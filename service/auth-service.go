package service

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/skip2/go-qrcode"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateToken() (token string) {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	var png []byte
	r.ParseForm()
	token := r.Form["token"][0]

	w.Header().Set("Content-Type", "image/png")
	png, err := qrcode.Encode("http://localhost:9000/qrcode/auth?token="+token, qrcode.Medium, 256)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("Length is %d bytes long\n", len(png))
	}

	w.Write(png)
}
