package model

import (
	"github.com/gofiber/fiber/v2"
	"pmkit/internal/pkg"
)

type ThreadLocal struct {
	LocalMap map[string]interface{}
}

func ThreadLocalWithUid(ctx *fiber.Ctx) *ThreadLocal {
	threadLocal := &ThreadLocal{
		LocalMap: make(map[string]interface{}),
	}
	threadLocal.Set("uid", pkg.GetCurrentUserId(ctx))
	return threadLocal
}

func (c *ThreadLocal) Set(name string, value interface{}) {
	c.LocalMap[name] = value
}

func (c *ThreadLocal) Get(name string) interface{} {
	value, found := c.LocalMap[name]
	if found {
		return value
	}
	return nil
}

func (c *ThreadLocal) GetUid() int64 {
	value := c.Get("uid")
	if value != nil {
		return value.(int64)
	}
	return 0
}

func (c *ThreadLocal) Remove(name string) {
	delete(c.LocalMap, name)
}

func (c *ThreadLocal) Clear() {
	for k := range c.LocalMap {
		delete(c.LocalMap, k)
	}
}
