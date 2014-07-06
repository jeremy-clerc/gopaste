package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
)

func TestGenerateID(t *testing.T) {
	rand.Seed(0)
    validID := regexp.MustCompile(
        fmt.Sprintf("^([a-zA-Z0-9]{%d})$", idLength))

    generatedID := GenerateID()
	if !validID.MatchString(generatedID) {
       t.Errorf("Generated ID for is not compliant to a-zA-Z0-9. Got: " + generatedID) 
    }
}
