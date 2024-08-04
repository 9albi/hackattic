package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"log/slog"
	"math/bits"
	"os"

	"github.com/9albi/hackattic/pkg/client"
)

type Problem struct {
	Block struct {
		Data  [][]interface{} `json:"data"`
		Nonce int             `json:"nonce,omitempty"`
	} `json:"block,omitempty"`
	Difficulty int `json:"difficulty,omitempty"`
}

type Solution struct {
	Nonce int `json:"nonce"`
}

func solve() error {
	token := os.Getenv("HACKATTIC_ACCESS_TOKEN")
	hackatticClient, err := client.NewHackatticClient("mini_miner", token)
	if err != nil {
		return err
	}

	var problem Problem
	err = hackatticClient.GetChallenge(&problem)
	if err != nil {
		return err
	}

	for {
		jsonData, err := json.Marshal(problem.Block)
		if err != nil {
			panic(err)
		}

		digest := sha256.Sum256(jsonData)

		difficulty := ComputeDifficulty(digest[:])

		if difficulty >= problem.Difficulty {
			slog.Info("hit", "nonce", problem.Block.Nonce, "difficulty", difficulty)
			break
		}

		problem.Block.Nonce++
	}

	solution := Solution{
		Nonce: problem.Block.Nonce,
	}

	result, err := hackatticClient.PostSolution(solution)
	if err != nil {
		return err
	}

	log.Print(string(result))
	return nil
}

func ComputeDifficulty(shaDigest []byte) int {
	difficulty := 0
	for _, part := range shaDigest {
		byteLeadingZeros := bits.LeadingZeros8(part)
		difficulty += byteLeadingZeros
		if byteLeadingZeros != 8 {
			break
		}
	}
	return difficulty
}

func main() {
	err := solve()
	if err != nil {
		panic(err)
	}
}
