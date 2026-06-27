package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"sudoku.reyel.cloud/sudoku"
)

type puzzleFile struct {
	Board [][]int `json:"board"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go <command> [args]")
		fmt.Fprintln(os.Stderr, "commands: generate, read")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		generateCmd()
	case "read":
		readCmd()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func generateCmd() {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	n := cmd.Int("n", 30, "number of given squares (17-81)")
	cmd.Parse(os.Args[2:])

	board := sudoku.Generate(*n)

	grid := make([][]int, 9)
	for r := range 9 {
		grid[r] = make([]int, 9)
		for c := range 9 {
			grid[r][c] = board[r*9+c]
		}
	}

	data := map[string]any{
		"board": grid,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	filename := fmt.Sprintf("puzzle-%d.json", time.Now().Unix())
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", filename)
}

func readCmd() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go read <puzzle-file>")
		os.Exit(1)
	}

	path := os.Args[2]

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	var p puzzleFile
	if err := json.Unmarshal(data, &p); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if len(p.Board) != 9 {
		fmt.Fprintln(os.Stderr, "invalid board: must be 9x9")
		os.Exit(1)
	}

	for r := range 9 {
		if len(p.Board[r]) != 9 {
			fmt.Fprintln(os.Stderr, "invalid board: must be 9x9")
			os.Exit(1)
		}

		for c := range 9 {
			if c%3 == 0 && c != 0 {
				fmt.Print("| ")
			}
			if p.Board[r][c] == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", p.Board[r][c])
			}
		}

		if r%3 == 2 && r != 8 {
			fmt.Println()
			fmt.Println("------+-------+------")
		} else {
			fmt.Println()
		}
	}
}
