/**
 * ChatTCPServer.go
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
    "sync/atomic"
)

func handleConnection(client_conn net.Conn, id int) {
    for {
        buffer := make([]byte, 1024)

        count, err := client_conn.Read(buffer)
        if err != nil {
            // fmt.Printf("client socket failed to read data from buffer\n")
            log.Fatal(err)
        }

        // check the header of the message
        fmt.Printf("Command %s\n",string(buffer[0]))
        switch string(buffer[0]) {
        case "1": 
            atomic.AddInt64(&Total_served_commands, 1)  // mutul exclusion
            client_conn.Write(bytes.ToUpper(buffer[1:count]))
        case "2":
            atomic.AddInt64(&Total_served_commands, 1)
            client_conn.Write([]byte(client_conn.RemoteAddr().String()))   
        case "3":
            atomic.AddInt64(&Total_served_commands, 1)
            served_count_string := strconv.FormatInt(atomic.LoadInt64(&Total_served_commands), 10)
            client_conn.Write([]byte(served_count_string)) 
        case "4":
            atomic.AddInt64(&Total_served_commands, 1)

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
            atomic.AddInt64(&TotalClient, -1)
            client_conn.Close()
            fmt.Printf("Client %d disconnected. Number of connected clients = %d\n", id, atomic.LoadInt64(&TotalClient))
            return
        default :
            fmt.Printf("Wrong option\n")
        }
    }
}


var Total_served_commands int64  // atomic variable
var LastestId = 0
var TotalClient int64  // atomic variable
var Time_start = time.Now()
var Listener net.Conn

func main() {
    serverPort := "44089"
    Total_served_commands = 0
    TotalClient = 0
 
    // create server socket
    Listener, err:= net.Listen("tcp", ":" + serverPort)
    if err != nil {
        // fmt.Printf("create server socket failed\n")
        log.Fatal(err)
    }
    fmt.Printf("Server is ready to receive on port %s\n", serverPort)

    // ctrl + c handling
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("Bye bye~")
        Listener.Close()
        close(c)
        os.Exit(0)
    }()

    // print the number of clients every one minute
    ticker := time.NewTicker(60 * time.Second)
    go func() {
        for {
            <-ticker.C
            fmt.Println("The number of conneted clients : ", atomic.LoadInt64(&TotalClient))
        }
    }()

    for {
        // create client socket
        conn, err:= Listener.Accept()
        if err != nil {
            // fmt.Printf("create client socket failed\n")
            log.Fatal(err)
        } else {
            atomic.AddInt64(&TotalClient, 1); // for mutual exclusion
            LastestId = LastestId + 1 
            fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
            fmt.Printf("Client %d connected. Number of connected clients = %d\n", LastestId, atomic.LoadInt64(&TotalClient))
        }
        
        // client socket do its work
        go handleConnection(conn, LastestId)
    }
}
 
 