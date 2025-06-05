package spine

// Keyframe represents a single keyframe for translation or rotation.
type Keyframe struct {
	Time  float64 `json:"time,omitempty"`
	X     float64 `json:"x,omitempty"`
	Y     float64 `json:"y,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// BoneAnimation holds keyframes for a bone.
type BoneAnimation struct {
	Translate []Keyframe `json:"translate,omitempty"`
	Rotate    []Keyframe `json:"rotate,omitempty"`
}

// Animation represents a single animation consisting of multiple bone animations.
type Animation struct {
	Bones map[string]BoneAnimation `json:"bones"`
}

// BoneTransform stores runtime transform values for a bone.
type BoneTransform struct {
	X        float64
	Y        float64
	Rotation float64
}

// sampleVec2 interpolates translation keyframes to get position at time t.
func sampleVec2(frames []Keyframe, t float64) (float64, float64) {
	if len(frames) == 0 {
		return 0, 0
	}
	if t <= frames[0].Time {
		return frames[0].X, frames[0].Y
	}
	for i := 0; i < len(frames)-1; i++ {
		f0 := frames[i]
		f1 := frames[i+1]
		if t < f1.Time {
			r := (t - f0.Time) / (f1.Time - f0.Time)
			x := f0.X + (f1.X-f0.X)*r
			y := f0.Y + (f1.Y-f0.Y)*r
			return x, y
		}
	}
	last := frames[len(frames)-1]
	return last.X, last.Y
}

// sampleValue interpolates rotation keyframes to get value at time t.
func sampleValue(frames []Keyframe, t float64) float64 {
	if len(frames) == 0 {
		return 0
	}
	if t <= frames[0].Time {
		return frames[0].Value
	}
	for i := 0; i < len(frames)-1; i++ {
		f0 := frames[i]
		f1 := frames[i+1]
		if t < f1.Time {
			r := (t - f0.Time) / (f1.Time - f0.Time)
			return f0.Value + (f1.Value-f0.Value)*r
		}
	}
	return frames[len(frames)-1].Value
}
