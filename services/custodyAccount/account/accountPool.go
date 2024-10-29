package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"trade/btlLog"
	"trade/models"
	cModels "trade/models/custodyModels"
	"trade/services/btldb"
)

// 定义用户池相关的错误信息
var (
	ErrMaxUserPoolReached = fmt.Errorf("用户池已满，无法添加新用户")
	ErrUserNotFound       = fmt.Errorf("用户不存在")
)

var pool *UserPool

// UserInfo 定义用户结构体
type UserInfo struct {
	User          *models.User
	Account       *models.Account
	LockAccount   *cModels.LockAccount
	PaymentMux    sync.Mutex // 支付锁
	RpcMux        sync.Mutex // RPC 锁
	LastActiveMux sync.Mutex // 上次活跃时间锁
	LastActive    time.Time  // 记录上次活跃时间
}

// UserPool 定义用户池结构体
type UserPool struct {
	users            map[string]*UserInfo // 存储用户信息的映射
	mutex            sync.RWMutex         // 读写锁，确保并发安全
	maxCapacity      int                  // 用户池最大容量
	inactiveDuration time.Duration        // 不活跃用户的最大持续时间
}

// 初始化用户池
func init() {
	pool = NewUserPool(30, 3*time.Minute)
	pool.StartCleanupScheduler(5 * time.Minute)
}

// NewUserPool 创建新的用户池
func NewUserPool(maxCapacity int, inactiveDuration time.Duration) *UserPool {
	if maxCapacity <= 0 {
		maxCapacity = 2000 // 默认容量
	}
	return &UserPool{
		users:            make(map[string]*UserInfo),
		maxCapacity:      maxCapacity,
		inactiveDuration: inactiveDuration,
	}
}

// RemoveUser 从用户池中删除用户
func (pool *UserPool) RemoveUser(userName string) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if _, exists := pool.users[userName]; exists {
		delete(pool.users, userName)
		btlLog.CUST.Info("用户 %s 已被删除\n", userName)
	} else {
		btlLog.CUST.Warning("警告: 用户 %s 不存在，无法删除\n", userName)
	}
}

// GetUser 根据用户ID获取用户
func (pool *UserPool) GetUser(userName string) (*UserInfo, bool) {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	user, exists := pool.users[userName]
	return user, exists
}

// CreateUser 创建新用户
func (pool *UserPool) CreateUser(userName string) (*UserInfo, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if existingUser, exists := pool.users[userName]; exists {
		return existingUser, nil
	}

	if len(pool.users) >= pool.maxCapacity {
		return nil, ErrMaxUserPoolReached
	}

	u, a, l, err := GetUserInfoFromDb(userName)
	if err != nil {
		btlLog.CUST.Error("错误: 获取用户 %s 信息时发生错误: %v\n", userName, err)
		return nil, err
	}

	newUser := &UserInfo{
		User:        u,
		Account:     a,
		LockAccount: l,
		LastActive:  time.Now(),
	}

	pool.users[userName] = newUser
	btlLog.CUST.Info("用户 %s 已被添加到用户池\n", userName)
	return newUser, nil
}

// ListUsers 列出所有用户
func (pool *UserPool) ListUsers() []*UserInfo {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	users := make([]*UserInfo, 0, len(pool.users))
	for _, user := range pool.users {
		users = append(users, user)
	}
	return users
}

// CleanupInactiveUsers 清理不活跃用户
func (pool *UserPool) CleanupInactiveUsers() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	now := time.Now()
	for userName, user := range pool.users {
		if userName == "admin" {
			continue
		}
		if now.Sub(user.LastActive) > pool.inactiveDuration {
			delete(pool.users, userName)
			btlLog.CUST.Info("用户 %s 因为不活跃被清理\n", userName)
		}
	}
}

// StartCleanupScheduler 启动定期清理不活跃用户的调度器
func (pool *UserPool) StartCleanupScheduler(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			pool.CleanupInactiveUsers()
		}
	}()
}

// UpdateUserActivity 更新用户活动时间
func (pool *UserPool) UpdateUserActivity(userName string) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if user, exists := pool.users[userName]; exists {
		user.LastActive = time.Now()
		btlLog.CUST.Info("用户 %s 的活跃时间已更新为 %s\n", userName, user.LastActive)
	} else {
		btlLog.CUST.Warning("警告: 用户 %s 不存在，无法更新活动时间\n", userName)
	}
}

// GetUserInfoFromDb 获取用户信息
func GetUserInfoFromDb(username string) (*models.User, *models.Account, *cModels.LockAccount, error) {
	// 获取用户信息
	user, err := btldb.ReadUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ErrUserNotFound
		}
		return nil, nil, nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	}
	if user.Status == 1 {
		return nil, nil, nil, errors.New("用户已被冻结")
	}
	// 获取Lit账户信息
	account := &models.Account{}
	account, err = GetAccountByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		account, err = CreateAccount(user)
		if err != nil {
			return nil, nil, nil, CustodyAccountCreateErr
		}
	}
	// 获取冻结账户信息
	lockAccount := &cModels.LockAccount{}
	lockAccount, err = GetLockAccountByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, fmt.Errorf("%w: %w", models.ReadDbErr, err)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果账户不存在，则创建账户
		lockAccount, err = CreateLockAccount(user)
		if err != nil {
			return nil, nil, nil, CustodyAccountCreateErr
		}
	}

	return user, account, lockAccount, nil
}
