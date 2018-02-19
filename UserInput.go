package main

import (
   "github.com/hajimehoshi/ebiten"
   "time"
)

const (
   HORIZONTAL_TICK = 1e8
   VERTICAL_TICK = 1e7
   SPACE_TICK = 2e8
   PAUSE_TICK = 1e8
)

var last_horizontal_tick,
    last_vertical_tick,
    last_space_tick,
    last_pause_tick int64

func time_gate(interval int64, last_tick *int64) bool {
   if time.Now().UnixNano() > (*last_tick + interval) {
      update_time_gate(last_tick)
      return true
   } else {
      return false
   }
}

func update_time_gate(last_tick *int64) {
   *last_tick = time.Now().UnixNano()
}

func input_left() bool {
   if ebiten.IsKeyPressed(ebiten.KeyLeft) &&
      time_gate(HORIZONTAL_TICK, &last_horizontal_tick) {
      return true
   }
   return false
}

func input_down() bool {
   if ebiten.IsKeyPressed(ebiten.KeyDown) &&
      time_gate(VERTICAL_TICK, &last_vertical_tick) {
      return true
   }
   return false
}

func input_right() bool {
   if ebiten.IsKeyPressed(ebiten.KeyRight) &&
      time_gate(HORIZONTAL_TICK, &last_horizontal_tick) {
      return true
   }
   return false
}

func input_space() bool {
   if ebiten.IsKeyPressed(ebiten.KeySpace) &&
      time_gate(SPACE_TICK, &last_space_tick) {
      return true
   }
   return false
}

func input_enter() bool {
   if ebiten.IsKeyPressed(ebiten.KeyEnter) &&
      time_gate(PAUSE_TICK, &last_pause_tick) {
      return true
   }
   return false
}

