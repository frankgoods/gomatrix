package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tm "github.com/buger/goterm"
)

var m sync.Mutex

var screenWidth = tm.Width() - 1
var screenHeight = tm.Height() - 1

var lettersOld [24]string
var lettersNew [24]string

func init() {
	for i := range lettersNew {
		lettersNew[i] = tm.Color(string('a'+i), tm.WHITE)
		lettersOld[i] = tm.Color(string('a'+i), tm.GREEN)
	}
}

var secretString = " I can only show you the door. You are the one that has to walk through it. "
var secretLetters []string

func init() {
	secretLetters = make([]string, len(secretString))
	for i := range secretLetters {
		secretLetters[i] = tm.Color(string(secretString[i]), tm.GREEN)
	}
}

var secretX int
var secretY = rand.Intn(screenHeight) + 1

const (
	CLEAR = 0
	PAINT = 1
)

func main() {
	if screenWidth < len(secretString) {
		fmt.Println("Sorry, your console window is too small to hold the matrix")
		return
	} else {
		secretX = rand.Intn(screenWidth-len(secretString)) + 1
	}

	fmt.Printf("\033[?25l") // hide cursor
	defer func() {
		fmt.Printf("\033[1;1H") // move cursor to 1,1
		fmt.Printf("\033[2J")   // clear the screen
		fmt.Printf("\033[?25h") // show cursor
	}()

	rand.Seed(time.Now().Unix())
	tm.Clear()

	for i := 0; i < screenWidth/7; i++ {
		newLetter()
	}

	// we want our deferred function to finish cleanup
	// so we catch terminating signals to let our program exit gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	ticker := time.NewTicker(20 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			m.Lock()
			tm.Flush()
			m.Unlock()
		case <-sig:
			ticker.Stop()
			return
		}
	}
}

func newLetter() {
	go func() {
		x, y, mode := newState()
		ticker := time.NewTicker(time.Duration(rand.Intn(50)+40) * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				m.Lock()
				tm.MoveCursor(x, y)
				if mode == PAINT {
					if secretPlace(x, y) {
						tm.Print(string(secretLetters[x-secretX]))
					} else {
						tm.Print(lettersOld[rand.Intn(len(lettersNew))])
					}
				}
				y++
				if y > screenHeight {
					x, y, mode = newState()
				}
				tm.MoveCursor(x, y)
				if mode == PAINT {
					tm.Print(lettersNew[rand.Intn(len(lettersNew))])
				} else if !secretPlace(x, y) {
					tm.Print(" ")
				}
				m.Unlock()
			}
		}
	}()
}

func secretPlace(x, y int) bool {
	return y == secretY && x >= secretX && x < secretX+len(secretString)
}

func newState() (x, y, mode int) {
	if rand.Intn(42) < 25 {
		mode = CLEAR
	} else {
		mode = PAINT
	}
	x = rand.Intn(screenWidth) + 1
	y = 1
	return
}
