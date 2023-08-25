package dao

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/app/model"
	"pmkit/internal/pkg"
	"pmkit/internal/pkg/database"
)

type RequirementCommentDao struct {
}

func (d *RequirementCommentDao) SaveComment(db *sqlx.Tx, comment *model.RequirementComment) (int64, error) {
	sql := "INSERT INTO `pk_requirement_comment` (`id`, `requirement_id`, `comment`, `create_by`, `create_time`) VALUES (:id, :requirement_id, :comment, :create_by, :create_time)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, comment)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementCommentDao) RemoveComment(db *sqlx.Tx, id int64) (int64, error) {
	sql := "UPDATE `pk_requirement_comment` SET `delete` = 1 WHERE `id` = ?"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.Exec(sql, id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementCommentDao) ListByRequirementId(requirementId int64) ([]*model.RequirementComment, error) {
	db := database.GetDBInstance()
	comments := make([]*model.RequirementComment, 0)
	sql := "SELECT `id`, `requirement_id`, `comment`, `delete`, `create_by`, `create_time` FROM `pk_requirement_comment` WHERE `requirement_id` = ? AND `delete` = 0"
	rows, err := db.Queryx(sql, requirementId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := &model.RequirementComment{}
		_ = rows.StructScan(c)
		comments = append(comments, c)
	}
	return comments, nil
}
