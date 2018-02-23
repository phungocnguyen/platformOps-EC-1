package models


type Control struct {
	Id          int
	ReqId       int
	CisId       string
	Category    string
	Requirement string
	Discussion  string
	CheckText   string
	FixText     string
	RowDesc     string
	BaselineId  int
}
