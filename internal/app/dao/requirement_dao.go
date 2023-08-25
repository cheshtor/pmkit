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

type RequirementDao struct {
}

func (d *RequirementDao) SaveRequirement(db *sqlx.Tx, requirement *model.Requirement) (int64, error) {
	sql := "INSERT INTO `pk_requirement` (`id`, `project_id`, `iteration_id`, `code`, `name`, `type`, `demander`, `priority`, `influence`, `owner`, `tracer`, `status`, `create_by`, `create_time`) VALUES (:id, :project_id, :iteration_id, :code, :name, :type, :demander, :priority, :influence, :owner, :tracer, :status, :create_by, :create_time)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, requirement)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementDao) EditRequirement(db *sqlx.Tx, requirement *model.Requirement) (int64, error) {
	setClause := "SET"
	if len(requirement.Name) != 0 {
		setClause += " `name` = :name,"
	}
	if requirement.Type.Value() != 0 {
		setClause += " `type` = :type,"
	}
	if len(requirement.Demander) != 0 {
		setClause += " `demander` = :demander,"
	}
	if requirement.Priority != 0 {
		setClause += " `priority` = :priority,"
	}
	if requirement.Influence.Value() != 0 {
		setClause += " `influence` = :influence,"
	}
	if requirement.Owner != 0 {
		setClause += " `owner` = :owner,"
	}
	if requirement.Tracer != 0 {
		setClause += " `tracer` = :tracer,"
	}
	if setClause == "SET" {
		return 0, errors.New("SQL 构造失败")
	}
	before, _ := strings.CutSuffix(setClause, ",")
	setClause = before
	sql := "UPDATE `pk_requirement` " + setClause + " WHERE `id` = :id"
	requirement.ModifiedTime = time.Now().UnixMilli()
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, requirement)
	if err != nil {
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (d *RequirementDao) ChangeStatus(db *sqlx.Tx, params map[string]interface{}) (int64, error) {
	sql := "UPDATE `pk_requirement` SET `status` = :status, `modified_time` = :modified_time, `modified_by` = :modified_by WHERE `id` = :id AND `delete` = 0"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, params)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementDao) RemoveRequirement(db *sqlx.Tx, params map[string]interface{}) (int64, error) {
	sql := "UPDATE `pk_requirement` SET `delete` = 1, `modified_time` = :modified_time, `modified_by` = :modified_by WHERE `id` = :id AND `delete` = 0"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, params)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementDao) FindByRequirementId(requirementId int64) (*model.Requirement, error) {
	db := database.GetDBInstance()
	requirement := &model.Requirement{}
	sql := "SELECT `id`, `project_id`, `iteration_id`, `code`, `name`, `type`, `demander`, `priority`, `influence`, `owner`, `tracer`, `status`, `delete`, `create_by`, `create_time`, `modified_by`, `modified_time` FROM `pk_requirement` WHERE `id` = ? AND `delete` = 0"
	err := db.Get(requirement, sql, requirementId)
	return requirement, err
}

func (d *RequirementDao) SearchList(projectId, iterationId int64, name string, requirementType int, demander string, priority int, influence int, owner, tracer int64, status int, offset, limit int64) ([]*model.Requirement, int64, error) {
	db := database.GetDBInstance()
	requirements := make([]*model.Requirement, 0)
	sql := "SELECT `id`, `project_id`, `iteration_id`, `code`, `name`, `type`, `demander`, `priority`, `influence`, `owner`, `tracer`, `status`, `delete`, `create_by`, `create_time`, `modified_by`, `modified_time` FROM `pk_requirement`"
	orderClause := " ORDER BY `create_time` DESC"
	limitClause := " LIMIT :offset, :limit"
	whereClause := ""
	if projectId != 0 {
		whereClause += " AND `project_id` = :project_id"
	}
	if iterationId != 0 {
		whereClause += " AND `iteration_id` = :iteration_id"
	}
	if len(name) != 0 {
		whereClause += " AND `name` LIKE CONCAT('%', :name, '%')"
	}
	if requirementType != 0 {
		whereClause += " AND `type` = :type"
	}
	if len(demander) != 0 {
		whereClause += " AND `demander` LIKE CONCAT('%', :demander, '%')"
	}
	if priority != 0 {
		whereClause += " AND `priority` = :priority"
	}
	if influence != 0 {
		whereClause += " AND `influence` = :influence"
	}
	if owner != 0 {
		whereClause += " AND `owner` = :owner"
	}
	if tracer != 0 {
		whereClause += " AND `tracer` = :tracer"
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
		"project_id":   projectId,
		"iteration_id": iterationId,
		"name":         name,
		"type":         requirementType,
		"demander":     demander,
		"priority":     priority,
		"influence":    influence,
		"owner":        owner,
		"tracer":       tracer,
		"status":       status,
		"offset":       offset,
		"limit":        limit,
	}
	funcName := pkg.GetRunningFuncName()
	log.Debugf("[%s] Execute SQL => %s\n", funcName, sql)
	rows, err := db.NamedQuery(sql, namedParams)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		r := &model.Requirement{}
		_ = rows.StructScan(r)
		requirements = append(requirements, r)
	}
	countSql := "SELECT COUNT(1) FROM `pk_requirement`" + whereClause
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
	return requirements, count, nil
}
