package ui2d

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

func loadSounds() *sounds {
	newSounds := &sounds{}

	footstepBase := "ui/assets/sounds/footstep0"
	for i := 0; i <= 9; i++ {
		filename := footstepBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			panic(err)
		}
		newSounds.footstep = append(newSounds.footstep, chunk)
	}
	doorOpenBase := "ui/assets/sounds/doorOpen_"
	for i := 1; i <= 2; i++ {
		filename := doorOpenBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			panic(err)
		}
		newSounds.doorOpen = append(newSounds.doorOpen, chunk)
	}
	doorCloseBase := "ui/assets/sounds/doorClose_"
	for i := 1; i <= 4; i++ {
		filename := doorCloseBase + strconv.Itoa(i) + ".ogg"
		chunk, err := mix.LoadWAV(filename)
		if err != nil {
			panic(err)
		}
		newSounds.doorClose = append(newSounds.doorClose, chunk)
	}

	return newSounds
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
