package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	filename := "foo"
	configName := "foo"
	version := "foo"
	runtimeStage := "foo"
	config, err := NewConfig(filename, configName, version, runtimeStage)
	if err == nil {
		t.Fail()
	}
	if config == nil {
		t.Fail()
	}
}
