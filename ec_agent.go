package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"platformOps-EC/models"
	"strings"
	"time"
)

type Control struct {
	title    string
	command  string
	output   string
	dateExe  string
	baseline string
}

func (c Control) SetOutput(output string) {
	c.output = output
}

func (c Control) GetTitle() string {
	return c.title
}

func (c Control) GetCommand() string {
	return c.command
}

func (c Control) GetOutput() string {
	return c.output
}

func (c Control) GetDateExe() string {
	return c.dateExe
}

func (c Control) GetBaseline() string {
	return c.baseline
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func getEC_Manifest(manifest string) []models.EC_Manifest {
	fmt.Println("- Parsing manifest", "[", manifest, "]")
	raw, err := ioutil.ReadFile(manifest)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []models.EC_Manifest
	errj := json.Unmarshal(raw, &c)
	if errj != nil {
		fmt.Println("error parsing json input", err)
	}
	return c
}

func Execute(output_buffer *bytes.Buffer, stack []*exec.Cmd) (err error) {
	var error_buffer bytes.Buffer
	pipe_stack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdin_pipe, stdout_pipe := io.Pipe()
		stack[i].Stdout = stdout_pipe
		stack[i].Stderr = &error_buffer
		stack[i+1].Stdin = stdin_pipe
		pipe_stack[i] = stdout_pipe
	}
	stack[i].Stdout = output_buffer
	stack[i].Stderr = &error_buffer
	if err := call(stack, pipe_stack); err != nil {
		log.Fatalln("Encounter Error", string(error_buffer.Bytes()), err)
	}
	return err
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

func main() {

	var input, output string

	flag.StringVar(&input, "i", "", "Input manifest json file. If missing, program will exit.")
	flag.StringVar(&output, "o", "output.txt", "Execution output location.")
	flag.Parse()

	if input == "" {
		fmt.Println("Missing input manifest. Program will exit.")
		os.Exit(1)
	}

	if output == "output.txt" {
		fmt.Println("Default to output.txt")

	}

	var controlO []Control

	baseline := getEC_Manifest(input)
	if len(baseline) < 1 {
		os.Exit(1)
	}

	fmt.Println("- Start executing commands")

	for _, manifest := range baseline {
		var b bytes.Buffer

		data := manifest.GetCommand()

		result := strings.Split(data, "|")
		array := make([]*exec.Cmd, len(result))
		for i := range result {
			s := strings.TrimSpace(result[i])
			commands := strings.Split(s, " ")
			args := commands[1:len(commands)]

			array[i] = exec.Command(commands[0], args...)

		}

		if err := Execute(&b,
			array); err != nil {
			log.Fatalln(err)
		}

		s := b.String()

		co := Control{title: manifest.GetTitle(),
			command:  manifest.GetCommand(),
			output:   s,
			dateExe:  DateTimeNow(),
			baseline: manifest.GetBaseline()}

		controlO = append(controlO, co)

		fmt.Println("- Done executing", "[", manifest.GetTitle(), "]")

	}
	writeToFile(controlO, output)
	fmt.Println("- Done writing to", "[", output, "]")

}

func writeToFile(baseline []Control, output string) {
	s_1 := "##################################"
	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	for i := range baseline {
		fmt.Fprintf(file, "\n%v", s_1)
		fmt.Fprintf(file, "\n%v", baseline[i].GetTitle())
		fmt.Fprintf(file, "\n%v", baseline[i].GetBaseline())
		fmt.Fprintf(file, "\n%v", baseline[i].GetDateExe())
		fmt.Fprintf(file, "\n%v", baseline[i].GetCommand())
		fmt.Fprintf(file, "\n%v\n", s_1)
		//fmt.Fprintf(file,"\n%v\n", baseline[i].GetCommand())
		fmt.Fprintf(file, "\n%v\n", baseline[i].GetOutput())
	}

}

func DateTimeNow() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
}

var Global_Version = "ec_agent_0.1"