package go_pager

import (
	"encoding/json"
	"os"
)

type Contact struct {
	Name string
	Number string
}

type Config struct {
    Contacts []Contact
}

func (config *Config) LoadFromFile(path string) {
	reader, _ := os.Open(path)
	decoder := json.NewDecoder(reader)
	decoder.Decode(&config) }
