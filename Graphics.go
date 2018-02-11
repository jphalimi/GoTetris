package main

import (
   "github.com/hajimehoshi/ebiten"
   "github.com/hajimehoshi/ebiten/ebitenutil"
   "github.com/hajimehoshi/ebiten/text"
	"github.com/golang/freetype/truetype"
   "golang.org/x/image/font"
   "image/color"
   "fmt"
   "path"
  	"io/ioutil"
   "strconv"
   "image"
   "time"
)

const (
   BLUE = 0
   RED = 1
   GREEN = 2
   YELLOW = 3
   PURPLE = 4
   PINK = 5
   GREY = 6
   WHITE = 7
   BLACK = 8
   NONE = 8
)

const (
   margin_h = 20
   margin_w = 20
)

type CoordF64 struct {
   x, y float64
}

type CoordInt struct {
   x, y int
}

type GBSurface struct {
   block *ebiten.Image
}

type Block struct {
   color int8
   pos CoordF64
}

type Square struct {
   surface *ebiten.Image
   block Block
   size CoordInt
}

type Graphics struct {
   block, background, splash, next_piece, level, score *ebiten.Image
   block_size CoordInt
   leftpanel, righttoppanel, rightbottompanel *Square
   game_offset CoordF64
   game_blocks [GRID_X][GRID_Y] Block
   next_blocks [4][4] Block
   ginit bool
   font font.Face
   last_splash_tick int64
}

func (g *Graphics) is_init() bool {
   return g.ginit
}

func new_image(filename string) *ebiten.Image {
   new_image, _,  err := ebitenutil.NewImageFromFile(
      path.Join(GetCurrentDir(),
         filename),
         ebiten.FilterNearest)

   if err != nil {
      panic(err)
   }

   return new_image
}

func (g *Graphics) init(screen *ebiten.Image) {
   w, h := screen.Size()
   g.init_panels(w, h)
   g.block = new_image("blocks.png")
   g.background = new_image("background.png")
   g.next_piece = new_image("next_piece.png")
   g.level = new_image("level.png")
   g.score = new_image("score.png")
   g.splash = new_image("splash.png")
   g.allocate_game_blocks(w, h)
   g.allocate_next_blocks(w, h)
   g.init_font()
   
   g.ginit = true
}

func (b *Block) init_block(offset CoordF64) {
   b.color = NONE
   b.pos = offset
}

func create_square(width, height float64, offset CoordF64, color int8) *Square {
   return &Square{
      surface: init_surface(int(width), int(height), color),
      block: Block{color, offset},
      size: CoordInt{int(width), int(height)}}
}

func (g *Graphics) init_panels(width, height int) {
   // Main game panel
   g.game_offset = CoordF64{margin_h / 2, margin_w / 2}
   leftpanel_width := float64(width / 2.0 - margin_w / 2.0)
   leftpanel_height := float64(height - margin_h)
   g.leftpanel = create_square(leftpanel_width, leftpanel_height, g.game_offset, WHITE)

   // Right panels
   width_lp := float64(width / 2.0 - margin_w)
   height_lp := float64(height / 2.0 - margin_h)
   woffset_lp := float64(width / 2.0 + margin_w / 2.0)
   g.righttoppanel = create_square(width_lp, height_lp,
      CoordF64{margin_h / 2.0, woffset_lp}, WHITE)
   g.rightbottompanel = create_square(width_lp, height_lp,
      CoordF64{float64(height / 2.0 + margin_h / 2.0), woffset_lp}, WHITE)
}

func (g *Graphics) allocate_game_blocks(width, height int) {
   // Allocate game blocks.
   offset := g.game_offset
   leftpanel_width := width / 2 - margin_w / 2
   leftpanel_height := height - margin_h
   step_x := float64(leftpanel_height) / float64(GRID_X)
   step_y := float64(leftpanel_width) / float64(GRID_Y)
   g.block_size.x = int(step_x)
   g.block_size.y = int(step_y)

   fmt.Printf("steps = %v %v %v %v %v %v\n", step_x, step_y, GRID_X, GRID_Y, leftpanel_height, leftpanel_width)
   for i := 0; i < GRID_X; i++ {
      for j := 0; j < GRID_Y; j++ {
         g.game_blocks[i][j].init_block(offset)
         offset.y += step_y
      }
      offset.x += step_x
      offset.y = g.game_offset.y
   }
}

func (g *Graphics) allocate_next_blocks(width, height int) {
   // Allocate next blocks.
   offset := g.righttoppanel.block.pos
   leftpanel_width := width / 2 - margin_w / 2
   leftpanel_height := height - margin_h
   step_x := float64(leftpanel_height) / float64(GRID_X)
   step_y := float64(leftpanel_width) / float64(GRID_Y)
   offset.x += 70
   offset.y += 70
   for i := 0; i < 4; i++ {
      for j := 0; j < 4; j++ {
         g.next_blocks[i][j].init_block(offset)
         offset.y += step_y
      }
      offset.x += step_x
      offset.y = g.righttoppanel.block.pos.y + 70
   }
}

func (g *Graphics) init_font() {
   // Font allocation
   f, err := ebitenutil.OpenFile(path.Join(GetCurrentDir(), "superstar-m54.ttf"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	tt, err := truetype.Parse(b)
	if err != nil {
		panic(err)
   }

	const dpi = 72
	g.font = truetype.NewFace(tt, &truetype.Options{
		Size:    30,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func (g *Graphics) drawScoreAndLevel(game *Game, screen *ebiten.Image) {
   offset := g.rightbottompanel.block.pos
   op := &ebiten.DrawImageOptions{}
   op.GeoM.Translate(offset.y, offset.x)
   screen.DrawImage(g.score, op)
   text.Draw(screen, strconv.Itoa(game.getScore()), g.font, 
      int(offset.y) + 10, int(offset.x) + 80, color.White);
   op = &ebiten.DrawImageOptions{}
   op.GeoM.Translate(offset.y, offset.x + 100)
   screen.DrawImage(g.level, op)
   text.Draw(screen, strconv.Itoa(game.getLevel()), g.font, 
      int(offset.y) + 10, int(offset.x) + 175, color.White);
}

func (g *Graphics) GetOffset(b *Block) int {
   switch b.color {
   case RED:      { return 0 }
   case YELLOW:   { return 1*g.block_size.y }
   case GREEN:    { return 2*g.block_size.y }
   case BLUE:     { return 3*g.block_size.y }
   case PURPLE:   { return 4*g.block_size.y }
   case PINK:     { return 5*g.block_size.y }
   case GREY:     { return 6*g.block_size.y }
   case NONE:     { return 0*g.block_size.y }
   }
   return 0
}

func GetColor(c int8) color.Color {
   switch c {
   case BLUE:   { return color.RGBA{0x24, 0x70, 0xca, 0xff}}
   case RED:    { return color.RGBA{0xbf, 0x37, 0x37, 0xff}}
   case GREEN:  { return color.RGBA{0x27, 0xaa, 0x2d, 0xff}}
   case YELLOW: { return color.RGBA{0xf1, 0xf1, 0x36, 0xff}}
   case PURPLE: { return color.RGBA{0xb4, 0x2c, 0xc4, 0xff}}
   case GREY:   { return color.RGBA{0x87, 0x87, 0x87, 0xff}}
   case WHITE:  { return color.RGBA{0xff, 0xff, 0xff, 0x44}}
   case BLACK:  { return color.RGBA{0x00, 0x00, 0x00, 0xDD}}
   default:     { return color.Black }
   }
}

func (g *Graphics) drawNextPiece(game *Game, screen *ebiten.Image) {
   offset := g.righttoppanel.block.pos
   op := &ebiten.DrawImageOptions{}
   op.GeoM.Translate(offset.y, offset.x)
   screen.DrawImage(g.next_piece, op)

   op = &ebiten.DrawImageOptions{}
   for i := 0; i < 4; i++ {
   	for j := 0; j < 4; j++ {
         if g.next_blocks[i][j].color == NONE {
            continue
         }
   		op.GeoM.Reset()
   		op.GeoM.Translate(
            g.next_blocks[i][j].pos.y,
            g.next_blocks[i][j].pos.x)
         block_offset := g.GetOffset(&g.next_blocks[i][j])
   		src := image.Rect(block_offset, 0,
            block_offset + int(g.block_size.y),
            int(g.block_size.x))
   		op.SourceRect = &src
         //fmt.Printf("Drawing piece %v;%v, SR{%v,%v,%v,%v}\n", g.game_blocks[i][j].size.x, g.game_blocks[i][j].size.y)
   		screen.DrawImage(g.block, op)
   	}
   }

   // Update the colors of the next blocks.
   for i := 0; i < 4; i++ {
      for j := 0; j < 4; j++ {
         if game.next_piece.grid[i][j] {
            g.next_blocks[i][j].setColor(game.next_piece.color)
         } else {
            g.next_blocks[i][j].resetColors()
         }
      }
   }
}

func (g *Graphics) draw_blocks(game *Game, screen *ebiten.Image) {
   op := &ebiten.DrawImageOptions{}
	for i := 0; i < GRID_X; i++ {
		for j := 0; j < GRID_Y; j++ {
         if g.game_blocks[i][j].color == NONE {
            continue
         }
			op.GeoM.Reset()
			op.GeoM.Translate(
            g.game_blocks[i][j].pos.y,
            g.game_blocks[i][j].pos.x)
         block_offset := g.GetOffset(&g.game_blocks[i][j])
			src := image.Rect(block_offset, 0,
            block_offset + int(g.block_size.y),
            int(g.block_size.x))
			op.SourceRect = &src
			screen.DrawImage(g.block, op)
		}
   }
}

func (g *Graphics) should_draw_starter() bool {
   cur_tick := time.Now().UnixNano() / (1e9) // 1s
   if cur_tick % 2 == 0 {
      return false
   } else {
      return true
   }
}

func (g *Graphics) drawSplash(game *Game, screen *ebiten.Image) {
   op := &ebiten.DrawImageOptions{}
   screen.DrawImage(g.splash, op)
   if input_space() {
      game.start_game()
   }
   if g.should_draw_starter() {
      text.Draw(screen, "Press Space when ready", g.font, 
         80, 440, color.White);
   }
}

func (g *Graphics) drawBackground(screen *ebiten.Image) {
   op := &ebiten.DrawImageOptions{}
   screen.DrawImage(g.background, op)
}

func (g *Graphics) draw(game *Game, screen *ebiten.Image) {
   g.drawBackground(screen)
   g.drawPanels(screen)
   g.alter_colors(game)

   g.draw_blocks(game, screen)

   g.drawScoreAndLevel(game, screen)
   g.drawNextPiece(game, screen)
}

func (g *Graphics) alter_colors(game *Game) {
   // Set colors for the heap.
   for i := 0; i < GRID_X; i++ {
      for j := 0; j < GRID_Y; j++ {
         if game.grid_set(i,j) {
            g.setColor(i, j, game.grid[i][j])
         } else {
            g.game_blocks[i][j].resetColors()
         }
      }
   }
   
   // Set colors for the current piece.
   start_x := game.cur_piece.grid_pos.x
   start_y := game.cur_piece.grid_pos.y
   piece_x := 0
   for i := start_x; i < start_x + 4; i++ {
      piece_y := 0
      for j := start_y; j < start_y+4; j++ {
         if game.cur_piece.grid[piece_x][piece_y] {
            g.setColor(i, j, game.cur_piece.color)
         }
         piece_y++
      }
      piece_x++
   }
}

func (g *Graphics) setColor(i, j int, c int8) {
   if (i >= 0 && i < GRID_X &&
       j >= 0 && i < GRID_X) {
      g.game_blocks[i][j].setColor(c)
   }
}

func (b *Block) resetColors() {
   b.color = NONE
}

func (b *Block) setColor(color int8) {
   b.color = color
}

func (g *Graphics) drawPanels(screen *ebiten.Image) {
   g.leftpanel.draw(screen)
   g.righttoppanel.draw(screen)
   g.rightbottompanel.draw(screen)
}

func (s *Square) draw(screen *ebiten.Image) {
   opts := &ebiten.DrawImageOptions{}
   opts.GeoM.Translate(s.block.pos.y, s.block.pos.x)
   s.surface.Fill(GetColor(s.block.color))
	src := image.Rect(0, 0, s.size.x, s.size.y)
	opts.SourceRect = &src
   screen.DrawImage(s.surface, opts)
}

func init_surface(width, height int,
                  color int8) *ebiten.Image {
   square, _ := ebiten.NewImage(width, height, ebiten.FilterNearest)
   square.Fill(GetColor(color))

   return square
}

//func 
