package syncpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LockPool_AllocAndFree(t *testing.T) {
	pool := NewLockPool(128, 64*1024, 2, 1024*1024)
	for i := 0; i < len(pool.classes); i++ {
		temp := make([][]byte, len(pool.classes[i].chunks))

		for j := 0; j < len(temp); j++ {
			mem := pool.Alloc(pool.classes[i].size)
			assert.Equal(t, cap(mem), pool.classes[i].size)
			temp[j] = mem
		}

		for j := 0; j < len(temp); j++ {
			pool.Free(temp[j])
		}
	}
}

func Test_LockPool_AllocSmall(t *testing.T) {
	pool := NewLockPool(128, 1024, 2, 1024)
	mem := pool.Alloc(64)
	assert.Equal(t, len(mem), 64)
	assert.Equal(t, cap(mem), 128)
	pool.Free(mem)
}

func Test_LockPool_AllocLarge(t *testing.T) {
	pool := NewLockPool(128, 1024, 2, 1024)
	mem := pool.Alloc(2048)
	assert.Equal(t, len(mem), 2048)
	assert.Equal(t, cap(mem), 2048)
	pool.Free(mem)
}

func Test_LockPool_DoubleFree(t *testing.T) {
	pool := NewLockPool(128, 1024, 2, 1024)
	mem := pool.Alloc(64)
	go func() {
		defer func() {
			assert.NotNil(t, recover())
		}()
		pool.Free(mem)
		pool.Free(mem)
	}()
}

func Test_LockPool_AllocSlow(t *testing.T) {
	pool := NewLockPool(128, 1024, 2, 1024)
	mem := pool.classes[len(pool.classes)-1].Pop()
	assert.Equal(t, cap(mem), 1024)

	mem = pool.Alloc(1024)
	assert.Equal(t, cap(mem), 1024)
}

func Benchmark_LockPool_AllocAndFree_128(b *testing.B) {
	pool := NewLockPool(128, 1024, 2, 64*1024)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(128))
		}
	})
}

func Benchmark_LockPool_AllocAndFree_256(b *testing.B) {
	pool := NewLockPool(128, 1024, 2, 64*1024)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(256))
		}
	})
}

func Benchmark_LockPool_AllocAndFree_512(b *testing.B) {
	pool := NewLockPool(128, 1024, 2, 64*1024)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(512))
		}
	})
}
