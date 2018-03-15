package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"platformOps-EC/converter"
	"platformOps-EC/models"
	"time"
	"strings"
	"platformOps-EC/services"
	"os/exec"
	"bytes"
)

/*
This is a av evidence collection agent:

Usage

-i Input file
-c configuration file
-m mode:
	- toJson (convert excel baseline to json manifest)
	- local (collect evidence using local input json manifest)
	- web (collect evidence using manifest from input endpoint url, send json result back to master)

*/

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

	body, err := ioutil.ReadAll(resp.Body)

	var baseline []models.ECManifest

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(body, &baseline)

	return baseline
}

func loadConfigiIntoSession(configFile string) map[string]string {
	fmt.Printf("- Loading configs [%v]\n", configFile)


	var config map[string]string
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Error loading the config: ", configFile)
		fmt.Printf("Error loading the config: %s/n ", configFile)
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	for k, v := range config {
		os.Setenv(k, v)
	}
	return config


}
func executeCommands(baseline []models.ECManifest) ([]models.ECManifestResult, []models.ECManifestResult) {
	var manifestErrors, manifestResults []models.ECManifestResult

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

		s := b.String()

		resultManifest := models.ECManifestResult{

			models.ECManifest{manifest.ReqId, manifest.Title, manifest.Command, manifest.Baseline},
			s,
			dateTimeNow()}

		manifestResults = append(manifestResults, resultManifest)

		if errorOutput != "" {
			errorManifest := models.ECManifestResult{
				models.ECManifest{manifest.ReqId, manifest.Title, manifest.Command, manifest.Baseline},
				errorOutput,
				dateTimeNow()}
			manifestErrors = append(manifestErrors, errorManifest)
		}
	}
	return manifestResults, manifestErrors
}
func dateTimeNow() string {
	return time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
}
func getErrorFileName(output string) string {
	return filepath.Join(filepath.Dir(output), "error_" + filepath.Base(output))
}

func main() {

	var input, output,config,   mode string


	fmt.Println("- Empowered by", models.ECVersion)

	flag.StringVar(&input, "i", "", "Input manifest json file. If missing, program will exit.")
	flag.StringVar(&output, "o", "output.txt", "Execution output location.")
	flag.StringVar(&config, "c", "config.toml", "External configuration location.")
	flag.StringVar(&mode, "m", "local", "Run as Web agent or local CLI agent. -m w as Web agent. Default local CLI agent. ")

	flag.Parse()

	env:= loadConfigiIntoSession(config)
	defer func() {
		os.Clearenv()
		for k, _ := range env{
			os.Unsetenv(k)
		}
	}()
	if input == "" {
		fmt.Println("Missing input manifest. Program will exit.")
		os.Exit(1)
	}

	if output == "output.txt" {
		fmt.Println("Default to output.txt")

	}

	switch mode {
	case "toJson":
		converter.ToJson(input, output)
	default:
		processManifest(input, output, mode)
	}

}

func processManifest(input string, output string, mode string, ) {

	var baseline []models.ECManifest

	if mode == "local" {
		baseline = getECManifest(input)
	} else if mode == "web" {
		baseline = getJsonManifestFromMaster(input)
	}

	if len(baseline) < 1 {
		fmt.Println("Baseline does not have controls.  Program will exit")
		os.Exit(1)
	}

	fmt.Println("- Start executing commands")

	manifestResults, manifestErrors := executeCommands(baseline)

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