// Copyright 2020 Marius van der Wijden
// This file is part of the fuzzy-vm library.
//
// The fuzzy-vm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The fuzzy-vm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the fuzzy-vm library. If not, see <http://www.gnu.org/licenses/>.

// Package fuzzer is the entry point for go-fuzz.
package fuzzer

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/holiman/goevmlab/fuzzing"

	"github.com/MariusVanDerWijden/FuzzyVM/generator"
)

var outputDir = "../out"

// Fuzz is the entry point for go-fuzz
func Fuzz(data []byte) int {
	testMaker, _ := generator.GenerateProgram(data)
	name := randTestName(data)
	// Execute the test and write out the resulting trace
	traceFile := setupTrace(name)
	defer traceFile.Close()
	testMaker.Fill(traceFile)
	// Save the test
	test := testMaker.ToGeneralStateTest(name)
	storeTest(test, name)
	return 0
}

func setupTrace(name string) *os.File {
	traceFile, err := os.OpenFile(fmt.Sprintf("./%v-trace.jsonl", name), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic("Could not write out trace file")
	}
	return traceFile
}

// storeTest saves a testcase to disk
func storeTest(test *fuzzing.GeneralStateTest, testName string) {
	path := fmt.Sprintf("%v/%v.json", outputDir, testName)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		panic("Could not open test file")
	}
	defer f.Close()
	// Write to file
	encoder := json.NewEncoder(f)
	if err = encoder.Encode(test); err != nil {
		panic("Could not encode state test")
	}
}

func randTestName(data []byte) string {
	var seedData [8]byte
	copy(seedData[:], data)
	seed := int64(binary.BigEndian.Uint64(seedData[:]))
	rand := rand.New(rand.NewSource(time.Now().UnixNano() ^ seed))
	return fmt.Sprintf("FuzzyVM-%v-%v", rand.Int31(), rand.Int31())
}