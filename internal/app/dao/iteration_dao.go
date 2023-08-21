package dao

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/app/model"
	"pmkit/internal/pkg"
	"pmkit/internal/pkg/database"
	"strings"
	"time"
)

type IterationDao struct {
}

func (d *IterationDao) SaveIteration(db *sqlx.Tx, iteration *model.Iteration) (int64, error) {
	sql := "INSERT INTO `pk_iteration` (`id`, `project_id`, `name`, `remark`, `start_date`, `end_date`, `create_by`, `create_time`) VALUES (:id, :project_id, :name, :remark, :start_date, :end_date, :create_by, :create_time)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, iteration)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *IterationDao) EditIteration(db *sqlx.Tx, iteration *model.Iteration) (int64, error) {
	setClause := "SET"
	if len(iteration.Name) != 0 {
		setClause += " `name` = :name, "
	}
	if len(iteration.Remark) != 0 {
		setClause += " `remark` = :remark, "
	}
	if iteration.StartDate != 0 {
		setClause += " `start_date` = :start_date, "
	}
	if iteration.EndDate != 0 {
		setClause += " `end_date` = :end_date, "
	}
	if iteration.ModifiedTime != 0 {
		setClause += " `modified_time` = :modified_time,"
	}
	if iteration.ModifiedBy != 0 {
		setClause += " `modified_by` = :modified_by,"
	}
	if setClause == "SET" {
		return 0, errors.New("SQL 构造失败")
	}
	before, _ := strings.CutSuffix(setClause, ",")
	setClause = before
	sql := "UPDATE `pk_iteration` " + setClause + " WHERE `id` = :id"
	iteration.ModifiedTime = time.Now().UnixMilli()
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, iteration)
	if err != nil {
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (d *IterationDao) RemoveIteration(db *sqlx.Tx, params map[string]interface{}) (int64, error) {
	sql := "UPDATE `pk_iteration` SET `delete` = 1, `modified_time` = :modified_time, `modified_by` = :modified_by WHERE `id` = :id AND `delete` = 0"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, params)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *IterationDao) FindById(id int64) (*model.Iteration, error) {
	db := database.GetDBInstance()
	iteration := &model.Iteration{}
	sql := "SELECT `id`, `project_id`, `name`, `remark`, `start_date`, `end_date`, `create_by`, `create_time`, `modified_by`, `modified_time` FROM `pk_iteration` WHERE `id` = ? AND `delete` = 0"
	err := db.Get(iteration, sql, id)
	return iteration, err
}

func (d *IterationDao) ListByProjectId(projectId int64) ([]*model.Iteration, error) {
	db := database.GetDBInstance()
	sql := "SELECT `id`, `project_id`, `name`, `remark`, `start_date`, `end_date`, `create_by`, `create_time`, `modified_by`, `modified_time` FROM `pk_iteration` WHERE `project_id` = ? AND `delete` = 0 ORDER BY `create_time` DESC"
	rows, err := db.Queryx(sql, projectId)
	if err != nil {
		return nil, err
	}
	var iterations = make([]*model.Iteration, 0)
	for rows.Next() {
		i := &model.Iteration{}
		_ = rows.StructScan(i)
		iterations = append(iterations, i)
	}
	return iterations, nil
}
