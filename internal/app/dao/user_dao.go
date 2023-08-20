package dao

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"pmkit/internal/app/model"
	"pmkit/internal/pkg"
	"pmkit/internal/pkg/database"
	"strings"
)

type UserDao struct {
}

func (d *UserDao) SaveUser(db *sqlx.Tx, user *model.User) (int64, error) {
	sql := "INSERT INTO `pk_user` (`id`, `phone`, `password`, `username`, `create_time`, `active`) VALUES (:id, :phone, :password, :username, :create_time, :active)"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, user)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (d *UserDao) EditUser(db *sqlx.Tx, user *model.User) (int64, error) {
	setClause := "SET"
	if len(user.Phone) != 0 {
		setClause += " `phone` = :phone,"
	}
	if len(user.Password) != 0 {
		setClause += " `password` = :password,"
	}
	if len(user.Username) != 0 {
		setClause += " `username` = :username,"
	}
	if setClause == "SET" {
		return 0, errors.New("SQL 构造失败")
	}
	before, _ := strings.CutSuffix(setClause, ",")
	setClause = before
	sql := "UPDATE `pk_user` " + setClause + " WHERE `id` = :id"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.NamedExec(sql, user)
	if err != nil {
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (d *UserDao) ActiveSwitch(db *sqlx.Tx, active bool, id int64) (int64, error) {
	sql := "UPDATE `pk_user` SET `active` = ? WHERE `id` = ?"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	result, err := db.Exec(sql, active, id)
	if err != nil {
		return -1, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (d *UserDao) FindById(id int64) (*model.User, error) {
	db := database.GetDBInstance()
	user := &model.User{}
	sql := "SELECT `id`, `phone`, `username`, `create_time`, `active` FROM `pk_user` WHERE `id` = ?"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	err := db.Get(user, sql, id)
	return user, err
}

func (d *UserDao) SearchList(phone string, username string, active interface{}, offset int64, limit int64) ([]*model.User, int64, error) {
	db := database.GetDBInstance()
	var users = make([]*model.User, 0)
	sql := "SELECT `id`, `phone`, `username`, `create_time`, `active` FROM `pk_user`"
	limitClause := " LIMIT :offset, :limit"
	whereClause := ""
	if len(phone) != 0 {
		whereClause += " AND `phone` LIKE CONCAT('%', :phone, '%')"
	}
	if len(username) != 0 {
		whereClause += " AND `username` LIKE CONCAT('%', :username, '%')"
	}
	if active != nil {
		whereClause += " AND `active` = :active"
	}
	if len(whereClause) != 0 {
		after, _ := strings.CutPrefix(whereClause, " AND")
		whereClause = " WHERE " + after
	}
	sql = sql + whereClause + limitClause
	namedParams := map[string]interface{}{
		"phone":    phone,
		"username": username,
		"active":   active,
		"offset":   offset,
		"limit":    limit,
	}
	funcName := pkg.GetRunningFuncName()
	log.Debugf("[%s] Execute SQL => %s\n", funcName, sql)
	rows, err := db.NamedQuery(sql, namedParams)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		u := &model.User{}
		_ = rows.StructScan(u)
		users = append(users, u)
	}
	countSql := "SELECT COUNT(1) FROM `pk_user`" + whereClause
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
	return users, count, nil
}

func (d *UserDao) CheckLogin(phone string, password string) (*model.User, error) {
	db := database.GetDBInstance()
	sql := "SELECT `id`, `phone`, `username`, `create_time`, `active` FROM `pk_user` WHERE `phone` = ? AND `password` = ?"
	log.Debugf("[%s] Execute SQL => %s\n", pkg.GetRunningFuncName(), sql)
	var user = new(model.User)
	err := db.Get(user, sql, phone, password)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}
