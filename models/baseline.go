package models

import ()

type Baseline struct {
        Id   int
        Name string
}

func (b *Baseline) SetId(id int) {
        b.Id = id
}

func (b Baseline) GetId() int {
        return b.Id
}

func (b *Baseline) SetName(name string) {
        b.Name = name
}

func (b Baseline) GetName() string {
        return b.Name
}
