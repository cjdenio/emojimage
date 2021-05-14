package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func main() {
	r := 0
	g := 0
	b := 0

	for r <= 10 {
		for g <= 10 {
			for b <= 10 {
				fmt.Printf("Trying %d-%d-%d...\n", r, g, b)
				fmt.Println("Generating image...")
				defer fmt.Println("")

				i := genImage(color.RGBA{uint8(float32(r) * 25.5), uint8(float32(g) * 25.5), uint8(float32(b) * 25.5), 255})
				resp, err := uploadEmoji(fmt.Sprintf("%02d%02d%02d", r, g, b), bytes.NewReader(i))
				fmt.Println(resp)
				if err != nil {
					fmt.Printf("Failed on %d-%d-%d, waiting 10 seconds... \n", r, g, b)
					time.Sleep(10 * time.Second)
					continue
				}

				b++
			}

			b = 0
			g++
		}

		g = 0
		r++
	}
}

func uploadEmoji(name string, data io.Reader) (string, error) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	go func() {
		w.WriteField("name", name)
		w.WriteField("mode", "data")
		w.WriteField("token", "TOKEN")

		fw, _ := w.CreateFormFile("image", "test.png")
		io.Copy(fw, data)

		w.Close()
		pw.Close()
	}()

	u, _ := url.Parse("https://hackclub.slack.com/api/emoji.add")
	resp, err := http.DefaultClient.Do(&http.Request{
		Method: http.MethodPost,
		URL:    u,
		Body:   pr,
		Header: http.Header{
			"Content-Type": {fmt.Sprintf("multipart/form-data; boundary=%s", w.Boundary())},
		},
	})
	if err != nil {
		return "", err
	}

	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Ok bool `json:"ok"`
	}

	json.Unmarshal(body, &parsed)

	if resp.StatusCode != 200 || !parsed.Ok {
		return string(body), errors.New("uh oh")
	}

	return string(body), nil
}

func genImage(c color.Color) []byte {
	m := image.NewRGBA(image.Rect(0, 0, 50, 50))
	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, m)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
