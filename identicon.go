package identicon

import (
	"crypto/sha512"
	"encoding/binary"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

// Code derives a code for use with Render.
func Code(str string) uint64 {
	buf := sha512.Sum512([]byte(str))
	return binary.BigEndian.Uint64(buf[56:])
}

type Settings struct {
	// TwoColor specifies if the identicon should be
	// generated using one or two colors.
	TwoColor bool

	// Alpha specifies the transparency of the generated identicon.
	Alpha uint8
}

// DefaultSettings returns a Settings object with the recommended settings.
func DefaultSettings() *Settings {
	return &Settings{
		TwoColor: true,
		Alpha:    255,
	}
}

// Render generates an identicon.
// code is a code derived by the Code function.
// totalSize specifies the total size in pixels. It is recommended that
// this is divisible by 3.
func Render(code uint64, totalSize int, settings *Settings) image.Image {
	penWidth := 0
	middleType := int(code & 0x03)
	middleInvert := code>>2&0x01 == 1
	cornerType := int(code >> 3 & 0x0f)
	cornerInvert := code>>7&0x01 == 1
	cornerTurn := int(code >> 8 & 0x03)
	sideType := int(code >> 10 & 0x0f)
	sideInvert := code>>14&0x01 == 1
	sideTurn := int(code >> 15 & 0x03)
	blue := code >> 17 & 0x1f
	green := code >> 22 & 0x1f
	red := code >> 27 & 0x1f
	secondRed := code >> 32 & 0x1f
	secondGreen := code >> 37 & 0x1f
	secondBlue := code >> 42 & 0x1f
	swapCross := code>>47&0x01 == 1
	middleType = middlePatchSet[middleType]
	foreColor := color.RGBA{R: uint8(red) << 3, G: uint8(green) << 3, B: uint8(blue) << 3, A: settings.Alpha}
	var secondColor color.RGBA
	if settings.TwoColor {
		secondColor = color.RGBA{R: uint8(secondRed) << 3, G: uint8(secondGreen) << 3, B: uint8(secondBlue) << 3, A: settings.Alpha}
	} else {
		secondColor = foreColor
	}
	var middleColor color.Color
	if swapCross {
		middleColor = foreColor
	} else {
		middleColor = secondColor
	}
	image := gg.NewContext(totalSize, totalSize)
	patchSize := float64(totalSize) / 3
	drawPatch(gg.Point{X: 1, Y: 1}, 0, middleInvert, middleType, image, patchSize, middleColor, penWidth)
	for i, p := range []gg.Point{{X: 1, Y: 0}, {X: 2, Y: 1}, {X: 1, Y: 2}, {X: 0, Y: 1}} {
		drawPatch(p, sideTurn+1+i, sideInvert, sideType, image, patchSize, foreColor, penWidth)
	}
	for i, p := range []gg.Point{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}} {
		drawPatch(p, cornerTurn+1+i, cornerInvert, cornerType, image, patchSize, secondColor, penWidth)
	}
	return image.Image()
}

func drawPatch(pos gg.Point, turn int, invert bool, type_ int, image *gg.Context, patchSize float64, foreColor color.Color, penWidth int) {
	path := pathSet[type_]
	turn %= 4
	image.Push()
	image.Translate(pos.X*patchSize+float64(penWidth)/2, pos.Y*patchSize+float64(penWidth)/2)
	image.RotateAbout(float64(turn)*math.Pi/2, patchSize/2, patchSize/2)
	for _, p := range path {
		image.LineTo(p.X/4*patchSize, p.Y/4*patchSize)
	}
	image.ClosePath()
	if invert {
		image.MoveTo(0, 0)
		image.LineTo(0, patchSize)
		image.LineTo(patchSize, patchSize)
		image.LineTo(patchSize, 0)
		image.ClosePath()
	}
	image.SetColor(foreColor)
	image.Fill()
	image.Pop()
}

var pathSet = [][]gg.Point{
	// [0] full square:
	{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
	// [1] right-angled triangle pointing top-left:
	{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 0, Y: 4}},
	// [2] upwardy triangle:
	{{X: 2, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
	// [3] left half of square, standing rectangle:
	{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 4}, {X: 0, Y: 4}},
	// [4] square standing on diagonale:
	{{X: 2, Y: 0}, {X: 4, Y: 2}, {X: 2, Y: 4}, {X: 0, Y: 2}},
	// [5] kite pointing topleft:
	{{X: 0, Y: 0}, {X: 4, Y: 2}, {X: 4, Y: 4}, {X: 2, Y: 4}},
	// [6] Sierpinski triangle, fractal triangles:
	{{X: 2, Y: 0}, {X: 4, Y: 4}, {X: 2, Y: 4}, {X: 3, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 4}, {X: 0, Y: 4}},
	// [7] sharp angled lefttop pointing triangle:
	{{X: 0, Y: 0}, {X: 4, Y: 2}, {X: 2, Y: 4}},
	// [8] small centered square:
	{{X: 1, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 3}, {X: 1, Y: 3}},
	// [9] two small triangles:
	{{X: 2, Y: 0}, {X: 4, Y: 0}, {X: 0, Y: 4}, {X: 0, Y: 2}, {X: 2, Y: 2}},
	// [10] small topleft square:
	{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}},
	// [11] downpointing right-angled triangle on bottom:
	{{X: 0, Y: 2}, {X: 4, Y: 2}, {X: 2, Y: 4}},
	// [12] uppointing right-angled triangle on bottom:
	{{X: 2, Y: 2}, {X: 4, Y: 4}, {X: 0, Y: 4}},
	// [13] small rightbottom pointing right-angled triangle on topleft:
	{{X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}},
	// [14] small lefttop pointing right-angled triangle on topleft:
	{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 2}},
	// [15] empty:
	{},
}

// get the [0] full square, [4] square standing on diagonale, [8] small centered square, or [15] empty tile:
var middlePatchSet = []int{0, 4, 8, 15}
