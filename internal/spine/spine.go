package spine

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/png"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

type Skeleton struct {
	Hash   string  `json:"hash"`
	Spine  string  `json:"spine"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Images string  `json:"images"`
	Audio  string  `json:"audio"`
}

type Bone struct {
	Name     string  `json:"name"`
	Parent   string  `json:"parent,omitempty"`
	Length   float64 `json:"length,omitempty"`
	Rotation float64 `json:"rotation,omitempty"`
	X        float64 `json:"x,omitempty"`
	Y        float64 `json:"y,omitempty"`
	Inherit  string  `json:"inherit,omitempty"`
}

type Slot struct {
	Name       string `json:"name"`
	Bone       string `json:"bone"`
	Attachment string `json:"attachment,omitempty"`
}

type Attachment struct {
	X        float64 `json:"x,omitempty"`
	Y        float64 `json:"y,omitempty"`
	Rotation float64 `json:"rotation,omitempty"`
	Width    float64 `json:"width,omitempty"`
	Height   float64 `json:"height,omitempty"`
}

type Skin struct {
	Name        string                           `json:"name"`
	Attachments map[string]map[string]Attachment `json:"attachments"`
}

type SpineData struct {
	Skeleton Skeleton `json:"skeleton"`
	Bones    []Bone   `json:"bones"`
	Slots    []Slot   `json:"slots"`
	Skins    []Skin   `json:"skins"`

	Img *ebiten.Image
}

func LoadSpineData(dataPath string) (*SpineData, error) {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var spineData SpineData
	if err := json.Unmarshal(data, &spineData); err != nil {
		return nil, err
	}
	imgPath := path.Join(path.Dir(dataPath), "hero-ess.png")
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	spineData.Img = ebiten.NewImageFromImage(img)

	return &spineData, nil
}
