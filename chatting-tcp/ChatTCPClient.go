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
    // "time"
    // "strings"
    "log"
    "os/signal"
    "syscall"
)

func updateMessages(client_conn net.Conn) {
    for {
        buffer := make([]byte, 1024)
        client_conn.Read(buffer)
        fmt.Printf("%s", string(buffer))
    }
}

 
func main() {
    // make TCP Connection with server
    serverName := "127.0.0.1"
    // serverName := "nsl2.cau.ac.kr"
    serverPort := "44089"
    argsWithProg := os.Args[1] // client nickname
    buffer := make([]byte, 4096)


    conn, err:= net.Dial("tcp", serverName+":"+serverPort)
    if err != nil {
        log.Fatal(err)
    }

    // send client's nickname to the server
    conn.Write([]byte(argsWithProg))
    fmt.Printf("my nickname %s\n", argsWithProg)
    // localAddr := conn.LocalAddr().(*net.TCPAddr)    
    // fmt.Printf("Client is running on port %d\n", localAddr.Port)
    // wait chat room attend permission
    n, err := conn.Read(buffer)
    if string(buffer[0]) == "0" {
        fmt.Printf("chatting room full. cannot connect\n")
        return
    } else if string(buffer[0]) == "1" {
        fmt.Printf("that nickname is already used by another user. cannot connect\n")
        return
    }

    // read welcome msg from the server and print
    fmt.Printf("%s\n", string(buffer[:n]))
             
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
 
    go updateMessages(conn)

    for {
        // buffer := make([]byte, 1024)
        message, err := bufio.NewReader(os.Stdin).ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        conn.Write([]byte(message))
 
         // add request type flag(header) (1,2,3,4) at the front of the body
    //      switch user_choice {
    //      case "1\n":
    //          fmt.Printf("Input lowercase sentence: ")
    //          input, err := bufio.NewReader(os.Stdin).ReadString('\n')
    //          if err != nil {
    //              log.Fatal(err)
    //          }
    //          input = "1" + input
 
    //          time_start := time.Now()
    //          conn.Write([]byte(input))
 
    //          conn.Read(buffer)
    //          time_elapsed := time.Since(time_start)
 
    //          fmt.Printf("Reply from server: %s", string(buffer))
    //          fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
 
    //      case "2\n":
    //          input := "2"
 
    //          time_start := time.Now()
    //          conn.Write([]byte(input))
 
    //          conn.Read(buffer)
    //          time_elapsed := time.Since(time_start)
 
    //          ip_port := strings.Split(string(buffer), ":")
 
    //          fmt.Printf("Reply from server: client IP = %s, port = %s\n", ip_port[0], ip_port[1])
    //          fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
 
    //      case "3\n":
    //          input := "3"
 
    //          time_start := time.Now()
    //          conn.Write([]byte(input))
 
    //          conn.Read(buffer)
    //          time_elapsed := time.Since(time_start)
 
    //          fmt.Printf("Reply from server: requests served = %s\n", string(buffer))
    //          fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)    
 
    //      case "4\n":
    //          input := "4"
 
    //          time_start := time.Now()
    //          conn.Write([]byte(input))
 
    //          conn.Read(buffer)
    //          time_elapsed := time.Since(time_start)
 
    //          fmt.Printf("Reply from server: run time = %s\n", string(buffer))
    //          fmt.Printf("RTT = %.3f ms\n\n", float64(time_elapsed.Nanoseconds())/1000000.0)
         
    //      case "5\n":
    //          input := "5"
    //          conn.Write([]byte(input))
 
    //          fmt.Printf("Bye bye~")
    //          conn.Close()
    //          return;
    //      default:
    //          fmt.Printf("Wrong option\n")
    //      }
    }
 
 }
 