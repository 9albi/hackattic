package main

import (
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"os"

	"github.com/9albi/hackattic/pkg/client"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"

	_ "image/png"
)

type Problem struct {
	ImageUrl string `json:"image_url"`
}

type Solution struct {
	Code string `json:"code"`
}

func solve() error {
	token := os.Getenv("HACKATTIC_ACCESS_TOKEN")
	hackatticClient, err := client.NewHackatticClient("reading_qr", token)
	if err != nil {
		return err
	}

	var problem Problem
	err = hackatticClient.GetChallenge(&problem)
	if err != nil {
		return err
	}

	fmt.Println(problem.ImageUrl)
	resp, err := http.Get(problem.ImageUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("received non 200 response Code")
	}

	image, err := png.Decode(resp.Body)

	// bounds := image.Bounds()
	// for i := 0; i < bounds.Dx(); i++ {
	// 	for j := 0; j < bounds.Dy(); j++ {
	// 		r, g, b, _ := image.At(i, j).RGBA()
	// 		if (r == 0) && (g == 0) && (b == 0) {
	// 			fmt.Println("found pixel", i, j)
	// 			return nil
	// 		}
	// 	}
	// }

	bmp, _ := gozxing.NewBinaryBitmapFromImage(image)

	qrReader := qrcode.NewQRCodeReader()
	qrResult, _ := qrReader.Decode(bmp, nil)

	solution := Solution{
		Code: qrResult.GetText(),
	}

	result, err := hackatticClient.PostSolution(solution)
	if err != nil {
		return err
	}

	fmt.Print(string(result))
	return nil
}

func main() {
	err := solve()
	if err != nil {
		panic(err)
	}
}
