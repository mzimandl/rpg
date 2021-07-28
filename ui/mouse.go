package ui

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type mouseState struct {
	leftButton, rightButton         bool
	prevLeftButton, prevRightButton bool
	x, y                            int32
	prevX, prevY                    int32
	lastClick                       time.Time
}

func NewMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()

	return &mouseState{
		x:           mouseX,
		y:           mouseY,
		leftButton:  !(leftButton == 0),
		rightButton: !(rightButton == 0),
	}
}

func (ms *mouseState) update() {
	ms.prevX, ms.prevY = ms.x, ms.y
	ms.prevLeftButton, ms.prevRightButton = ms.leftButton, ms.rightButton

	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	ms.leftButton = (mouseButtonState&sdl.ButtonLMask() != 0)
	ms.rightButton = (mouseButtonState&sdl.ButtonRMask() != 0)
	ms.x, ms.y = mouseX, mouseY

	if ms.leftUnclicked() {
		ms.lastClick = time.Now()
	}
}

func (ms *mouseState) leftDoubleClicked() bool {
	return ms.leftButton && !ms.prevLeftButton && time.Now().Sub(ms.lastClick).Milliseconds() < 250
}

func (ms *mouseState) leftClicked() bool {
	return ms.leftButton && !ms.prevLeftButton
}

func (ms *mouseState) leftUnclicked() bool {
	return !ms.leftButton && ms.prevLeftButton
}

func (ms *mouseState) rightClicked() bool {
	return ms.rightButton && !ms.prevRightButton
}

func (ms *mouseState) rightUnclicked() bool {
	return !ms.rightButton && ms.prevRightButton
}

func (ms *mouseState) onRect(rect *sdl.Rect) bool {
	return rect.HasIntersection(&sdl.Rect{ms.x, ms.y, 1, 1})
}
