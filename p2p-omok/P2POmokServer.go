/**
 * P2POmokServer.go
 * NAME : Seokwon Yoon
 * ID : 20174089
 **/

package main

import (
    // "bytes"
    "fmt"
    "net"
    "time"
    "strconv"
    "log"
    "os"
    "os/signal"
    "syscall"
    "strings"
)

// for debugging
func printUserDatabase() {
    fmt.Println("-----------------------------------------------")
    for ipport, user := range UserDatabase {
		fmt.Println("   ip/port:", ipport, "=>", "user:%s", user)
	}
}
func duplicatedNickname(nick string) bool {
    for _, user := range UserDatabase {
        if user.nickname == nick {
            return true;
        }
    }
    return false;
}

type Client struct {
    Conn net.Conn
    nickname string
}

// code 1 : duplicated nickname
// code 3 : ordinary message
// code 5 : close connection
// code 7 : \list
// code 6 : \dm
// code 5 : \exit
// code 8 : \ver
// code 9 : \rtt
var UserDatabase = make(map[string]*Client)   // ip,port : user
var Listener net.Conn
var ServerPort = "54089"
var Time_start = time.Now()
var SoftwareVersion = "1.0.0"

func main() {
    // create server socket
    Listener, err:= net.Listen("tcp", ":" + ServerPort)
    if err != nil {
        // fmt.Printf("create server socket failed\n")
        log.Fatal(err)
    } else {
        fmt.Printf("Server is ready to receive on port %s\n", ServerPort)
    }

    // ctrl + c handling
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("gg~")
        
        // close all the clients' connection
        for _, user := range UserDatabase {
            user.Conn.Write([]byte("5"))
            // user.Conn.Write([]byte("[Server has been shut down]"))
            user.Conn.Close()
        }
        Listener.Close()
        close(c)
        os.Exit(0)
    }()

    oddPlayer := ""
    evenPlayer := ""
    oddUdpIp := ""
    evenUdpIp := ""
    oddUdpPort := ""
    evenUdpPort := ""

    for {
        // create client socket
        clientConn, err:= Listener.Accept()
        if err != nil {
            // fmt.Printf("create client socket failed\n")
            log.Fatal(err)
        }
        
        // save user info at the UserDatabase
        _, exists := UserDatabase[clientConn.RemoteAddr().String()]
        if !exists {
            nicknameAndUdpAddr := make([]byte, 1024)
            n, _ := clientConn.Read(nicknameAndUdpAddr)
            // fmt.Printf("user nickname is  %s\n", string(nickname[:32]))
            s := strings.Split(string(nicknameAndUdpAddr[:n]), " ")
            nick := s[0]
            udpAddr := s[1]

            // check whether nickname is already taken
            if duplicatedNickname(nick) {
                fmt.Printf("duplicated nickname\n",)
                clientConn.Write([]byte("1"))
                clientConn.Close()  
                continue
            }

            // add user to the database
            UserDatabase[clientConn.RemoteAddr().String()] = &Client{clientConn, string(nick)}

            welcomeMessage := "welcome "+string(nick)+" to p2p-omok server at 165.194.35.202:" + ServerPort + "\n"

            // send welcome message
            if (len(UserDatabase) % 2  == 0) {
                evenUdpPort = udpAddr
                evenPlayer = clientConn.RemoteAddr().String()
                s := strings.Split(clientConn.RemoteAddr().String(), ":")
                evenUdpIp = s[0]
                welcomeMessage = welcomeMessage + UserDatabase[oddPlayer].nickname + " is waiting for you (" + oddUdpIp + ":" + oddUdpPort + ").\n" + UserDatabase[oddPlayer].nickname + " plays first.\n"
                clientConn.Write([]byte("3" + welcomeMessage))
                UserDatabase[oddPlayer].Conn.Write([]byte("3" + UserDatabase[evenPlayer].nickname + " joined (" + evenUdpIp + ":" + evenUdpPort + "). you play first.\n"))

            } else {
                oddUdpPort = udpAddr
                oddPlayer = clientConn.RemoteAddr().String()
                s := strings.Split(clientConn.RemoteAddr().String(), ":")
                oddUdpIp = s[0]
                welcomeMessage = welcomeMessage + "waiting for an opponent\n"
                clientConn.Write([]byte("3" + welcomeMessage))
            }
            log := string(nick) + " udp addr " + udpAddr + ". There are " + strconv.Itoa(len(UserDatabase)) + " users connected."
            fmt.Printf("%s\n",log)
                
            // close TCP connection
            if (len(UserDatabase) % 2  == 0) {
                // wait for ack message from client
                buffer := make([]byte, 1024)
                UserDatabase[oddPlayer].Conn.Read(buffer)
                
                UserDatabase[oddPlayer].Conn.Write([]byte("5" + UserDatabase[evenPlayer].nickname + " " + evenUdpIp + ":" + evenUdpPort + " 1"))
                UserDatabase[evenPlayer].Conn.Write([]byte("5" + UserDatabase[oddPlayer].nickname + " " + oddUdpIp + ":" + oddUdpPort + " 0"))

                UserDatabase[oddPlayer].Conn.Close()
                UserDatabase[evenPlayer].Conn.Close()

                delete(UserDatabase, string(oddPlayer));
                delete(UserDatabase, string(evenPlayer));
                // printUserDatabase()
            }
        }
    }
}
 
 