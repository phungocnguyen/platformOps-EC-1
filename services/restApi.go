package services

import (
	"bytes"
	"ec/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"platformOps-EC/models"
	"time"
)

func PostECResultsToMaster(url string, ecResults []models.ECResult) string {

	var jsonStr = []byte(models.ToJson(ecResults))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Printf("Error posting to master %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func SendResultToMaster(url string, ecResults []models.ECResult) string {
	var jsonStr = []byte(models.ToJson(ecResults))

	var myClient = &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Date", services.DateTimeNow())
	req.Header.Add("User-Agent", models.ECVersion)
	resp, err := myClient.Do(req)
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func GetManifestFromMaster(url string) []models.ECManifest {

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
