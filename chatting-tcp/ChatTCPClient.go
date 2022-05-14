/**
 * ChatTCPClient.go
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
    // "strings"
    "log"
    "os/signal"
    "syscall"
)

func updateMessages(client_conn net.Conn) {
    for {
        buffer := make([]byte, 1024)
        client_conn.Read(buffer)

        switch string(buffer[0]) {
        case "5":
            fmt.Println("gg~")
            client_conn.Close()
            os.Exit(0)
        case "2":
            fmt.Println("\"I hate professor\" is not not allowed\ngg~")
            client_conn.Close()
            os.Exit(0)
        case "3":
            fmt.Printf("%s", string(buffer[1:]))
        case "9":
            time_elapsed := time.Since(time_start)
            fmt.Printf("RTT = %.3f ms\n", float64(time_elapsed.Nanoseconds())/1000000.0)
        default:
            // fmt.Println("wrong code " + string(buffer[0]))
        }
    }
}

// code 0 : chatting room is full
// code 1 : nickname is duplicated
// code 2 : "i hate professor" detected
// code 3 : cahtting message
// code 5 : server has been quitted
// code 7 : \list
// code 6 : \dm
// code 5 : \exit
// code 8 : \ver
// code 9 : \rtt
var time_start = time.Now()

func main() {
    // make TCP Connection with server
    // serverName := "127.0.0.1"
    serverName := "nsl2.cau.ac.kr"
    serverPort := "44089"
    argsWithProg := os.Args[1] // client nickname
    buffer := make([]byte, 4096)


    conn, err:= net.Dial("tcp", serverName+":"+serverPort)
    if err != nil {
        log.Fatal(err)
    }

    // send client's nickname to the server
    conn.Write([]byte(argsWithProg))
    // fmt.Printf("my nickname %s\n", argsWithProg)
    // localAddr := conn.LocalAddr().(*net.TCPAddr)    
    // fmt.Printf("Client is running on port %d\n", localAddr.Port)

    // wait chat room attend permission
    _, err = conn.Read(buffer)
    if string(buffer[0]) == "0" {
        fmt.Printf("chatting room full. cannot connect\n")
        return
    } else if string(buffer[0]) == "1" {
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
 
    go updateMessages(conn)

    for {
        // buffer := make([]byte, 1024)
        message, err := bufio.NewReader(os.Stdin).ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        if message == "\n" {
            continue
        }
        // \command cases
        // code 7 : \list
        // code 6 : \dm
        // code 5 : \exit
        // code 8 : \ver
        // code 9 : \rtt
        if len(message) >= 3 && message[0:3] == "\\dm" {
            conn.Write([]byte("6" + message))
            // fmt.Println("6" + message)
        } else if message[0] == '\\' {
            switch message {
            case "\\list\n":
                conn.Write([]byte("7"))
            case "\\exit\n":
                fmt.Println("gg~")
                conn.Write([]byte("5"))
                conn.Close()
                os.Exit(0)
            case "\\ver\n":
                conn.Write([]byte("8"))
            case "\\rtt\n":
                time_start = time.Now()
                conn.Write([]byte("9"))
            default:
                fmt.Println("invalid command")
            }
        } else {
            conn.Write([]byte("3" + message))
        }
    }
 
 }
 