package main

import (
   //"fmt"
   "github.com/hajimehoshi/ebiten"
)

func input_left() bool {
   if ebiten.IsKeyPressed(ebiten.KeyLeft) {
      //fmt.Printf("Left\n")
      return true
   }
   return false
}

func input_down() bool {
   if ebiten.IsKeyPressed(ebiten.KeyDown) {
      //fmt.Printf("Down\n")
      return true
   }
   return false
}

func input_right() bool {
   if ebiten.IsKeyPressed(ebiten.KeyRight) {
      //fmt.Printf("Right\n")
      return true
   }
   return false
}

func input_space() bool {
   if ebiten.IsKeyPressed(ebiten.KeySpace) {
      //fmt.Printf("Space\n")
      return true
   }
   return false
}

func input_enter() bool {
   if ebiten.IsKeyPressed(ebiten.KeyEnter) {
      //fmt.Printf("Enter\n")
      return true
   }
   return false
}

