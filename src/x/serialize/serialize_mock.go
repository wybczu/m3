// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/x/serialize (interfaces: TagEncoder,TagEncoderPool,TagDecoder,TagDecoderPool,MetricTagsIterator,MetricTagsIteratorPool)

// Copyright (c) 2021 Uber Technologies, Inc.
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

// Package serialize is a generated GoMock package.
package serialize

import (
	"reflect"

	"github.com/m3db/m3/src/x/checked"
	"github.com/m3db/m3/src/x/ident"

	"github.com/golang/mock/gomock"
)

// MockTagEncoder is a mock of TagEncoder interface.
type MockTagEncoder struct {
	ctrl     *gomock.Controller
	recorder *MockTagEncoderMockRecorder
}

// MockTagEncoderMockRecorder is the mock recorder for MockTagEncoder.
type MockTagEncoderMockRecorder struct {
	mock *MockTagEncoder
}

// NewMockTagEncoder creates a new mock instance.
func NewMockTagEncoder(ctrl *gomock.Controller) *MockTagEncoder {
	mock := &MockTagEncoder{ctrl: ctrl}
	mock.recorder = &MockTagEncoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagEncoder) EXPECT() *MockTagEncoderMockRecorder {
	return m.recorder
}

// Data mocks base method.
func (m *MockTagEncoder) Data() (checked.Bytes, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Data")
	ret0, _ := ret[0].(checked.Bytes)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Data indicates an expected call of Data.
func (mr *MockTagEncoderMockRecorder) Data() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockTagEncoder)(nil).Data))
}

// Encode mocks base method.
func (m *MockTagEncoder) Encode(arg0 ident.TagIterator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Encode indicates an expected call of Encode.
func (mr *MockTagEncoderMockRecorder) Encode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockTagEncoder)(nil).Encode), arg0)
}

// Finalize mocks base method.
func (m *MockTagEncoder) Finalize() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Finalize")
}

// Finalize indicates an expected call of Finalize.
func (mr *MockTagEncoderMockRecorder) Finalize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockTagEncoder)(nil).Finalize))
}

// Reset mocks base method.
func (m *MockTagEncoder) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockTagEncoderMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockTagEncoder)(nil).Reset))
}

// MockTagEncoderPool is a mock of TagEncoderPool interface.
type MockTagEncoderPool struct {
	ctrl     *gomock.Controller
	recorder *MockTagEncoderPoolMockRecorder
}

// MockTagEncoderPoolMockRecorder is the mock recorder for MockTagEncoderPool.
type MockTagEncoderPoolMockRecorder struct {
	mock *MockTagEncoderPool
}

// NewMockTagEncoderPool creates a new mock instance.
func NewMockTagEncoderPool(ctrl *gomock.Controller) *MockTagEncoderPool {
	mock := &MockTagEncoderPool{ctrl: ctrl}
	mock.recorder = &MockTagEncoderPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagEncoderPool) EXPECT() *MockTagEncoderPoolMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockTagEncoderPool) Get() TagEncoder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(TagEncoder)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockTagEncoderPoolMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTagEncoderPool)(nil).Get))
}

// Init mocks base method.
func (m *MockTagEncoderPool) Init() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init")
}

// Init indicates an expected call of Init.
func (mr *MockTagEncoderPoolMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockTagEncoderPool)(nil).Init))
}

// Put mocks base method.
func (m *MockTagEncoderPool) Put(arg0 TagEncoder) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Put", arg0)
}

// Put indicates an expected call of Put.
func (mr *MockTagEncoderPoolMockRecorder) Put(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTagEncoderPool)(nil).Put), arg0)
}

// MockTagDecoder is a mock of TagDecoder interface.
type MockTagDecoder struct {
	ctrl     *gomock.Controller
	recorder *MockTagDecoderMockRecorder
}

// MockTagDecoderMockRecorder is the mock recorder for MockTagDecoder.
type MockTagDecoderMockRecorder struct {
	mock *MockTagDecoder
}

// NewMockTagDecoder creates a new mock instance.
func NewMockTagDecoder(ctrl *gomock.Controller) *MockTagDecoder {
	mock := &MockTagDecoder{ctrl: ctrl}
	mock.recorder = &MockTagDecoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagDecoder) EXPECT() *MockTagDecoderMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTagDecoder) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockTagDecoderMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTagDecoder)(nil).Close))
}

// Current mocks base method.
func (m *MockTagDecoder) Current() ident.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Current")
	ret0, _ := ret[0].(ident.Tag)
	return ret0
}

// Current indicates an expected call of Current.
func (mr *MockTagDecoderMockRecorder) Current() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Current", reflect.TypeOf((*MockTagDecoder)(nil).Current))
}

// CurrentIndex mocks base method.
func (m *MockTagDecoder) CurrentIndex() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentIndex")
	ret0, _ := ret[0].(int)
	return ret0
}

// CurrentIndex indicates an expected call of CurrentIndex.
func (mr *MockTagDecoderMockRecorder) CurrentIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentIndex", reflect.TypeOf((*MockTagDecoder)(nil).CurrentIndex))
}

// Duplicate mocks base method.
func (m *MockTagDecoder) Duplicate() ident.TagIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Duplicate")
	ret0, _ := ret[0].(ident.TagIterator)
	return ret0
}

// Duplicate indicates an expected call of Duplicate.
func (mr *MockTagDecoderMockRecorder) Duplicate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Duplicate", reflect.TypeOf((*MockTagDecoder)(nil).Duplicate))
}

// Err mocks base method.
func (m *MockTagDecoder) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockTagDecoderMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockTagDecoder)(nil).Err))
}

// Len mocks base method.
func (m *MockTagDecoder) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockTagDecoderMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockTagDecoder)(nil).Len))
}

// Next mocks base method.
func (m *MockTagDecoder) Next() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockTagDecoderMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockTagDecoder)(nil).Next))
}

// Remaining mocks base method.
func (m *MockTagDecoder) Remaining() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remaining")
	ret0, _ := ret[0].(int)
	return ret0
}

// Remaining indicates an expected call of Remaining.
func (mr *MockTagDecoderMockRecorder) Remaining() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remaining", reflect.TypeOf((*MockTagDecoder)(nil).Remaining))
}

// Reset mocks base method.
func (m *MockTagDecoder) Reset(arg0 checked.Bytes) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset", arg0)
}

// Reset indicates an expected call of Reset.
func (mr *MockTagDecoderMockRecorder) Reset(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockTagDecoder)(nil).Reset), arg0)
}

// Rewind mocks base method.
func (m *MockTagDecoder) Rewind() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Rewind")
}

// Rewind indicates an expected call of Rewind.
func (mr *MockTagDecoderMockRecorder) Rewind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rewind", reflect.TypeOf((*MockTagDecoder)(nil).Rewind))
}

// MockTagDecoderPool is a mock of TagDecoderPool interface.
type MockTagDecoderPool struct {
	ctrl     *gomock.Controller
	recorder *MockTagDecoderPoolMockRecorder
}

// MockTagDecoderPoolMockRecorder is the mock recorder for MockTagDecoderPool.
type MockTagDecoderPoolMockRecorder struct {
	mock *MockTagDecoderPool
}

// NewMockTagDecoderPool creates a new mock instance.
func NewMockTagDecoderPool(ctrl *gomock.Controller) *MockTagDecoderPool {
	mock := &MockTagDecoderPool{ctrl: ctrl}
	mock.recorder = &MockTagDecoderPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagDecoderPool) EXPECT() *MockTagDecoderPoolMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockTagDecoderPool) Get() TagDecoder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(TagDecoder)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockTagDecoderPoolMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTagDecoderPool)(nil).Get))
}

// Init mocks base method.
func (m *MockTagDecoderPool) Init() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init")
}

// Init indicates an expected call of Init.
func (mr *MockTagDecoderPoolMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockTagDecoderPool)(nil).Init))
}

// Put mocks base method.
func (m *MockTagDecoderPool) Put(arg0 TagDecoder) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Put", arg0)
}

// Put indicates an expected call of Put.
func (mr *MockTagDecoderPoolMockRecorder) Put(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTagDecoderPool)(nil).Put), arg0)
}

// MockMetricTagsIterator is a mock of MetricTagsIterator interface.
type MockMetricTagsIterator struct {
	ctrl     *gomock.Controller
	recorder *MockMetricTagsIteratorMockRecorder
}

// MockMetricTagsIteratorMockRecorder is the mock recorder for MockMetricTagsIterator.
type MockMetricTagsIteratorMockRecorder struct {
	mock *MockMetricTagsIterator
}

// NewMockMetricTagsIterator creates a new mock instance.
func NewMockMetricTagsIterator(ctrl *gomock.Controller) *MockMetricTagsIterator {
	mock := &MockMetricTagsIterator{ctrl: ctrl}
	mock.recorder = &MockMetricTagsIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricTagsIterator) EXPECT() *MockMetricTagsIteratorMockRecorder {
	return m.recorder
}

// Bytes mocks base method.
func (m *MockMetricTagsIterator) Bytes() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes.
func (mr *MockMetricTagsIteratorMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockMetricTagsIterator)(nil).Bytes))
}

// Close mocks base method.
func (m *MockMetricTagsIterator) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockMetricTagsIteratorMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMetricTagsIterator)(nil).Close))
}

// Current mocks base method.
func (m *MockMetricTagsIterator) Current() ([]byte, []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Current")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	return ret0, ret1
}

// Current indicates an expected call of Current.
func (mr *MockMetricTagsIteratorMockRecorder) Current() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Current", reflect.TypeOf((*MockMetricTagsIterator)(nil).Current))
}

// Err mocks base method.
func (m *MockMetricTagsIterator) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockMetricTagsIteratorMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockMetricTagsIterator)(nil).Err))
}

// Next mocks base method.
func (m *MockMetricTagsIterator) Next() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockMetricTagsIteratorMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockMetricTagsIterator)(nil).Next))
}

// NumTags mocks base method.
func (m *MockMetricTagsIterator) NumTags() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumTags")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumTags indicates an expected call of NumTags.
func (mr *MockMetricTagsIteratorMockRecorder) NumTags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumTags", reflect.TypeOf((*MockMetricTagsIterator)(nil).NumTags))
}

// Reset mocks base method.
func (m *MockMetricTagsIterator) Reset(arg0 []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset", arg0)
}

// Reset indicates an expected call of Reset.
func (mr *MockMetricTagsIteratorMockRecorder) Reset(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockMetricTagsIterator)(nil).Reset), arg0)
}

// TagValue mocks base method.
func (m *MockMetricTagsIterator) TagValue(arg0 []byte) ([]byte, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TagValue", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// TagValue indicates an expected call of TagValue.
func (mr *MockMetricTagsIteratorMockRecorder) TagValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TagValue", reflect.TypeOf((*MockMetricTagsIterator)(nil).TagValue), arg0)
}

// MockMetricTagsIteratorPool is a mock of MetricTagsIteratorPool interface.
type MockMetricTagsIteratorPool struct {
	ctrl     *gomock.Controller
	recorder *MockMetricTagsIteratorPoolMockRecorder
}

// MockMetricTagsIteratorPoolMockRecorder is the mock recorder for MockMetricTagsIteratorPool.
type MockMetricTagsIteratorPoolMockRecorder struct {
	mock *MockMetricTagsIteratorPool
}

// NewMockMetricTagsIteratorPool creates a new mock instance.
func NewMockMetricTagsIteratorPool(ctrl *gomock.Controller) *MockMetricTagsIteratorPool {
	mock := &MockMetricTagsIteratorPool{ctrl: ctrl}
	mock.recorder = &MockMetricTagsIteratorPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricTagsIteratorPool) EXPECT() *MockMetricTagsIteratorPoolMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMetricTagsIteratorPool) Get() MetricTagsIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(MetricTagsIterator)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockMetricTagsIteratorPoolMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMetricTagsIteratorPool)(nil).Get))
}

// Init mocks base method.
func (m *MockMetricTagsIteratorPool) Init() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init")
}

// Init indicates an expected call of Init.
func (mr *MockMetricTagsIteratorPoolMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockMetricTagsIteratorPool)(nil).Init))
}

// Put mocks base method.
func (m *MockMetricTagsIteratorPool) Put(arg0 MetricTagsIterator) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Put", arg0)
}

// Put indicates an expected call of Put.
func (mr *MockMetricTagsIteratorPoolMockRecorder) Put(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockMetricTagsIteratorPool)(nil).Put), arg0)
}
