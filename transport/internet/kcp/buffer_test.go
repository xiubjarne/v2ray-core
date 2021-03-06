package kcp_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
)

func TestBuffer(t *testing.T) {
	assert := assert.On(t)

	b := NewBuffer()

	for i := 0; i < NumDistro; i++ {
		x := b.Allocate()
		assert.Pointer(x).IsNotNil()
		x.Release()
	}
	assert.Pointer(b.Allocate()).IsNil()
	b.Release()
}

func TestSingleRelease(t *testing.T) {
	assert := assert.On(t)

	b := NewBuffer()
	x := b.Allocate()
	x.Release()

	y := b.Allocate()
	assert.Pointer(y.Value).IsNotNil()

	b.Release()
	y.Release()

	z := b.Allocate()
	assert.Pointer(z).IsNil()
}
