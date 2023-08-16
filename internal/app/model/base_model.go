package model

type BaseModel struct {
	Id           int64 `db:"id" json:"id"`
	CreateBy     int64 `db:"create_by" json:"createBy"`
	CreateTime   int64 `db:"create_time" json:"createTime"`
	ModifiedBy   int64 `db:"modified_by" json:"modifiedBy"`
	ModifiedTime int64 `db:"modified_time" json:"modifiedTime"`
}
