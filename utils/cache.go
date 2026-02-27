/*
Package utils OI缓存管理

主要功能：
- NewOICacheManager() *OICacheManager                                    // 创建OI缓存管理器
- (m *OICacheManager) Get(symbol string) *indicators.OICache             // 获取缓存
- (m *OICacheManager) Update(symbol string, oi float64, timestamp int64) // 更新缓存
- (m *OICacheManager) GetAll() map[string]*indicators.OICache            // 获取所有缓存
*/
package utils

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// OICache OI缓存结构（避免循环依赖，在这里重新定义）
type OICache struct {
	Symbol     string    // 交易对
	History    []float64 // 历史OI值（从新到旧，最多5个）
	Timestamps []int64   // 对应的时间戳
}

// OICacheManager OI缓存管理器
type OICacheManager struct {
	caches map[string]*OICache
	mu     sync.RWMutex
	maxSize int // 每个symbol最多保存的历史记录数
}

// NewOICacheManager 创建OI缓存管理器
func NewOICacheManager(maxSize int) *OICacheManager {
	if maxSize <= 0 {
		maxSize = 5 // 默认保存5个历史记录
	}
	
	Info("创建OI缓存管理器", zap.Int("max_size", maxSize))
	
	return &OICacheManager{
		caches:  make(map[string]*OICache),
		maxSize: maxSize,
	}
}

// Get 获取指定交易对的OI缓存
func (m *OICacheManager) Get(symbol string) *OICache {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	cache, exists := m.caches[symbol]
	if !exists {
		return nil
	}
	
	return cache
}

// Update 更新指定交易对的OI缓存
func (m *OICacheManager) Update(symbol string, oi float64, timestamp int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	cache, exists := m.caches[symbol]
	if !exists {
		// 创建新缓存
		cache = &OICache{
			Symbol:     symbol,
			History:    []float64{},
			Timestamps: []int64{},
		}
		m.caches[symbol] = cache
	}
	
	// 添加新值到开头
	cache.History = append([]float64{oi}, cache.History...)
	cache.Timestamps = append([]int64{timestamp}, cache.Timestamps...)
	
	// 保持最大数量
	if len(cache.History) > m.maxSize {
		cache.History = cache.History[:m.maxSize]
		cache.Timestamps = cache.Timestamps[:m.maxSize]
	}
	
	Debug("更新OI缓存",
		zap.String("symbol", symbol),
		zap.Float64("oi", oi),
		zap.Int("history_count", len(cache.History)),
	)
}

// GetAll 获取所有缓存
func (m *OICacheManager) GetAll() map[string]*OICache {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// 返回副本，避免外部修改
	result := make(map[string]*OICache, len(m.caches))
	for symbol, cache := range m.caches {
		result[symbol] = cache
	}
	
	return result
}

// Clear 清空指定交易对的缓存
func (m *OICacheManager) Clear(symbol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.caches, symbol)
	Info("清空OI缓存", zap.String("symbol", symbol))
}

// ClearAll 清空所有缓存
func (m *OICacheManager) ClearAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.caches = make(map[string]*OICache)
	Info("清空所有OI缓存")
}

// GetCacheCount 获取缓存数量
func (m *OICacheManager) GetCacheCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return len(m.caches)
}

// GetSymbols 获取所有已缓存的交易对
func (m *OICacheManager) GetSymbols() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	symbols := make([]string, 0, len(m.caches))
	for symbol := range m.caches {
		symbols = append(symbols, symbol)
	}
	
	return symbols
}

// IsExpired 检查缓存是否过期
// maxAge: 最大缓存时间（秒）
func (m *OICacheManager) IsExpired(symbol string, maxAge int64) bool {
	cache := m.Get(symbol)
	if cache == nil || len(cache.Timestamps) == 0 {
		return true
	}
	
	// 检查最新数据的时间戳
	latestTimestamp := cache.Timestamps[0]
	currentTimestamp := time.Now().Unix()
	
	return (currentTimestamp - latestTimestamp) > maxAge
}

// CleanExpired 清理过期缓存
// maxAge: 最大缓存时间（秒）
func (m *OICacheManager) CleanExpired(maxAge int64) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	currentTimestamp := time.Now().Unix()
	cleaned := 0
	
	for symbol, cache := range m.caches {
		if len(cache.Timestamps) == 0 {
			delete(m.caches, symbol)
			cleaned++
			continue
		}
		
		latestTimestamp := cache.Timestamps[0]
		if (currentTimestamp - latestTimestamp) > maxAge {
			delete(m.caches, symbol)
			cleaned++
		}
	}
	
	if cleaned > 0 {
		Info("清理过期OI缓存", zap.Int("count", cleaned))
	}
	
	return cleaned
}

// GetStats 获取缓存统计信息
func (m *OICacheManager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	totalRecords := 0
	for _, cache := range m.caches {
		totalRecords += len(cache.History)
	}
	
	return map[string]interface{}{
		"symbol_count":  len(m.caches),
		"total_records": totalRecords,
		"max_size":      m.maxSize,
	}
}
