/**
 * MultiClientTCPServer.go
 * NAME : Seokwon Yoon
 * ID : 20174089
 **/

package main

import (
    "bytes"
    "fmt"
    "net"
    "time"
    "strconv"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func handleConnection(client_conn net.Conn, id int) {
    // if main go routine ends, close all the client conns
    // c := make(chan os.Signal)
    // signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    // go func() {
    //     <-c
    //     Listener.Close()
    //     // 열려있는 애들 다 닫자.
    //     fmt.Println("--------client closed")
    //     os.Exit(0)
    // }()

    for {
        buffer := make([]byte, 1024)

        count, err := client_conn.Read(buffer)
        if err != nil {
            fmt.Printf("client socket failed to read data from buffer\n")
            log.Fatal(err)
        }
        // check the header of the message
        fmt.Printf("Command %s\n",string(buffer[0]))
        switch string(buffer[0]) {
        case "1": 
            Total_served_commands += 1
            client_conn.Write(bytes.ToUpper(buffer[1:count]))
        case "2":
            Total_served_commands += 1
            client_conn.Write([]byte(client_conn.RemoteAddr().String()))   
        case "3":
            Total_served_commands += 1
            served_count_string := strconv.Itoa(Total_served_commands)
            client_conn.Write([]byte(served_count_string)) 
        case "4":
            Total_served_commands += 1

            time_elapsed := time.Since(Time_start)

            // result should be HH:MM:SS format
            hh_s := strconv.Itoa(int(time_elapsed.Hours()))
            if len(hh_s) == 1 {
                hh_s = "0" + hh_s
            }
            mm_s := strconv.Itoa(int(time_elapsed.Minutes())% 60)
            if len(mm_s) == 1 {
                mm_s = "0" + mm_s
            }
            ss_s := strconv.Itoa(int(time_elapsed.Seconds())% 60)
            if len(ss_s) == 1 {
                ss_s = "0" + ss_s
            }
            result := hh_s + ":" + mm_s + ":" + ss_s

            client_conn.Write([]byte(result))
        case "5":
            // close client connection
            TotalClient = TotalClient - 1;
            client_conn.Close()
            fmt.Printf("Client %d disconnected. Number of connected clients = %d\n", id, TotalClient)
            return
        default :
            fmt.Printf("Wrong option\n")
        }
    }
    // close client conn 
}

var Total_served_commands = 0
var LastestId = 0
var TotalClient = 0
var Time_start = time.Now()
var Listener net.Conn

func main() {
    serverPort := "44089"
 
    // create server socket
    Listener, err:= net.Listen("tcp", ":" + serverPort)
    if err != nil {
        fmt.Printf("create server socket failed\n")
        log.Fatal(err)
    }
    fmt.Printf("Server is ready to receive on port %s\n", serverPort)

    // ctrl + c handling
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        Listener.Close()
        fmt.Println("Bye bye~")
        os.Exit(0)
    }()

    for {
        // create client socket
        conn, err:= Listener.Accept()
        if err != nil {
            // fmt.Printf("create client socket failed\n")
            log.Fatal(err)
        } else {
            TotalClient = TotalClient + 1;
            LastestId = LastestId + 1;
            fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
            fmt.Printf("Client %d connected. Number of connected clients = %d\n", LastestId, TotalClient)
        }
        
        // client socket do works
        go handleConnection(conn, LastestId)
    }
}
 
 