package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type keyboardState struct {
	current  []uint8
	previous []uint8
}

func NewKeyboardState() *keyboardState {
	current := sdl.GetKeyboardState()
	previous := make([]uint8, len(current))
	for i, v := range current {
		previous[i] = v
	}
	return &keyboardState{current, previous}
}

func (ks *keyboardState) pressed(scancode uint8) bool {
	return ks.current[scancode] != 0 && ks.previous[scancode] == 0
}

func (ks *keyboardState) hold(scancode uint8) bool {
	return ks.current[scancode] != 0
}

func (ks *keyboardState) update() {
	for i, v := range ks.current {
		ks.previous[i] = v
	}
	ks.current = sdl.GetKeyboardState()
}
