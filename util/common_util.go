package util

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func BaseJsonEncode(data interface{}) string {
	mjson, _ := json.Marshal(data)
	mString := string(mjson)
	return mString
}

func GetFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		return f.Name()
	}
	return ""
}

//hashCode
func Hash2(s string) int64 {
	var h int64 = 0
	ln := len(s)
	for i := 0; i < ln; i++ {
		h = 31*h + int64(s[i])
	}
	if h < 0 {
		return -h
	}
	return h
}

//判断是不是在数组组
func VailIsInInt(intArray []int, i int) bool {
	for _, typ := range intArray {
		if typ == i {
			return true
		}
	}
	return false
}

//判断是不是在数组组
func VailIsInInt64(intArray []int64, i int64) bool {
	for _, typ := range intArray {
		if typ == i {
			return true
		}
	}
	return false
}

//判断是不是在数组组
func VailIsInStr(intArray []string, i string) bool {
	for _, typ := range intArray {
		if typ == i {
			return true
		}
	}
	return false
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return int64(rand.Intn(int(max-min))) + min
}

func Removal(strs []string) []string {
	data := []string{}
	lsMap := make(map[string]bool)
	for _, str := range strs {
		if _, ok := lsMap[str]; ok {
			continue
		}
		data = append(data, str)
	}
	return data
}

//uuids去重
func RemovalUuid(strs []string) []string {
	data := []string{}
	lsMap := make(map[string]bool)
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" || len(str) != 32 {
			continue
		}
		if _, ok := lsMap[str]; ok {
			continue
		}
		data = append(data, str)
	}
	return data
}
