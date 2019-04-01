package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type NavetimeRes struct {
	Type     string             `json:"type"`
	Features []NavetimeFeatures `json:"features"`
}

type NavetimeFeatures struct {
	Type     string            `json:"type"`
	Geometry *NavetimeGeometry `json:"geometry"`
}

type NavetimeGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

func main() {

	f, e := os.Open("data.json")
	if e != nil {
		fmt.Println(e)
		return
	}
	buf := bufio.NewReader(f)
	line, _ := buf.ReadString('\n')
	line = strings.TrimSpace(line)

	navetimeRes := &NavetimeRes{}
	json.Unmarshal([]byte(line), &navetimeRes)

	send_body, _ := json.Marshal(navetimeRes)

	fmt.Println(string(send_body))
}
