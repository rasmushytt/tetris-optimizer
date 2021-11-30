package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

var mySquare [][]string

func main() {
	args := os.Args[1]
	if args != "" {
		file, err := os.Open(args)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		defer func() {
			if err = file.Close(); err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}()
		var myArray = ReadInput(file)
		Solve(myArray)
		PrintSolution()
	}
}

func ReadInput(file io.Reader) [][4][4]string {
	var tetrominoArray [][4][4]string
	var tetromino [4][4]string
	scanner := bufio.NewScanner(file)
	i, in, flag, alpha := 0, 0, true, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for scanner.Scan() {
		str := scanner.Text()
		if str == "" {
			if flag {
				flag = false
				continue
			} else {
				fmt.Println("ERROR")
				os.Exit(0)
			}
		}
		var arr [4]string
		if len(str) != 4 {
			fmt.Println("ERROR")
			os.Exit(0)
		}
		for ind := range arr {
			if rune(str[ind]) == '.' {
				arr[ind] = "."
			} else if rune(str[ind]) == '#' {
				arr[ind] = string(alpha[in])
			} else {
				fmt.Println("ERROR")
				os.Exit(0)
			}
		}
		tetromino[i] = arr
		i++
		if i == 4 {
			flag = true
			i = 0
			in++
			if !CheckTetromino(tetromino) {
				fmt.Println("ERROR")
				os.Exit(0)
			}
			tetromino = OptimizeTetromino(tetromino)
			tetrominoArray = append(tetrominoArray, tetromino)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return tetrominoArray
}

func OptimizeTetromino(tetromino [4][4]string) [4][4]string {
	//optimzes tetromino
	i := 0
	for {
		zeroes := 0
		for j := 0; j < 4; j++ {
			if tetromino[i][j] == "." {
				zeroes++
			}
		}
		if zeroes == 4 { //if row is all zeroes, shift by 1 row to top
			tetromino = ShiftVertical(tetromino)
			continue
		}
		break
	}
	for {
		zeroes := 0
		for j := 0; j < 4; j++ {
			if tetromino[j][i] == "." {
				zeroes++
			}
		}
		if zeroes == 4 { //if col is all zeroes, shift by 1 col to left
			tetromino = ShiftHorizontal(tetromino)
			continue
		}
		break
	}
	return tetromino
}

func ShiftVertical(tetromino [4][4]string) [4][4]string {
	//shifts tetromino row by 1
	temp := tetromino[0]
	tetromino[0] = tetromino[1]
	tetromino[1] = tetromino[2]
	tetromino[2] = tetromino[3]
	tetromino[3] = temp
	return tetromino
}
func ShiftHorizontal(tetromino [4][4]string) [4][4]string {
	//shifts tetromino col by 1
	tetromino = Transpose(tetromino)
	tetromino = ShiftVertical(tetromino)
	tetromino = Transpose(tetromino)
	return tetromino
}

func Transpose(slice [4][4]string) [4][4]string {
	//transpose tetromino
	xl := len(slice[0])
	yl := len(slice)
	var result [4][4]string
	for i := range result {
		result[i] = [4]string{}
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func CheckTetromino(tetromino [4][4]string) bool {
	c, d := 0, 0

	for a, elem := range tetromino {
		for b, elem2 := range elem {
			if elem2 != "." {
				d++
				if a+1 < 4 && tetromino[a+1][b] != "." {
					c++
				}
				if a-1 >= 0 && tetromino[a-1][b] != "." {
					c++
				}
				if b+1 < 4 && tetromino[a][b+1] != "." {
					c++
				}
				if b-1 >= 0 && tetromino[a][b-1] != "." {
					c++
				}
			}
		}
	}
	if d != 4 {
		return false
	}
	if c == 6 || c == 8 {
		return true
	}
	return false
}

func Solve(tetrominoes [][4][4]string) [][]string {
	//initial board starts with dimmension 4*4, if we can't place all tetrominoes
	//increase size by 1 and initialize board
	l := int(math.Ceil(math.Sqrt(float64(4 * len(tetrominoes)))))
	mySquare = InitSquare(l)
	for !BacktrackSolver(tetrominoes, 0) {
		l++
		mySquare = InitSquare(l)
	}
	return mySquare
}
func InitSquare(n int) [][]string {
	//initializes a square
	var Square [][]string
	var row []string
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			row = append(row, ".")
		}
		Square = append(Square, row)
		row = []string{}
	}
	return Square
}

func BacktrackSolver(tetrominoes [][4][4]string, n int) bool {
	if n == len(tetrominoes) { //base condition when all tetrominoes are placed, board is solved
		return true
	}
	for i := 0; i < len(mySquare); i++ {
		for j := 0; j < len(mySquare); j++ {
			if CheckInsert(i, j, tetrominoes[n]) { //check if we can place current tetrominoe on the board anywhere
				Insert(i, j, tetrominoes[n]) // if we can place it at this location, check if we can place another piece
				if BacktrackSolver(tetrominoes, n+1) {
					return true
				}
				Remove(i, j, tetrominoes[n]) //if the next piece can't be placed, backtrack
			}
		}
	} // if we can't place tetro anywhere, return false
	return false
}
func CheckInsert(i, j int, tetro [4][4]string) bool { //check if we can place piece at current location
	for a := 0; a < 4; a++ {
		for b := 0; b < 4; b++ {
			if tetro[a][b] != "." {
				if i+a == len(mySquare) || j+b == len(mySquare) || mySquare[i+a][j+b] != "." {
					return false
				}
			}
		}

	}
	return true
}

func Insert(i, j int, tetro [4][4]string) { // insert piece and when all 4 pieces "#" are placed, no need to place '.'
	a, b, c := 0, 0, 0
	for a < 4 {
		for b < 4 {
			if tetro[a][b] != "." {
				c++
				mySquare[i+a][j+b] = tetro[a][b]
				if c == 4 {
					break
				}
			}
			b++
		}
		b = 0
		a++
	}
}

func Remove(i, j int, tetro [4][4]string) { //remove piece at current location
	a, b, c := 0, 0, 0
	for a < 4 {
		for b < 4 {
			if tetro[a][b] != "." {
				if c == 4 {
					break
				}
				mySquare[i+a][j+b] = "."
			}
			b++
		}
		b = 0
		a++
	}
}

func PrintSolution() {
	for i := range mySquare {
		for j := range mySquare {
			fmt.Printf("%s  ", mySquare[i][j])
		}
		fmt.Printf("\n")
	}
}
