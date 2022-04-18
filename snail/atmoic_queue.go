package container

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type esCache struct {
	putNo uint32 // 在round(putNo/capaciity)次循环中写入
	getNo uint32 // 在round(getNo/capaciity)次循环中读取
	value interface{}
}

// lock free queue
type EsQueue struct {
	capaciity uint32
	capMod    uint32
	putPos    uint32 // put指针，总是增加的，溢出之后可能会小于getpos。
	getPos    uint32
	cache     []esCache
}

func NewQueue(capaciity uint32) *EsQueue {
	q := new(EsQueue)
	q.capaciity = minQuantity(capaciity)
	q.capMod = q.capaciity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]esCache, q.capaciity)
	for i := range q.cache {
		cache := &q.cache[i]
		cache.getNo = uint32(i)
		cache.putNo = uint32(i)
	}
	cache := &q.cache[0]
	cache.getNo = q.capaciity
	cache.putNo = q.capaciity
	return q
}

func (q *EsQueue) String() string {
	getPos := atomic.LoadUint32(&q.getPos)
	putPos := atomic.LoadUint32(&q.putPos)
	return fmt.Sprintf("Queue{capaciity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capaciity, q.capMod, putPos, getPos)
}

func (q *EsQueue) Capaciity() uint32 {
	return q.capaciity
}

func (q *EsQueue) Quantity() uint32 {
	var putPos, getPos uint32
	var quantity uint32
	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		quantity = putPos - getPos
	} else {
		quantity = q.capMod + (putPos - getPos)
	}

	return quantity
}

// put queue functions
func (q *EsQueue) Put(val interface{}) (ok bool, quantity uint32) {
	var putPos, putPosNew, getPos, posCnt uint32
	var cache *esCache
	capMod := q.capMod

	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	// 已超容量，失败。
	if posCnt >= capMod-1 {
		runtime.Gosched()
		return false, posCnt
	}

	// put总是增加的。
	putPosNew = putPos + 1
	// 先占坑，获取本轮的本坑写入权。
	if !atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
		// 被抢了，失败。
		// 这里Gosched一下，因为客户端失败后会循环调用这个函数。
		runtime.Gosched()
		return false, posCnt
	}

	// MOD(X, M) = X&(M-1)  //M必须是2的N次方 //q.capMod = q.capaciity - 1
	// M-1 最后的几个bit位必定是1，所以可以用来取模。
	// 取模，获得实际地址。
	cache = &q.cache[putPosNew&capMod]

	// 队列满，又没有消费者的时候，会for循环。runtime.Gosched()会让出时间。
	for {
		getNo := atomic.LoadUint32(&cache.getNo)
		putNo := atomic.LoadUint32(&cache.putNo)
		// putPosNew == putNo，检查cache的putNo是不是本轮，防止覆盖写入（即下一轮的写入）。
		// getNo == putNo，检查是否被读过，如果没有不能写入。
		if putPosNew == putNo && getNo == putNo {
			// 队列没满，可以写入。
			// 由于读之前会判断有没有写入，所以不会出现同时读写的情况。
			cache.value = val
			// 允许下一轮写入，和putPosNew == putNo这个判断相互斥。
			atomic.AddUint32(&cache.putNo, q.capaciity)
			return true, posCnt + 1
		} else {
			runtime.Gosched()
		}
	}
}

// get queue functions
func (q *EsQueue) Get() (val interface{}, ok bool, quantity uint32) {
	var putPos, getPos, getPosNew, posCnt uint32
	var cache *esCache
	capMod := q.capMod

	putPos = atomic.LoadUint32(&q.putPos)
	getPos = atomic.LoadUint32(&q.getPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt < 1 {
		runtime.Gosched()
		return nil, false, posCnt
	}

	getPosNew = getPos + 1
	// 获取读权。
	if !atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
		// 获取失败。
		runtime.Gosched()
		return nil, false, posCnt
	}

	cache = &q.cache[getPosNew&capMod]

	// 队列空，又没有生产者的时候，会for循环。runtime.Gosched()会让出时间。
	// 可以考虑加sleep
	for {
		getNo := atomic.LoadUint32(&cache.getNo)
		putNo := atomic.LoadUint32(&cache.putNo)
		// getPosNew == getNo。防止同时读。
		// getNo == putNo-q.capaciity。读取前，再次检测是否被写入，如果没写入，说明没有新数据，不读。putNo-q.capaciity表示上一轮的写。
		if getPosNew == getNo && getNo == putNo-q.capaciity {
			val = cache.value
			// 允许下一轮读。
			atomic.AddUint32(&cache.getNo, q.capaciity)
			return val, true, posCnt - 1
		} else {
			runtime.Gosched()
		}
	}
}

// round 到最近的2的倍数 --> 应该是2的N次方
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
