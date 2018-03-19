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
	"path"
	"path/filepath"
	"platformOps-EC/converter"
	"platformOps-EC/models"
	"platformOps-EC/services"
	"strings"
)

/*
This is a av evidence collection agent:

Usage

-i Input file
-o Output file
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

	return services.GetManifestFromMaster(url)
}

func CollectEvidence(baseline []models.ECManifest) []models.ECResult {

	var ecResults []models.ECResult
	for _, manifest := range baseline {
		var errorOutputs, resultOutputs []string

		data := manifest.Command
		if len(data) > 0 {
			fmt.Printf("- Executing [%v]\n", manifest.Title)
			for c := range data {
				var b bytes.Buffer

				result := strings.Split(data[c], "|")
				array := make([]*exec.Cmd, len(result))
				for i := range result {
					s := strings.TrimSpace(result[i])
					commands := strings.Split(s, " ")
					args := commands[1:]
					services.WrapperCliVarsToEnvVars(args)
					array[i] = exec.Command(commands[0], args...)
				}

				errorOutput := services.Execute(&b, array)

				resultOutputs = append(resultOutputs, b.String())

				if errorOutput != "" {
					errorOutputs = append(errorOutputs, errorOutput)
				}
			}

		}
		resultManifest := models.ECResult{
			models.ECManifest{manifest.ReqId, manifest.Title, manifest.Command, manifest.BaselineUid, manifest.ControlUid},
			services.GetHostNameExec(),
			resultOutputs,
			errorOutputs,
			services.DateTimeNow()}

		ecResults = append(ecResults, resultManifest)

	}
	return ecResults
}

func writeToFile(baseline []models.ECResult, output string, resultType string, isJson bool) {
	hashString := "##################################"
	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	switch isJson {
	case true:
		fmt.Fprint(file, models.ToJson(baseline))

	case false:
		for i := range baseline {
			fmt.Fprintf(file, "\n%v", hashString)
			fmt.Fprintf(file, "\nVersion:  %v", models.ECVersion)
			fmt.Fprintf(file, "\nReq Id:   %v", baseline[i].ReqId)
			fmt.Fprintf(file, "\nTitle:    %v", baseline[i].Title)
			fmt.Fprintf(file, "\nBaseline: %v", baseline[i].BaselineUid)
			fmt.Fprintf(file, "\nControl: %v", baseline[i].ControlUid)
			fmt.Fprintf(file, "\nDate Exc: %v", baseline[i].DateExe)
			fmt.Fprintf(file, "\nHost Exc: %v", baseline[i].HostExec)
			fmt.Fprintf(file, "\n%v:", "Command")
			for c := range baseline[i].Command {
				fmt.Fprintf(file, "\n        [%v]", baseline[i].Command[c])
			}
			fmt.Fprintf(file, "\n%v\n", hashString)

			switch resultType {
			case "stdOutput":
				fmt.Fprintf(file, "\n%v\n", strings.Join(baseline[i].StdOutput, "\n"))
			case "stdErrOutput":
				fmt.Fprintf(file, "\n%v\n", strings.Join(baseline[i].StdErrOutput, "\n"))
			case "both":
				fmt.Fprintf(file, "\n%v\n", strings.Join(baseline[i].StdOutput, "\n"))
				fmt.Fprintf(file, "\n%v\n", strings.Join(baseline[i].StdErrOutput, "\n"))
			}
		}
	}
}

func getFileName(output string, outputType string) string {
	var fileName string
	switch outputType {
	case "error":
		fileName = filepath.Join(filepath.Dir(output), outputType+"_"+filepath.Base(output))
	case "json":
		ext := path.Ext(output)
		fileName = output[0:len(output)-len(ext)] + ".json"
	default:
		fileName = output
	}

	return fileName
}

func main() {

	var input, output, config, mode string

	fmt.Println("- Empowered by", models.ECVersion)

	flag.StringVar(&input, "i", "", "Input manifest json file. If missing, program will exit.")
	flag.StringVar(&output, "o", "output.txt", "Execution output location.")
	flag.StringVar(&config, "c", "config.toml", "External configuration location.")
	flag.StringVar(&mode, "m", "local", "Run as Web agent or local CLI agent. -m w as Web agent. Default local CLI agent. ")

	flag.Parse()

	env := services.LoadConfig(config)
	defer func() {
		os.Clearenv()
		for k, _ := range env {
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
		configMap := services.LoadConfig(config)
		services.SetEnvConfig(configMap)
		processManifest(input, output, mode)
		services.UnsetEnvConfig(configMap)

	}

}

func processManifest(input string, output string, mode string) {

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

	ecResults := CollectEvidence(baseline)

	// write result to output file
	writeToFile(ecResults, output, "stdOutput", false)
	fmt.Printf("- Done writing to [%v]\n", output)

	// write error to error output file
	errorFile := getFileName(output, "error")
	writeToFile(ecResults, errorFile, "stdErrOutput", false)
	fmt.Printf("- Done writing error to [%v]\n", errorFile)

	// write both StdOut & StdErrOut to  json file
	bothFile := getFileName(output, "json")
	writeToFile(ecResults, bothFile, "", true)
	fmt.Printf("- Done writing json to [%v]\n", bothFile)

}
