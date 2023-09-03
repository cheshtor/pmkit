package service

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/app/dao"
	"pmkit/internal/app/model"
	"pmkit/internal/pkg"
	"pmkit/internal/pkg/database"
	"time"
)

var requirementDao dao.RequirementDao
var contentDao dao.RequirementContentDao
var commentDao dao.RequirementCommentDao
var trackDao dao.RequirementTrackDao

type RequirementService struct {
}

func (s *RequirementService) SaveRequirement(wholeRequirement *model.WholeRequirement) (*model.Requirement, error) {
	requirement := wholeRequirement.SeparateRequirement()
	content := wholeRequirement.SeparateContent()
	err := database.Run(func(db *sqlx.Tx) error {
		requirement.Id = pkg.GetId()
		requirement.CreateTime = time.Now().UnixMilli()
		affectedRows, err := requirementDao.SaveRequirement(db, requirement)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("保存需求信息失败")
		}
		affectedRows, err = contentDao.SaveContent(db, requirement.Id, content.Content)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("保存需求详情失败")
		}
		return nil
	})
	return requirement, err
}

func (s *RequirementService) EditRequirement(t *model.ThreadLocal, wholeRequirement *model.WholeRequirement) (bool, error) {
	requirement := wholeRequirement.SeparateRequirement()
	content := wholeRequirement.SeparateContent()
	err := database.Run(func(db *sqlx.Tx) error {
		requirement.ModifiedTime = time.Now().UnixMilli()
		requirement.CreateBy = t.GetUid()
		affectedRows, err := requirementDao.EditRequirement(db, requirement)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("更新需求信息失败")
		}
		affectedRows, err = contentDao.EditContent(db, requirement.Id, content.Content)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("更新需求详情失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *RequirementService) ChangeRequirementStatus(t *model.ThreadLocal, id int64, status model.RequirementStatus) (bool, error) {
	requirement, _ := s.GetRequirementById(id)
	canTransfer := requirement.Status.TransferTo(status)
	if !canTransfer {
		return false, errors.New("此需求不能变更为目标状态")
	}
	err := database.Run(func(db *sqlx.Tx) error {
		params := make(map[string]interface{})
		params["id"] = id
		params["status"] = status.Value()
		params["modified_by"] = t.GetUid()
		params["modified_time"] = time.Now().UnixMilli()
		affectedRows, err := requirementDao.ChangeStatus(db, params)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("更新需求状态失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *RequirementService) RemoveRequirement(t *model.ThreadLocal, id int64) (bool, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		params := make(map[string]interface{})
		params["id"] = id
		params["modified_by"] = t.GetUid()
		params["modified_time"] = time.Now().UnixMilli()
		affectedRows, err := requirementDao.RemoveRequirement(db, params)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("删除需求失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *RequirementService) GetRequirementById(id int64) (*model.WholeRequirement, error) {
	requirement, err := requirementDao.FindByRequirementId(id)
	if err != nil {
		return nil, fmt.Errorf("需求不存在。ID：%d", id)
	}
	content, _ := contentDao.FindByRequirementId(id)
	whole := new(model.WholeRequirement)
	whole.Merge(requirement, content)
	return whole, err
}

func (s *RequirementService) SearchRequirement(cond *model.Requirement, pageNo, pageSize int64) (*model.Page, error) {
	calcedPageNo, offset, calcedPageSize := pkg.ResolvePage(pageNo, pageSize)
	list, count, err := requirementDao.SearchList(cond.ProjectId, cond.IterationId, cond.Name, cond.Type.Value(), cond.Demander, cond.Priority, cond.Influence.Value(), cond.Owner, cond.Tracer, cond.Status.Value(), offset, calcedPageSize)
	if err != nil {
		return nil, err
	}
	var page = new(model.Page)
	page.PageNo = calcedPageNo
	page.PageSize = calcedPageSize
	page.Rows = make([]interface{}, len(list))
	for index, item := range list {
		page.Rows[index] = item
	}
	page.TotalCount = count
	return page, nil
}

func (s *RequirementService) SaveRequirementTrack(track *model.RequirementTrack) (*model.RequirementTrack, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		track.Id = pkg.GetId()
		track.CreateTime = time.Now().UnixMilli()
		affectedRows, err := trackDao.SaveTrack(db, track)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("保存需求状态变更记录失败")
		}
		return nil
	})
	return track, err
}

func (s *RequirementService) GetAllRequirementTracks(requirementId int64) ([]*model.RequirementTrack, error) {
	return trackDao.ListByRequirementId(requirementId)
}

func (s *RequirementService) SaveRequirementComment(comment *model.RequirementComment) (*model.RequirementComment, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		comment.Id = pkg.GetId()
		comment.CreateTime = time.Now().UnixMilli()
		affectedRows, err := commentDao.SaveComment(db, comment)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("保存需求评论失败")
		}
		return nil
	})
	return comment, err
}

func (s *RequirementService) RemoveRequirementComment(commentId int64) (bool, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		affectedRows, err := commentDao.RemoveComment(db, commentId)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("删除需求评论失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *RequirementService) GetAllRequirementComment(requirementId int64) ([]*model.RequirementComment, error) {
	return commentDao.ListByRequirementId(requirementId)
}
