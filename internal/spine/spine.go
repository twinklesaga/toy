package spine

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/png"
	"os"
	"path"
	"strings"

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

type Spine struct {
	Skeleton Skeleton
	Bones    []Bone
	Slots    []Slot
	Skins    []Skin

	Atlas *Atlas
	Image *ebiten.Image
}

type spineData struct {
	Skeleton Skeleton `json:"skeleton"`
	Bones    []Bone   `json:"bones"`
	Slots    []Slot   `json:"slots"`
	Skins    []Skin   `json:"skins"`
}

func LoadSpineData(dataPath string) (*Spine, error) {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var spineData spineData
	if err := json.Unmarshal(data, &spineData); err != nil {
		return nil, err
	}

	atlasPath := path.Join(path.Dir(dataPath), strings.Replace(path.Base(dataPath), ".json", ".atlas", -1))

	atlas, err := parseAtlas(atlasPath)
	if err != nil {
		return nil, err
	}

	spine := &Spine{
		Skeleton: spineData.Skeleton,
		Bones:    spineData.Bones,
		Slots:    spineData.Slots,
		Skins:    spineData.Skins,
		Atlas:    &atlas,
	}

	imgPath := path.Join(path.Dir(dataPath), atlas.ImageName)
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	spine.Image = ebiten.NewImageFromImage(img)

	return spine, nil
}
