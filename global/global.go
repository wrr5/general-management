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
	// 立即检查并刷新一次
	if tm.ShouldRefresh() {
		tm.refresh()
	}

	// 每天固定时间刷新（凌晨0:05）
	go func() {
		for {
			now := time.Now()
			// 计算今天凌晨0:05
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 5, 0, 0, now.Location())
			if now.After(today) {
				// 如果已经过了今天0:05，就计算明天0:05
				today = today.Add(24 * time.Hour)
			}
			duration := today.Sub(now)

			log.Printf("距离下次token刷新还有: %v", duration)
			time.Sleep(duration)

			tm.refresh()
		}
	}()
}

// ForceRefresh 强制刷新token
func (tm *TokenManager) ForceRefresh() {
	log.Println("强制刷新token...")
	tm.refresh()
}
