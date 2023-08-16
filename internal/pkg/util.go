package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"strings"
)

func GenMd5(plainText string) string {
	if len(plainText) == 0 {
		return ""
	}
	hash := md5.New()
	hash.Write([]byte(plainText))
	return hex.EncodeToString(hash.Sum(nil))
}

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

func GetRunningFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func GetCurrentUserId(c *fiber.Ctx) int64 {
	localsUid := c.Locals("uid")
	if localsUid != nil {
		return int64(localsUid.(int))
	}
	return 0
}

func SetCurrentUserId(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	if len(token) > 0 {
		userId, err := GetUserIdFromToken(token)
		if err != nil {
			return err
		}
		c.Locals("uid", userId)
	} else {
		c.Locals("uid", 1690713268944900096)
	}
	return nil
}

func FixRequestURL(url string) string {
	ok := strings.Contains(url, "?")
	if ok {
		url = url[:strings.Index(url, "?")]
	}
	return url
}
