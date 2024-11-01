package main

import (
    "syscall/js"
    "math"
    "time"
)

const (
    canvasWidth  = 800
    canvasHeight = 600
    paddleHeight = 100
    paddleWidth  = 10
    ballSize     = 10

    playerSpeed  = 5
    ballSpeed    = 5
)

type Game struct {
    ballX, ballY      float64
    ballVelX, ballVelY float64
    player1Y, player2Y float64
}

var game = Game{
    ballX:   canvasWidth / 2,
    ballY:   canvasHeight / 2,
    ballVelX: ballSpeed,
    ballVelY: ballSpeed,
    player1Y: canvasHeight / 2,
    player2Y: canvasHeight / 2,
}

// JavaScript canvas API interface
var ctx js.Value

// Initialize the game
func initGame() {
    game.ballX = canvasWidth / 2
    game.ballY = canvasHeight / 2
    game.ballVelX = ballSpeed
    game.ballVelY = ballSpeed
    game.player1Y = canvasHeight / 2
    game.player2Y = canvasHeight / 2
}

// Poll the key inputs from JavaScript
func getKeyInput() float64 {
    // Get the window object
    window := js.Global()

    // Access the keysPressed object
    keysPressed := window.Get("keysPressed")
    if keysPressed.IsUndefined() {
        return 0 // No input for this frame
    }

    yDirection := 0.0

    // If ArrowUp is pressed, move up
    if keysPressed.Get("ArrowUp").Bool() {
        yDirection = -playerSpeed
    }

    // If ArrowDown is pressed, move down
    if keysPressed.Get("ArrowDown").Bool() {
        yDirection = playerSpeed
    }

    return yDirection
}

// Game logic - called every frame
func updateGame() {
    // Ball movement
    game.ballX += game.ballVelX
    game.ballY += game.ballVelY

    // Ball collision with top and bottom walls
    if game.ballY-ballSize <= 0 || game.ballY+ballSize >= canvasHeight {
        game.ballVelY = -game.ballVelY
    }

    // Ball collision with player paddles
    if game.ballX-ballSize <= paddleWidth {
        if game.ballY >= game.player1Y-paddleHeight/2 && game.ballY <= game.player1Y+paddleHeight/2 {
            game.ballVelX = -game.ballVelX
        }
    }
    if game.ballX+ballSize >= canvasWidth-paddleWidth {
        if game.ballY >= game.player2Y-paddleHeight/2 && game.ballY <= game.player2Y+paddleHeight/2 {
            game.ballVelX = -game.ballVelX
        }
    }

    // Player 1 (user) paddle movement based on key input
    direction := getKeyInput()
    game.player1Y += direction

    // Clamp player 1's paddle within the screen bounds
    if game.player1Y-paddleHeight/2 < 0 {
        game.player1Y = paddleHeight / 2
    } else if game.player1Y+paddleHeight/2 > canvasHeight {
        game.player1Y = canvasHeight - paddleHeight/2
    }

    // Update player 2 AI (automatic movement)
    if math.Abs(game.ballY-game.player2Y) > playerSpeed {
        if game.ballY > game.player2Y {
            game.player2Y += playerSpeed
        } else {
            game.player2Y -= playerSpeed
        }
    }

    // Check scoring
    if game.ballX+ballSize >= canvasWidth || game.ballX-ballSize <= 0 {
        initGame()  // Reset game
    }
}

// Render the game UI
func renderGame() {
    ctx.Call("clearRect", 0, 0, canvasWidth, canvasHeight)

    // Draw player paddles
    ctx.Call("fillRect", 0, game.player1Y-paddleHeight/2, paddleWidth, paddleHeight)
    ctx.Call("fillRect", canvasWidth-paddleWidth, game.player2Y-paddleHeight/2, paddleWidth, paddleHeight)

    // Draw ball
    ctx.Call("beginPath")
    ctx.Call("arc", game.ballX, game.ballY, ballSize, 0, 2*math.Pi)
    ctx.Call("fill")
}

// Main game loop
func gameLoop() {
    updateGame()
    renderGame()
    time.AfterFunc(time.Second/60, gameLoop)  // Keep running at ~60 FPS
}

// Entry point
func main() {
    // Set up the canvas context
    doc := js.Global().Get("document")
    canvas := doc.Call("getElementById", "gameCanvas")
    ctx = canvas.Call("getContext", "2d")

    // Initialize the game state
    initGame()

    // Start game loop
    gameLoop()

    // Prevent the main function from exiting immediately
    quit := make(chan struct{})
    <-quit
}