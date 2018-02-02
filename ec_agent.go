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

	fmt.Println("- Empowered by", models.EC_version)

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

	var manifest_results []models.EC_Manifest_Result

	baseline := getEC_Manifest(input)
	if len(baseline) < 1 {
		os.Exit(1)
	}

	fmt.Println("- Start executing commands")

	for _, manifest := range baseline {
		var b bytes.Buffer

		data := manifest.Command

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

		result_manifest := models.EC_Manifest_Result{
		                    models.EC_Manifest{manifest.Title, manifest.Command, manifest.Baseline},
			                s,
			                DateTimeNow()}

		manifest_results = append(manifest_results, result_manifest)

		fmt.Println("- Done executing", "[", manifest.Title, "]")

	}
	writeToFile(manifest_results, output)
	fmt.Println("- Done writing to", "[", output, "]")

}

func writeToFile(baseline []models.EC_Manifest_Result, output string) {
	s_1 := "##################################"
	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	for i := range baseline {
		fmt.Fprintf(file, "\n%v", s_1)
		fmt.Fprintf(file, "\nTitle:    %v", baseline[i].Title)
		fmt.Fprintf(file, "\nBaseline: %v", baseline[i].Baseline)
		fmt.Fprintf(file, "\nDate Exc: %v", baseline[i].DateExe)
		fmt.Fprintf(file, "\nCommand:  %v", baseline[i].Command)
		fmt.Fprintf(file, "\nVersion:  %v", models.EC_version)
		fmt.Fprintf(file, "\n%v\n", s_1)
		fmt.Fprintf(file, "\n%v\n", baseline[i].Output)
	}

}

func DateTimeNow() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
}

