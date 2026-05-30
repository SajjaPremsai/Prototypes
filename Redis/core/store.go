package core

import (
	"time"
)

var store map[string]*Obj

type Obj struct{
	value interface{}
	ExpireAt int64
}

func init(){
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64) *Obj{
	var expireAt int64 = -1
	if durationMs > 0{
		expireAt = time.Now().UnixMilli() + durationMs
	}

	return &Obj{
		value: value,
		ExpireAt: expireAt,
	}
}

func Put(k string, obj * Obj){
	store[k] = obj                             
}

func Get(k string) *Obj{
	return store[k]
}