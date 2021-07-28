package ui

import (
	"math/rand"
	"strconv"

	"github.com/veandco/go-sdl2/mix"
)

type sounds struct {
	doorOpen  []*mix.Chunk
	doorClose []*mix.Chunk
	footstep  []*mix.Chunk
}

func (ui *ui) loadAudio() {
	var err error
	ui.music, err = mix.LoadMUS("ui/assets/dungeon002.ogg")
	if err != nil {
		panic(err)
	}

	ui.sounds = &sounds{}
	footstepBase := "ui/assets/sounds/footstep0"
	for i := 0; i <= 9; i++ {
		filename := footstepBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			panic(err)
		}
		ui.sounds.footstep = append(ui.sounds.footstep, chunk)
	}
	doorOpenBase := "ui/assets/sounds/doorOpen_"
	for i := 1; i <= 2; i++ {
		filename := doorOpenBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			panic(err)
		}
		ui.sounds.doorOpen = append(ui.sounds.doorOpen, chunk)
	}
	doorCloseBase := "ui/assets/sounds/doorClose_"
	for i := 1; i <= 4; i++ {
		filename := doorCloseBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			panic(err)
		}
		ui.sounds.doorClose = append(ui.sounds.doorClose, chunk)
	}

}

func (s *sounds) Free() {
	for _, chunk := range s.doorOpen {
		chunk.Free()
	}
	for _, chunk := range s.footstep {
		chunk.Free()
	}
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}
