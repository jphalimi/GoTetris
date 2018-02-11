package main

import (
   "fmt"
   "github.com/hajimehoshi/ebiten"
   "github.com/hajimehoshi/ebiten/ebitenutil"
   "time"
)

var g Game
var graphics Graphics

func init() {
   g.init()
}

func show_fps(screen *ebiten.Image, up_t, dr_t int64) {
   msg := fmt.Sprintf(`FPS: %0.2f, up = %v, dr = %v`, ebiten.CurrentFPS(), up_t, dr_t)
   ebitenutil.DebugPrint(screen, msg)
}

func update(screen *ebiten.Image) error {
   if !graphics.is_init() {
      fmt.Printf("Init Graphics...\n")
      graphics.init(screen)
   }

   if g.is_started() {
      // Display game.
      up_s := time.Now().UnixNano()
      g.update()
      up_t := time.Now().UnixNano() - up_s
      dr_s := time.Now().UnixNano()
      graphics.draw(&g, screen)
      dr_t := time.Now().UnixNano() - dr_s
      show_fps(screen, up_t, dr_t)
   } else {
      // Display splash.
      dr_s := time.Now().UnixNano()
      graphics.drawSplash(&g, screen)
      dr_t := time.Now().UnixNano() - dr_s
      show_fps(screen, 0, dr_t)
   }

   return nil
}

func main() {
   if err := ebiten.Run(update, 480, 480, 1, "Tetris"); err != nil {
      panic(err)
   }
}

