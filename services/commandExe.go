package services

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Execute(outputBuffer *bytes.Buffer, stack []*exec.Cmd) (errorOutput string) {
	var errorBuffer bytes.Buffer
	pipeStack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdinPipe, stdoutPipe := io.Pipe()
		stack[i].Stdout = stdoutPipe
		stack[i].Stderr = &errorBuffer
		stack[i+1].Stdin = stdinPipe
		pipeStack[i] = stdoutPipe
	}
	stack[i].Stdout = outputBuffer
	stack[i].Stderr = &errorBuffer

	if err := call(stack, pipeStack); err != nil {
		fmt.Printf("Encountered Error %v, %v\n", string(errorBuffer.Bytes()), err)
		errorOutput = fmt.Sprintf("%v\n", err)
	}

	if errorOutput != "" {
		return fmt.Sprintf("%v\n%v", errorOutput, string(errorBuffer.Bytes()))
	}

	return

}

func call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}

func GetHostNameExec() string {
	hostname, err := os.Hostname()

	if err != nil {
		log.Printf("Error getting hostname %v", err)
		return "NA"
	}

	return hostname
}

func SetEnvConfig(config map[string]string) {
	for k, v := range config {
		os.Setenv(k, v)
	}
}

func UnsetEnvConfig(config map[string]string) {
	for k := range config {
		os.Unsetenv(k)
	}
}

func PrintEnv(config map[string]string) {
	for k := range config {
		fmt.Printf("Key=%v, Value=%v\n", k, os.Getenv(k))
	}
}

func LoadConfig(configFile string) (config map[string]string) {
	fmt.Printf("- Loading configs [%v]\n", configFile)
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Error loading the config: ", err)
		fmt.Printf("Error loading the config: %s/n ", configFile)
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return
}

func WrapperCliVarsToEnvVars(args []string) {

	for k := range args {
		if atIdx := strings.Index(args[k], "$"); atIdx == 0 {
			key := strings.Replace(args[k], "$", "", 1)

			if val, ok := os.LookupEnv(key); ok {
				args[k] = val
			}
		}
	}
}

func PrintAllEnv() {
	fmt.Println(os.Environ())
}
