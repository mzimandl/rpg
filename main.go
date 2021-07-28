package main

import (
	"rpg/game"
	"rpg/ui"
	"runtime"
	"sync"
)

func main() {
	g := game.NewGame() // problematic multiple window view, something with sdl and threads
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		runtime.LockOSThread() // SDL has to stay on one thread
		ui := ui.NewUI(g.InputChan, g.LevelChan)
		ui.Run()
		ui.Destroy()
		wg.Done()
	}()
	g.Run()
	wg.Wait()
	ui.Destroy()
}
