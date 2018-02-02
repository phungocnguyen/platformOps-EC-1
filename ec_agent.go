package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"platformOps-EC/models"
	"platformOps-EC/services"
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



func main() {

	fmt.Println("- Empowered by",Global_Version)

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

		if err := services.Execute(&b,
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
		fmt.Fprintf(file, "\nTitle:    %v", baseline[i].GetTitle())
		fmt.Fprintf(file, "\nBaseline: %v", baseline[i].GetBaseline())
		fmt.Fprintf(file, "\nDate Exc: %v", baseline[i].GetDateExe())
		fmt.Fprintf(file, "\nCommand:  %v", baseline[i].GetCommand())
		fmt.Fprintf(file, "\nVersion:  %v", Global_Version)
		fmt.Fprintf(file, "\n%v\n", s_1)
		fmt.Fprintf(file, "\n%v\n", baseline[i].GetOutput())
	}

}

func DateTimeNow() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
}

var Global_Version = "ec_agent_v.0.2"