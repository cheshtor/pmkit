package model

type User struct {
	Id         int64  `db:"id" json:"id"`
	Phone      string `db:"phone" json:"phone"`
	Password   string `db:"password" json:"password"`
	Username   string `db:"username" json:"username"`
	CreateTime int64  `db:"create_time" json:"createTime"`
	Active     bool   `db:"active" json:"active"`
}
