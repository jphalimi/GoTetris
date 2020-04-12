package main

import (
   "fmt"
   "math/rand"
)

const (
   STRAIGHT = 0
   TO_THE_RIGHT = 1
   UPSIDE_DOWN = 2
   TO_THE_LEFT = 3
)

type PieceGrid [4][4]bool

type Piece struct {
   grid_pos CoordInt
   grid PieceGrid
   color int8
   direction int
   debug_name string
}

func (p *Piece) init() {
   p.grid_pos = CoordInt{0,3}
   p.direction = STRAIGHT
   p.color = int8(rand.Intn(7))
}

func (p *Piece) init_square() {
   p.grid = PieceGrid{
      {false,false,false,false},
      {false,true,true,false},
      {false,true,true,false},
      {false,false,false,false}}
   p.debug_name = "Square"
}

func (p *Piece) init_lshape_left() {
   p.grid = PieceGrid{
      {false,true,false,false},
      {false,true,false,false},
      {false,true,true,false},
      {false,false,false,false}}
   p.debug_name = "LShape left"
}

func (p *Piece) init_lshape_right() {
   p.grid = PieceGrid{
      {false,false,true,false},
      {false,false,true,false},
      {false,true,true,false},
      {false,false,false,false}}
   p.debug_name = "LShape Right"
}

func (p *Piece) init_sshape_left() {
   p.grid = PieceGrid{
      {false,true,false,false},
      {false,true,true,false},
      {false,false,true,false},
      {false,false,false,false}}
   p.debug_name = "SShape Left"
}

func (p *Piece) init_sshape_right() {
   p.grid = PieceGrid{
      {false,false,true,false},
      {false,true,true,false},
      {false,true,false,false},
      {false,false,false,false}}
   p.debug_name = "SShape Right"
}

func (p *Piece) init_tshape() {
   p.grid = PieceGrid{
      {false,false,true,false},
      {false,true,true,false},
      {false,false,true,false},
      {false,false,false,false}}
   p.debug_name = "TShape"
}

func (p *Piece) init_stick() {
   p.grid = PieceGrid{
      {false,false,true,false},
      {false,false,true,false},
      {false,false,true,false},
      {false,false,true,false}}
   p.debug_name = "Stick"
}

func (p *Piece) name() string {
   return p.debug_name
}

func (p *Piece) rotate() {
   var new_grid PieceGrid
   for i := 0; i < 4; i++ {
      for j := 0; j < 4; j++ {
         new_grid[j][3-i] = p.grid[i][j]
      }
   }

   for i := 0; i < 4; i++ {
      for j := 0; j < 4; j++ {
         p.grid[i][j] = new_grid[i][j]
      }
   }

   // Update direction
   p.direction++
   p.direction %= 4

   //fmt.Printf("Piece rotated, dir = %v.\n", p.direction)

   p.adjust_pos_if_necessary()
}

func (p *Piece) adjust_pos_if_necessary() {
   if p.debug_name == "SShape Left" || p.debug_name == "SShape Right" {
      //fmt.Printf("Updating grid pos %v...\n", p.grid_pos)
      switch p.direction {
         case TO_THE_LEFT: {
            p.grid_pos.y++
         }
         case STRAIGHT: {
            p.grid_pos.y--
            p.grid_pos.x--
         }
      }
      //fmt.Printf("Updating grid pos %v, dir = %v...\n", p.grid_pos, p.direction)
   }
   if p.debug_name == "TShape" {
      switch p.direction {
         case UPSIDE_DOWN: {
            p.grid_pos.y++
            p.grid_pos.x--
         }
         case STRAIGHT: {
            p.grid_pos.y--
         }
      }
   }
   if p.debug_name == "Stick" {
      switch p.direction {
         case TO_THE_LEFT: {
            p.grid_pos.y--
         }
         case UPSIDE_DOWN: {
            p.grid_pos.y++
            p.grid_pos.x--
         }
      }
   }
}

func (p *Piece) fall() {
   //fmt.Printf("Falling (%v,%v) to (%v,%v)\n", p.grid_pos.x, p.grid_pos.y, p.grid_pos.x+1, p.grid_pos.y)
   p.grid_pos.x++
}

func (p *Piece) move_right() {
   //fmt.Printf("-> (%v,%v) to (%v,%v)\n", p.grid_pos.x, p.grid_pos.y, p.grid_pos.x, p.grid_pos.y+1)
   p.grid_pos.y++
}

func (p *Piece) move_left() {
   //fmt.Printf("<- (%v,%v) to (%v,%v)\n", p.grid_pos.x, p.grid_pos.y, p.grid_pos.x, p.grid_pos.y-1)
   p.grid_pos.y--
}

func (p *Piece) is_set(i, j int) bool {
   return p.grid[i][j]
}

func (p *Piece) get_color() int8 {
   return p.color
}

func (p *Piece) vector_bottom() [4]int {
   var vec[4] int
   for i := 0; i < 4; i++ {
      found := false
      for j := 0; j < 4; j++ {
         if p.grid[j][i] {
            vec[i] = j
            found = true
         }
      }
      if !found {
         vec[i] = -1
      }
   }
   return vec
}

func (p *Piece) vector_left() [4]int {
   var vec[4] int
   for i := 0; i < 4; i++ {
      found := false
      for j := 3; j >= 0; j-- {
         if p.grid[i][j] {
            vec[i] = j
            found = true
         }
      }
      if !found {
         vec[i] = -1
      }
   }
   return vec
}

func (p *Piece) vector_right() [4]int {
   var vec[4] int
   for i := 0; i < 4; i++ {
      found := false
      for j := 0; j < 4; j++ {
         if p.grid[i][j] {
            vec[i] = j
            found = true
         }
      }
      if !found {
         vec[i] = -1
      }
   }
   return vec
}

