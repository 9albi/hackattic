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

const HackatticProblemURL = "https://hackattic.com/challenges/%s/%s?access_token=%s%s"

type hackatticClient struct {
	httpClient *http.Client
	challenge  string
	token      string
}

func NewHackatticClient(challenge, token string) (*hackatticClient, error) {
	if token == "" {
		return nil, fmt.Errorf("missing token")
	}

	if challenge == "" {
		return nil, fmt.Errorf("missing challenge name")
	}

	return &hackatticClient{
		&http.Client{},
		challenge,
		token,
	}, nil
}

func (c *hackatticClient) GetChallenge(out interface{}) error {
	formattedURL := fmt.Sprintf(HackatticProblemURL, c.challenge, "problem", c.token, "")

	req, err := http.NewRequest("GET", formattedURL, nil)
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

func (c *hackatticClient) PostSolution(solution interface{}) ([]byte, error) {
	formattedURL := fmt.Sprintf(HackatticProblemURL,
		c.challenge,
		"solve",
		c.token,
		"&playground=1",
	)

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
