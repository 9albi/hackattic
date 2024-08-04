package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/9albi/hackattic/pkg/client"
)

type Problem struct {
	BytesString string `json:"bytes"`
}

type Solution struct {
	Int             int32   `json:"int"`
	Uint            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float32 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

func solve() error {
	token := os.Getenv("HACKATTIC_ACCESS_TOKEN")
	hackatticClient, err := client.NewHackatticClient("help_me_unpack", token)
	if err != nil {
		return err
	}

	var problem Problem
	err = hackatticClient.GetChallenge(&problem)
	if err != nil {
		return err
	}

	// decode base64 encoded input
	decodedBytes, err := base64.StdEncoding.DecodeString(problem.BytesString)
	if err != nil {
		return err
	}

	// parse Int
	intBytes := decodedBytes[0:4]

	var intValue int32
	intValue += int32(intBytes[0])
	intValue += int32(intBytes[1]) << 8
	intValue += int32(intBytes[2]) << 16
	intValue += int32(intBytes[3]) << 24

	fmt.Println("int", intValue)

	// parse Uint
	uintBytes := decodedBytes[4:8]

	var uintValue uint32
	uintValue += uint32(uintBytes[0])
	uintValue += uint32(uintBytes[1]) << 8
	uintValue += uint32(uintBytes[2]) << 16
	uintValue += uint32(uintBytes[3]) << 24

	fmt.Println("uint", uintValue)

	// parse Short
	shortBytes := decodedBytes[8:12]

	var shortValue int16

	shortValue += int16(shortBytes[0])
	shortValue += int16(shortBytes[1]) << 8
	shortValue += int16(shortBytes[2]) << 16
	shortValue += int16(shortBytes[3]) << 24

	fmt.Println("short", shortValue)

	// parse Float
	floatBytes := decodedBytes[12:16]

	// TODO:
	// sign := floatBytes[0] >> 7
	//
	// exponent := floatBytes[0] << 1
	//
	// fmt.Println("exponennt", exponent)
	bits := binary.LittleEndian.Uint32(floatBytes)
	floatValue := math.Float32frombits(bits)

	fmt.Println("float", floatValue)

	// parse Double
	doubleBytes := decodedBytes[16:24]

	doubleBits := binary.LittleEndian.Uint64(doubleBytes)
	doubleValue := math.Float64frombits(doubleBits)

	fmt.Println("double", doubleValue)

	// parse Double (BigEndian)
	bigEndianDoubleBytes := decodedBytes[24:32]

	bigEndianDoubleBits := binary.BigEndian.Uint64(bigEndianDoubleBytes)
	bigEndianDoubleValue := math.Float64frombits(bigEndianDoubleBits)

	fmt.Println("bigEndian double", bigEndianDoubleValue)

	solution := Solution{
		Int:             intValue,
		Uint:            uintValue,
		Short:           shortValue,
		Float:           floatValue,
		Double:          doubleValue,
		BigEndianDouble: bigEndianDoubleValue,
	}

	result, err := hackatticClient.PostSolution(solution)
	if err != nil {
		return err
	}

	log.Print(string(result))

	return nil
}

func main() {
	err := solve()
	if err != nil {
		panic(err)
	}
}
