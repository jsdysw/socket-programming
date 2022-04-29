/**
 * EasyUDPClient.go
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
)

func main() {
    // serverName := "127.0.0.1"
    serverName := "nsl2.cau.ac.kr"
    serverPort := "34089"

    pconn, err:= net.ListenPacket("udp", ":")
    if err != nil {
        log.Fatal(err)
    }

    localAddr := pconn.LocalAddr().(*net.UDPAddr)
    fmt.Printf("Client is running on port %d\n", localAddr.Port)

    server_addr, err := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
    if err != nil {
        log.Fatal(err)
    }

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

        // print Menu at conosle
        fmt.Printf("<Menu>\n")
        fmt.Printf("1) Convert text to UPPER-case\n")
        fmt.Printf("2) get my IP address and port number\n")
        fmt.Printf("3) get server request count\n")
        fmt.Printf("4) get server running time\n")
        fmt.Printf("5) exit\n")
        
        // get user's choice
        fmt.Printf("Input option: ")
        user_choice, err := bufio.NewReader(os.Stdin).ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("\n")

        // add request type flag (1,2,3,4) at the front of the body
        switch user_choice {
        case "1\n":
            fmt.Printf("Input lowercase sentence: ")
            input, err := bufio.NewReader(os.Stdin).ReadString('\n')
            if err != nil {
                log.Fatal(err)
            }

            input = "1" + input

            time_start := time.Now()
            pconn.WriteTo([]byte(input), server_addr)

            pconn.ReadFrom(buffer)
            time_elapsed := time.Since(time_start)

            fmt.Printf("Reply from server: %s", string(buffer))
            fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)

        case "2\n":
            input := "2"

            time_start := time.Now()
            pconn.WriteTo([]byte(input), server_addr)

            pconn.ReadFrom(buffer)
            time_elapsed := time.Since(time_start)

            ip_port := strings.Split(string(buffer), ":")

            fmt.Printf("Reply from server: client IP = %s, port = %s\n", ip_port[0], ip_port[1])
            fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
        
        case "3\n":
            input := "3"

            time_start := time.Now()
            pconn.WriteTo([]byte(input), server_addr)

            pconn.ReadFrom(buffer)
            time_elapsed := time.Since(time_start)

            fmt.Printf("Reply from server: requests served = %s\n", string(buffer))
            fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)

        case "4\n":
            input := "4"
            
            time_start := time.Now()
            pconn.WriteTo([]byte(input), server_addr)

            pconn.ReadFrom(buffer)
            time_elapsed := time.Since(time_start)

            fmt.Printf("Reply from server: run time = %s\n", string(buffer))
            fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)

        case "5\n":
            input := "5"
            pconn.WriteTo([]byte(input), server_addr)

            fmt.Printf("Bye bye~\n")
            pconn.Close()
            return;
        default :
            fmt.Printf("Wrong option\n")
        }
    }
  
}