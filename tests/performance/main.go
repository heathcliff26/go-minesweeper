package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
)

const workerThreads = 10

func main() {
	difficulty, _ := minesweeper.NewCustomDifficulty(700, 99, 99)
	measureLoop("NewGameWithSafePos", 100, func(res chan time.Duration) {
		measure(res, func() {
			minesweeper.NewGameWithSafePos(difficulty, minesweeper.NewPos(50, 50))
		})
	})
	measureLoop("NewGameWithSafeArea", 100, func(res chan time.Duration) {
		measure(res, func() {
			minesweeper.NewGameWithSafeArea(difficulty, minesweeper.NewPos(50, 50))
		})
	})
	measureLoop("AssistedMode", 100, measureAssistedMode)
}

// Run the measurement multiple times, running in parallel as set by workerThreads
func measureLoop(name string, iterations int, f func(res chan time.Duration)) string {
	fmt.Printf("Measuring %s using %d iterations\n", name, iterations)

	rest := iterations % workerThreads

	measurements := make([]time.Duration, iterations)
	res := make(chan time.Duration, workerThreads)
	i := 0

	for x := 0; x < iterations/workerThreads; x++ {
		for y := 0; y < workerThreads; y++ {
			go f(res)
		}
		for y := 0; y < workerThreads; y++ {
			measurements[i] = <-res
			i++
		}
	}
	for y := 0; y < rest; y++ {
		go f(res)
	}
	for y := 0; y < rest; y++ {
		measurements[i] = <-res
		i++
	}

	fmt.Printf("Finished measuring %s, calculating result\n\n", name)
	max, min, average := calculateStats(measurements)
	output := fmt.Sprintf("Name: %s\nMax: %s\nMin: %s\nAverage: %s\nMeasurements:\n", name, dString(max), dString(min), dString(average))
	for i := 0; i < iterations; i++ {
		output += dString(measurements[i])
		if (i+1)%workerThreads == 0 {
			output += "\n"
		} else {
			output += " "
		}
	}
	if !(iterations%workerThreads == 0) {
		output += "\n"
	}
	fmt.Print(output + "\n")
	return output
}

func measure(res chan time.Duration, f func()) {
	start := time.Now()
	f()
	res <- time.Since(start)
}

func calculateStats(measurements []time.Duration) (max, min, average time.Duration) {
	var sum time.Duration = 0

	if len(measurements) == 0 {
		return
	}
	min = measurements[0]
	for _, m := range measurements {
		sum += m
		if m > max {
			max = m
		}
		if m < min {
			min = m
		}
	}
	average = sum / time.Duration(len(measurements))
	return
}

func dString(d time.Duration) string {
	return strconv.FormatInt(d.Microseconds(), 10) + "µs"
}

func measureAssistedMode(res chan time.Duration) {
	save, err := minesweeper.LoadSave("testdata/game.sav")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	game := save.Game()

	f, err := os.ReadFile("testdata/positions.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var positions []minesweeper.Pos
	err = json.Unmarshal(f, &positions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	start := time.Now()
	for _, pos := range positions {
		s := game.CheckField(pos)
		s.ObviousMines()
		s.ObviousSafePos()
		if s.GameOver || s.GameWon {
			fmt.Println("Something unexpected happened, the game finished")
		}
	}
	res <- time.Since(start)
}
