package spine

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Rect 구조체: 사각 영역 정보를 저장
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Sprite 구조체: 개별 스프라이트 정보를 저장
type Sprite struct {
	Name    string
	Bounds  Rect // 아틀라스 내 위치와 크기
	Offsets Rect // 원본 이미지의 오프셋 및 크기
	Rotated bool // 회전 여부 (rotate: 90 이 있으면 true)
}

// Atlas 구조체: 아틀라스 전체 정보를 저장
type Atlas struct {
	ImageName            string   // 이미지 파일 이름
	Width, Height        int      // 아틀라스 전체 크기
	FilterMin, FilterMag string   // 필터 종류 (Minification, Magnification)
	PremultipliedAlpha   bool     // 프리멀티플라이드 알파 여부
	Sprites              []Sprite // 스프라이트 리스트
}

func (a *Atlas) FindSprite(name string) (Sprite, error) {
	for _, sprite := range a.Sprites {
		if sprite.Name == name {
			return sprite, nil
		}
	}
	return Sprite{}, errors.New("Sprite not found")
}

// 문자열 형태 "x,y,width,height"를 Rect 구조체로 파싱하는 헬퍼 함수
func parseRect(s string) (Rect, error) {
	// 앞뒤 공백 제거 후 콤마로 분리
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 4 {
		return Rect{}, fmt.Errorf("rect parsing fail: %s", s)
	}
	// 각 부분을 정수로 변환
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Rect{}, err
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Rect{}, err
	}
	w, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return Rect{}, err
	}
	h, err := strconv.Atoi(strings.TrimSpace(parts[3]))
	if err != nil {
		return Rect{}, err
	}
	return Rect{X: x, Y: y, Width: w, Height: h}, nil
}

// 아틀라스 파일을 파싱하는 함수
func parseAtlas(filename string) (Atlas, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Atlas{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // 빈 줄은 무시
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return Atlas{}, err
	}

	if len(lines) < 4 {
		return Atlas{}, fmt.Errorf("아틀라스 파일 형식 오류: 최소 4줄의 헤더 정보가 필요합니다")
	}

	// Atlas 구조체 채우기
	atlas := Atlas{}
	atlas.ImageName = lines[0]

	// 두 번째 줄: size: w,h
	if strings.HasPrefix(lines[1], "size:") {
		sizeData := strings.TrimSpace(strings.TrimPrefix(lines[1], "size:"))
		// "w,h" 형식을 분할하여 정수 변환
		sizeParts := strings.Split(sizeData, ",")
		if len(sizeParts) == 2 {
			atlas.Width, _ = strconv.Atoi(strings.TrimSpace(sizeParts[0]))
			atlas.Height, _ = strconv.Atoi(strings.TrimSpace(sizeParts[1]))
		}
	} else {
		return Atlas{}, fmt.Errorf("size 정보 누락 또는 형식 오류")
	}

	// 세 번째 줄: filter: type1,type2
	if strings.HasPrefix(lines[2], "filter:") {
		filterData := strings.TrimSpace(strings.TrimPrefix(lines[2], "filter:"))
		filterParts := strings.Split(filterData, ",")
		if len(filterParts) == 2 {
			atlas.FilterMin = strings.TrimSpace(filterParts[0])
			atlas.FilterMag = strings.TrimSpace(filterParts[1])
		} else {
			// 콤마로 두 부분이 아니면 하나만 적용
			atlas.FilterMin = strings.TrimSpace(filterData)
			atlas.FilterMag = ""
		}
	} else {
		return Atlas{}, fmt.Errorf("filter 정보 누락 또는 형식 오류")
	}

	// 네 번째 줄: pma: true/false
	if strings.HasPrefix(lines[3], "pma:") {
		pmaData := strings.TrimSpace(strings.TrimPrefix(lines[3], "pma:"))
		atlas.PremultipliedAlpha = (pmaData == "true")
	} else {
		return Atlas{}, fmt.Errorf("pma 정보 누락 또는 형식 오류")
	}

	// 5번째 줄부터 스프라이트 목록 파싱
	i := 4
	for i < len(lines) {
		name := lines[i]
		i++
		// bounds 라인 파싱
		if i >= len(lines) || !strings.HasPrefix(lines[i], "bounds:") {
			return Atlas{}, fmt.Errorf("스프라이트 '%s'의 bounds 정보 누락", name)
		}
		boundsData := strings.TrimSpace(strings.TrimPrefix(lines[i], "bounds:"))
		i++
		// offsets 라인 파싱
		if i >= len(lines) || !strings.HasPrefix(lines[i], "offsets:") {
			return Atlas{}, fmt.Errorf("스프라이트 '%s'의 offsets 정보 누락", name)
		}
		offsetsData := strings.TrimSpace(strings.TrimPrefix(lines[i], "offsets:"))
		i++
		// rotate (선택 사항)
		rotated := false
		if i < len(lines) && strings.HasPrefix(lines[i], "rotate:") {
			// "rotate: 90"인 경우 회전으로 처리
			rotated = true
			i++
		}

		// Rect 문자열을 Rect 구조체로 파싱
		boundsRect, err := parseRect(boundsData)
		if err != nil {
			return Atlas{}, fmt.Errorf("'%s' bounds 파싱 오류: %v", name, err)
		}
		offsetsRect, err := parseRect(offsetsData)
		if err != nil {
			return Atlas{}, fmt.Errorf("'%s' offsets 파싱 오류: %v", name, err)
		}

		// Sprite 구조체 생성 후 Atlas에 추가
		sprite := Sprite{
			Name:    name,
			Bounds:  boundsRect,
			Offsets: offsetsRect,
			Rotated: rotated,
		}
		atlas.Sprites = append(atlas.Sprites, sprite)
	}

	return atlas, nil
}
