package model

type ProjectStatus int

const (
	// Waiting 待建
	Waiting ProjectStatus = iota + 1
	// InProcess 在建
	InProcess
	// Pause 暂停
	Pause
	// Stop 终止（项目未结项而停止）
	Stop
	// Finish 结项
	Finish
)

func (ps ProjectStatus) String() string {
	return [...]string{"待建", "在建", "暂停", "终止", "结项"}[ps]
}

type Project struct {
	BaseModel
	Name        string        `db:"name" json:"name"`
	Employer    string        `db:"employer" json:"employer"`
	Contractor  string        `db:"contractor" json:"contractor"`
	Supervisor  string        `db:"supervisor" json:"supervisor"`
	Executor    string        `db:"executor" json:"executor"`
	Description string        `db:"description" json:"description"`
	StartDate   int64         `db:"start_date" json:"startDate"`
	EndDate     int64         `db:"end_date" json:"endDate"`
	Status      ProjectStatus `db:"status" json:"status"`
	Delete      bool          `db:"delete" json:"delete"`
}
