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


func channelForFirstPlayerMsg(client *Client) {
    for {
        // read clients chat msg
        buffer := make([]byte, 1024)
        client.Conn.Read(buffer)
        
        switch string(buffer[0]) {
        case "5": // close connection
            // leaving message
            m := "<" + client.nickname + "> left. There are " + strconv.Itoa(len(UserDatabase)-1) + " users now\n"
            fmt.Print(m)

            delete(UserDatabase, string(client.Conn.RemoteAddr().String()));
            client.Conn.Close()
            // printUserDatabase()
            return
        case "3":
            // fmt.Printf("matching finished\n")
        default :
            fmt.Printf(string(buffer))
            // fmt.Printf("matching done\n")
            return
        }
        
    }
}

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
var oddPlayer = ""
var evenPlayer = ""
var oddUdpIp = ""
var evenUdpIp = ""
var oddUdpPort = ""
var evenUdpPort = ""

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
            s := strings.Split(string(nicknameAndUdpAddr[1:n]), " ")
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

            log := string(nick) + " joined from " + clientConn.RemoteAddr().String() + ". UDP port " + udpAddr + ".\n" + strconv.Itoa(len(UserDatabase)) + " user connected,"
            fmt.Printf("%s",log)

            // send welcome message
            if (len(UserDatabase) % 2  == 0) {
                evenUdpPort = udpAddr
                evenPlayer = clientConn.RemoteAddr().String()
                s := strings.Split(clientConn.RemoteAddr().String(), ":")
                evenUdpIp = s[0]
                welcomeMessage = welcomeMessage + UserDatabase[oddPlayer].nickname + " is waiting for you (" + oddUdpIp + ":" + oddUdpPort + ").\n"
                clientConn.Write([]byte("3" + welcomeMessage))

                buffer := make([]byte, 1024)
                clientConn.Read(buffer)

                UserDatabase[evenPlayer].Conn.Write([]byte("3" + UserDatabase[oddPlayer].nickname + " plays first.\n"))
                UserDatabase[oddPlayer].Conn.Write([]byte("3" + UserDatabase[evenPlayer].nickname + " joined (" + evenUdpIp + ":" + evenUdpPort + "). you play first.\n"))
            
                // buffer := make([]byte, 1024)
                // print("read start\n")
                clientConn.Read(buffer)
                // print("read done\n")

                UserDatabase[evenPlayer].Conn.Write([]byte("5" + UserDatabase[oddPlayer].nickname + " " + oddUdpIp + ":" + oddUdpPort + " 0"))
                UserDatabase[oddPlayer].Conn.Write([]byte("5" + UserDatabase[evenPlayer].nickname + " " + evenUdpIp + ":" + evenUdpPort + " 1"))

                UserDatabase[oddPlayer].Conn.Close()
                UserDatabase[evenPlayer].Conn.Close()

                log := " notifying " + UserDatabase[oddPlayer].nickname + " and " + UserDatabase[evenPlayer].nickname + ".\n" + UserDatabase[oddPlayer].nickname + " and " + UserDatabase[evenPlayer].nickname + " disconnected."
                fmt.Printf("%s\n",log)

                delete(UserDatabase, string(oddPlayer));
                delete(UserDatabase, string(evenPlayer));

                // printUserDatabase()
            } else {
                oddUdpPort = udpAddr
                oddPlayer = clientConn.RemoteAddr().String()
                s := strings.Split(clientConn.RemoteAddr().String(), ":")
                oddUdpIp = s[0]
                welcomeMessage = welcomeMessage + "waiting for an opponent\n"
                clientConn.Write([]byte("3" + welcomeMessage))
                go channelForFirstPlayerMsg(UserDatabase[clientConn.RemoteAddr().String()])

                log := " waiting for another"
                fmt.Printf("%s\n",log)
            }
        }
    }
}
 
 