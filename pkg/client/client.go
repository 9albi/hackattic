package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type (
	Option func(*hackatticClient)
)

func WithEnv(c *hackatticClient) {
	c.token = os.Getenv("HACKATTIC_ACCESS_TOKEN")
}

const HackatticProblemURL = "https://hackattic.com/challenges/%s/problem?access_token=%s"

type hackatticClient struct {
	challenge  string
	token      string
	httpClient *http.Client
}

func NewHackatticClient(challenge, token string) (*hackatticClient, error) {
	if token == "" {
		return nil, fmt.Errorf("missing token")
	}

	if challenge == "" {
		return nil, fmt.Errorf("missing challenge name")
	}

	return &hackatticClient{
		challenge,
		token,
		&http.Client{},
	}, nil
}

func (c *hackatticClient) PostSolution(solution interface{}) ([]byte, error) {
	formattedURL := fmt.Sprintf(HackatticProblemURL, c.challenge, c.token)

	jsonBytes, err := json.Marshal(solution)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", formattedURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *hackatticClient) GetChallenge(out interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf(HackatticProblemURL, c.challenge, c.token), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBody, out)
}
