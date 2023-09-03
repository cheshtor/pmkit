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

// TransferTo 需求状态转换校验
func (rs RequirementStatus) TransferTo(targetStatus RequirementStatus) bool {
	currentValue := rs.Value()
	targetValue := targetStatus.Value()
	// 目标状态编码必须在有效范围内
	if targetValue < 1 || targetValue > Delivery.Value() {
		return false
	}
	// 无转换
	if targetValue == currentValue {
		return true
	}
	// 目标状态编码小于当前状态编码，只有当前状态是 Verify 和 Checking 时允许
	if targetValue < currentValue && (currentValue != Verify.Value() || currentValue != Checking.Value()) {
		return false
	}
	// 状态之间严格递进，不允许跨越
	if targetValue-currentValue != 1 {
		return false
	}
	return true
}

// Requirement 需求
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

// RequirementContent 需求描述
type RequirementContent struct {
	RequirementId int64  `db:"requirement_id" json:"requirementId"`
	Content       string `db:"content" json:"content"`
}

// WholeRequirement 包含需求基本信息和需求描述详情的实体
type WholeRequirement struct {
	*Requirement
	*RequirementContent
}

func (wr *WholeRequirement) SeparateRequirement() *Requirement {
	r := new(Requirement)
	r.Id = wr.Id
	r.ProjectId = wr.ProjectId
	r.IterationId = wr.IterationId
	r.Code = wr.Code
	r.Name = wr.Name
	r.Type = wr.Type
	r.Demander = wr.Demander
	r.Priority = wr.Priority
	r.Influence = wr.Influence
	r.Owner = wr.Owner
	r.Tracer = wr.Tracer
	r.Status = wr.Status
	r.Delete = wr.Delete
	r.CreateTime = wr.CreateTime
	r.CreateBy = wr.CreateBy
	r.ModifiedTime = wr.ModifiedTime
	r.ModifiedBy = wr.ModifiedBy
	return r
}

func (wr *WholeRequirement) SeparateContent() *RequirementContent {
	c := new(RequirementContent)
	c.RequirementId = wr.Id
	c.Content = wr.Content
	return c
}

func (wr *WholeRequirement) Merge(requirement *Requirement, content *RequirementContent) {
	if requirement == nil {
		wr.Id = requirement.Id
		wr.ProjectId = requirement.ProjectId
		wr.IterationId = requirement.IterationId
		wr.Code = requirement.Code
		wr.Name = requirement.Name
		wr.Type = requirement.Type
		wr.Demander = requirement.Demander
		wr.Priority = requirement.Priority
		wr.Influence = requirement.Influence
		wr.Owner = requirement.Owner
		wr.Tracer = requirement.Tracer
		wr.Status = requirement.Status
		wr.Delete = requirement.Delete
		wr.CreateTime = requirement.CreateTime
		wr.CreateBy = requirement.CreateBy
		wr.ModifiedTime = requirement.ModifiedTime
		wr.ModifiedBy = requirement.ModifiedBy
	}
	if content != nil {
		wr.RequirementId = requirement.Id
		wr.Content = content.Content
	}
}

// RequirementTrack 需求状态追踪记录
type RequirementTrack struct {
	Id            int64             `db:"id" json:"id"`
	RequirementId int64             `db:"requirement_id" json:"requirementId"`
	Status        RequirementStatus `db:"status" json:"status"`
	CreateBy      int64             `db:"create_by" json:"createBy"`
	CreateTime    int64             `db:"create_time" json:"createTime"`
}

// RequirementComment 需求评论
type RequirementComment struct {
	Id            int64  `db:"id" json:"id"`
	RequirementId int64  `db:"requirement_id" json:"requirementId"`
	Comment       string `db:"comment" json:"comment"`
	Delete        bool   `db:"delete" json:"delete"`
	CreateBy      int64  `db:"create_by" json:"createBy"`
	CreateTime    int64  `db:"create_time" json:"createTime"`
}
