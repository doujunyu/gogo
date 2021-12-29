package job

import (
	"time"
)

type Cache struct {
	CacheChan chan int         `Testing:"生成管道,添加和编辑用管道操作"`
	Data      map[string]*CacheValue `Testing:"操作的数据"`
}

// CacheKey 管道数据格式
type CacheKey struct {
	WorkType int         `Testing:"操作方式：1=添加,2=编辑,3=删除,4=增加时间,5=减少时间"`
	Name     string      `Testing:"缓存名"`
	Data     *CacheValue `Testing:"缓存数据"`
}

// CacheValue 缓存数据格式
type CacheValue struct {
	Data interface{}
	Time int64
}

// NewCache 创建一个结构体
func NewCache() *Cache {
	return &Cache{
		CacheChan: make(chan int,1),
		Data:      make(map[string]*CacheValue),
	}
}

// Set 新增缓存,时间=0表示长期存储,无论有没有数据都直接赋值
func (c *Cache) Set(name string, data interface{}, second int64) bool {
	c.CacheChan <- 1
	c.Data[name] = &CacheValue{
		Time: time.Now().Unix() + second,
		Data: data,
	}
	<-c.CacheChan
	return true
}

// SetCover Set 新增缓存,时间=0表示长期存储,如果数据存在则不会进行覆盖,并返回false
func (c *Cache) SetCover(name string, data interface{}, second int64) bool {
	c.CacheChan <- 1
	if !c.Exists(name) {
		<-c.CacheChan
		return false
	}
	c.Data[name] = &CacheValue{
		Time: time.Now().Unix() + second,
		Data: data,
	}
	<-c.CacheChan
	return true
}
// SetDataCover 重置缓存数据
func (c *Cache) SetDataCover(name string, data interface{}) bool {
	c.CacheChan <- 1
	if !c.Exists(name) {
		<-c.CacheChan
		return false
	}
	c.Data[name].Data = data
	<-c.CacheChan
	return true
}

// SetTimeCover 重置缓存时间
func (c *Cache) SetTimeCover(name string, second int64) bool {
	c.CacheChan <- 1
	if !c.Exists(name) {
		<-c.CacheChan
		return false
	}
	c.Data[name].Time = time.Now().Unix() + second
	<-c.CacheChan
	return true
}

// SetTimeIncrease 增加缓存时间
func (c *Cache) SetTimeIncrease(name string, second int64) bool {
	c.CacheChan <- 1
	if !c.Exists(name) {
		<-c.CacheChan
		return false
	}
	c.Data[name].Time +=  second
	<-c.CacheChan
	return true
}

// SetTimeDecrease 减少缓存时间
func (c *Cache) SetTimeDecrease(name string, second int64) bool {
	c.CacheChan <- 1
	if !c.Exists(name) {
		<-c.CacheChan
		return false
	}
	c.Data[name].Time -=  second
	<-c.CacheChan
	return true
}

// Del 删除
func (c *Cache) Del(name string) bool {
	c.CacheChan <- 1
	if !c.Exists(name) {
		<-c.CacheChan
		return true
	}
	delete(c.Data, name)
	<-c.CacheChan
	return true
}

// Get 获取缓存数据
func (c *Cache) Get(name string) interface{} {
	if !c.Exists(name) {
		return nil
	}
	return c.Data[name].Data
}

// GetTime 获取缓存时间
func (c *Cache) GetTime(name string) int64 {
	if !c.Exists(name) {
		return 0
	}
	return c.Data[name].Time - time.Now().Unix()
}

// Exists 判断数据是否存在
func (c *Cache) Exists(name string) bool {
	if c.Data[name] == nil || time.Now().Unix() > c.Data[name].Time {
		return false
	}
	return true
}

// ChanLongTime 清理过期的缓存
func (c *Cache) ChanLongTime() {
	for {
		for name, value := range c.Data {
			time.Sleep(time.Second)
			if value.Time < time.Now().Unix() {
				c.Del(name)
			}
		}
	}
}
