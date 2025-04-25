// #nosec -- This is a helper command. Should only be used for development.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
)

var input, output string
var usage bool

func init() {
	flag.StringVar(&input, "i", "", "Input file")
	flag.StringVar(&output, "o", "", "Output file")
	flag.BoolVar(&usage, "h", false, "Show helptext")
}

func main() {
	flag.Parse()
	if usage {
		flag.Usage()
		os.Exit(0)
	}
	if input == "" {
		fmt.Println("You need to specify an input file with -i <file>")
		os.Exit(1)
	}
	if output == "" {
		fmt.Println("You need to specify an output file with -o <file>")
		os.Exit(1)
	}
	f, err := os.ReadFile(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reg := regexp.MustCompile(`time=.* level=INFO msg="Checking field" pos="\((\d*), (\d*)\)"\n`)

	positions := reg.FindAllStringSubmatch(string(f), -1)
	p := make([]minesweeper.Pos, len(positions))
	for i, pos := range positions {
		x, err := strconv.Atoi(pos[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		y, err := strconv.Atoi(pos[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		p[i] = minesweeper.NewPos(x, y)
	}
	fmt.Printf("Found %d positions in log\n", len(p))

	buf, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(output, buf, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Success!!!")
}
