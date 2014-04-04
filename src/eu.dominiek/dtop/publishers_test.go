package main

import (
	"fmt"
	"testing"
)

func TestCleanDistro(t *testing.T) {
	assertEquals("Ubuntu 13.04", clean_distro("Ubuntu 13.04"))
	assertEquals("Ubuntu 13.04", clean_distro("\"Ubuntu 13.04\""))
}

func assertEquals(expected string, actual string) {
	if expected != actual {
		panic(fmt.Sprintf("Expected '%s' but was '%s'.", expected, actual))
	}
}
