package model

type ProjectStatus int

const (
	// Waiting 待建
	Waiting ProjectStatus = 1
	// InProcess 在建
	InProcess ProjectStatus = 2
	// Pause 暂停
	Pause ProjectStatus = 3
	// Stop 终止（项目未结项而停止）
	Stop ProjectStatus = 4
	// Finish 结项
	Finish ProjectStatus = 5
)

func (ps ProjectStatus) Value() int {
	switch ps {
	case Waiting:
		return 1
	case InProcess:
		return 2
	case Pause:
		return 3
	case Stop:
		return 4
	case Finish:
		return 5
	default:
		return 0
	}
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
