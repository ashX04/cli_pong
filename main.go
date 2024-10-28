package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	paddleHeight      = 4
	paddleWidth       = 1
	ballSize          = 1
	initialBallSpeedX = 1
	initialBallSpeedY = 1
)

// Game represents the state of the Pong game
type Game struct {
	width      int
	height     int
	paddle1    float64 // left paddle Y position
	paddle2    float64 // right paddle Y position
	ballX      float64
	ballY      float64
	ballSpeedX float64
	ballSpeedY float64
	score1     int
	score2     int
}

// Initialize the game state
func initialModel() Game {
	return Game{
		width:      80,
		height:     24,
		paddle1:    10,
		paddle2:    10,
		ballX:      40,
		ballY:      12,
		ballSpeedX: initialBallSpeedX,
		ballSpeedY: initialBallSpeedY,
	}
}

func (g Game) Init() tea.Cmd {
	return tick
}

// Update handles game logic and input
func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return g, tea.Quit
		case "w":
			if g.paddle1 > 0 {
				g.paddle1--
			}
		case "s":
			if g.paddle1 < float64(g.height-paddleHeight) {
				g.paddle1++
			}
		case "up":
			if g.paddle2 > 0 {
				g.paddle2--
			}
		case "down":
			if g.paddle2 < float64(g.height-paddleHeight) {
				g.paddle2++
			}
		}
	case tickMsg:
		// Update ball position
		g.ballX += g.ballSpeedX
		g.ballY += g.ballSpeedY

		// Ball collision with top and bottom walls
		if g.ballY <= 0 || g.ballY >= float64(g.height-1) {
			g.ballSpeedY = -g.ballSpeedY
		}

		// Ball collision with paddles
		if g.ballX <= float64(paddleWidth) && g.ballY >= g.paddle1 && g.ballY <= g.paddle1+float64(paddleHeight) {
			g.ballSpeedX = -g.ballSpeedX
		}
		if g.ballX >= float64(g.width-paddleWidth-1) && g.ballY >= g.paddle2 && g.ballY <= g.paddle2+float64(paddleHeight) {
			g.ballSpeedX = -g.ballSpeedX
		}

		// Score points
		if g.ballX <= 0 {
			g.score2++
			g.resetBall()
		}
		if g.ballX >= float64(g.width) {
			g.score1++
			g.resetBall()
		}

		return g, tick
	}

	return g, nil
}

func (g *Game) resetBall() {
	g.ballX = float64(g.width) / 2
	g.ballY = float64(g.height) / 2
	g.ballSpeedX = initialBallSpeedX
	g.ballSpeedY = initialBallSpeedY
}

// View renders the game
func (g Game) View() string {
	// Create the game board
	board := make([][]string, g.height)
	for i := range board {
		board[i] = make([]string, g.width)
		for j := range board[i] {
			board[i][j] = " "
		}
	}

	// Draw paddles
	for i := 0; i < paddleHeight; i++ {
		if int(g.paddle1)+i < g.height {
			board[int(g.paddle1)+i][0] = "█"
		}
		if int(g.paddle2)+i < g.height {
			board[int(g.paddle2)+i][g.width-1] = "█"
		}
	}

	// Draw ball
	if int(g.ballY) >= 0 && int(g.ballY) < g.height && int(g.ballX) >= 0 && int(g.ballX) < g.width {
		board[int(g.ballY)][int(g.ballX)] = "●"
	}

	// Convert board to string
	var s string
	s += fmt.Sprintf("Score: %d - %d\n", g.score1, g.score2)
	for i := range board {
		for j := range board[i] {
			s += board[i][j]
		}
		s += "\n"
	}
	s += "\nControls: W/S for left paddle, ↑/↓ for right paddle, Q to quit"

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(s)
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second / 30)
	return tickMsg{}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
