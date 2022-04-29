/**
 * EasyUDPServer.go
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

    serverPort := "34089"

    pconn, err:= net.ListenPacket("udp", ":"+serverPort)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server is ready to receive on port %s\n", serverPort)

    // ctrl + c handling
    c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
        <-c
        fmt.Println("Bye bye~")
        pconn.Close()
		os.Exit(0)
	}()
    
    for {
        buffer := make([]byte, 1024)

        count, r_addr, err:= pconn.ReadFrom(buffer)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Connection request from %s\n", r_addr.String())
        fmt.Printf("Command %s\n",string(buffer[0]))

        // first character of the message are used as header
        switch string(buffer[0]) {
        case "1":
            total_served_commands += 1
            pconn.WriteTo(bytes.ToUpper(buffer[1:count]), r_addr)      
        case "2":
            total_served_commands += 1
            pconn.WriteTo([]byte(r_addr.String()), r_addr)   
        case "3":
            total_served_commands += 1
            served_count_string := strconv.Itoa(total_served_commands)
            pconn.WriteTo([]byte(served_count_string), r_addr) 
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

            // fmt.Println(result)
            pconn.WriteTo([]byte(result), r_addr)
        case "5":
            // client has left, do nothing
            continue
        default :
            fmt.Printf("Wrong option\n")
        }

    }
}

