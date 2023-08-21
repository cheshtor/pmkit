package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yourbasic/bit"
	"pmkit/internal/app/model"
	"pmkit/internal/app/service"
	"pmkit/internal/pkg"
)

var iterationService service.IterationService

type IterationController struct {
}

func (ic *IterationController) RestController() string {
	return "/v1/iteration"
}

func (ic *IterationController) GetIteration() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/get", permissions, func(c *fiber.Ctx) error {
		id := c.QueryInt("iterationId", 0)
		if id == 0 {
			return fmt.Errorf("找不到指定迭代。ID：%d", id)
		}
		iteration, err := iterationService.GetIterationById(int64(id))
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(iteration))
	}
}

func (ic *IterationController) ListByProjectId() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/listByProjectId", permissions, func(c *fiber.Ctx) error {
		projectId := int64(c.QueryInt("projectId", 0))
		if projectId == 0 {
			return fmt.Errorf("找不到指定项目。ID：%d", projectId)
		}
		iterations, err := iterationService.ListByProjectId(projectId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(iterations))
	}
}

func (ic *IterationController) AddIteration() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/add", permissions, func(c *fiber.Ctx) error {
		var iteration = new(model.Iteration)
		err := c.BodyParser(iteration)
		if err != nil {
			return err
		}
		threadLocal := model.ThreadLocalWithUid(c)
		iteration, err = iterationService.SaveIteration(threadLocal, iteration)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(iteration))
	}
}

func (ic *IterationController) UpdateIteration() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/edit", permissions, func(c *fiber.Ctx) error {
		var iteration = new(model.Iteration)
		err := c.BodyParser(iteration)
		if err != nil {
			return err
		}
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := iterationService.EditIteration(threadLocal, iteration)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (ic *IterationController) RemoveIteration() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/remove", permissions, func(c *fiber.Ctx) error {
		iterationId := c.QueryInt("iterationId", 0)
		if iterationId == 0 {
			return fmt.Errorf("找不到指定迭代。ID：%d", iterationId)
		}
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := iterationService.RemoveIteration(threadLocal, int64(iterationId))
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}
