package main

import (
   "fmt"
   "github.com/hajimehoshi/ebiten"
   "github.com/hajimehoshi/ebiten/ebitenutil"
   "time"
)

var g Game
var graphics Graphics
var init_knob bool

func init() {
   g.init()
}

func update(screen *ebiten.Image) error {
   if !graphics.is_init() {
      fmt.Printf("Init Graphics...\n")
      graphics.init(screen)
   }

   up_s := time.Now().UnixNano()
   g.update()
   up_t := time.Now().UnixNano() - up_s
   dr_s := time.Now().UnixNano()
   graphics.draw(&g, screen)
   dr_t := time.Now().UnixNano() - dr_s
   msg := fmt.Sprintf(`FPS: %0.2f, up = %v, dr = %v`, ebiten.CurrentFPS(), up_t, dr_t)
   ebitenutil.DebugPrint(screen, msg)
   return nil
}

func main() {
   init_knob = false
   if err := ebiten.Run(update, 480, 480, 1, "Tetris"); err != nil {
      panic(err)
   }
}

