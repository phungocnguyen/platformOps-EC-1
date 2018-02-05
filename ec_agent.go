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
	"path/filepath"
	"platformOps-EC/models"
	"platformOps-EC/services"
	"strings"
	"time"
)


func getECManifest(manifest string) []models.ECManifest {
	fmt.Printf("- Parsing manifest [%v]\n", manifest)
	raw, err := ioutil.ReadFile(manifest)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []models.ECManifest
	errUnmarshal := json.Unmarshal(raw, &c)
	if errUnmarshal != nil {
		fmt.Println("error parsing json input", err)
	}
	return c
}



func main() {

	fmt.Println("- Empowered by", models.ECVersion)

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

	var manifestResults []models.ECManifestResult

	var manifestErrors []models.ECManifestResult

	baseline := getECManifest(input)
	if len(baseline) < 1 {
		os.Exit(1)
	}

	fmt.Println("- Start executing commands")

	for _, manifest := range baseline {
		var b bytes.Buffer

		data := manifest.Command

		fmt.Printf("- Executing [%v]\n", manifest.Title)

		result := strings.Split(data, "|")
		array := make([]*exec.Cmd, len(result))
		for i := range result {
			s := strings.TrimSpace(result[i])
			commands := strings.Split(s, " ")
			args := commands[1:len(commands)]

			array[i] = exec.Command(commands[0], args...)

		}

		errorOutput := services.Execute(&b, array)

		/*if err != nil {
			log.Fatalln(err)
		}
		*/

		s := b.String()

		resultManifest := models.ECManifestResult{
		                    models.ECManifest{manifest.ReqId,manifest.Title, manifest.Command, manifest.Baseline},
			                s,
			                dateTimeNow()}

		manifestResults = append(manifestResults, resultManifest)

		if errorOutput!= "" {
			errorManifest := models.ECManifestResult{
				models.ECManifest{manifest.ReqId, manifest.Title, manifest.Command, manifest.Baseline},
				errorOutput,
				dateTimeNow()}
			manifestErrors = append(manifestErrors, errorManifest)
		}



	}

	writeToFile(manifestResults, output)
	fmt.Printf("- Done writing to [%v]\n", output)

	if len(manifestErrors) > 0 {
		errorFile := getErrorFileName(output)
		writeToFile(manifestErrors, errorFile)
		fmt.Printf("- Done writing error to [%v]\n", errorFile)
	}
}

func writeToFile(baseline []models.ECManifestResult, output string) {
	hashString := "##################################"
	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	for i := range baseline {
		fmt.Fprintf(file, "\n%v", hashString)
		fmt.Fprintf(file, "\nReq Id:   %v", baseline[i].ReqId)
		fmt.Fprintf(file, "\nTitle:    %v", baseline[i].Title)
		fmt.Fprintf(file, "\nBaseline: %v", baseline[i].Baseline)
		fmt.Fprintf(file, "\nDate Exc: %v", baseline[i].DateExe)
		fmt.Fprintf(file, "\nCommand:  %v", baseline[i].Command)
		fmt.Fprintf(file, "\nVersion:  %v", models.ECVersion)
		fmt.Fprintf(file, "\n%v\n", hashString)
		fmt.Fprintf(file, "\n%v\n", baseline[i].Output)
	}

}

func dateTimeNow() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
}

func getErrorFileName(output string) string{
	return filepath.Join(filepath.Dir(output), "error_"+filepath.Base(output))
}

