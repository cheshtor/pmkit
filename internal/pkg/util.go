package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"strings"
)

// GenMd5 计算入参的 MD5 值
func GenMd5(plainText string) string {
	if len(plainText) == 0 {
		return ""
	}
	hash := md5.New()
	hash.Write([]byte(plainText))
	return hex.EncodeToString(hash.Sum(nil))
}

// ResolvePage 将页码和数据量修正合法并计算出数据查询时的偏移量
func ResolvePage(pageNo, pageSize int64) (int64, int64, int64) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize
	return pageNo, offset, pageSize
}

// GetRunningFuncName 获取当前执行函数的全名（包含文件路径）
func GetRunningFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// GetCurrentUserId 从上下文中获取当前用户 ID
func GetCurrentUserId(c *fiber.Ctx) int64 {
	localsUid := c.Locals("uid")
	if localsUid != nil {
		return int64(localsUid.(int))
	}
	return 0
}

// SetCurrentUserId 从请求头中获取 Token。如果 Token 存在且合法则解析出当前用户 ID。
func SetCurrentUserId(c *fiber.Ctx) (int64, error) {
	configs = GetConfig()
	uid := configs.GetInt64("app.default-uid")
	token := c.Get("Authorization", "")
	if len(token) > 0 {
		userId, err := GetUserIdFromToken(token)
		if err != nil {
			return -1, err
		}
		uid = userId
	}
	c.Locals("uid", uid)
	return uid, nil
}

// CleanRequestURL 从请求路径中截掉 querystring
func CleanRequestURL(url string) string {
	ok := strings.Contains(url, "?")
	if ok {
		url = url[:strings.Index(url, "?")]
	}
	return url
}
