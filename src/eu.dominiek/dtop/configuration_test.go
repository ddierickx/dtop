package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDTopConfigurationSerializationRoundtrip(t *testing.T) {
	testUser := NewDTopUser("ho", "dor")
	users := []DTopUser{*testUser}
	cfg := NewDTopConfiguration("name", "description", users, "static", 12345)
	jsonBytes, _ := SerializeDTopConfigurationToJson(cfg)
	deserializedCfg, _ := DeserializeDTopConfigurationFromJson(jsonBytes)

	if cfg.IsAuth() != (len(users) > 0) {
		panic(fmt.Sprintf("IsAuth call returned an invalid result, expected %t but was %t.", cfg.IsAuth, (len(users) > 0)))
	}

	if !reflect.DeepEqual(cfg, deserializedCfg) {
		panic("Serialization/deserialization failure.")
	}
}

func checkValidity(cfg *DTopConfiguration, valid bool) {
	if result, _ := cfg.IsValid(); result != valid {
		panic(fmt.Sprintf("Expected configuration to be valid=%t, but was %t", valid, result))
	}
}

func TestValidateDTopConfiguration(t *testing.T) {
	testUser := NewDTopUser("ho", "dor")
	users := []DTopUser{*testUser}
	checkValidity(NewDTopConfiguration("name", "description", users, "/tmp", 12345), true)
	checkValidity(NewDTopConfiguration("", "description", users, "/tmp", 12345), false)
	checkValidity(NewDTopConfiguration("name", "", users, "/tmp", 12345), false)
	checkValidity(NewDTopConfiguration("name", "description", users, "/tmp", 0), false)
	checkValidity(NewDTopConfiguration("name", "description", users, "phony", 8080), false)
}
