package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ECManifest struct {
	ReqId       int      `json:"reqId"`
	Title       string   `json:"title"`
	Command     []string `json:"command"`
	BaselineUid string   `json:"baseline"`
	ControlUid  string   `json:"control"`
}

type ECManifestResult struct {
	ECManifest
	Output  string `json:"output"`
	DateExe string `json:"dateExe"`
}

type ECResult struct {
	ECManifest
	HostExec     string   `json:"host"`
	StdOutput    []string `json:"stdOutput"`
	StdErrOutput []string `json:"stdErrOutput"`
	DateExe      string   `json:"dateExe"`
}

type BatchSubmision struct {
	Id           int
	BatchUid     string `json:"batchUid"`
	DateSubmit   string `json:"dateSubmit"`
	TimeSubmit   string `json:"timeSubmit"`
	UserSubmit   string `json:"userSubmit"`
	ResultSubmit []ECResult
}

func (p ECManifest) ToString() string {
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

func ToObject(jsonStr string, p interface{}) {

	raw := []byte(jsonStr)
	err := json.Unmarshal(raw, &p)

	if err != nil {
		log.Printf("Error converting json string to object %v", err)
	}

}
