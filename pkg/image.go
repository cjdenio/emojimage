package pkg

import (
	"fmt"
	"image/png"
	"io"
	"math"

	"github.com/nfnt/resize"
)

func GenerateImage(img io.Reader, size uint) (string, error) {
	src, err := png.Decode(img)
	if err != nil {
		return "", err
	}

	resized := resize.Resize(size, size, src, resize.NearestNeighbor)

	message := ""

	for y := 0; y < int(size); y++ {
		for x := 0; x < int(size); x++ {
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
