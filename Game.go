package main

import (
   "math/rand"
   "time"
   "fmt"
)

const (
   GRID_X = 20
   GRID_Y = 10
   GAME_INITIAL_SPEED = 5e8
   MOVE_STEP = 8e7
   LEVEL_STEP = 5e7
   ROTATION_STEP = 1e8
   INPUT_SPEED = 0
)

type Game struct {
   score int
   grid [GRID_X][GRID_Y]int8
   cur_piece *Piece
   next_piece *Piece
   speed int64
   level int
   last_tick int64
   last_move_tick, last_rotation_tick int64
   picked_new_piece bool
   game_started, game_over bool
}

func (g *Game) print() {
   //fmt.Printf("score: %v\n", g.score)
}

func (g *Game) init() {
   rand.Seed(time.Now().UnixNano())
   g.resetScore()
   g.level = 1
   g.speed = GAME_INITIAL_SPEED // In Nano-seconds.
   g.last_tick = time.Now().UnixNano()
   g.last_move_tick = 0
   g.last_rotation_tick = 0
   g.cur_piece = g.pick_piece()
   g.next_piece = g.pick_piece()
   g.game_over = false
   g.game_started = false

   // Init grid.
   for i := 0; i < GRID_X; i++ {
      for j := 0; j < GRID_Y; j++ {
         g.grid[i][j] = -1
      }
   }
}

func (g *Game) is_started() bool {
   return g.game_started
}

func (g *Game) start_game() {
   g.game_started = true
}

func (g *Game) getScore() int {
   return g.score
}

func (g *Game) getLevel() int {
   return g.level
}

func (g *Game) grid_set(i, j int) bool {
   return g.grid[i][j] != -1
}

func (g *Game) addScore(load int) {
   g.score += load
   scores_to_level := map[int]int {
      300: 2,
      800: 3,
      1300: 4,
      2000: 5,
      4000: 6,
      6000: 7,
      8000: 8,
      10000: 9,
      15000: 10}
   for k, v := range scores_to_level {
      if g.score > k && v > g.level {
         g.setLevel(v)
      }
   }
}

func (g *Game) setLevel(level int) {
   if g.level > 10 {
      fmt.Printf("Level 10 is max!\n")
      return
   }
   g.level = level
   g.speed = GAME_INITIAL_SPEED - LEVEL_STEP * int64(g.level - 1)
   fmt.Printf("Level updated to %v, speed is %v\n", g.level, g.speed)
}

func (g *Game) resetScore() {
   g.score = 0
}

func (g *Game) end() bool {
   return false
}

func (g *Game) pick_piece() *Piece {
   new_piece := &Piece{}
   new_piece.init()
   //switch piece_type := 6; piece_type {
   switch piece_type := rand.Intn(7); piece_type {
      case 0: new_piece.init_square()
      case 1: new_piece.init_lshape_left()
      case 2: new_piece.init_lshape_right()
      case 3: new_piece.init_sshape_left()
      case 4: new_piece.init_sshape_right()
      case 5: new_piece.init_tshape()
      case 6: new_piece.init_stick()
      default: fmt.Println("Invalid piece type: %v\n", piece_type)
   }

   //fmt.Printf("Picked piece %v\n", new_piece.name())
   g.picked_new_piece = true
   return new_piece
}

func (g* Game) should_update_state() bool {
   cur_tick := time.Now().UnixNano()
   if cur_tick < (g.last_tick + g.speed) {
      return false
   } else {
      return true
   }
}

func (g* Game) update_tick() {
   g.last_tick = time.Now().UnixNano()
}

func (g *Game) update_rotation_tick() {
   g.last_rotation_tick = time.Now().UnixNano()
}

func (g *Game) update_move_tick() {
   g.last_move_tick = time.Now().UnixNano()
}

func (g *Game) piece_to_heap() {
   grid_x := g.cur_piece.grid_pos.x
   grid_y := g.cur_piece.grid_pos.y
   piece_x := 0
   //fmt.Printf("piece_to_heap: (%v,%v)\n", grid_x, grid_y)
   for i := grid_x; i < grid_x+4; i++ {
      piece_y := 0
      for j := grid_y; j < grid_y+4; j++ {
         if g.cur_piece.is_set(piece_x, piece_y) {
            //fmt.Printf("setting grid(%v,%v)\n", i, j)
            g.grid[i][j] = g.cur_piece.get_color()
         }
         piece_y++
      }
      piece_x++
   }
}

func (g *Game) remove_line(x int) {
   for i := x; i > 0; i-- {
      for j := 0; j < GRID_Y; j++ {
         g.grid[i][j] = g.grid[i-1][j]
      }
   }
}

func (g *Game) clean_heap() {
   nb_lines_removed := 0
   for i := 0; i < GRID_X; i++ {
      full_line := true
      for j := 0; j < GRID_Y; j++ {
         if g.grid[i][j] == -1 {
            full_line = false
            break
         }
      }
      if full_line {
         g.remove_line(i)
         nb_lines_removed++
      }
   }
   g.addScore(nb_lines_removed * nb_lines_removed * 100)
}

func (g *Game) collision_bottom(p *Piece) bool {
   pos := p.grid_pos
   vec := p.vector_bottom()

   for i := 0; i < 4; i++ {
      if vec[i] == -1 {
         continue
      }
      piece_x := vec[i] + pos.x
      piece_y := pos.y + i
      grid_x := piece_x + 1
      grid_y := piece_y
      if piece_y >= GRID_Y || piece_y < 0 || grid_x < 0 {
         continue
      }
      //fmt.Printf("cb? %v %v %v %v %v %v %v\n", vec, pos, piece_x, piece_y, grid_x, grid_y, i)
      if piece_x >= (GRID_X-1) || g.grid[grid_x][grid_y] != -1 {
         //fmt.Printf("Bottom collision detected on (%v,%v), Pos(%v,%v), vec(%v,%v,%v,%v)\n", grid_x, grid_y, pos.x, pos.y, vec[0], vec[1], vec[2], vec[3])
         return true
      }
   }

   return false
}

func (g *Game) collision_left(p *Piece) bool {
   pos := p.grid_pos
   vec := p.vector_left()

   for i := 0; i < 4; i++ {
      if vec[i] == -1 {
         continue
      }
      piece_x := pos.x + i
      piece_y := vec[i] + pos.y
      grid_x := piece_x
      grid_y := piece_y - 1
      if grid_x >= GRID_X || grid_x < 0 {
         break
      }

      //fmt.Printf("cl? %v %v %v %v %v %v %v\n", vec, pos, piece_x, piece_y, grid_x, grid_y, i)
      if piece_y <= 0 || g.grid[grid_x][grid_y] != -1 {
         //fmt.Printf("Left collision detected on (%v,%v), Pos(%v,%v), vec(%v,%v,%v,%v)\n", grid_x, grid_y, pos.x, pos.y, vec[0], vec[1], vec[2], vec[3])
         return true
      }
   }

   return false
}

func (g *Game) collision_right(p *Piece) bool {
   pos := p.grid_pos
   vec := p.vector_right()

   for i := 0; i < 4; i++ {
      // if piece does not have any collision for this coordinate.
      if vec[i] == -1 {
         continue
      }
      piece_x := pos.x + i
      piece_y := vec[i] + pos.y
      grid_x := piece_x
      grid_y := piece_y + 1
      if grid_x >= GRID_X || grid_x < 0 {
         break
      }

      //fmt.Printf("cr? %v %v %v %v %v %v %v\n", vec, pos, piece_x, piece_y, grid_x, grid_y, i)
      if piece_y >= (GRID_Y-1) || g.grid[grid_x][grid_y] != -1 {
         //fmt.Printf("Right collision detected on (%v,%v), Pos(%v,%v), vec(%v,%v,%v,%v)\n", grid_x, grid_y, pos.x, pos.y, vec[0], vec[1], vec[2], vec[3])
         return true
      }
   }

   return false
}

func (g *Game) should_move() bool {
   cur_tick := time.Now().UnixNano()
   if cur_tick < (g.last_move_tick + MOVE_STEP) {
      return false
   } else {
      return true
   }
}

func (g *Game) should_rotate() bool {
   cur_tick := time.Now().UnixNano()
   if cur_tick < (g.last_rotation_tick + ROTATION_STEP) {
      return false
   } else {
      return true
   }
}

func (g *Game) handle_input() {
   if g.should_move() && input_left() && !g.collision_left(g.cur_piece) {
      g.cur_piece.move_left()
      g.update_move_tick()
   }
   if g.should_move() && input_right() && !g.collision_right(g.cur_piece) {
      g.cur_piece.move_right()
      g.update_move_tick()
   }
   if input_down() && !g.collision_bottom(g.cur_piece) {
      g.cur_piece.fall()
   }
   if g.should_rotate() && input_space() {
      g.rotate_piece_if_possible()
      g.update_rotation_tick()
   }
}

func (g *Game) rotate_piece_if_possible() {
   // Create a copy of the cur piece.
   piece_rotated := *g.cur_piece

   // Rotate it.
   piece_rotated.rotate()

   // If it does not collide, update cur_piece.
   if !g.collides(&piece_rotated) {
      g.cur_piece = &piece_rotated
   }
}

func (g *Game) collides(p *Piece) bool {
   return g.collision_right(p) ||
          g.collision_left(p) ||
          g.collision_bottom(p)
}

func (g* Game) update() {
   g.handle_input()
   if g.should_update_state() {
      //fmt.Printf("Updating game now! %v\n", g.last_tick)

      if g.collision_bottom(g.cur_piece) {
         g.addScore(10)
         g.piece_to_heap()
         g.clean_heap()
         g.cur_piece = g.next_piece
         if g.collides(g.cur_piece) {
            fmt.Printf("Game Over\n")
            g.game_over = true
         }
         g.next_piece = g.pick_piece()
      } else {
         g.cur_piece.fall()
      }

      g.update_tick()
   }
}
