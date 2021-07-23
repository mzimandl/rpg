package ui2d

import (
	"math/rand"

	"github.com/veandco/go-sdl2/mix"
)

type sounds struct {
	doorOpen []*mix.Chunk
	footstep []*mix.Chunk
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}
