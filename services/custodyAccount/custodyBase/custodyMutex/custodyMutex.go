package custodyMutex

import (
	"sync"
	"sync/atomic"
	"time"
)

type CustodyMutex struct {
	lock       *sync.Mutex
	expiryTime int64 // 过期时间（UNIX 时间戳）
}

var (
	custodyLocks sync.Map  // 使用 sync.Map 替代普通 map，以支持并发读写
	cleanupOnce  sync.Once // 确保清理函数只被调用一次
)

const expiryInterval = 5 * time.Minute // 过期锁的检查间隔

// GetCustodyMutex 获取一个 CustodyMutex，如果不存在则创建一个
func GetCustodyMutex(key string) *sync.Mutex {
	// 启动清理过期锁的 goroutine（仅一次）
	cleanupOnce.Do(func() {
		go clearExpiredLocks()
	})

	// 当前时间
	timeNow := time.Now()

	// 创建新的 CustodyMutex 实例，并设置过期时间
	newMutex := &CustodyMutex{
		lock:       &sync.Mutex{},
		expiryTime: timeNow.Add(expiryInterval).Unix(),
	}

	// 尝试加载或存储锁
	actual, loaded := custodyLocks.LoadOrStore(key, newMutex)
	if loaded {
		// 返回已存在的锁，并更新时间
		existingMutex := actual.(*CustodyMutex)
		existingMutex.UpdateExpiryTime(timeNow) // 更新过期时间
		return existingMutex.lock
	}

	// 返回新创建的锁
	return newMutex.lock
}

// clearExpiredLocks 每5分钟检查并清理过期锁
func clearExpiredLocks() {
	ticker := time.NewTicker(expiryInterval)
	defer ticker.Stop()

	for range ticker.C {
		cleanExpiredLocks()
	}
}

// cleanExpiredLocks 检查并删除过期的 custodyLocks
func cleanExpiredLocks() {
	now := time.Now().Unix()
	custodyLocks.Range(func(key, value interface{}) bool {
		cm := value.(*CustodyMutex)
		// 判断锁是否过期
		if atomic.LoadInt64(&cm.expiryTime) < now {
			custodyLocks.Delete(key) // 删除过期锁
		}
		return true // 继续迭代
	})
}

// UpdateExpiryTime 更新 CustodyMutex 的过期时间
func (cm *CustodyMutex) UpdateExpiryTime(timeNow time.Time) {
	newExpiry := timeNow.Add(expiryInterval).Unix()
	atomic.StoreInt64(&cm.expiryTime, newExpiry) // 线程安全地更新过期时间
}
