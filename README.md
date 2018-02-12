# GoTetris
GoTetris is simple Tetris game written in Go for fun and profit. GoTetris uses <a href="https://github.com/hajimehoshi/ebiten">Ebiten</a>, a simple 2D Go library. The game has currently been tessted on MacOS only, although it should work on Windows and Linux as well. Feel free to report any issues you may find along your route.

<img src="https://raw.githubusercontent.com/jphalimi/GoTetris/master/resources/screenshot_splash.png" alt="Splash screen" width="250" height="250" /> <img src="https://raw.githubusercontent.com/jphalimi/GoTetris/master/resources/screenshot.png" alt="Splash screen" width="250" height="250" />

There are no binary releases available at this time. Please follow the guide below to build the game yourself.

# Prerequisites
1. Have Go installed on your machine. Go the <a href="https://golang.org/doc/install">official website</a> if you don't yet.
2. You will need to set the GOPATH environment variable setup appropriately as well.
3. For Linux users: Install `xorg-dev` and `freeglut3-dev` packages.

# Setup & build
First, you will need to download the game dependencies. Please run the following command in your local project folder:
```
go get
```
Now, you can build the game:
```
go build
```

# Play
Now you can launch the generated binary:
```
./GoTetris
```
Have fun! :-)

# Credits
The "blocks.png" image comes from the "blocks" example from the <a href="https://github.com/hajimehoshi/ebiten">Ebiten</a> project.
