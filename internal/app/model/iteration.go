package model

type Iteration struct {
	BaseModel
	ProjectId int64  `db:"project_id" json:"projectId"`
	Name      string `db:"name" json:"name"`
	Remark    string `db:"remark" json:"remark"`
	StartDate int64  `db:"start_date" json:"startDate"`
	EndDate   int64  `db:"end_date" json:"endDate"`
	Delete    bool   `db:"delete" json:"delete"`
}
