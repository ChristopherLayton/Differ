package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cell struct {
	Value   int
	Derived string
	Origin  int
	Trans   int
}

type Diff struct {
	Index int
	Op    string
}

func createGrid(rows, columns int) [][]Cell {
	grid := make([][]Cell, rows)
	for i := 0; i < rows; i++ {
		newRow := make([]Cell, columns)
		grid[i] = newRow
	}
	return grid
}

func LevenshteinDiff(str1, str2 string) []Diff {


	grid  := createGrid(len(str2) + 1,len(str1) + 1)

	for i := 0; i < len(grid[0]); i++ {
		grid[0][i] = Cell{Value: i, Derived: "N"}
	}

	for i := 0; i < len(grid); i++ {
		grid[i][0] = Cell{Value: i, Derived: "N"}
	}
	grid   = fillOutGrid(str1, str2, grid)
    printGrid(grid)
	diffs := traverseGrid(str1, str2, grid)
	return diffs

}

func main() {
	str1 := os.Args[1]
	str2 := os.Args[2]
	diffs := LevenshteinDiff(str1, str2)
	str1Arr := strings.Split(str1, "")
	str2Arr := strings.Split(str2, "")
	for _, diff := range diffs {
		if diff.Op == "+" {
			//Color string2
			str2Arr[diff.Index] = "\x1b[32m" + str2Arr[diff.Index] + "\x1b[0m"
			continue
		}
		if diff.Op == "-" {
			//Color string1
			str1Arr[diff.Index] = "\x1b[31m" + str1Arr[diff.Index] + "\x1b[0m"
			continue
		}
		if diff.Op == "/" {
			//Color string1
			str1Arr[diff.Index] = "\x1b[7m" + str1Arr[diff.Index] + "\x1b[0m"
			continue
		}
	}

	str1 = strings.Join(str1Arr, "")
	str2 = strings.Join(str2Arr, "")
	fmt.Printf(str1 + "\n")
	fmt.Printf(str2 + "\n")
}

func traverseGrid(str1 string, str2 string, grid [][]Cell) []Diff {
	var strDiffs []Diff
	startCell := grid[len(str2)][len(str1)]
	num := startCell.Value
	var recurse func(row int, col int, cCell Cell, grid [][]Cell) string
	recurse = func(row int, col int, cCell Cell, grid [][]Cell) string {
		if num <= 0 {
			return ""
		}
		dec := false
		var nextCell Cell
		if cCell.Derived == "U" {
			nextCell = grid[row-1][col]
			if cCell.Value != nextCell.Value {
				strDiffs = append(strDiffs, Diff{Op: "+", Index: cCell.Trans})
				dec = true
			}
			row = row - 1
		} else if cCell.Derived == "D" {
			nextCell = grid[row-1][col-1]
			if cCell.Value != nextCell.Value {
				strDiffs = append(strDiffs, Diff{Op: "/", Index: cCell.Origin})
				dec = true
			}
			row = row - 1
			col = col - 1
		} else if cCell.Derived == "L" {
			nextCell = grid[row][col-1]
			if cCell.Value != nextCell.Value {
				strDiffs = append(strDiffs, Diff{Op: "-", Index: cCell.Origin})
				dec = true
			}
			col = col - 1
		}

		if dec {
			num = num - 1
		}
		return recurse(row, col, nextCell, grid)
	}
	recurse(len(str2), len(str1), startCell, grid)
	return strDiffs

}

func fillOutGrid(original string, transformed string, grid [][]Cell) [][]Cell {
	for i := 1; i < len(grid[0]); i++ {
		for j := 1; j < len(grid); j++ {
			cost := 1
			if original[i-1] == transformed[j-1] {
				cost = 0
			}
			derive, val := minimum(grid[j-1][i].Value+1, grid[j-1][i-1].Value+cost, grid[j][i-1].Value+1)
			grid[j][i] = Cell{Value: val, Derived: derive, Origin: i - 1, Trans: j - 1}
		}
	}
	return grid
}

func minimum(up int, diag int, left int) (string, int) {

	if up <= left && up <= diag {
		return "U", up
	}

	if left <= up && left <= diag {
		return "L", left
	}

	if diag <= up && diag <= left {
		return "D", diag
	}

	return "", 0
}

func printGrid(grid [][]Cell) {
	maxVal := grid[len(grid)-1][len(grid[0])-1].Value
	var output string
	for _, val := range grid {
		output = output + "|"
		for _, cell := range val {
			var cel string
			if cell.Value < 10 && maxVal >= 10 {
				cel = " "
			}
			cel = cel + strconv.Itoa(cell.Value) + cell.Derived + "|"
			output = output + cel
		}
		output = output + "\n"
	}
	fmt.Println(output)
}
