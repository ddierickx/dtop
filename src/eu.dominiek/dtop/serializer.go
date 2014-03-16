package main

import(
	"encoding/json"
)

// Server->client serialization protocol function template.
type EventSerializer func(Event) []byte

// Default JSON implementation of EventSerializer.
func jsonEventSerializer(event Event) []byte {
	serialized, err := json.Marshal(event)

	if err != nil {
		panic("Serialization error: " + err.Error())
	}

	return serialized
}