/**
 * P2POmokClient.go
 * NAME : Seokwon Yoon
 * ID : 20174089
 **/

package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
    "strings"
    "log"
    "os/signal"
    "syscall"   
    "os/exec"
    "runtime"
    "regexp"
    "strconv"
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
    clearMap["darwin"] = func() {
        cmd := exec.Command("clear") //darwin example, its tested
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

func gameStart() {
    opponent_addr, _ := net.ResolveUDPAddr("udp", OpponentAddr)
    print("send to ")
    print(opponent_addr.String())

    // print board
    board = Board{}
    x, y, count, win, cnt := -1, -1, 0, 0, 0
    for i := 0; i < Row; i++ {
        var tempRow []int
        for j := 0; j < Col; j++ {
            tempRow = append(tempRow, 0)
        }
        board = append(board, tempRow)
    }
    printBoard(board)

    for {
        userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
        pconn.WriteTo([]byte(userInput), opponent_addr)
        if err != nil {
            log.Fatal(err)
        }
        if userInput == "\n" {
            continue
        }

        if userInput[0] == '\\' {
            if userInput[1] == '\\' {
                if (TurnFlag == turn) {
                    println("not your turn")
                    continue
                } 
                re := regexp.MustCompile("[0-9]+")
                result := re.FindAllString(userInput, -1)
                if (len(result) == 2) {
                    x, _ = strconv.Atoi(result[0])
                    y, _ = strconv.Atoi(result[1])
                    // println(result[0])
                    // println(result[1])
                    cnt = 2
                } else {
                    cnt = 0
                }
                
            } else {
                switch userInput {
                case "\\gg\n":
                    fmt.Println("gg~")
                    // conn.Write([]byte("7"))
                case "\\exit\n":
                    fmt.Println("gg~")
                    // conn.Write([]byte("5"))
                    // conn.Close()
                    // os.Exit(0)
                default:
                    fmt.Println("invalid command")
                }
            }
        } else {
            // send message to opponent
            pconn.WriteTo([]byte("3" + userInput), opponent_addr)
            continue
        }

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
        } else {
            pconn.WriteTo([]byte("9" + userInput), opponent_addr)
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
}

func readOppnentMessage() {
    buffer := make([]byte, 1024)
    for {
        count, _, _:= pconn.ReadFrom(buffer)
            
        switch string(buffer[0]) {
        case "3":
            fmt.Printf("%s> %s", OpponentNick, string(buffer[1:count]))         
        case "9":
            fmt.Printf("stone %s", string(buffer[1:count]))    
            re := regexp.MustCompile("[0-9]+")
            result := re.FindAllString(string(buffer[1:count]), -1)
            x, _ := strconv.Atoi(result[0])
            y, _ := strconv.Atoi(result[1])  

            if turn == 0 {
                board[x][y] = 1
            } else {
                board[x][y] = 2
            }
            clear()
            printBoard(board)
            win := checkWin(board, x, y)
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
        default:
            // fmt.Println("wrong code " + string(buffer[0]))
        }
    }
}
func waitForOpponent(client_conn net.Conn) {
    buffer := make([]byte, 1024)

    for {
        n, _ := client_conn.Read(buffer)

        switch string(buffer[0]) {
        case "5":
            // fmt.Printf("close server connection, opponent info : %s\n", string(buffer[1:]))
            s := strings.Split(string(buffer[1:n]), " ")
            OpponentNick = s[0]
            OpponentAddr = s[1]
            TurnFlag, _ = strconv.Atoi(s[2])
            // fmt.Printf("%s %s %d\n", OpponentNick, OpponentAddr, TurnFlag)
            client_conn.Close()
            return
        case "3":
            fmt.Printf("%s", string(buffer[1:]))
            client_conn.Write([]byte("3 ack"))
        default:
            // fmt.Println("wrong code " + string(buffer[0]))
        }
    }
}

// code 1 : nickname is duplicated
// code 3 : oridinary message
// code 5 : close connection 
// code 7 : \gg
// code 6 : \exit
// code 9 : \\ <x> <y>
var time_start = time.Now()
var OpponentNick = ""
var OpponentAddr = ""
var myAddr = ""
var TurnFlag = 0
var turn = 0
var pconn net.PacketConn
var board = Board{}


func main() {
    // make TCP Connection with server
    serverName := "127.0.0.1"
    // serverName := "nsl2.cau.ac.kr"
    serverPort := "54089"
    argsWithProg := os.Args[1] // client nickname
    buffer := make([]byte, 4096)

    conn, err:= net.Dial("tcp", serverName+":"+serverPort)
    if err != nil {
        log.Fatal(err)
    }
    pconn, _ = net.ListenPacket("udp", ":")
    localAddr := pconn.LocalAddr().(*net.UDPAddr)

    // send client's nickname and udp address to the server
    conn.Write([]byte(argsWithProg + " " + strconv.Itoa(localAddr.Port)))

    // wait for permission
    _, err = conn.Read(buffer)
    if string(buffer[0]) == "1" {
        fmt.Printf("that nickname is already used by another user. cannot connect\n")
        return
    }

    if string(buffer[0]) == "3" {
        // read welcome msg from the server and print
        fmt.Printf("%s", string(buffer[1:]))
    } 

    // ctrl + c handling
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("gg~")
        conn.Write([]byte("5"))
        conn.Close()
        os.Exit(0)
    }()

    waitForOpponent(conn)
    fmt.Println("\n")

    println("my udp address : ", localAddr.String())
    println("opponent udp address : ", OpponentAddr)
    
    go readOppnentMessage()
    gameStart()
}
 