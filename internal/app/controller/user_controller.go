package controller

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yourbasic/bit"
	"pmkit/internal/app/model"
	"pmkit/internal/app/service"
	"pmkit/internal/pkg"
)

var userService service.UserService

type UserController struct {
}

func (uc *UserController) RestController() string {
	return "/v1/user"
}

func (uc *UserController) CheckLogin() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/login", permissions, func(c *fiber.Ctx) error {
		var user = new(model.User)
		err := c.BodyParser(user)
		if err != nil {
			return err
		}
		token, err := userService.CheckLogin(user.Phone, user.Password)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(token))
	}
}

func (uc *UserController) GetUser() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/get", permissions, func(c *fiber.Ctx) error {
		id := c.QueryInt("id", 0)
		if id == 0 {
			return fmt.Errorf("找不到指定用户。ID：%d", id)
		}
		user, err := userService.GetUserById(int64(id))
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(user))
	}
}

func (uc *UserController) ListUser() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/list", permissions, func(c *fiber.Ctx) error {
		pageNo := int64(c.QueryInt("pageNo", 1))
		pageSize := int64(c.QueryInt("pageSize", 10))
		var condition model.User
		err := c.BodyParser(&condition)
		if err != nil {
			return err
		}
		list, err := userService.GetUserList(condition, pageNo, pageSize)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(list))
	}
}

func (uc *UserController) AddUser() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/add", permissions, func(c *fiber.Ctx) error {
		var user = new(model.User)
		err := c.BodyParser(user)
		if err != nil {
			return err
		}
		success, err := userService.SaveUser(user)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (uc *UserController) UpdateUser() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/update", permissions, func(c *fiber.Ctx) error {
		var user = new(model.User)
		err := c.BodyParser(user)
		if err != nil {
			return err
		}
		success, err := userService.EditUser(user)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (uc *UserController) ActiveSwitch() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/switch", permissions, func(c *fiber.Ctx) error {
		id := int64(c.QueryInt("id", 0))
		if id == 0 {
			return fmt.Errorf("找不到指定用户。ID：%d", id)
		}
		activeParam := c.Query("active", "")
		if len(activeParam) == 0 || (activeParam != "enable" && activeParam != "disable") {
			return errors.New("用户目标状态不合法")
		}
		active := activeParam == "enable"
		success, err := userService.ActiveSwitch(id, active)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}
