package spine

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/png"
	"math"
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

	Atlas   *Atlas
	Image   *ebiten.Image
	BoneMap map[string]*Bone

	BoneState        map[string]BoneTransform
	Animations       map[string]Animation
	CurrentAnimation string
	time             float64
}

func (s *Spine) findBoneTransform(name string) (*Bone, float64, float64, float64, error) {
	var fx, fy, fr float64

	cur := name
	for len(cur) > 0 {
		b, ok := s.BoneMap[cur]
		if !ok {
			return nil, 0, 0, 0, errors.New("no bone found")
		}
		state := s.BoneState[cur]
		fx += b.X + state.X
		fy += b.Y + state.Y
		fr += b.Rotation + state.Rotation
		cur = b.Parent
	}

	b := s.BoneMap[name]
	return b, fx, -fy, fr, nil
}

func (s *Spine) getAttachment(slotName, attachmentName string) (Attachment, bool) {
	if len(s.Skins) == 0 {
		return Attachment{}, false
	}
	skin := s.Skins[0]
	if slotMap, ok := skin.Attachments[slotName]; ok {
		if att, ok := slotMap[attachmentName]; ok {
			return att, true
		}
	}
	return Attachment{}, false
}

func (s *Spine) Draw(screen *ebiten.Image, x, y float64) {
	for _, slot := range s.Slots {
		bone, bx, by, br, err := s.findBoneTransform(slot.Bone)
		if err != nil || bone == nil {
			continue
		}

		att, ok := s.getAttachment(slot.Name, slot.Attachment)
		if !ok {
			continue
		}

		sprite, err := s.Atlas.FindSprite(slot.Attachment)
		if err != nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(sprite.Offsets.Width)/2, -float64(sprite.Offsets.Height)/2)
		if sprite.Rotated {
			op.GeoM.Rotate(math.Pi / 2)
		}
		op.GeoM.Rotate((br + att.Rotation) * math.Pi / 180)
		op.GeoM.Translate(x+bx+att.X, y+by-att.Y)

		sub := s.Image.SubImage(toImageRect(sprite)).(*ebiten.Image)
		screen.DrawImage(sub, op)
	}
}

func toImageRect(sprite Sprite) image.Rectangle {
	if sprite.Rotated {
		return image.Rect(sprite.Bounds.X, sprite.Bounds.Y,
			sprite.Bounds.X+sprite.Bounds.Height,
			sprite.Bounds.Y+sprite.Bounds.Width)
	}

	return image.Rect(sprite.Bounds.X, sprite.Bounds.Y,
		sprite.Bounds.X+sprite.Bounds.Width,
		sprite.Bounds.Y+sprite.Bounds.Height)
}

func (s *Spine) FindBone(name string) (*Bone, error) {
	for _, b := range s.Bones {
		if b.Name == name {
			return &b, nil
		}
	}
	return nil, errors.New("no bone found")
}

func (s *Spine) FindBonePos(name string) (*Bone, float64, float64, error) {
	var fx = 0.0
	var fy = 0.0

	cur := name
	for len(cur) > 0 {
		b, ok := s.BoneMap[cur]
		if !ok {
			return nil, 0.0, 0.0, errors.New("no bone found")
		}
		st := s.BoneState[cur]
		fx += b.X + st.X
		fy += b.Y + st.Y
		cur = b.Parent
	}
	b, _ := s.BoneMap[name]
	return b, fx, -fy, nil
}

type spineData struct {
	Skeleton   Skeleton             `json:"skeleton"`
	Bones      []Bone               `json:"bones"`
	Slots      []Slot               `json:"slots"`
	Skins      []Skin               `json:"skins"`
	Animations map[string]Animation `json:"animations"`
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

	spine.BoneMap = make(map[string]*Bone)
	spine.BoneState = make(map[string]BoneTransform)
	for i := range spine.Bones {
		b := &spine.Bones[i]
		spine.BoneMap[b.Name] = b
		spine.BoneState[b.Name] = BoneTransform{}
	}

	spine.Animations = spineData.Animations

	return spine, nil
}

// SetAnimation sets the current animation by name.
func (s *Spine) SetAnimation(name string) {
	if _, ok := s.Animations[name]; ok {
		s.CurrentAnimation = name
		s.time = 0
	}
}

// Update advances the current animation by dt seconds.
func (s *Spine) Update(dt float64) {
	if s.CurrentAnimation == "" {
		return
	}
	anim, ok := s.Animations[s.CurrentAnimation]
	if !ok {
		return
	}
	s.time += dt
	for name, ba := range anim.Bones {
		state := s.BoneState[name]
		if len(ba.Translate) > 0 {
			x, y := sampleVec2(ba.Translate, s.time)
			state.X = x
			state.Y = y
		}
		if len(ba.Rotate) > 0 {
			r := sampleValue(ba.Rotate, s.time)
			state.Rotation = r
		}
		s.BoneState[name] = state
	}
}
