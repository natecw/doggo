package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const url = "https://dog.ceo/api/breeds/image/random"

type response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	a := app.New()
	w := a.NewWindow("Doggos")

	err, file := nextDoggo()
	check(err)

	img := canvas.NewImageFromFile(file)
	w.SetContent(img)
	w.Resize(fyne.NewSize(640, 480))

	w.ShowAndRun()
}

func nextDoggo() (error, string) {
	res, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		return err, ""
	}
	defer res.Body.Close()

	doggoReference := new(response)
	err = json.NewDecoder(res.Body).Decode(&doggoReference)
	if err != nil {
		return err, ""
	}

	file, err := os.CreateTemp(".", "")
	if err != nil {
		return err, ""
	}

	pictureUrl := strings.ReplaceAll(doggoReference.Message, "\\", "")
	res, err = http.Get(pictureUrl)
	if err != nil {
		return err, ""
	}
	defer res.Body.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err, ""
	}
	defer file.Close()

	return nil, file.Name()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
