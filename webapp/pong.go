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
    player1X, player1Y float64
    player2Y float64
}

var game = Game{
    ballX:   canvasWidth / 2,
    ballY:   canvasHeight / 2,
    ballVelX: ballSpeed,
    ballVelY: ballSpeed,
    player1X: 0,
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

func getKeyInput() (float64, float64) {
    // Get the global window object
    window := js.Global()

    // Ensure keysPressed exists
    keysPressed := window.Get("keysPressed")
    if keysPressed.IsUndefined() {
        return 0, 0 // No input for this frame
    }

    // Y-axis movement (up/down)
    yDirection := 0.0
    if keysPressed.Get("ArrowUp").Bool() {
        yDirection = -playerSpeed
    }
    if keysPressed.Get("ArrowDown").Bool() {
        yDirection = playerSpeed
    }

    // X-axis movement (left/right)
    xDirection := 0.0
    if keysPressed.Get("ArrowLeft").Bool() {
        xDirection = -playerSpeed
    }
    if keysPressed.Get("ArrowRight").Bool() {
        xDirection = playerSpeed
    }

    // Return both X-axis and Y-axis direction
    return xDirection, yDirection
}

// Game logic - called every frame
func updateGame() {
    // Ball movement (same as before)
    game.ballX += game.ballVelX
    game.ballY += game.ballVelY

    // Ball collision with top and bottom walls
    if game.ballY-ballSize <= 0 || game.ballY+ballSize >= canvasHeight {
        game.ballVelY = -game.ballVelY
    }

    // Ball collision with player 1's paddle
    // Check if the ball is within the X range (player1X to player1X + paddleWidth)
    // AND within the Y range (player1Y - paddleHeight/2 to player1Y + paddleHeight / 2)
    if game.ballX-ballSize <= game.player1X + paddleWidth && game.ballX+ballSize >= game.player1X {
        if game.ballY >= game.player1Y-paddleHeight/2 && game.ballY <= game.player1Y+paddleHeight/2 {
            game.ballVelX = -game.ballVelX
        }
    }

    // Ball collision with player 2's (AI) paddle
    if game.ballX+ballSize >= canvasWidth-paddleWidth {
        if game.ballY >= game.player2Y-paddleHeight/2 && game.ballY <= game.player2Y+paddleHeight/2 {
            game.ballVelX = -game.ballVelX
        }
    }

    // Player 1 (user) paddle movement (horizontal and vertical)
    xDir, yDir := getKeyInput()
    game.player1Y += yDir
    game.player1X += xDir // Allow player 1 to move left/right

    // Clamp player 1's paddle within the screen bounds
    if game.player1Y-paddleHeight/2 < 0 {
        game.player1Y = paddleHeight / 2
    } else if game.player1Y+paddleHeight/2 > canvasHeight {
        game.player1Y = canvasHeight - paddleHeight/2
    }

    // Clamp player 1's X position (left side to middle of the screen)
    if game.player1X < 0 {
        game.player1X = 0
    } else if game.player1X+paddleWidth > canvasWidth/2 {
        game.player1X = canvasWidth/2 - paddleWidth
    }

    // Update player 2 AI (only Y-axis movement, unchanged)
    if math.Abs(game.ballY-game.player2Y) > playerSpeed {
        if game.ballY > game.player2Y {
            game.player2Y += playerSpeed
        } else {
            game.player2Y -= playerSpeed
        }
    }

    // Check scoring/reset (same as before)
    if game.ballX+ballSize >= canvasWidth || game.ballX-ballSize <= 0 {
        initGame()  // Reset game
    }
}

// Render the game UI
func renderGame() {
    // Clear the canvas
    ctx.Call("clearRect", 0, 0, canvasWidth, canvasHeight)

    // Draw player 1's paddle (with both X and Y)
    ctx.Call("fillRect", game.player1X, game.player1Y-paddleHeight/2, paddleWidth, paddleHeight)

    // Draw player 2's paddle (only Y)
    ctx.Call("fillRect", canvasWidth-paddleWidth, game.player2Y-paddleHeight/2, paddleWidth, paddleHeight)

    // Draw the ball
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