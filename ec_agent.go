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
	"net/http"
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

func getJsonManifestFromMaster(url string) []models.ECManifest {

	var myClient = &http.Client{Timeout: 10 * time.Second}

	resp, err := myClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//decoder := json.NewDecoder(resp.Body)
	//fmt.Println(decoder.Decode(&baseline))

	body, err := ioutil.ReadAll(resp.Body)

	var baseline []models.ECManifest

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(body, &baseline)
	return baseline
}

func CollectEvidence (baseline []models.ECManifest)  ([]models.ECManifestResult, []models.ECManifestResult)  {

	var manifestResults, manifestErrors []models.ECManifestResult
	for _, manifest := range baseline {
		var b bytes.Buffer
		var errorOutput string
		s := ""
		data := manifest.Command
		if data != "" {
			fmt.Printf("- Executing [%v]\n", manifest.Title)

			result := strings.Split(data, "|")
			array := make([]*exec.Cmd, len(result))
			for i := range result {
				s := strings.TrimSpace(result[i])
				commands := strings.Split(s, " ")
				args := commands[1:len(commands)]

				array[i] = exec.Command(commands[0], args...)

			}

			errorOutput = services.Execute(&b, array)

			/*if err != nil {
				log.Fatalln(err)
			}
			*/

			s = b.String()
		}
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
	return manifestResults, manifestErrors
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


func main() {

	fmt.Println("- Empowered by", models.ECVersion)

	var input, output, mode string

	flag.StringVar(&input,  "i", "", "Input manifest json file. If missing, program will exit.")
	flag.StringVar(&output, "o", "output.txt", "Execution output location.")
	flag.StringVar(&mode,   "m", "local", "Run as Web agent or local CLI agent. -m w as Web agent. Default local CLI agent. ")
	flag.Parse()

	if input == "" {
		fmt.Println("Missing input manifest. Program will exit.")
		os.Exit(1)
	}

	if output == "output.txt" {
		fmt.Println("Default to output.txt")

	}

	// populate baseline by requesting manifest from master
	// or input manifest at command line argument.
	var baseline []models.ECManifest
	if mode == "local" {
		baseline = getECManifest(input)
	} else {
		baseline = getJsonManifestFromMaster(input)
	}
	if len(baseline) < 1 {
		os.Exit(1)
	}

	fmt.Println("- Start executing commands")
	//collecting evidence by executing commands
	manifestResults, manifestErrors := CollectEvidence(baseline)

	// write result to output file
	writeToFile(manifestResults, output)
	fmt.Printf("- Done writing to [%v]\n", output)

	// write error to error output file
	if len(manifestErrors) > 0 {
		errorFile := getErrorFileName(output)
		writeToFile(manifestErrors, errorFile)
		fmt.Printf("- Done writing error to [%v]\n", errorFile)
	}
}

