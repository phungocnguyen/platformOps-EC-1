package models

type Baseline struct {
	Id   int
	Uid  string
	Name string
}

type Command struct {
	Id        int
	Cmd       string
	ExeOrder  int
	ControlId string
}
