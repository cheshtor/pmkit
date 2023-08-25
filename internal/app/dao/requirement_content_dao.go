package dao

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/pkg"
)

type RequirementContentDao struct {
}

func (d *RequirementContentDao) SaveContent(db *sqlx.Tx, requirementId int64, content string) (int64, error) {
	sql := "INSERT INTO `pk_requirement_content` (`requirement_id`, `content`) VALUES (?, ?)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.Exec(sql, requirementId, content)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementContentDao) EditContent(db *sqlx.Tx, requirementId int64, content string) (int64, error) {
	sql := "UPDATE `pk_requirement_content` SET `content` = ? WHERE `requirement_id` = ?"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.Exec(sql, content, requirementId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
