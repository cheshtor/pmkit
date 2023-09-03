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

var projectService service.ProjectService

type ProjectController struct {
}

func (pc *ProjectController) RestController() string {
	return "/v1/project"
}

func (pc *ProjectController) GetProject() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/get", permissions, func(c *fiber.Ctx) error {
		id := c.QueryInt("projectId", 0)
		if id == 0 {
			return fmt.Errorf("找不到指定项目。ID：%d", id)
		}
		project, err := projectService.GetProjectById(int64(id))
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(project))
	}
}

func (pc *ProjectController) ListProject() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/list", permissions, func(c *fiber.Ctx) error {
		pageNo := int64(c.QueryInt("pageNo", 1))
		pageSize := int64(c.QueryInt("pageSize", 10))
		var condition model.Project
		err := c.BodyParser(&condition)
		if err != nil {
			return err
		}
		list, err := projectService.GetProjectList(&condition, pageNo, pageSize)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(list))
	}
}

func (pc *ProjectController) AddProject() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/add", permissions, func(c *fiber.Ctx) error {
		var project = new(model.Project)
		err := c.BodyParser(project)
		if err != nil {
			return err
		}
		uid := pkg.GetCurrentUserId(c)
		project.CreateBy = uid
		project, err = projectService.SaveProject(project)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(project))
	}
}

func (pc *ProjectController) UpdateProject() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/edit", permissions, func(c *fiber.Ctx) error {
		var project = new(model.Project)
		err := c.BodyParser(project)
		if err != nil {
			return err
		}
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := projectService.EditProject(threadLocal, project)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (pc *ProjectController) ChangeStatus() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/status", permissions, func(c *fiber.Ctx) error {
		projectId := int64(c.QueryInt("projectId", 0))
		if projectId == 0 {
			return fmt.Errorf("找不到指定项目。ID：%d", projectId)
		}
		targetStatus := model.ProjectStatus(c.QueryInt("targetStatus", 0))
		timestamp := int64(c.QueryInt("timestamp", 0))
		if timestamp == 0 && (targetStatus == model.InProcess || targetStatus == model.Pause || targetStatus == model.Stop || targetStatus == model.Finish) {
			return errors.New("项目要变更的目标状态缺少参数：日期")
		}
		threadLocal := model.ThreadLocalWithUid(c)
		threadLocal.Set("timestamp", timestamp)
		success, err := projectService.ChangeProjectStatus(threadLocal, projectId, targetStatus)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (pc *ProjectController) RemoveProject() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/remove", permissions, func(c *fiber.Ctx) error {
		projectId := int64(c.QueryInt("projectId", 0))
		if projectId == 0 {
			return fmt.Errorf("找不到指定项目。ID：%d", projectId)
		}
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := projectService.RemoveProject(threadLocal, projectId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}
