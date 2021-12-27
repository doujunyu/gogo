package cache

import (
	"fmt"
	"sync"
	"time"
)
type cache struct {
	Time int64
	Data interface{}
}

var CacheData map[string]*cache
var CacheDataOnce sync.Once

// CacheReady 建立链接,实现单例模式
func init()  {
	CacheDataOnce.Do(func() {
		data := make(map[string]*cache)
		CacheData = data
	})
	//return CacheData
}

// Set 新增缓存,时间=0表示长期存储
func Set(name string,data interface{},second int64) bool {
	if second != 0{
		second += time.Now().Unix()
	}
	mapData := &cache{
		Time: second ,
		Data: data,
	}
	CacheData[name] = mapData
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
	CacheData[name].Time = second
	return true
}

// SetTimeIncrease 增加缓存时间
func SetTimeIncrease(name string,second int64) bool{
	if Exists(name) {
		return false
	}
	CacheData[name].Time += second
	return true
}

// SetTimeDecrease 减少缓存时间
func SetTimeDecrease(name string,second int64) bool{
	if Exists(name) {
		return false
	}
	CacheData[name].Time -= second
	return true
}

// Get 获取缓存数据
func Get(name string) interface{}{
	//if Exists(name) {
	//	return 0
	//}
	return CacheData[name].Data
}

// GetTime 获取缓存时间
func GetTime (name string) int64{
	if Exists(name) {
		return 0
	}
	return CacheData[name].Time
}

// Exists 判断数据是否存在
func Exists (name string) bool {
	if CacheData[name] == nil || time.Now().Unix() > CacheData[name].Time {
		fmt.Println(CacheData[name],CacheData[name].Time)
		return false
	}
	return true
}



