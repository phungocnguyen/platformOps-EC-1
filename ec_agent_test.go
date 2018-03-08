package main

import (
	"github.com/BurntSushi/toml"
	"github.com/codeskyblue/go-sh"
	"log"
	"strings"
	"testing"
)

func TestCommandExecutionWithVariables(t *testing.T) {

	session := sh.NewSession()
	session.SetEnv("BUILD_ID", "123")

	out, _ := session.Command("bash", "-c", "echo $BUILD_ID").Output()

	if strings.Compare(strings.TrimSuffix(string(out), "\n"), "123") != 0 {

		t.Errorf("expected 123 got %s ", string(out))

	}
}

func TestLoadConfigIntoSession(t *testing.T) {
	session := sh.NewSession()
	var config map[string]string
	configFile := "test-data/ec-config.toml"

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}
	loadConfigiIntoSession(session, configFile)

	for k, v := range config {
		if session.Env[k] != v {
			t.Errorf("expected session value %s for key %s, but got %s ", v, k, session.Env[k])
		}
	}

}
