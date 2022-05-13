/**
 * ChatTCPServer.go
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
    // "sync/atomic"
)

func handleConnection(client *Client) {
    for {
        // read clients chat msg
        buffer := make([]byte, 1024)
        _, err := client.Conn.Read(buffer)
        if err != nil {
            // fmt.Printf("client socket failed to read data from buffer\n")
            log.Fatal(err)
        }
        // fmt.Printf("client's message %s\n", string(buffer))

        // share all the message to connected users
        for ipport, user := range UserDatabase {
            if (ipport != string(client.Conn.RemoteAddr().String())) {
                m := client.nickname + "> " + string(buffer)
                fmt.Println(m)
                user.Conn.Write([]byte(m))
            }
        }

        // // check the header of the message
        // fmt.Printf("Command %s\n",string(buffer[0]))
        // switch string(buffer[0]) {
        // case "1": 
        //     atomic.AddInt64(&Total_served_commands, 1)  // mutul exclusion
        //     client_conn.Write(bytes.ToUpper(buffer[1:count]))
        // case "2":
        //     atomic.AddInt64(&Total_served_commands, 1)
        //     client_conn.Write([]byte(client_conn.RemoteAddr().String()))   
        // case "3":
        //     atomic.AddInt64(&Total_served_commands, 1)
        //     served_count_string := strconv.FormatInt(atomic.LoadInt64(&Total_served_commands), 10)
        //     client_conn.Write([]byte(served_count_string)) 
        // case "4":
        //     atomic.AddInt64(&Total_served_commands, 1)

        //     time_elapsed := time.Since(Time_start)

        //     // result should be HH:MM:SS format
        //     hh_s := strconv.Itoa(int(time_elapsed.Hours()))
        //     if len(hh_s) == 1 {
        //         hh_s = "0" + hh_s
        //     }
        //     mm_s := strconv.Itoa(int(time_elapsed.Minutes())% 60)
        //     if len(mm_s) == 1 {
        //         mm_s = "0" + mm_s
        //     }
        //     ss_s := strconv.Itoa(int(time_elapsed.Seconds())% 60)
        //     if len(ss_s) == 1 {
        //         ss_s = "0" + ss_s
        //     }
        //     result := hh_s + ":" + mm_s + ":" + ss_s

        //     client_conn.Write([]byte(result))
        // case "5":
        //     // close client connection
        //     atomic.AddInt64(&TotalClient, -1)
        //     client_conn.Close()
        //     // fmt.Printf("Client %d disconnected. Number of connected clients = %d\n", id, atomic.LoadInt64(&TotalClient))
        //     return
        // default :
        //     fmt.Printf("Wrong option\n")
        // }
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

// code 0 : chatting room is full
// code 1 : duplicated nickname
var UserDatabase map[string]*Client  // ip,port : user
var Listener net.Conn
var ServerPort string
var Time_start = time.Now()

func main() {
    ServerPort := "44089"
    UserDatabase = make(map[string]*Client) 
 
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
        fmt.Println("Bye bye~")
        // erase all the clients' connection
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
        // check whether chatting room is full
        if (len(UserDatabase) >= 2) {
            fmt.Printf("room is full\n",)
            clientConn.Write([]byte("0"))
            clientConn.Close()
            continue
        }
        
        // save user info at the UserDatabase
        _, exists := UserDatabase[clientConn.RemoteAddr().String()]
        if !exists {
            nickname := make([]byte, 1024)
            clientConn.Read(nickname)
            // check whether nickname is already taken
            if duplicatedNickname(string(nickname[:32])) {
                fmt.Printf("duplicated nickname\n",)
                clientConn.Write([]byte("1"))
                clientConn.Close()  
                continue
            }
            // fmt.Printf("user nickname is  %s\n", string(nickname[:32]))
            UserDatabase[clientConn.RemoteAddr().String()] = &Client{clientConn, string(nickname[:32])}
            // send welcome message
            welcomeMessage := "[Welcome "+string(nickname[:32])+" to CAU network class chat room at 165.194.35.202:" + ServerPort + "]\n" +"[There are " + strconv.Itoa(len(UserDatabase)) + " users connected.]" 
            clientConn.Write([]byte(welcomeMessage))
            log := string(nickname[:32]) + " joined from " + clientConn.RemoteAddr().String() + ". There are " + strconv.Itoa(len(UserDatabase)) + " users connected"
            fmt.Printf("%s\n",log)
        }
        // fmt.Printf("Connection request from %s\n", clientConn.RemoteAddr().String())
        printUserDatabase()
        // client socket do its work
        go handleConnection(UserDatabase[clientConn.RemoteAddr().String()])
        
    }
}
 
 