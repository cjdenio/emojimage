package main

import (
	"fmt"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/slack-go/slack"
)

const CHANNEL = "C021WGW2LSF"

type EventPayload struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Event     struct {
		Type    string `json:"type"`
		Channel string `json:"channel_id"`
		File    struct {
			ID string `json:"id"`
		} `json:"file"`
	} `json:"event"`
}

func main() {
	r := gin.Default()

	r.POST("/slack/events", func(c *gin.Context) {
		var payload EventPayload
		c.BindJSON(&payload)

		if payload.Type == "url_verification" {
			c.String(200, payload.Challenge)
			return
		}

		c.String(200, "")

		if payload.Event.Channel != CHANNEL {
			return
		}

		client := slack.New(os.Getenv("SLACK_TOKEN"))

		f, _, _, err := client.GetFileInfo(payload.Event.File.ID, 100, 1)
		if err != nil {
			log.Println(err)
			return
		}

		share := f.Shares.Public[CHANNEL][0]

		if share.ThreadTs != "" {
			return
		}

		reportErr := func(message string) {
			client.PostMessage(CHANNEL, slack.MsgOptionText(fmt.Sprintf("Something went wrong: `%s`", message), false), slack.MsgOptionTS(share.Ts))
		}

		if f.Filetype != "png" {
			reportErr("File must be a PNG")
			return
		}

		u, _ := url.Parse(f.URLPrivate)
		resp, err := http.DefaultClient.Do(&http.Request{
			Method: http.MethodGet,
			URL:    u,
			Header: http.Header{
				"Authorization": {fmt.Sprintf("Bearer %s", os.Getenv("SLACK_TOKEN"))},
			},
		})
		if err != nil {
			reportErr(err.Error())
			return
		}

		image, err := emojImage(resp.Body)
		if err != nil {
			reportErr(err.Error())
			return
		}

		_, err = client.UploadFile(slack.FileUploadParameters{
			Channels:        []string{CHANNEL},
			Content:         image,
			ThreadTimestamp: share.Ts,
			InitialComment:  "Copy the following to your clipboard, then paste it somewhere in Slack!",
		})
		if err != nil {
			reportErr(err.Error())
			return
		}

		client.PostMessage(CHANNEL, slack.MsgOptionText(image, false), slack.MsgOptionTS(share.Ts))
	})

	r.Run("0.0.0.0:3000")
}

func emojImage(img io.Reader) (string, error) {
	src, err := png.Decode(img)
	if err != nil {
		return "", nil
	}

	resized := resize.Resize(20, 20, src, resize.NearestNeighbor)

	message := ""

	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			c := resized.At(x, y)
			r, g, b, a := c.RGBA()
			if a < 16384 {
				message += ":blank:"
				continue
			}

			message += fmt.Sprintf(":%02d%02d%02d:", int(math.Round(float64(r)/6553.5)), int(math.Round(float64(g)/6553.5)), int(math.Round(float64(b)/6553.5)))
		}

		message += "\n"
	}

	return message, nil
}
