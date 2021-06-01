package util

import (
	"github.com/hashicorp/go-version"
	"strings"
)

//版本判断  appVersion现在的版本  和 checkVersion的对比
// (注意checkVersion要带比较符号，
// 如验证  6.1.9是否大于等于 6.1.8  appVersion="6.1.9" checkVersion=">=6.1.8" )
func VersionCompare(appVersion string, checkVersion string) bool {
	if checkVersion == "" {
		return true
	}
	if strings.Contains(appVersion, "_") {
		appVersion = strings.Replace(appVersion, "_", ".", 10)
	}
	compare, err := version.NewConstraint(string(checkVersion))
	if err != nil {
		return false
	}
	appV, err := version.NewVersion(appVersion)
	if err != nil {
		return false
	}
	return compare.Check(appV)
}

//灰度
func VersionGrayscale(src string) bool {
	if src == "intelmgtv" {
		return false
	}
	return true
}
