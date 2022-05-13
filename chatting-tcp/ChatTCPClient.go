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
    "strings"
    "log"
    "os/signal"
    "syscall"
)
 
func main() {
    // make TCP Connection with server
    serverName := "127.0.0.1"
    // serverName := "nsl2.cau.ac.kr"
    serverPort := "44089"
    argsWithProg := os.Args[1]
 
    conn, err:= net.Dial("tcp", serverName+":"+serverPort)
    if err != nil {
        log.Fatal(err)
    }
     
    // send client's nickname to the server
    conn.Write([]byte(argsWithProg))
    // fmt.Printf("my nickname %s\n", argsWithProg)
    localAddr := conn.LocalAddr().(*net.TCPAddr)
    // fmt.Printf("Client is running on port %d\n", localAddr.Port)
    
    // read welcome msg from the server and print
    buffer := make([]byte, 1024)
    conn.Read(buffer)
    fmt.Printf("%s\n", string(buffer))
             

    // ctrl + c handling
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("Bye bye~")
        input := "5"
        conn.Write([]byte(input))
        conn.Close()
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
 
         // add request type flag(header) (1,2,3,4) at the front of the body
         switch user_choice {
         case "1\n":
             fmt.Printf("Input lowercase sentence: ")
             input, err := bufio.NewReader(os.Stdin).ReadString('\n')
             if err != nil {
                 log.Fatal(err)
             }
             input = "1" + input
 
             time_start := time.Now()
             conn.Write([]byte(input))
 
             conn.Read(buffer)
             time_elapsed := time.Since(time_start)
 
             fmt.Printf("Reply from server: %s", string(buffer))
             fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
 
         case "2\n":
             input := "2"
 
             time_start := time.Now()
             conn.Write([]byte(input))
 
             conn.Read(buffer)
             time_elapsed := time.Since(time_start)
 
             ip_port := strings.Split(string(buffer), ":")
 
             fmt.Printf("Reply from server: client IP = %s, port = %s\n", ip_port[0], ip_port[1])
             fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
 
         case "3\n":
             input := "3"
 
             time_start := time.Now()
             conn.Write([]byte(input))
 
             conn.Read(buffer)
             time_elapsed := time.Since(time_start)
 
             fmt.Printf("Reply from server: requests served = %s\n", string(buffer))
             fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)    
 
         case "4\n":
             input := "4"
 
             time_start := time.Now()
             conn.Write([]byte(input))
 
             conn.Read(buffer)
             time_elapsed := time.Since(time_start)
 
             fmt.Printf("Reply from server: run time = %s\n", string(buffer))
             fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
         
         case "5\n":
             input := "5"
             conn.Write([]byte(input))
 
             fmt.Printf("Bye bye~")
             conn.Close()
             return;
         default:
             fmt.Printf("Wrong option\n")
         }
     }
 
 }
 