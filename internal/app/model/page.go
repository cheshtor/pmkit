package model

type Page struct {
	PageNo     int64         `json:"pageNo"`
	PageSize   int64         `json:"pageSize"`
	TotalCount int64         `json:"totalCount"`
	Rows       []interface{} `json:"rows"`
}

type Paging struct {
	PageNo    int64       `json:"pageNo"`
	PageSize  int64       `json:"pageSize"`
	Condition interface{} `json:"condition"`
}
