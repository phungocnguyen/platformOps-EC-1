package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCommandExecutionWithVariables(t *testing.T) {

	os.Setenv("BUILD_ID", "123")

	out, _ := exec.Command("echo", os.ExpandEnv("$BUILD_ID")).Output()

	if strings.Compare(strings.TrimSuffix(string(out), "\n"), "123") != 0 {

		t.Errorf("expected 123 got %s ", string(out))

	}
}

func TestLoadConfigIntoSession(t *testing.T) {

	var config map[string]string
	configFile := "test-data/ec-config.toml"

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	loadConfigIntoSession(configFile)

	for k, v := range config {
		if os.Getenv(k) != v {
			t.Errorf("expected session value %s for key %s, but got %s ", v, k, os.Getenv(k))

		}
	}

}
