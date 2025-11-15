package global

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type TokenManager struct {
	token      string
	lastUpdate time.Time
	lock       sync.RWMutex
}

var (
	DB *gorm.DB
	TM = &TokenManager{}
)

// Init 初始化token管理器
func (tm *TokenManager) Init() {
	tm.refresh()
}

// refresh 刷新token
func (tm *TokenManager) refresh() {
	token := GetToken() // 调用你的GetToken函数

	tm.lock.Lock()
	defer tm.lock.Unlock()

	tm.token = token
	tm.lastUpdate = time.Now()

	log.Printf("token刷新成功, 更新时间: %s", tm.lastUpdate.Format("2006-01-02 15:04:05"))
}

// Get 获取当前token
func (tm *TokenManager) Get() string {
	tm.lock.RLock()
	defer tm.lock.RUnlock()
	return tm.token
}

// ShouldRefresh 检查是否需要刷新
func (tm *TokenManager) ShouldRefresh() bool {
	tm.lock.RLock()
	defer tm.lock.RUnlock()
	return tm.token == "" || time.Since(tm.lastUpdate) > 24*time.Hour
}

// StartAutoRefresh 启动自动刷新
func (tm *TokenManager) StartAutoRefresh() {
	// 每天凌晨刷新
	go func() {
		for {
			now := time.Now()
			// 计算到明天凌晨的时间
			next := now.Add(24 * time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 5, 0, 0, next.Location()) // 凌晨0点5分
			duration := next.Sub(now)

			time.Sleep(duration)
			tm.refresh()
		}
	}()

	// 备用：每小时检查一次，防止意外情况
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if tm.ShouldRefresh() {
				log.Println("检测到token需要刷新，执行刷新...")
				tm.refresh()
			}
		}
	}()
}

// ForceRefresh 强制刷新token
func (tm *TokenManager) ForceRefresh() {
	log.Println("强制刷新token...")
	tm.refresh()
}
