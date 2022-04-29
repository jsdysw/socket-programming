/**
 * EasyTCPServer.go
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

 
func main() {
    time_start := time.Now()

    total_served_commands := 0
    serverPort := "44089"
 
    listener, err:= net.Listen("tcp", ":" + serverPort)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server is ready to receive on port %s\n", serverPort)

    var client_conn net.Conn

    // ctrl + c handling
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        listener.Close()
        fmt.Println("Bye bye~")
        os.Exit(0)
        return
    }()

    conn, err:= listener.Accept()
    if err != nil {
        log.Fatal(err)
    } else {
        client_conn = conn
        fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
    }
    
    for {
        // ctrl + c handling
        c := make(chan os.Signal)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        go func() {
            <-c
            client_conn.Close()
            listener.Close()
            fmt.Println("Bye bye~")
            os.Exit(0)
            return
        }()
        
        buffer := make([]byte, 1024)
            
        count, err := client_conn.Read(buffer)
        if err != nil {
            log.Fatal(err)
        }
        // check the header of the message
        fmt.Printf("Command %s\n",string(buffer[0]))
        switch string(buffer[0]) {
        case "1": 
            total_served_commands += 1
            client_conn.Write(bytes.ToUpper(buffer[1:count]))
        case "2":
            total_served_commands += 1
            client_conn.Write([]byte(client_conn.RemoteAddr().String()))   
        case "3":
            total_served_commands += 1
            served_count_string := strconv.Itoa(total_served_commands)
            client_conn.Write([]byte(served_count_string)) 
        case "4":
            total_served_commands += 1
            time_elapsed := time.Since(time_start)

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
            // close client conn and wait anotherå
            client_conn.Close()
            // fmt.Printf("waiting for another client\n")
            conn, err:= listener.Accept()
            if err != nil {
                log.Fatal(err)
            } else {
                client_conn = conn
                fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
            }
        default :
            fmt.Printf("Wrong option\n")
        }
    }
}
 
 