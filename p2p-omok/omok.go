package main

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "time"
)

const (
    Row = 10
    Col = 10
)

type Board [][]int

func printBoard(b Board) {
    fmt.Print("   ")
    for j := 0; j < Col; j++ {
        fmt.Printf("%2d", j)
    }

    fmt.Println()
    fmt.Print("  ")
    for j := 0; j < 2*Col+3; j++ {
        fmt.Print("-")
    }

    fmt.Println()

    for i := 0; i < Row; i++ {
        fmt.Printf("%d |", i)
        for j := 0; j < Col; j++ {
            c := b[i][j]
            if c == 0 {
                fmt.Print(" +")
            } else if c == 1 {
                fmt.Print(" 0")
            } else if c == 2 {
                fmt.Print(" @")
            } else {
                fmt.Print(" |")
            }
        }

        fmt.Println(" |")
    }

    fmt.Print("  ")
    for j := 0; j < 2*Col+3; j++ {
        fmt.Print("-")
    }

    fmt.Println()
}

func checkWin(b Board, x, y int) int {
    lastStone := b[x][y]
    startX, startY, endX, endY := x, y, x, y

    // Check X
    for startX-1 >= 0 && b[startX-1][y] == lastStone {
        startX--
    }
    for endX+1 < Row && b[endX+1][y] == lastStone {
        endX++
    }

    if endX-startX+1 >= 5 {
        return lastStone
    }

    // Check Y
    startX, startY, endX, endY = x, y, x, y
    for startY-1 >= 0 && b[x][startY-1] == lastStone {
        startY--
    }
    for endY+1 < Row && b[x][endY+1] == lastStone {
        endY++
    }

    if endY-startY+1 >= 5 {
        return lastStone
    }

    // Check Diag 1
    startX, startY, endX, endY = x, y, x, y
    for startX-1 >= 0 && startY-1 >= 0 && b[startX-1][startY-1] == lastStone {
        startX--
        startY--
    }
    for endX+1 < Row && endY+1 < Col && b[endX+1][endY+1] == lastStone {
        endX++
        endY++
    }

    if endY-startY+1 >= 5 {
        return lastStone
    }

    // Check Diag 2
    startX, startY, endX, endY = x, y, x, y
    for startX-1 >= 0 && endY+1 < Col && b[startX-1][endY+1] == lastStone {
        startX--
        endY++
    }
    for endX+1 < Row && startY-1 >= 0 && b[endX+1][startY-1] == lastStone {
        endX++
        startY--
    }

    if endY-startY+1 >= 5 {
        return lastStone
    }

    return 0
}

func clear() {
    fmt.Printf("%s", runtime.GOOS)

    clearMap := make(map[string]func()) //Initialize it
    clearMap["linux"] = func() {
        cmd := exec.Command("clear") //Linux example, its tested
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clearMap["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
        cmd.Stdout = os.Stdout
        cmd.Run()
    }

    value, ok := clearMap[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
    if ok {                             //if we defined a clearMap func for that platform:
        value() //we execute it
    } else { //unsupported platform
        panic("Your platform is unsupported! I can't clearMap terminal screen :(")
    }
}

func main() {
    board := Board{}
    x, y, turn, count, win := -1, -1, 0, 0, 0
    for i := 0; i < Row; i++ {
        var tempRow []int
        for j := 0; j < Col; j++ {
            tempRow = append(tempRow, 0)
        }
        board = append(board, tempRow)
    }

    clear()
    printBoard(board)

    for {
        print("please enter \"x y\" coordinate >> ")
        cnt, _ := fmt.Scanf("%d %d ", &x, &y)

        if cnt != 2 {
            fmt.Println("error, must enter x y!")
            time.Sleep(1 * time.Second)
            continue
        } else if x < 0 || y < 0 || x >= Row || y >= Col {
            fmt.Println("error, out of bound!")
            time.Sleep(1 * time.Second)
            continue
        } else if board[x][y] != 0 {
            fmt.Println("error, already used!")
            time.Sleep(1 * time.Second)
            continue
        }

        if turn == 0 {
            board[x][y] = 1
        } else {
            board[x][y] = 2
        }

        clear()
        printBoard(board)

        win = checkWin(board, x, y)
        if win != 0 {
            fmt.Printf("player %d wins!\n", win)
            break
        }

        count += 1
        if count == Row*Col {
            fmt.Printf("draw!\n")
            break
        }

        turn = (turn + 1) % 2
    }

    return
}

