package collager

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"golang.org/x/image/webp"
)

type SaveTo struct {
	Name string
	Type FileType
}

type FileType int

const (
	PNG FileType = 0
	JPG FileType = 1
)

var (
	ErrNoName           = errors.New("name value is empty")
	ErrNoImages         = errors.New("you must add at least one image to use collage")
	ErrTypeNotSupported = errors.New("file type not supported")
)

type Collager struct {
	Images []image.Image
}

func NewCollager() *Collager {
	return &Collager{}
}

func (c *Collager) FromDetect(b []byte) (err error) {
	detected := http.DetectContentType(b)
	switch detected {
	case "image/webp":
		return c.FromWebp(b)
	case "image/jpeg":
		return c.FromJPG(b)
	case "image/png":
		return c.FromPNG(b)
	default:
		return ErrTypeNotSupported
	}
}

func (c *Collager) FromWebp(b []byte) (err error) {
	img, err := webp.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	c.Images = append(c.Images, img)
	return
}

func (c *Collager) FromBytes(b []byte) (err error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	c.Images = append(c.Images, img)
	return
}

func (c *Collager) FromPNG(b []byte) (err error) {
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	c.Images = append(c.Images, img)
	return
}

func (c *Collager) FromJPG(b []byte) (err error) {
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	c.Images = append(c.Images, img)
	return
}

func (c *Collager) SaveTo(rgba *image.RGBA, s *SaveTo) (err error) {
	if s.Name == "" {
		return ErrNoName
	}
	file, err := os.Create(s.Name)
	if err != nil {
		return err
	}
	switch s.Type {
	case PNG:
		err = png.Encode(file, rgba)
	case JPG:
		err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
	default:
		return ErrTypeNotSupported
	}
	if err != nil {
		return err
	}
	return
}

func findOptimalSize(imgs []image.Image) (x, y int) {
	for _, v := range imgs {
		if v.Bounds().Dx() > x {
			x = v.Bounds().Dx()
		}
	}

	for _, v := range imgs {
		if v.Bounds().Dy() > y {
			y = v.Bounds().Dy()
		}
	}

	return
}

func (c *Collager) Collage(x, y int, st ...*SaveTo) (ret *image.RGBA, err error) {
	if len(c.Images) == 0 {
		return ret, ErrNoImages
	}
	imgX, imgY := findOptimalSize(c.Images)
	maxX, maxY := x*imgX, y*imgY
	ret = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxX, maxY}})

	var dx, dy = 0, 0
	for _, img := range c.Images {
		rec := img.Bounds()
		draw.Draw(ret, image.Rect(dx, dy, dx+rec.Dx(), dy+rec.Dy()), img, image.Point{0, 0}, draw.Src)
		dx += rec.Dx()
		if dx >= maxX {
			dx = 0
			dy += rec.Dy()
		}
	}

	if len(st) > 0 {
		return ret, c.SaveTo(ret, st[0])
	}

	return
}
