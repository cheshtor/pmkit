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

var userDao dao.UserDao

type UserService struct {
}

func (s *UserService) SaveUser(user *model.User) (bool, error) {
	err := database.Run(func(db *sqlx.Tx) error {
		user.Id = pkg.GetId()
		user.Password = pkg.GenMd5(user.Password)
		user.CreateTime = time.Now().UnixMilli()
		user.Active = true
		affectedRows, err := userDao.SaveUser(db, user)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return fmt.Errorf("保存用户信息失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *UserService) EditUser(user *model.User) (bool, error) {
	_, _ = s.GetUserById(user.Id)
	err := database.Run(func(db *sqlx.Tx) error {
		if user.Password != "" {
			user.Password = pkg.GenMd5(user.Password)
		}
		affectedRows, err := userDao.EditUser(db, user)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return fmt.Errorf("更新用户信息失败")
		}
		return nil
	})
	return err == nil, err
}

func (s *UserService) ActiveSwitch(id int64, active bool) (bool, error) {
	exists, _ := s.GetUserById(id)
	if exists.Active == active {
		return true, nil
	}
	err := database.Run(func(db *sqlx.Tx) error {
		affectedRows, err := userDao.ActiveSwitch(db, active, id)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			status := "冻结"
			if active {
				status = "激活"
			}
			return fmt.Errorf("%s用户失败", status)
		}
		return nil
	})
	return err == nil, err
}

func (s *UserService) GetUserById(id int64) (*model.User, error) {
	user, err := userDao.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在。ID：%d", id)
	}
	return user, nil
}

func (s *UserService) GetUserList(condition *model.User, pageNo int64, pageSize int64) (*model.Page, error) {
	calcedPageNo, offset, calcedPageSize := pkg.ResolvePage(pageNo, pageSize)
	list, count, err := userDao.SearchList(condition.Phone, condition.Username, condition.Active, offset, calcedPageSize)
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

func (s *UserService) CheckLogin(phone string, password string) (string, error) {
	password = pkg.GenMd5(password)
	user, err := userDao.CheckLogin(phone, password)
	if err != nil {
		return "", err
	}
	if !user.Active {
		return "", errors.New("账号已被禁用，请联系管理员。")
	}
	return pkg.GenToken(user.Id)
}
