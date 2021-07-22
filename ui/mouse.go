package ui2d

import "github.com/veandco/go-sdl2/sdl"

type mouseState struct {
	leftButton, rightButton         bool
	prevLeftButton, prevRightButton bool
	x, y                            int32
	prevX, prevY                    int32
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

func (mouseState *mouseState) update() {
	mouseState.prevX, mouseState.prevY = mouseState.x, mouseState.y
	mouseState.prevLeftButton, mouseState.prevRightButton = mouseState.leftButton, mouseState.rightButton

	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	mouseState.leftButton = (mouseButtonState&sdl.ButtonLMask() != 0)
	mouseState.rightButton = (mouseButtonState&sdl.ButtonRMask() != 0)
	mouseState.x, mouseState.y = mouseX, mouseY
}

func (mouseState *mouseState) leftClicked() bool {
	return mouseState.leftButton && !mouseState.prevLeftButton
}

func (mouseState *mouseState) leftUnclicked() bool {
	return !mouseState.leftButton && mouseState.prevLeftButton
}

func (mouseState *mouseState) rightClicked() bool {
	return mouseState.rightButton && !mouseState.prevRightButton
}

func (mouseState *mouseState) rightUnclicked() bool {
	return !mouseState.rightButton && mouseState.prevRightButton
}

func (mouseState *mouseState) onArea(rect *sdl.Rect) bool {
	return rect.HasIntersection(&sdl.Rect{mouseState.x, mouseState.y, 1, 1})
}
