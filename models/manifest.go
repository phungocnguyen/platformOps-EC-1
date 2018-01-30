package models

import (
	"encoding/json"
	"fmt"
	"os"
)

type EC_Manifest struct {
	Title    string `json:"title"`
	Command  string `json:"command"`
	Baseline string `json:"baseline"`
}

func (c EC_Manifest) GetCommand() string {
	return c.Command
}

func (c EC_Manifest) GetTitle() string {
	return c.Title
}

func (c EC_Manifest) GetBaseline() string {
	return c.Baseline
}

func (p EC_Manifest) ToString() string {
	return ToJson(p)
}

func ToJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}