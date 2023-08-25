package dao

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/app/model"
	"pmkit/internal/pkg"
	"pmkit/internal/pkg/database"
)

type RequirementTrackDao struct {
}

func (d *RequirementTrackDao) SaveTrack(db *sqlx.Tx, track *model.RequirementTrack) (int64, error) {
	sql := "INSERT INTO `pk_requirement_track` (`id`, `requirement_id`, `status`, `create_by`, `create_time`) VALUES (:id, :requirement_id, :status, :create_by, :create_time)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, track)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *RequirementTrackDao) ListByRequirementId(requirementId int64) ([]*model.RequirementTrack, error) {
	db := database.GetDBInstance()
	tracks := make([]*model.RequirementTrack, 0)
	sql := "SELECT `id`, `requirement_id`, `status`, `create_by`, `create_time` FROM `pk_requirement_track` WHERE `requirement_id` = ? ORDER BY `create_time`"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	rows, err := db.Queryx(sql, requirementId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := &model.RequirementTrack{}
		_ = rows.StructScan(t)
		tracks = append(tracks, t)
	}
	return tracks, nil
}
