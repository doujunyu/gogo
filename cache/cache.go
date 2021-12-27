package cache

import (
	"sync"
	"time"
)
type cache struct {
	Time int64
	Data interface{}
}

var cacheData map[string]*cache
var cacheDataOnce sync.Once

// CacheReady 建立链接,实现单例模式
func init() {
	cacheDataOnce.Do(func() {
		data := make(map[string]*cache)
		cacheData = data
	})
}

// Set 新增缓存,时间=0表示长期存储
func Set(name string,second int64,data interface{}) bool {
	if second != 0{
		second += time.Now().Unix()
	}
	mapData := &cache{
		Time: second ,
		Data: data,
	}
	cacheData[name] = mapData
	return true
}

// SetTimeCover 重置缓存时间
func SetTimeCover(name string,second int64) bool{
	if !Exists(name) {
		return false
	}
	if second != 0{
		second += time.Now().Unix()
	}
	cacheData[name].Time = second
	return true
}

// SetTimeIncrease 增加缓存时间
func SetTimeIncrease(name string,second int64) bool{
	if Exists(name) {
		return false
	}
	cacheData[name].Time += second
	return true
}

// SetTimeDecrease 减少缓存时间
func SetTimeDecrease(name string,second int64) bool{
	if Exists(name) {
		return false
	}
	cacheData[name].Time -= second
	return true
}

// Get 获取缓存数据
func Get(name string) interface{}{
	if Exists(name) {
		return 0
	}
	return cacheData[name].Data
}

// GetTime 获取缓存时间
func GetTime (name string) int64{
	if Exists(name) {
		return 0
	}
	return cacheData[name].Time
}

// Exists 判断数据是否存在
func Exists (name string) bool {
	if cacheData[name] == nil || time.Now().Unix() > cacheData[name].Time {
		return false
	}
	return true
}



