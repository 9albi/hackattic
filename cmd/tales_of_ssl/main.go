package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/9albi/hackattic/pkg/client"
)

type Problem struct {
	PrivateKey   string `json:"private_key"`
	RequiredData struct {
		Domain        string `json:"domain"`
		Serial_number string `json:"serial_number"`
		Country       string `json:"country"`
	} `json:"required_data"`
}

type Solution struct {
	Certificate string `json:"certificate"`
}

func solve() error {
	token := os.Getenv("HACKATTIC_ACCESS_TOKEN")
	hackatticClient, err := client.NewHackatticClient("tales_of_ssl", token)
	if err != nil {
		return err
	}

	var problem Problem
	err = hackatticClient.GetChallenge(&problem)
	if err != nil {
		return err
	}

	// parse serial_number from hex string to int
	serialNumberInt := big.Int{}
	serialNumberInt.SetString(problem.RequiredData.Serial_number, 0)

	// the private_key is base64 encoded
	decodedPrivateKey, _ := base64.StdEncoding.DecodeString(problem.PrivateKey)

	// parse the private_key as rsa private key
	priv, err := x509.ParsePKCS1PrivateKey(decodedPrivateKey)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Get(fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fields=cca2", problem.RequiredData.Country))
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type RESTResult []struct {
		CountryCode string `json:"cca2"`
	}

	var restResult RESTResult
	if err := json.Unmarshal(respBody, &restResult); err != nil {
		return err
	}

	if len(restResult) == 0 {
		return fmt.Errorf("couldn't find informations for country %s", problem.RequiredData.Country)
	}

	countryCode := restResult[0].CountryCode
	fmt.Printf("Country: %s, Country Code: %s\n", problem.RequiredData.Country, countryCode)

	template := x509.Certificate{
		SerialNumber: &serialNumberInt,
		Subject: pkix.Name{
			Organization: []string{problem.RequiredData.Country},
			Country:      []string{countryCode},
			CommonName:   problem.RequiredData.Domain,
		},
		IsCA: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certificateEncoded := base64.StdEncoding.EncodeToString(derBytes)

	solution := Solution{
		Certificate: string(certificateEncoded),
	}

	result, err := hackatticClient.PostSolution(solution)
	if err != nil {
		return err
	}

	fmt.Println(string(result))
	return nil
}

func main() {
	err := solve()
	if err != nil {
		panic(err)
	}
}
