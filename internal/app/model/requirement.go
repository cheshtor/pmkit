package model

// RequirementType 需求类型
type RequirementType int

const (
	Business    RequirementType = 1
	Extension   RequirementType = 2
	Improvement RequirementType = 3
)

func (rt RequirementType) Value() int {
	switch rt {
	case Business:
		return 1
	case Extension:
		return 2
	case Improvement:
		return 3
	default:
		return 0
	}
}

// RequirementInfluence 需求对系统的影响程度
type RequirementInfluence int

const (
	High   RequirementInfluence = 1
	Middle RequirementInfluence = 2
	Low    RequirementInfluence = 3
	None   RequirementInfluence = 4
)

func (ri RequirementInfluence) Value() int {
	switch ri {
	case High:
		return 1
	case Middle:
		return 2
	case Low:
		return 3
	case None:
		return 4
	default:
		return 0
	}
}

// RequirementStatus 需求状态
type RequirementStatus int

const (
	// Unclaimed 待认领（需求池中无人认领的需求）
	Unclaimed RequirementStatus = 1
	// Checking 待审核（审核通过的需求才能最终实施）
	Checking RequirementStatus = 2
	// Schedule 待排期
	Schedule RequirementStatus = 3
	// WaitCoding 待实施
	WaitCoding RequirementStatus = 4
	// Coding 实施中
	Coding RequirementStatus = 5
	// WaitDelivery 待交付
	WaitDelivery RequirementStatus = 6
	// Verify 待验收（验收不通过则复制当前需求作为子需求，重新排期）
	Verify RequirementStatus = 7
	// Delivery 已交付（验收通过）
	Delivery RequirementStatus = 8
)

func (rs RequirementStatus) Value() int {
	switch rs {
	case Unclaimed:
		return 1
	case Checking:
		return 2
	case Schedule:
		return 3
	case WaitCoding:
		return 4
	case Coding:
		return 5
	case WaitDelivery:
		return 6
	case Verify:
		return 7
	case Delivery:
		return 8
	default:
		return 0
	}
}

type Requirement struct {
	BaseModel
	ProjectId   int64                `db:"project_id" json:"projectId"`
	IterationId int64                `db:"iteration_id" json:"iterationId"`
	Code        string               `db:"code" json:"code"`
	Name        string               `db:"name" json:"name"`
	Type        RequirementType      `db:"type" json:"type"`
	Demander    string               `db:"demander" json:"demander"`
	Priority    int                  `db:"priority" json:"priority"`
	Influence   RequirementInfluence `db:"influence" json:"influence"`
	Owner       int64                `db:"owner" json:"owner"`
	Tracer      int64                `db:"tracer" json:"tracer"`
	Status      RequirementStatus    `db:"status" json:"status"`
	Delete      bool                 `db:"delete" json:"delete"`
}

type RequirementContent struct {
	RequirementId int64  `db:"requirement_id" json:"requirementId"`
	Content       string `db:"content" json:"content"`
}
