package main

import (
	"rpg/game"
	ui2d "rpg/ui"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	g := game.NewGame() // problematic multiple window view, something with sdl and threads
	wg.Add(1)
	go func() {
		runtime.LockOSThread() // SDL has to stay on one thread
		ui := ui2d.NewUI(g.InputChan, g.LevelChan)
		ui.Run()
		ui.Destroy()
		wg.Done()
	}()
	g.Run()

	wg.Wait()
	ui2d.Destroy()
}
