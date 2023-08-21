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

var iterationDao dao.IterationDao

type IterationService struct {
}

func (s *IterationService) SaveIteration(t *model.ThreadLocal, iteration *model.Iteration) (*model.Iteration, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		iteration.Id = pkg.GetId()
		iteration.CreateTime = time.Now().UnixMilli()
		iteration.CreateBy = t.GetUid()
		affectedRows, err := iterationDao.SaveIteration(db, iteration)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("保存迭代信息失败")
		}
		return nil
	})
	return iteration, err
}

func (s *IterationService) EditIteration(t *model.ThreadLocal, iteration *model.Iteration) (bool, error) {
	_, _ = s.GetIterationById(iteration.Id)
	err := database.Run(func(db *sqlx.Tx) error {
		iteration.ModifiedTime = time.Now().UnixMilli()
		iteration.ModifiedBy = t.GetUid()
		affectedRows, err := iterationDao.EditIteration(db, iteration)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("更新迭代信息失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *IterationService) RemoveIteration(t *model.ThreadLocal, id int64) (bool, error) {
	_, _ = s.GetIterationById(id)
	err := database.Run(func(db *sqlx.Tx) error {
		params := make(map[string]interface{})
		params["id"] = id
		params["modified_time"] = time.Now().UnixMilli()
		params["modified_by"] = t.GetUid()
		affectedRows, err := iterationDao.RemoveIteration(db, params)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("删除迭代失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *IterationService) GetIterationById(id int64) (*model.Iteration, error) {
	iteration, err := iterationDao.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("迭代不存在。ID：%d", id)
	}
	return iteration, nil
}

func (s *IterationService) ListByProjectId(projectId int64) ([]*model.Iteration, error) {
	projectService := new(ProjectService)
	_, _ = projectService.GetProjectById(projectId)
	return iterationDao.ListByProjectId(projectId)
}
