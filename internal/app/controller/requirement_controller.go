package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yourbasic/bit"
	"pmkit/internal/app/model"
	"pmkit/internal/app/service"
	"pmkit/internal/pkg"
)

var requirementService service.RequirementService

type RequirementController struct {
}

func (rc *RequirementController) RestController() string {
	return "/v1/requirement"
}

func (rc *RequirementController) GetRequirement() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/get", permissions, func(c *fiber.Ctx) error {
		requirementId := int64(c.QueryInt("requirementId", 0))
		if requirementId == 0 {
			return fmt.Errorf("找不到指定需求。ID：%d", requirementId)
		}
		wholeRequirement, err := requirementService.GetRequirementById(requirementId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(wholeRequirement))
	}
}

func (rc *RequirementController) ListRequirement() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/list", permissions, func(c *fiber.Ctx) error {
		pageNo := int64(c.QueryInt("pageNo", 1))
		pageSize := int64(c.QueryInt("pageSize", 10))
		var condition model.Requirement
		err := c.BodyParser(&condition)
		if err != nil {
			return err
		}
		list, err := requirementService.SearchRequirement(&condition, pageNo, pageSize)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(list))
	}
}

func (rc *RequirementController) AddRequirement() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/add", permissions, func(c *fiber.Ctx) error {
		var wholeRequirement = new(model.WholeRequirement)
		err := c.BodyParser(wholeRequirement)
		if err != nil {
			return err
		}
		uid := pkg.GetCurrentUserId(c)
		wholeRequirement.CreateBy = uid
		requirement, err := requirementService.SaveRequirement(wholeRequirement)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(requirement))
	}
}

func (rc *RequirementController) UpdateRequirement() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/edit", permissions, func(c *fiber.Ctx) error {
		var wholeRequirement = new(model.WholeRequirement)
		err := c.BodyParser(wholeRequirement)
		if err != nil {
			return err
		}
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := requirementService.EditRequirement(threadLocal, wholeRequirement)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (rc *RequirementController) ChangeRequirementStatus() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/status", permissions, func(c *fiber.Ctx) error {
		requirementId := int64(c.QueryInt("requirementId", 0))
		if requirementId == 0 {
			return fmt.Errorf("找不到指定需求。ID：%d", requirementId)
		}
		targetStatus := model.RequirementStatus(c.QueryInt("targetStatus", 0))
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := requirementService.ChangeRequirementStatus(threadLocal, requirementId, targetStatus)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (rc *RequirementController) RemoveRequirement() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/remove", permissions, func(c *fiber.Ctx) error {
		requirementId := int64(c.QueryInt("requirementId", 0))
		if requirementId == 0 {
			return fmt.Errorf("找不到指定需求。ID：%d", requirementId)
		}
		threadLocal := model.ThreadLocalWithUid(c)
		success, err := requirementService.RemoveRequirement(threadLocal, requirementId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (rc *RequirementController) AddTrack() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/track/add", permissions, func(c *fiber.Ctx) error {
		var track = new(model.RequirementTrack)
		err := c.BodyParser(track)
		if err != nil {
			return err
		}
		uid := pkg.GetCurrentUserId(c)
		track.CreateBy = uid
		track, err = requirementService.SaveRequirementTrack(track)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(track))
	}
}

func (rc *RequirementController) ListTrack() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/track/list", permissions, func(c *fiber.Ctx) error {
		requirementId := int64(c.QueryInt("requirementId", 0))
		if requirementId == 0 {
			return fmt.Errorf("找不到指定需求。ID：%d", requirementId)
		}
		tracks, err := requirementService.GetAllRequirementTracks(requirementId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(tracks))
	}
}

func (rc *RequirementController) AddComment() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodPost, "/comment/add", permissions, func(c *fiber.Ctx) error {
		var comment = new(model.RequirementComment)
		err := c.BodyParser(comment)
		if err != nil {
			return err
		}
		uid := pkg.GetCurrentUserId(c)
		comment.CreateBy = uid
		comment, err = requirementService.SaveRequirementComment(comment)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(comment))
	}
}

func (rc *RequirementController) RemoveComment() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/comment/remove", permissions, func(c *fiber.Ctx) error {
		commentId := int64(c.QueryInt("commentId", 0))
		if commentId == 0 {
			return fmt.Errorf("找不到指定评论。ID：%d", commentId)
		}
		success, err := requirementService.RemoveRequirementComment(commentId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(success))
	}
}

func (rc *RequirementController) ListComment() (string, string, *bit.Set, func(*fiber.Ctx) error) {
	var permissions = bit.New()
	return fiber.MethodGet, "/comment/list", permissions, func(c *fiber.Ctx) error {
		requirementId := int64(c.QueryInt("requirementId", 0))
		if requirementId == 0 {
			return fmt.Errorf("找不到指定需求。ID：%d", requirementId)
		}
		comments, err := requirementService.GetAllRequirementComment(requirementId)
		if err != nil {
			return err
		}
		return c.JSON(pkg.Success(comments))
	}
}
