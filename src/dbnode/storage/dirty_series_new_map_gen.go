// Copyright (c) 2019 Uber Technologies, Inc.
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

// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package storage

import (
	"github.com/m3db/m3/src/x/ident"
	"github.com/m3db/m3/src/x/pool"

	"github.com/cespare/xxhash"
)

// dirtySeriesMapOptions provides options used when created the map.
type dirtySeriesMapOptions struct {
	InitialSize int
	KeyCopyPool pool.BytesPool
}

// newDirtySeriesMap returns a new byte keyed map.
func newDirtySeriesMap(opts dirtySeriesMapOptions) *dirtySeriesMap {
	var (
		copyFn     dirtySeriesMapCopyFn
		finalizeFn dirtySeriesMapFinalizeFn
	)
	if pool := opts.KeyCopyPool; pool == nil {
		copyFn = func(k idAndBlockStart) idAndBlockStart {
			return idAndBlockStart{
				id:         ident.BytesID(append([]byte(nil), k.id.Bytes()...)),
				blockStart: k.blockStart,
			}
		}
	} else {
		copyFn = func(k idAndBlockStart) idAndBlockStart {
			bytes := k.id.Bytes()
			keyLen := len(bytes)
			pooled := pool.Get(keyLen)[:keyLen]
			copy(pooled, bytes)
			return idAndBlockStart{
				id:         ident.BytesID(pooled),
				blockStart: k.blockStart,
			}
		}
		finalizeFn = func(k idAndBlockStart) {
			if slice, ok := k.id.(ident.BytesID); ok {
				pool.Put(slice)
			}
		}
	}
	return _dirtySeriesMapAlloc(_dirtySeriesMapOptions{
		hash: func(k idAndBlockStart) dirtySeriesMapHash {
			hash := uint64(7)
			hash = 31*hash + xxhash.Sum64(k.id.Bytes())
			hash = 31*hash + uint64(k.blockStart)
			return dirtySeriesMapHash(hash)
		},
		equals: func(x, y idAndBlockStart) bool {
			return x.id.Equal(y.id) && x.blockStart == y.blockStart
		},
		copy:        copyFn,
		finalize:    finalizeFn,
		initialSize: opts.InitialSize,
	})
}