package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"atomicgo.dev/cursor"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

type Cell struct {
	char        rune
	highlighted bool
}

type Matrix struct {
	columns       [][]Cell
	width, height int
}

func (m *Matrix) init() {
	m.width, m.height = screen.Size()

	for i := 0; i < m.width; i++ {
		if i%2 != 0 {
			continue
		}

		var column []Cell
		for j := 0; j < m.height; j++ {
			randIndex := rand.Intn(len(chars))
			char := rune(chars[randIndex])
			cell := Cell{char: char, highlighted: false}

			column = append(column, cell)
		}

		m.columns = append(m.columns, column)
	}
}

func (m *Matrix) update() {
	const invertChance = 10

	// update the first row
	for i := 0; i < len(m.columns); i++ {
		previous := m.columns[i][1]
		flip := rand.Intn(invertChance) == 0

		if !flip {
			m.columns[i][0].highlighted = previous.highlighted
		} else {
			m.columns[i][0].highlighted = !previous.highlighted
		}
	}

	// update the rest of the rows
	for i := 0; i < len(m.columns); i++ {
		for j := len(m.columns[i]) - 1; j > 0; j-- {
			previous := m.columns[i][j-1]
			m.columns[i][j].highlighted = previous.highlighted
		}
	}
}

func (m *Matrix) draw() {
	screen.MoveTopLeft()
	width, height := screen.Size()

	// screen size changed
	if width != m.width || height != m.height {
		m.columns = nil
		m.init()
	}

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if c%2 != 0 {
				fmt.Print(" ")
				continue
			}

			cell := m.columns[c/2][r]

			if cell.highlighted {
				color.Set(color.FgHiGreen)
			} else {
				color.Set(color.FgBlack)
			}

			fmt.Printf("%c", cell.char)
		}
	}
}

func catchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cursor.Show()
		screen.Clear()
		os.Exit(0)
	}()
}

func init() {
	color.Set(color.FgHiGreen)
	cursor.Hide()
	screen.Clear()
	screen.MoveTopLeft()
	catchInterrupt()
}

func main() {
	var m Matrix
	m.init()

	for {
		m.update()
		m.draw()
		time.Sleep(time.Millisecond * 50)
	}
}
