package models

type Baseline struct {
        Id   int
        Name string
}
type Command struct {
        Id              int
        Cmd             string
        ExeOrder        int
        ControlId       int
}