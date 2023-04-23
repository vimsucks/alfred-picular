package main

import (
	"alfred-picular/picular"
	"errors"
	"fmt"
	aw "github.com/deanishe/awgo"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path"
)

const ImageDirPath = "./images"

var wf = aw.New()

func run() {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = args[0]
	}
	if query == "" {
		return
	}

	picularResponse, err := picular.SearchColor(query)
	if err != nil {
		wf.Fatal("search color failed: " + err.Error())
	}

	if err = makeImageDirIfNotExists(); err != nil {
		wf.Fatal("create image dir failed: " + err.Error())
	}

	for _, c := range picularResponse.Colors {
		hex := c.Color
		rgba, err := parseHexColorFast(hex)
		if err != nil {
			wf.Fatal("parse hex color failed: " + hex)
		}
		filePath, err := createColorImageIfNotExists(hex, rgba)
		if err != nil {
			wf.Fatal("create image file failed: " + err.Error())
		}
		wf.NewItem(hex).
			Arg(hex).
			Subtitle(fmt.Sprintf("rgb(%d,%d,%d)", rgba.R, rgba.G, rgba.B)).
			Icon(&aw.Icon{
				Value: filePath,
				Type:  aw.IconTypeImage,
			}).
			Valid(true)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

// createColorImageIfNotExists create a simple color image in jpg format, return file path
func createColorImageIfNotExists(hex string, rgba color.RGBA) (string, error) {
	imgPath := path.Join(ImageDirPath, hex+".jpg")
	if _, err := os.Stat(imgPath); err == nil {
		return imgPath, nil
	}

	rect := image.Rect(0, 0, 1, 1)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{C: rgba}, image.Point{}, draw.Src)
	out, err := os.Create(imgPath)
	if err != nil {
		return "", err
	}
	var opt jpeg.Options
	opt.Quality = 80
	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		return "", err
	}
	return imgPath, nil
}

func makeImageDirIfNotExists() error {
	return os.MkdirAll(ImageDirPath, os.ModePerm)
}

var errInvalidFormat = errors.New("invalid format")

// parseHexColorFast parse hex color
// example: RED #B40404 -> color.RGBA{R: 180, G: 4, B: 4, A: 255}
func parseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
