package service

import (
	"errors"
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
