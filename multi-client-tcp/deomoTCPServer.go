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

type Client struct {
    Conn net.Conn
    Count  int
}
 
func main() {
    time_start := time.Now()

    serverPort := "44089"
 
    listener, err:= net.Listen("tcp", ":" + serverPort)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server is ready to receive on port %s\n", serverPort)
    
    database := make(map[string]*Client) 

    for {
        // make tcp connection with client
        c1 := make(chan string, 1)
        go func() {
            conn, err:= listener.Accept()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
            // server remembers client info as map
            _, exists := database[conn.RemoteAddr().String()]
            if !exists {
                database[conn.RemoteAddr().String()] = &Client{conn, 0}
                // val = database[conn.RemoteAddr().String()]
            }
            c1 <- "result 1"
        }()

        // ctrl + c handling
        c := make(chan os.Signal)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        go func() {
            <-c
            for key, element := range database {
                delete(database, key);
                element.Conn.Close()
            }
            listener.Close()
            fmt.Println("Bye bye~")
            os.Exit(0)
            return
        }()

        // wait new client connection for 1 second 
        select {
        case res := <-c1:
            fmt.Println(res)
        case <-time.After(time.Second * 1):
            // fmt.Println("timeout 1")
        }
        
        for key, element := range database {
            buffer := make([]byte, 1024)
            
            count, err := element.Conn.Read(buffer)
            if err != nil {
                log.Fatal(err)
            }

            // check the header of the message
            fmt.Printf("Command %s\n",string(buffer[0]))

            switch string(buffer[0]) {
            case "1": 
                element.Count += 1
                element.Conn.Write(bytes.ToUpper(buffer[1:count]))
            case "2":
                element.Count += 1
                element.Conn.Write([]byte(key))   
            case "3":
                element.Count += 1
                served_count_string := strconv.Itoa(element.Count)
                element.Conn.Write([]byte(served_count_string)) 
            case "4":
                element.Count += 1
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
    
                element.Conn.Write([]byte(result))
            case "5":
                // delete client info at the database
                delete(database, key);
                element.Conn.Close()
            default :
                fmt.Printf("Wrong option\n")
            }

        }
    }
}
 
 