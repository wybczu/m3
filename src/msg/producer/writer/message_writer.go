// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package writer

import (
	"container/list"
	"errors"
	"math"
	"sync"
	"time"
	stdunsafe "unsafe"

	"github.com/uber-go/tally"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/m3db/m3/src/msg/producer"
	"github.com/m3db/m3/src/msg/protocol/proto"
	"github.com/m3db/m3/src/x/clock"
	"github.com/m3db/m3/src/x/instrument"
	"github.com/m3db/m3/src/x/retry"
	"github.com/m3db/m3/src/x/unsafe"
)

// MessageRetryNanosFn returns the message backoff time for retry in nanoseconds.
type MessageRetryNanosFn func(writeTimes int) int64

var (
	errInvalidBackoffDuration = errors.New("invalid backoff duration")
	errFailAllConsumers       = errors.New("could not write to any consumer")
	errNoWriters              = errors.New("no writers")
)

const _recordMessageDelayEvery = 4 // keep it a power of two value to keep modulo fast

type messageWriterMetrics struct {
	withoutConsumerScope       bool
	scope                      tally.Scope
	opts                       instrument.TimerOptions
	writeSuccess               tally.Counter
	allConsumersWriteError     tally.Counter
	noWritersError             tally.Counter
	writeAfterCutoff           tally.Counter
	writeBeforeCutover         tally.Counter
	messageAcked               tally.Counter
	messageClosed              tally.Counter
	messageDroppedBufferFull   tally.Counter
	messageDroppedTTLExpire    tally.Counter
	messageRetry               tally.Counter
	messageConsumeLatency      tally.Timer
	messageWriteDelay          tally.Timer
	scanBatchLatency           tally.Timer
	scanTotalLatency           tally.Timer
	writeSuccessLatency        tally.Histogram
	writeErrorLatency          tally.Histogram
	enqueuedMessages           tally.Counter
	dequeuedMessages           tally.Counter
	processedWrite             tally.Counter
	processedClosed            tally.Counter
	processedNotReady          tally.Counter
	processedTTL               tally.Counter
	processedAck               tally.Counter
	processedDrop              tally.Counter
	forcedFlush                tally.Counter
	forcedFlushTimeout         tally.Counter
	forcedFlushFailedOne       tally.Counter
	forcedFlushFailedAll       tally.Counter
	forcedFlushLatency         tally.Histogram
	forcedFlushSingleConsumer  tally.Counter
	forcedFlushNotEnoughBuffer tally.Counter
}

func (m *messageWriterMetrics) withConsumer(consumer string) *messageWriterMetrics {
	if m.withoutConsumerScope {
		return m
	}
	return newMessageWriterMetricsWithConsumer(m.scope, m.opts, consumer, false)
}

func newMessageWriterMetrics(
	scope tally.Scope,
	opts instrument.TimerOptions,
	withoutConsumerScope bool,
) *messageWriterMetrics {
	return newMessageWriterMetricsWithConsumer(scope, opts, "unknown", withoutConsumerScope)
}

func newMessageWriterMetricsWithConsumer(
	scope tally.Scope,
	opts instrument.TimerOptions,
	consumer string,
	withoutConsumerScope bool,
) *messageWriterMetrics {
	consumerScope := scope
	if !withoutConsumerScope {
		consumerScope = scope.Tagged(map[string]string{"consumer": consumer})
	}
	return &messageWriterMetrics{
		withoutConsumerScope: withoutConsumerScope,
		scope:                scope,
		opts:                 opts,
		writeSuccess:         consumerScope.Counter("write-success"),
		allConsumersWriteError: consumerScope.
			Tagged(map[string]string{"error-type": "all-consumers"}).
			Counter("write-error"),
		noWritersError: consumerScope.
			Tagged(map[string]string{"error-type": "no-writers"}).
			Counter("write-error"),
		writeAfterCutoff: consumerScope.
			Tagged(map[string]string{"reason": "after-cutoff"}).
			Counter("invalid-write"),
		writeBeforeCutover: consumerScope.
			Tagged(map[string]string{"reason": "before-cutover"}).
			Counter("invalid-write"),
		messageAcked:  consumerScope.Counter("message-acked"),
		messageClosed: consumerScope.Counter("message-closed"),
		messageDroppedBufferFull: consumerScope.Tagged(
			map[string]string{"reason": "buffer-full"},
		).Counter("message-dropped"),
		messageDroppedTTLExpire: consumerScope.Tagged(
			map[string]string{"reason": "ttl-expire"},
		).Counter("message-dropped"),
		messageRetry:          consumerScope.Counter("message-retry"),
		messageConsumeLatency: instrument.NewTimer(consumerScope, "message-consume-latency", opts),
		messageWriteDelay:     instrument.NewTimer(consumerScope, "message-write-delay", opts),
		scanBatchLatency:      instrument.NewTimer(consumerScope, "scan-batch-latency", opts),
		scanTotalLatency:      instrument.NewTimer(consumerScope, "scan-total-latency", opts),
		writeSuccessLatency: consumerScope.Histogram("write-success-latency",
			tally.MustMakeExponentialDurationBuckets(time.Millisecond*10, 2, 15)),
		writeErrorLatency: consumerScope.Histogram("write-error-latency",
			tally.MustMakeExponentialDurationBuckets(time.Millisecond*10, 2, 15)),
		enqueuedMessages: consumerScope.Counter("message-enqueue"),
		dequeuedMessages: consumerScope.Counter("message-dequeue"),
		processedWrite: consumerScope.
			Tagged(map[string]string{"result": "write"}).
			Counter("message-processed"),
		processedClosed: consumerScope.
			Tagged(map[string]string{"result": "closed"}).
			Counter("message-processed"),
		processedNotReady: consumerScope.
			Tagged(map[string]string{"result": "not-ready"}).
			Counter("message-processed"),
		processedTTL: consumerScope.
			Tagged(map[string]string{"result": "ttl"}).
			Counter("message-processed"),
		processedAck: consumerScope.
			Tagged(map[string]string{"result": "ack"}).
			Counter("message-processed"),
		processedDrop: consumerScope.
			Tagged(map[string]string{"result": "drop"}).
			Counter("message-processed"),
		forcedFlush:          consumerScope.Counter("forced-flush"),
		forcedFlushTimeout:   consumerScope.Counter("forced-flush-timeout"),
		forcedFlushFailedOne: consumerScope.Counter("forced-flush-failed-one"),
		forcedFlushFailedAll: consumerScope.Counter("forced-flush-failed-all"),
		forcedFlushLatency: consumerScope.Histogram(
			"forced-flush-latency",
			tally.MustMakeExponentialDurationBuckets(time.Millisecond*10, 2, 15),
		),
		forcedFlushSingleConsumer:  consumerScope.Counter("forced-flush-single-consumer"),
		forcedFlushNotEnoughBuffer: consumerScope.Counter("forced-flush-not-enough-buffer"),
	}
}

type messageWriter struct {
	sync.RWMutex

	replicatedShardID   uint64
	mPool               *messagePool
	opts                Options
	nextRetryAfterNanos MessageRetryNanosFn
	encoder             proto.Encoder
	numConnections      int

	msgID            uint64
	queue            *list.List
	consumerWriters  []consumerWriter
	iterationIndexes []int
	acks             *acks
	cutOffNanos      int64
	cutOverNanos     int64
	messageTTLNanos  int64
	msgsToWrite      []*message
	isClosed         bool
	doneCh           chan struct{}
	wg               sync.WaitGroup
	// metrics can be updated when a consumer instance changes, so must be guarded with RLock
	metrics      atomic.UnsafePointer //  *messageWriterMetrics
	nextFullScan time.Time
	lastNewWrite *list.Element

	nowFn clock.NowFn
}

func newMessageWriter(
	replicatedShardID uint64,
	mPool *messagePool,
	opts Options,
	m *messageWriterMetrics,
) *messageWriter {
	if opts == nil {
		opts = NewOptions()
	}
	nowFn := time.Now
	mw := &messageWriter{
		replicatedShardID:   replicatedShardID,
		mPool:               mPool,
		opts:                opts,
		nextRetryAfterNanos: opts.MessageRetryNanosFn(),
		encoder:             proto.NewEncoder(opts.EncoderOptions()),
		numConnections:      opts.ConnectionOptions().NumConnections(),
		msgID:               0,
		queue:               list.New(),
		acks:                newAckHelper(opts.InitialAckMapSize()),
		cutOffNanos:         0,
		cutOverNanos:        0,
		msgsToWrite:         make([]*message, 0, opts.MessageQueueScanBatchSize()),
		isClosed:            false,
		doneCh:              make(chan struct{}),
		nowFn:               nowFn,
	}
	mw.metrics.Store(stdunsafe.Pointer(m))
	return mw
}

// Write writes a message, messages not acknowledged in time will be retried.
// New messages will be written in order, but retries could be out of order.
func (w *messageWriter) Write(rm *producer.RefCountedMessage) {
	var (
		nowNanos = w.nowFn().UnixNano()
		msg      = w.newMessage()
		metrics  = w.Metrics()
	)
	w.Lock()
	if !w.isValidWriteWithLock(nowNanos, metrics) {
		w.Unlock()
		w.close(msg)
		return
	}
	rm.IncRef()
	w.msgID++
	meta := metadata{
		metadataKey: metadataKey{
			shard: w.replicatedShardID,
			id:    w.msgID,
		},
	}
	msg.Set(meta, rm, nowNanos)
	w.acks.add(meta, msg)
	// Make sure all the new writes are ordered in queue.
	metrics.enqueuedMessages.Inc(1)
	if w.lastNewWrite != nil {
		w.lastNewWrite = w.queue.InsertAfter(msg, w.lastNewWrite)
	} else {
		w.lastNewWrite = w.queue.PushFront(msg)
	}
	w.Unlock()
}

func (w *messageWriter) isValidWriteWithLock(nowNanos int64, metrics *messageWriterMetrics) bool {
	if w.opts.IgnoreCutoffCutover() {
		return true
	}

	if w.cutOffNanos > 0 && nowNanos >= w.cutOffNanos {
		metrics.writeAfterCutoff.Inc(1)
		return false
	}
	if w.cutOverNanos > 0 && nowNanos < w.cutOverNanos {
		metrics.writeBeforeCutover.Inc(1)
		return false
	}

	return true
}

func (w *messageWriter) write(
	consumerWriters []consumerWriter,
	metrics *messageWriterMetrics,
	m *message,
) error {
	m.IncReads()
	m.SetSentAt(w.nowFn().UnixNano())
	msg, isValid := m.Marshaler()
	if !isValid {
		m.DecReads()
		return nil
	}
	// The write function is accessed through only one thread,
	// so no lock is required for encoding.
	err := w.encoder.Encode(msg)
	m.DecReads()
	if err != nil {
		return err
	}
	var (
		// NB(r): Always select the same connection index per shard.
		connIndex = int(w.replicatedShardID % uint64(w.numConnections))
		writeData = w.encoder.Bytes()
	)

	cw := w.chooseConsumerWriter(
		consumerWriters,
		connIndex,
		len(writeData),
	)

	start := w.nowFn().UnixNano()
	if err := cw.Write(connIndex, writeData); err != nil {
		metrics.writeErrorLatency.RecordDuration(time.Duration(w.nowFn().UnixNano() - start))
		metrics.allConsumersWriteError.Inc(1)
		return errFailAllConsumers
	}

	metrics.writeSuccess.Inc(1)
	return nil
}

// Ack acknowledges the metadata.
func (w *messageWriter) Ack(meta metadata) bool {
	if acked, expectedProcessNanos := w.acks.ack(meta); acked {
		m := w.Metrics()
		m.messageConsumeLatency.Record(time.Duration(w.nowFn().UnixNano() - expectedProcessNanos))
		m.messageAcked.Inc(1)
		return true
	}
	return false
}

// Init initialize the message writer.
func (w *messageWriter) Init() {
	w.wg.Add(1)
	go func() {
		w.scanMessageQueueUntilClose()
		w.wg.Done()
	}()
}

func (w *messageWriter) scanMessageQueueUntilClose() {
	var (
		interval = w.opts.MessageQueueNewWritesScanInterval()
		jitter   = time.Duration(
			// approx ~40 days of jitter at millisecond precision - more than enough
			unsafe.Fastrandn(uint32(interval.Milliseconds())),
		) * time.Millisecond
	)
	// NB(cw): Add some jitter before the tick starts to reduce
	// some contention between all the message writers.
	time.Sleep(jitter)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.scanMessageQueue()
		case <-w.doneCh:
			return
		}
	}
}

func (w *messageWriter) scanMessageQueue() {
	w.RLock()
	e := w.queue.Front()
	w.lastNewWrite = nil
	isClosed := w.isClosed
	w.RUnlock()

	var (
		nowFn            = w.nowFn
		msgsToWrite      []*message
		beforeScan       = nowFn()
		beforeBatchNanos = beforeScan.UnixNano()
		batchSize        = w.opts.MessageQueueScanBatchSize()
		consumerWriters  []consumerWriter
		fullScan         = isClosed || beforeScan.After(w.nextFullScan)
		m                = w.Metrics()
		scanMetrics      scanBatchMetrics
		skipWrites       bool
	)
	defer scanMetrics.record(m)
	for e != nil {
		w.Lock()
		e, msgsToWrite = w.scanBatchWithLock(e, beforeBatchNanos, batchSize, fullScan, &scanMetrics)
		consumerWriters = w.consumerWriters
		w.Unlock()
		if !fullScan && len(msgsToWrite) == 0 {
			m.scanBatchLatency.Record(time.Duration(nowFn().UnixNano() - beforeBatchNanos))
			// If this is not a full scan, abort after the iteration batch
			// that no new messages were found.
			break
		}
		if skipWrites {
			m.scanBatchLatency.Record(time.Duration(nowFn().UnixNano() - beforeBatchNanos))
			continue
		}
		if err := w.writeBatch(consumerWriters, m, msgsToWrite); err != nil {
			// When we can't write to any consumer writer, skip the writes in this scan
			// to avoid meaningless attempts but continue to clean up the queue.
			skipWrites = true
		}
		nowNanos := nowFn().UnixNano()
		m.scanBatchLatency.Record(time.Duration(nowNanos - beforeBatchNanos))
		beforeBatchNanos = nowNanos
	}
	afterScan := nowFn()
	m.scanTotalLatency.Record(afterScan.Sub(beforeScan))
	if fullScan {
		w.nextFullScan = afterScan.Add(w.opts.MessageQueueFullScanInterval())
	}
}

func (w *messageWriter) writeBatch(
	consumerWriters []consumerWriter,
	metrics *messageWriterMetrics,
	messages []*message,
) error {
	if len(consumerWriters) == 0 {
		// Not expected in a healthy/valid placement.
		metrics.noWritersError.Inc(int64(len(messages)))
		return errNoWriters
	}
	delay := metrics.messageWriteDelay
	nowFn := w.nowFn
	for i := range messages {
		start := nowFn().UnixNano()
		if err := w.write(consumerWriters, metrics, messages[i]); err != nil {
			return err
		}
		if i%_recordMessageDelayEvery == 0 {
			now := nowFn().Unix()
			delay.Record(time.Duration(now - messages[i].ExpectedProcessAtNanos()))
			metrics.writeSuccessLatency.RecordDuration(time.Duration(now - start))
		}
	}
	return nil
}

// scanBatchWithLock iterates the message queue with a lock. It returns after
// visited enough elements. So it holds the lock for less time and allows new
// writes to be unblocked.
func (w *messageWriter) scanBatchWithLock(
	start *list.Element,
	nowNanos int64,
	batchSize int,
	fullScan bool,
	scanMetrics *scanBatchMetrics,
) (*list.Element, []*message) {
	var (
		iterated int
		next     *list.Element
	)
	metrics := w.Metrics()
	w.msgsToWrite = w.msgsToWrite[:0]
	for e := start; e != nil; e = next {
		iterated++
		if iterated > batchSize {
			break
		}
		next = e.Next()
		m := e.Value.(*message)
		if w.isClosed {
			scanMetrics[_processedClosed]++
			// Simply ack the messages here to mark them as consumed for this
			// message writer, this is useful when user removes a consumer service
			// during runtime that may be unhealthy to consume the messages.
			// So that the unacked messages for the unhealthy consumer services
			// do not stay in memory forever.
			// NB: The message must be added to the ack map to be acked here.
			w.acks.ack(m.Metadata())
			w.removeFromQueueWithLock(e, m, metrics)
			scanMetrics[_messageClosed]++
			continue
		}
		if m.RetryAtNanos() >= nowNanos {
			scanMetrics[_processedNotReady]++
			if !fullScan {
				// If this is not a full scan, bail after the first element that
				// is not a new write.
				break
			}
			continue
		}
		// If the message exceeded its allowed ttl of the consumer service,
		// remove it from the buffer.
		if w.messageTTLNanos > 0 && m.InitNanos()+w.messageTTLNanos <= nowNanos {
			scanMetrics[_processedTTL]++
			// There is a chance the message was acked right before the ack is
			// called, in which case just remove it from the queue.
			if acked, _ := w.acks.ack(m.Metadata()); acked {
				scanMetrics[_messageDroppedTTLExpire]++
			}
			w.removeFromQueueWithLock(e, m, metrics)
			continue
		}
		if m.IsAcked() {
			scanMetrics[_processedAck]++
			w.removeFromQueueWithLock(e, m, metrics)
			continue
		}
		if m.IsDroppedOrConsumed() {
			scanMetrics[_processedDrop]++
			// There is a chance the message could be acked between m.Acked()
			// and m.IsDroppedOrConsumed() check, in which case we should not
			// mark it as dropped, just continue and next tick will remove it
			// as acked.
			if m.IsAcked() {
				continue
			}
			w.acks.remove(m.Metadata())
			w.removeFromQueueWithLock(e, m, metrics)
			scanMetrics[_messageDroppedBufferFull]++
			continue
		}
		m.IncWriteTimes()
		writeTimes := m.WriteTimes()
		m.SetRetryAtNanos(w.nextRetryAfterNanos(writeTimes) + nowNanos)
		if writeTimes > 1 {
			scanMetrics[_messageRetry]++
		}
		scanMetrics[_processedWrite]++
		w.msgsToWrite = append(w.msgsToWrite, m)
	}
	return next, w.msgsToWrite
}

// Close closes the writer.
// It should block until all buffered messages have been acknowledged.
func (w *messageWriter) Close() {
	w.Lock()
	if w.isClosed {
		w.Unlock()
		return
	}
	w.isClosed = true
	w.Unlock()
	// NB: Wait until all messages cleaned up then close.
	w.waitUntilAllMessageRemoved()
	close(w.doneCh)
	w.wg.Wait()
}

func (w *messageWriter) waitUntilAllMessageRemoved() {
	// The message writers are being closed sequentially, checking isEmpty()
	// before always waiting for the first tick can speed up Close() a lot.
	if w.isEmpty() {
		return
	}
	ticker := time.NewTicker(w.opts.CloseCheckInterval())
	defer ticker.Stop()

	for range ticker.C {
		if w.isEmpty() {
			return
		}
	}
}

func (w *messageWriter) isEmpty() bool {
	w.RLock()
	l := w.queue.Len()
	w.RUnlock()
	return l == 0
}

// ReplicatedShardID returns the replicated shard id.
func (w *messageWriter) ReplicatedShardID() uint64 {
	return w.replicatedShardID
}

func (w *messageWriter) CutoffNanos() int64 {
	w.RLock()
	res := w.cutOffNanos
	w.RUnlock()
	return res
}

func (w *messageWriter) SetCutoffNanos(nanos int64) {
	w.Lock()
	w.cutOffNanos = nanos
	w.Unlock()
}

func (w *messageWriter) CutoverNanos() int64 {
	w.RLock()
	res := w.cutOverNanos
	w.RUnlock()
	return res
}

func (w *messageWriter) SetCutoverNanos(nanos int64) {
	w.Lock()
	w.cutOverNanos = nanos
	w.Unlock()
}

func (w *messageWriter) MessageTTLNanos() int64 {
	w.RLock()
	res := w.messageTTLNanos
	w.RUnlock()
	return res
}

func (w *messageWriter) SetMessageTTLNanos(nanos int64) {
	w.Lock()
	w.messageTTLNanos = nanos
	w.Unlock()
}

// AddConsumerWriter adds a consumer writer.
func (w *messageWriter) AddConsumerWriter(cw consumerWriter) {
	w.Lock()
	newConsumerWriters := make([]consumerWriter, 0, len(w.consumerWriters)+1)
	newConsumerWriters = append(newConsumerWriters, w.consumerWriters...)
	newConsumerWriters = append(newConsumerWriters, cw)

	w.iterationIndexes = make([]int, len(newConsumerWriters))
	for i := range w.iterationIndexes {
		w.iterationIndexes[i] = i
	}
	w.consumerWriters = newConsumerWriters
	w.Unlock()
}

// RemoveConsumerWriter removes the consumer writer for the given address.
func (w *messageWriter) RemoveConsumerWriter(addr string) {
	w.Lock()
	newConsumerWriters := make([]consumerWriter, 0, len(w.consumerWriters)-1)
	for _, cw := range w.consumerWriters {
		if cw.Address() == addr {
			continue
		}
		newConsumerWriters = append(newConsumerWriters, cw)
	}

	w.iterationIndexes = make([]int, len(newConsumerWriters))
	for i := range w.iterationIndexes {
		w.iterationIndexes[i] = i
	}
	w.consumerWriters = newConsumerWriters
	w.Unlock()
}

// Metrics returns the metrics. These are dynamic and change if downstream consumer instance changes.
func (w *messageWriter) Metrics() *messageWriterMetrics {
	return (*messageWriterMetrics)(w.metrics.Load())
}

// SetMetrics sets the metrics
//
// This allows changing the labels of the metrics when the downstream consumer instance changes.
func (w *messageWriter) SetMetrics(m *messageWriterMetrics) {
	w.metrics.Store(stdunsafe.Pointer(m))
}

// QueueSize returns the number of messages queued in the writer.
func (w *messageWriter) QueueSize() int {
	return w.acks.size()
}

func (w *messageWriter) newMessage() *message {
	return w.mPool.Get()
}

func (w *messageWriter) removeFromQueueWithLock(e *list.Element, m *message, metrics *messageWriterMetrics) {
	w.queue.Remove(e)
	metrics.dequeuedMessages.Inc(1)
	w.close(m)
}

func (w *messageWriter) close(m *message) {
	m.Close()
	w.mPool.Put(m)
}

type acks struct {
	mtx  sync.Mutex
	acks map[uint64]*message
}

// nolint: unparam
func newAckHelper(size int) *acks {
	return &acks{
		acks: make(map[uint64]*message, size),
	}
}

func (a *acks) add(meta metadata, m *message) {
	a.mtx.Lock()
	a.acks[meta.metadataKey.id] = m
	a.mtx.Unlock()
}

func (a *acks) remove(meta metadata) {
	a.mtx.Lock()
	delete(a.acks, meta.metadataKey.id)
	a.mtx.Unlock()
}

// ack processes the ack. returns true if the message was not already acked. additionally returns the expected
// processing time for lag calculations.
func (a *acks) ack(meta metadata) (bool, int64) {
	a.mtx.Lock()
	m, ok := a.acks[meta.metadataKey.id]
	if !ok {
		a.mtx.Unlock()
		// Acking a message that is already acked, which is ok.
		return false, 0
	}

	delete(a.acks, meta.metadataKey.id)
	a.mtx.Unlock()

	expectedProcessAtNanos := m.ExpectedProcessAtNanos()
	m.Ack()

	return true, expectedProcessAtNanos
}

func (a *acks) size() int {
	a.mtx.Lock()
	l := len(a.acks)
	a.mtx.Unlock()
	return l
}

type metricIdx byte

const (
	_messageClosed metricIdx = iota
	_messageDroppedBufferFull
	_messageDroppedTTLExpire
	_messageRetry
	_processedAck
	_processedClosed
	_processedDrop
	_processedNotReady
	_processedTTL
	_processedWrite
	_lastMetricIdx
)

type scanBatchMetrics [_lastMetricIdx]int32

func (m *scanBatchMetrics) record(metrics *messageWriterMetrics) {
	m.recordNonzeroCounter(_messageClosed, metrics.messageClosed)
	m.recordNonzeroCounter(_messageDroppedBufferFull, metrics.messageDroppedBufferFull)
	m.recordNonzeroCounter(_messageDroppedTTLExpire, metrics.messageDroppedTTLExpire)
	m.recordNonzeroCounter(_messageRetry, metrics.messageRetry)
	m.recordNonzeroCounter(_processedAck, metrics.processedAck)
	m.recordNonzeroCounter(_processedClosed, metrics.processedClosed)
	m.recordNonzeroCounter(_processedDrop, metrics.processedDrop)
	m.recordNonzeroCounter(_processedNotReady, metrics.processedNotReady)
	m.recordNonzeroCounter(_processedTTL, metrics.processedTTL)
	m.recordNonzeroCounter(_processedWrite, metrics.processedWrite)
}

func (m *scanBatchMetrics) recordNonzeroCounter(idx metricIdx, c tally.Counter) {
	if m[idx] > 0 {
		c.Inc(int64(m[idx]))
	}
}

// NextRetryNanosFn creates a MessageRetryNanosFn based on the retry options.
func NextRetryNanosFn(retryOpts retry.Options) func(int) int64 {
	var (
		jitter              = retryOpts.Jitter()
		backoffFactor       = retryOpts.BackoffFactor()
		initialBackoff      = retryOpts.InitialBackoff()
		maxBackoff          = retryOpts.MaxBackoff()
		initialBackoffFloat = float64(initialBackoff)
	)

	// inlined and specialized retry function that does not have any state that needs to be kept
	// between tries
	return func(writeTimes int) int64 {
		backoff := initialBackoff.Nanoseconds()
		if writeTimes >= 1 {
			backoffFloat64 := initialBackoffFloat * math.Pow(backoffFactor, float64(writeTimes-1))
			backoff = int64(backoffFloat64)
		}
		// Validate the value of backoff to make sure Fastrandn() does not panic and
		// check for overflow from the exponentiation op - unlikely, but prevents weird side effects.
		halfInMicros := (backoff / 2) / int64(time.Microsecond)
		if jitter && backoff >= 2 && halfInMicros < math.MaxUint32 {
			// Jitter can be only up to ~1 hour in microseconds, but it's not a limitation here
			jitterInMicros := unsafe.Fastrandn(uint32(halfInMicros))
			jitterInNanos := time.Duration(jitterInMicros) * time.Microsecond
			halfInNanos := time.Duration(halfInMicros) * time.Microsecond
			backoff = int64(halfInNanos + jitterInNanos)
		}
		// Clamp backoff to maxBackoff
		if maxBackoff := maxBackoff.Nanoseconds(); backoff > maxBackoff {
			backoff = maxBackoff
		}
		return backoff
	}
}

// StaticRetryNanosFn creates a MessageRetryNanosFn based on static config.
func StaticRetryNanosFn(backoffDurations []time.Duration) (MessageRetryNanosFn, error) {
	if len(backoffDurations) == 0 {
		return nil, errInvalidBackoffDuration
	}
	backoffInt64s := make([]int64, 0, len(backoffDurations))
	for _, b := range backoffDurations {
		backoffInt64s = append(backoffInt64s, int64(b))
	}
	return func(writeTimes int) int64 {
		retry := writeTimes - 1
		l := len(backoffInt64s)
		if retry < l {
			return backoffInt64s[retry]
		}
		return backoffInt64s[l-1]
	}, nil
}

func (w *messageWriter) chooseConsumerWriter(
	consumerWriters []consumerWriter,
	connIndex int,
	writeLen int,
) consumerWriter {
	if len(consumerWriters) == 1 {
		w.Metrics().forcedFlushSingleConsumer.Inc(1)
		return consumerWriters[0]
	}

	// find the consumer writer with the max available buffer.
	max, maxBuf := w.getConsumerWriterWithMaxBuffer(consumerWriters, connIndex)

	// if the available buffer is able to accommodate the write, return the consumer writer.
	// This means that the consumer writer will not be blocked on the write.
	if maxBuf >= writeLen {
		return max
	}

	m := w.Metrics()
	m.forcedFlush.Inc(1)

	startTs := w.nowFn().UnixNano()
	// Since we are not able to find a consumer writer that can accommodate the write,
	// we initiate a forced flush on all available the consumer writers.
	// The first one to return will be the chosen as the least loaded consumer writer.
	// Note that doing a forced operation on all consumer writers is fine since, a Write()
	// will anyway invoke a forced Flush(). But the downside of simply invoking a write
	// is that the entire consumer writer will be blocked in that process.
	// Therefore it makes sense to initiate a forced Flush() on all available consumer
	// writers and wait for the first one to return. This way, we can utilize the connections
	// to the replicas if available in a more efficient manner.
	doneCh := make(chan int, len(consumerWriters))
	// intentionally leave the doneCh open to avoid panics in case a forcedFlush finishes afte
	// this function returns.
	w.beginForcedFlush(doneCh, consumerWriters, connIndex)

	// wait for first consumer writer to finish.
	cw := w.waitForForcedFlush(doneCh, consumerWriters)
	if cw != nil {
		max = cw
		if cw.AvailableBuffer(connIndex) < writeLen {
			// The consumer writer should have enough buffer to accommodate the write.
			// if not, log and emit a metric.
			m.forcedFlushNotEnoughBuffer.Inc(1)
			w.opts.InstrumentOptions().Logger().Info(
				"forced flush, still not enough buffer",
				zap.String("consumer", cw.Address()),
			)
		}
	}

	m.forcedFlushLatency.RecordDuration(time.Duration(w.nowFn().UnixNano() - startTs))

	// return the consumer writer with the max buffer or the consumer writer that
	// returned first from the forced flush operation.
	return max
}

func (w *messageWriter) beginForcedFlush(
	doneCh chan<- int,
	consumerWriters []consumerWriter,
	connIndex int,
) {
	m := w.Metrics()
	for i := range consumerWriters {
		i := i
		go func(idx int) {
			if err := consumerWriters[idx].ForcedFlush(connIndex); err != nil {
				m.forcedFlushFailedOne.Inc(1)
				doneCh <- -1
				return
			}
			doneCh <- idx
		}(i)
	}
}

func (w *messageWriter) getConsumerWriterWithMaxBuffer(
	consumerWriters []consumerWriter,
	connIndex int,
) (consumerWriter, int) {
	max := consumerWriters[0]
	maxBufSize := consumerWriters[0].AvailableBuffer(connIndex)
	for i := 1; i < len(consumerWriters); i++ {
		bufSize := consumerWriters[i].AvailableBuffer(connIndex)
		if bufSize > maxBufSize {
			max = consumerWriters[i]
			maxBufSize = bufSize
		}
	}

	return max, maxBufSize
}

// waitForForcedFlush returns the first consumerWriter to complete
// the forced flush operation or nil if all consumer writers failed / timed out.
func (w *messageWriter) waitForForcedFlush(
	doneCh <-chan int,
	consumerWriters []consumerWriter,
) consumerWriter {
	var cw consumerWriter
	m := w.Metrics()
	// wait for the first consumer writer to return.
	// In case both the consumer writers are blocked for more than forcedFlushTimeout time,
	// we will short circuit and return nil.
	t := time.NewTicker(w.opts.ConnectionOptions().ForcedFlushTimeout())
	defer t.Stop()

waitLoop:
	for range len(consumerWriters) {
		select {
		case idx := <-doneCh:
			if idx == -1 {
				// received an error from a consumer writer.
				// wait for success or failure from the rest.
				continue waitLoop
			}
			cw = consumerWriters[idx]
			break waitLoop // break from the loop as soon as we get the first consumer writer to return.
		case <-t.C:
			// if no consumer writer returns within the timeout, return the max consumer writer.
			m.forcedFlushTimeout.Inc(1)
			break waitLoop
		}
	}

	if cw == nil {
		m.forcedFlushFailedAll.Inc(1)
	}

	return cw
}
