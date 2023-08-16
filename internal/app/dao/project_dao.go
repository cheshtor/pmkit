package dao

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/app/model"
	"pmkit/internal/pkg"
)

type ProjectDao struct {
}

func (d *ProjectDao) SavaProject(db *sqlx.Tx, project *model.Project) (int64, error) {
	sql := "INSERT INTO `pk_project` (`id`, `name`, `employer`, `contractor`, `supervisor`, `executor`, `description`, `start_date`, `status`, `create_by`, `create_time`) VALUES (:id, :name, :employer, :contractor, :supervisor, :executor, :description, :start_date, :status, :create_by, :create_time)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, project)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *ProjectDao) EditProject(db *sqlx.Tx, project *model.Project) (int64, error) {
	return 0, nil
}

func (d *ProjectDao) FindById(id int64) (*model.Project, error) {
	return nil, nil
}

func (d *ProjectDao) SearchList(name, employer, contractor, supervisor, executor string, status int) ([]*model.Project, int64, error) {
	return nil, 0, nil
}
