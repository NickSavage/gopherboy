package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// Define structs to match the JSON structure
type CPUTest struct {
	Name    string       `json:"name"`
	Initial CPUState     `json:"initial"`
	Final   CPUState     `json:"final"`
	Cycles  []CycleEvent `json:"cycles"`
}

type CPUState struct {
	A   uint8    `json:"a"`
	B   uint8    `json:"b"`
	C   uint8    `json:"c"`
	D   uint8    `json:"d"`
	E   uint8    `json:"e"`
	F   uint8    `json:"f"`
	H   uint8    `json:"h"`
	L   uint8    `json:"l"`
	PC  uint16   `json:"pc"`
	SP  uint16   `json:"sp"`
	RAM [][2]int `json:"ram"`
}

type CycleEvent []interface{}

// Function to parse the JSON
func ParseGameBoyTest(jsonData []byte) ([]CPUTest, error) {
	var tests []CPUTest
	err := json.Unmarshal(jsonData, &tests)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func TestCPU(t *testing.T) {
	cpu := InitCPU()
	log.Printf("CPU initialized %v", cpu)
}

func LoadAllTests(dirPath string) (map[string][]CPUTest, error) {
	allTests := make(map[string][]CPUTest)
	var loadErrors []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("Error accessing path %s: %v", path, err))
			return nil // Continue with other files
		}

		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("Error reading file %s: %v", path, err))
			return nil // Continue with other files
		}

		var tests []CPUTest
		if err := json.Unmarshal(data, &tests); err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("Error parsing JSON from file %s: %v", path, err))
			return nil // Continue with other files
		}

		fileName := filepath.Base(path)
		fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
		allTests[fileName] = tests

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	if len(loadErrors) > 0 {
		for _, err := range loadErrors {
			fmt.Printf("Warning: %s\n", err)
		}
	}

	return allTests, nil
}

func RunTest(test CPUTest, t *testing.T) {
	cpu := InitCPU()
	cpu.Registers[RegA] = test.Initial.A
	cpu.Registers[RegB] = test.Initial.B
	cpu.Registers[RegC] = test.Initial.C
	cpu.Registers[RegD] = test.Initial.D
	cpu.Registers[RegE] = test.Initial.E
	cpu.Flags.SetValue(test.Initial.F)
	cpu.Registers[RegH] = test.Initial.H
	cpu.Registers[RegL] = test.Initial.L
	cpu.PC = test.Initial.PC - 1
	cpu.SP = test.Initial.SP
	for _, ram := range test.Initial.RAM {
		cpu.Memory[ram[0]] = uint8(ram[1])
	}

	// run

	cpu.ParseNextOpcode()

	// check

	if cpu.Registers[RegA] != test.Final.A {
		log.Printf("Test %v: A register mismatch: expected %v, got %v", test.Name, test.Final.A, cpu.Registers[RegA])
		t.Errorf("Test %v: A register mismatch: expected %v, got %v", test.Name, test.Final.A, cpu.Registers[RegA])
	}
	if cpu.Registers[RegB] != test.Final.B {
		log.Printf("Test %v: B register mismatch: expected %v, got %v", test.Name, test.Final.B, cpu.Registers[RegB])
		t.Errorf("Test %v: B register mismatch: expected %v, got %v", test.Name, test.Final.B, cpu.Registers[RegB])
	}
	if cpu.Registers[RegC] != test.Final.C {
		log.Printf("Test %v: C register mismatch: expected %v, got %v", test.Name, test.Final.C, cpu.Registers[RegC])
		t.Errorf("Test %v: C register mismatch: expected %v, got %v", test.Name, test.Final.C, cpu.Registers[RegC])
	}
	if cpu.Registers[RegD] != test.Final.D {
		log.Printf("Test %v: D register mismatch: expected %v, got %v", test.Name, test.Final.D, cpu.Registers[RegD])
		t.Errorf("Test %v: D register mismatch: expected %v, got %v", test.Name, test.Final.D, cpu.Registers[RegD])
	}
	if cpu.Registers[RegE] != test.Final.E {
		log.Printf("Test %v: E register mismatch: expected %v, got %v", test.Name, test.Final.E, cpu.Registers[RegE])
		t.Errorf("Test %v: E register mismatch: expected %v, got %v", test.Name, test.Final.E, cpu.Registers[RegE])
	}
	if cpu.Registers[RegF] != test.Final.F {
		log.Printf("Test %v: F register mismatch: expected %v, got %v", test.Name, test.Final.F, cpu.Registers[RegF])
		t.Errorf("Test %v: F register mismatch: expected %v, got %v", test.Name, test.Final.F, cpu.Registers[RegF])
	}
	if cpu.Registers[RegH] != test.Final.H {
		log.Printf("Test %v: H register mismatch: expected %v, got %v", test.Name, test.Final.H, cpu.Registers[RegH])
		t.Errorf("Test %v: H register mismatch: expected %v, got %v", test.Name, test.Final.H, cpu.Registers[RegH])
	}
	if cpu.Registers[RegL] != test.Final.L {
		log.Printf("Test %v: L register mismatch: expected %v, got %v", test.Name, test.Final.L, cpu.Registers[RegL])
		t.Errorf("Test %v: L register mismatch: expected %v, got %v", test.Name, test.Final.L, cpu.Registers[RegL])
	}
	if cpu.PC != test.Final.PC-1 {
		log.Printf("Test %v: PC mismatch: expected %v, got %v", test.Name, test.Final.PC-1, cpu.PC)
		t.Errorf("Test %v: PC mismatch: expected %v, got %v", test.Name, test.Final.PC-1, cpu.PC)
	}
	if cpu.SP != test.Final.SP {
		log.Printf("Test %v: SP mismatch: expected %v, got %v", test.Name, test.Final.SP, cpu.SP)
		t.Errorf("Test %v: SP mismatch: expected %v, got %v", test.Name, test.Final.SP, cpu.SP)
	}
	for i, ram := range test.Final.RAM {
		if cpu.Memory[ram[0]] != uint8(ram[1]) {
			log.Printf("Test %v: RAM mismatch at index %v: expected %v, got %v", test.Name, i, ram[1], cpu.Memory[ram[0]])
			t.Errorf("Test %v: RAM mismatch at index %v: expected %v, got %v", test.Name, i, ram[1], cpu.Memory[ram[0]])
		}
	}
}

func TestLoadAllTests(t *testing.T) {
	files, _ := LoadAllTests("tests/cpu_tests")
	for _, tests := range files {
		for _, test := range tests {
			RunTest(test, t)
		}
	}
}
