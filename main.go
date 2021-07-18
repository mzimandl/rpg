package main

import (
	"rpg/game"
	ui2d "rpg/ui"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	g := game.NewGame(1, "game/maps/level1.map") // problematic multiple window view, something with sdl and threads
	for i := range g.LevelChans {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()

			runtime.LockOSThread() // SDL has to stay on one thread
			ui := ui2d.NewUI(g.InputChan, g.LevelChans[i])
			ui.Run()
			ui.Destroy()
		}(i)
	}
	g.Run()

	wg.Wait()
	ui2d.Destroy()
}
