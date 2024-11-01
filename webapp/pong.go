package main

import (
    "syscall/js"
    "math"
    "math/rand"
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

    // Randomize the starting direction of the ball
    angle := rand.Float64()*120 + 30 // Angle between 30 and 150 degrees
    radAngle := angle * (math.Pi / 180.0) // Convert degrees to radians

    // Set initial velocity based on the angle
    game.ballVelX = ballSpeed * math.Cos(radAngle)  // X component of velocity
    game.ballVelY = ballSpeed * math.Sin(radAngle)  // Y component of velocity

    // Randomize whether the ball starts moving to the left or right
    if rand.Float64() > 0.5 {
        game.ballVelX = -game.ballVelX
    }

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
    if game.ballX-ballSize <= game.player1X+paddleWidth && game.ballX+ballSize >= game.player1X {
        if game.ballY >= game.player1Y-paddleHeight/2 && game.ballY <= game.player1Y+paddleHeight/2 {
            game.ballVelX = -game.ballVelX

            // Add spin/bounce variation based on where it hits the paddle
            game.ballVelY += (game.ballY - game.player1Y) * 0.05  // Fine-tune this effect
        }
    }

    // Ball collision with player 2's (AI) paddle
    if game.ballX+ballSize >= canvasWidth-paddleWidth && game.ballX-ballSize <= canvasWidth { // Right side collision check
        if game.ballY >= game.player2Y-paddleHeight/2 && game.ballY <= game.player2Y+paddleHeight/2 {
            game.ballVelX = -game.ballVelX // Bounce the ball back
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

    // Update player 2 AI paddle movement with randomness (only Y)
    // Introduce randomness to AI movements
    randomFactor := rand.Float64() * 2  // Full randomness factor
    aiSpeed := playerSpeed * 0.7 * randomFactor  // AI will be slower than player

    // AI reaction delay (set it to 1 to make the AI react immediately)
    reactionDelay := 1  
    if rand.Intn(reactionDelay) == 0 {  
        if math.Abs(game.ballY - game.player2Y) > aiSpeed {
            if game.ballY > game.player2Y {
                game.player2Y += aiSpeed
            } else {
                game.player2Y -= aiSpeed
            }
        }
    }

    // Check scoring/reset
    if game.ballX+ballSize >= canvasWidth || game.ballX-ballSize <= 0 {
        initGame()  // Reset game if the ball hits the right or left boundaries
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
    rand.Seed(time.Now().UnixNano()) // Seed randomness with current time
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