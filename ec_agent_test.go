package main

import (
	"testing"
	"os/exec"
	"log"
	"fmt"
	"bytes"
)

func TestCommandExecutionWithVariables(t *testing.T){

		var1 := "FOO=BAR"
		foo:="bar"
		cmd := exec.Command("echo", foo)
		cmd.Env =[]string{var1}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %v\n", out.String())


}
