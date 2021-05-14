package main

import (
	"bytes"
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
	"os"
)

func main() {
	for r := 0; r <= 5; r++ {
		for g := 0; g <= 5; g++ {
			for b := 0; b <= 5; b++ {
				if r < 3 || (g < 1 && r == 3) || (b < 4 && g == 1 && r == 3) {
					continue
				}

				fmt.Printf("Trying %d-%d-%d...\n", r, g, b)

				f, err := os.Open(fmt.Sprintf("img/%d%d%d.png", r, g, b))
				if err != nil {
					fmt.Println("skipping")
					continue
				}

				result, err := uploadEmoji(fmt.Sprintf("color_%d_%d_%d", r, g, b), f)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(result)

				// i := genImage(color.RGBA{uint8(r * 51), uint8(g * 51), uint8(b * 51), 255})
				// err := os.WriteFile(fmt.Sprintf("img/%d%d%d.png", r, g, b), i, 0777)
				// if err != nil {
				// 	log.Fatal(err)
				// }
			}
		}
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

	if resp.StatusCode != 200 {
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
