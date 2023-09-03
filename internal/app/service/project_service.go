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

var projectDao dao.ProjectDao

type ProjectService struct {
}

func (s *ProjectService) SaveProject(project *model.Project) (*model.Project, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		project.Id = pkg.GetId()
		project.CreateTime = time.Now().UnixMilli()
		affectedRows, err := projectDao.SavaProject(db, project)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("保存项目信息失败")
		}
		return nil
	})
	return project, err
}

func (s *ProjectService) EditProject(t *model.ThreadLocal, project *model.Project) (bool, error) {
	_, _ = s.GetProjectById(project.Id)
	err := database.Run(func(db *sqlx.Tx) error {
		project.ModifiedTime = time.Now().UnixMilli()
		project.ModifiedBy = t.GetUid()
		affectedRows, err := projectDao.EditProject(db, project)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("更新项目信息失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *ProjectService) ChangeProjectStatus(t *model.ThreadLocal, id int64, status model.ProjectStatus) (bool, error) {
	exists, _ := s.GetProjectById(id)
	err := database.Run(func(db *sqlx.Tx) error {
		params := make(map[string]interface{})
		params["id"] = id
		params["status"] = status
		params["start_date"] = exists.StartDate
		params["end_date"] = exists.EndDate
		if status == model.InProcess {
			params["start_date"] = t.Get("timestamp").(int64)
		}
		if status == model.Pause || status == model.Stop || status == model.Finish {
			params["end_date"] = t.Get("timestamp").(int64)
		}
		params["modified_time"] = time.Now().UnixMilli()
		params["modified_by"] = t.GetUid()
		affectedRows, err := projectDao.EditStatus(db, params)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("更新项目状态失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *ProjectService) RemoveProject(t *model.ThreadLocal, id int64) (bool, error) {
	_, _ = s.GetProjectById(id)
	err := database.Run(func(db *sqlx.Tx) error {
		params := make(map[string]interface{})
		params["id"] = id
		params["modified_time"] = time.Now().UnixMilli()
		params["modified_by"] = t.GetUid()
		affectedRows, err := projectDao.RemoveProject(db, params)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("删除项目失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *ProjectService) GetProjectById(id int64) (*model.Project, error) {
	project, err := projectDao.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("项目不存在。ID：%d", id)
	}
	return project, nil
}

func (s *ProjectService) GetProjectList(condition *model.Project, pageNo int64, pageSize int64) (*model.Page, error) {
	calcedPageNo, offset, calcedPageSize := pkg.ResolvePage(pageNo, pageSize)
	list, count, err := projectDao.SearchList(condition.Name, condition.Employer, condition.Contractor, condition.Supervisor, condition.Executor, condition.StartDate, condition.EndDate, condition.Status.Value(), offset, calcedPageSize)
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
