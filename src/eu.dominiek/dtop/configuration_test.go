package main

import (
	"testing"
	"fmt"
	"reflect"
)

func TestDTopConfigurationSerializationRoundtrip(t *testing.T) {
	testUser := NewDTopUser("ho", "dor")
	users := []DTopUser { *testUser }
    cfg := NewDTopConfiguration("name", "description", users, "static", 12345)
    jsonBytes, _ := SerializeDTopConfigurationToJson(cfg)
    deserializedCfg, _ := DeserializeDTopConfigurationFromJson(jsonBytes)
 
    if !reflect.DeepEqual(cfg, deserializedCfg) {
    	panic("Serialization/deserialization failure.")
    }
}

func checkValidity(cfg *DTopConfiguration, valid bool) {
	if result, _ := ValidateDTopConfiguration(cfg); result != valid {
		panic(fmt.Sprintf("Expected configuration to be valid=%t, but was %t", valid, result))
	}
}

func TestValidateDTopConfiguration(t *testing.T) {
	testUser := NewDTopUser("ho", "dor")
	users := []DTopUser { *testUser }
    checkValidity(NewDTopConfiguration("name", "description", users, "/tmp", 12345), true)
    checkValidity(NewDTopConfiguration("", "description", users, "/tmp", 12345), false)
    checkValidity(NewDTopConfiguration("name", "", users, "/tmp", 12345), false)
    checkValidity(NewDTopConfiguration("name", "description", users, "/tmp", 0), false)
    checkValidity(NewDTopConfiguration("name", "description", users, "phony", 8080), false)
}