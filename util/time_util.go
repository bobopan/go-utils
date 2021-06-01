package util

import (
	"fmt"
	"strconv"
	"time"
)

var cstZone = time.FixedZone("CST", 8*3600)

const (
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
	YYYY_MM_DD_00_00_00 = "2006-01-02 00:00:00"
	YYYY_MM_DD_HH_MM    = "2006年01月02日 15:04"
	ORIGIN_DATE         = "2006-01-02 15:04:05"
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM_DD_CH       = "2006年01月02日"
	YYYY__MM__DD        = "2006_01_02"
	TIME_HOUR           = 3600
	YYYYMMDD            = "20060102"
	YYYYMMDDHHIISS      = "20060102150405"
)

var DateTimeFormatMap = map[string]string{
	"Y": "2006",
	"y": "06",
	"m": "01",
	"d": "02",
	"H": "15",
	"i": "04",
	"s": "05",
}

//今日剩余秒数
func GetTimeLeftToday() int {
	timeStr := time.Now().In(cstZone).Format(YYYY_MM_DD_CH)
	t, _ := time.Parse(YYYY_MM_DD_CH, timeStr)
	timeNumber := t.Unix()
	return int(time.Now().In(cstZone).Unix() - timeNumber)
}

func GetAfterDateStr(sec int, format string) string {
	t := time.Now().In(cstZone).Unix()
	t = t - int64(sec)
	return time.Unix(t, 0).In(cstZone).Format(format)
}

//获取当前格式化时间
func GetDateFormat(dateStr string) string {
	data := ""
	for _, v := range dateStr { // i 是字符的字节位置，v 是字符的拷贝
		s := fmt.Sprintf("%c", v)
		if item, ok := DateTimeFormatMap[s]; ok {
			data = data + item
		} else {
			data = data + s
		}
	}
	return time.Now().In(cstZone).Format(data)
}

//获取当前格式化时间
func GetDateIntFormat(dateStr string) int64 {
	data := ""
	for _, v := range dateStr { // i 是字符的字节位置，v 是字符的拷贝
		s := fmt.Sprintf("%c", v)
		if item, ok := DateTimeFormatMap[s]; ok {
			data = data + item
		}
	}
	t, err := strconv.ParseInt(time.Now().In(cstZone).Format(data), 10, 64)
	if err != nil {
		return 0
	}
	return int64(t)
}

//获取昨天的日期
func GetYesDateInt() int64 {
	nTime := time.Now().In(cstZone)
	yesTime := nTime.AddDate(0, 0, -1)
	logDay := yesTime.Format("20060102")
	t, err := strconv.ParseInt(logDay, 10, 64)
	if err != nil {
		return 0
	}
	return int64(t)
}

// 获取当天是星期几
func GetWeekDayInt() int {
	nTime := time.Now().In(cstZone)
	week := int(nTime.Weekday())
	return week
}
