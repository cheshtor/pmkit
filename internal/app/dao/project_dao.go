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
	setClause := "SET"
	if len(project.Name) != 0 {
		setClause += " `name` = :name,"
	}
	if len(project.Employer) != 0 {
		setClause += " `employer` = :employer,"
	}
	if len(project.Contractor) != 0 {
		setClause += " `contractor` = :contractor,"
	}
	if len(project.Supervisor) != 0 {
		setClause += " `supervisor` = :supervisor,"
	}
	if len(project.Executor) != 0 {
		setClause += " `executor` = :executor,"
	}
	if len(project.Description) != 0 {
		setClause += " `description` = :description,"
	}
	if project.StartDate != 0 {
		setClause += " `start_date` = :start_date,"
	}
	if project.ModifiedTime != 0 {
		setClause += " `modified_time` = :modified_time,"
	}
	if project.ModifiedBy != 0 {
		setClause += " `modified_by` = :modified_by,"
	}
	if setClause == "SET" {
		return 0, errors.New("SQL 构造失败")
	}
	before, _ := strings.CutSuffix(setClause, ",")
	setClause = before
	sql := "UPDATE `pk_project` " + setClause + " WHERE `id` = :id"
	project.ModifiedTime = time.Now().UnixMilli()
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, project)
	if err != nil {
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (d *ProjectDao) EditStatus(db *sqlx.Tx, params map[string]interface{}) (int64, error) {
	sql := "UPDATE `pk_project` SET `status` = :status, `start_date` = :start_date, `end_date` = :end_date, `modified_time` = :modified_time, `modified_by` = :modified_by WHERE `id` = :id AND `delete` = 0"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, params)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *ProjectDao) RemoveProject(db *sqlx.Tx, params map[string]interface{}) (int64, error) {
	sql := "UPDATE `pk_project` SET `delete` = 1, `modified_time` = :modified_time, `modified_by` = :modified_by WHERE `id` = :id AND `delete` = 0"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, params)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *ProjectDao) FindById(id int64) (*model.Project, error) {
	db := database.GetDBInstance()
	project := &model.Project{}
	sql := "SELECT `id`, `name`, `employer`, `contractor`, `supervisor`, `executor`, `description`, `start_date`, `end_date`, `status`, `delete`, `create_by`, `create_time`, `modified_by`, `modified_time` FROM `pk_project` WHERE `id` = ? AND `delete` = 0"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	err := db.Get(project, sql, id)
	return project, err
}

func (d *ProjectDao) SearchList(name, employer, contractor, supervisor, executor string, startDate, endDate int64, status int, offset, limit int64) ([]*model.Project, int64, error) {
	db := database.GetDBInstance()
	var projects = make([]*model.Project, 0)
	sql := "SELECT `id`, `name`, `employer`, `contractor`, `supervisor`, `executor`, `description`, `start_date`, `end_date`, `status`, `delete`, `create_by`, `create_time`, `modified_by`, `modified_time` FROM `pk_project`"
	orderClause := " ORDER BY `create_time` DESC"
	limitClause := " LIMIT :offset, :limit"
	whereClause := ""
	if len(name) != 0 {
		whereClause += " AND `name` LIKE CONCAT('%', :name, '%')"
	}
	if len(employer) != 0 {
		whereClause += " AND `employer` LIKE CONCAT('%', :employer, '%')"
	}
	if len(contractor) != 0 {
		whereClause += " AND `contractor` LIKE CONCAT('%', :contractor, '%')"
	}
	if len(supervisor) != 0 {
		whereClause += " AND `supervisor` LIKE CONCAT('%', :supervisor, '%')"
	}
	if len(executor) != 0 {
		whereClause += " AND `executor` LIKE CONCAT('%', :executor, '%')"
	}
	if startDate != 0 {
		whereClause += " AND `start_date` >= :start_date"
	}
	if endDate != 0 {
		whereClause += " AND `end_date` <= :end_date"
	}
	if status != 0 {
		whereClause += " AND `status` = :status"
	}
	if len(whereClause) != 0 {
		after, _ := strings.CutPrefix(whereClause, " AND")
		whereClause = " WHERE " + after
	}
	sql = sql + whereClause + orderClause + limitClause
	namedParams := map[string]interface{}{
		"name":       name,
		"employer":   employer,
		"contractor": contractor,
		"supervisor": supervisor,
		"executor":   executor,
		"start_date": startDate,
		"end_date":   endDate,
		"status":     status,
		"offset":     offset,
		"limit":      limit,
	}
	funcName := pkg.GetRunningFuncName()
	log.Debugf("[%s] Execute SQL => %s\n", funcName, sql)
	rows, err := db.NamedQuery(sql, namedParams)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		p := &model.Project{}
		_ = rows.StructScan(p)
		projects = append(projects, p)
	}
	countSql := "SELECT COUNT(1) FROM `pk_project`" + whereClause
	log.Debugf("[%s] Execute SQL => %s\n", funcName, countSql)
	res, err := db.NamedQuery(countSql, namedParams)
	if err != nil {
		return nil, 0, err
	}
	var count int64
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			return nil, 0, err
		}
	}
	return projects, count, nil
}
