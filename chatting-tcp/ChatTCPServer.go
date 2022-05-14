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
    "strings"
)

func handleConnection(client *Client) {
    for {
        // read clients chat msg
        buffer := make([]byte, 1024)
        n, err := client.Conn.Read(buffer)
        if err != nil {
            // fmt.Printf("client socket failed to read data from buffer\n")
            log.Fatal(err)
        }
        // fmt.Printf("client's message %s\n", string(buffer))

        switch string(buffer[0]) {
        case "5": // close connection
            // leaving message
            m := "<" + client.nickname + "> left. There are " + strconv.Itoa(len(UserDatabase)-1) + " users now\n"
            fmt.Print(m)

            delete(UserDatabase, string(client.Conn.RemoteAddr().String()));
            client.Conn.Close()
            // printUserDatabase()
            // share leaving msg to connected users
            for _, user := range UserDatabase {
                user.Conn.Write([]byte("3" + m))
            }
            return

        case "3": // ordinary message
            // detect "i hate professor"
            warning := ""
            violation := false
            if strings.Contains(string(bytes.ToUpper(buffer[1:])), "I HATE PROFESSOR") {
                // fmt.Println("i hate professor has detected")
                delete(UserDatabase, string(client.Conn.RemoteAddr().String()));
                client.Conn.Write([]byte("2"))
                warning = "[" + client.nickname + " is disconnected. There are " + strconv.Itoa(len(UserDatabase)) + " users in the chat room.]\n"
                violation = true
            }
            // share all the message to connected users
            for ipport, user := range UserDatabase {
                if ipport != string(client.Conn.RemoteAddr().String()) {
                    m := client.nickname + "> " + string(buffer[1:n])
                    user.Conn.Write([]byte("3" + m))
                }
                if violation {
                    user.Conn.Write([]byte(warning))
                }
            }
            if violation {
                fmt.Print(warning)
                client.Conn.Close()
                return;
            }
        // \command cases
        case "7":  // \list
            m := ""
            for ipport, user := range UserDatabase {
                split := strings.Split(ipport, ":")
                m += "<" + user.nickname + ", " + split[0] + ", " + split[1] +">\n"                
            }
            client.Conn.Write([]byte("3" + m))
        case "8":  // \ver
            client.Conn.Write([]byte("3" + SoftwareVersion + "\n"))
        case "9":  // \rtt
            client.Conn.Write([]byte("9"))
        case "6":  // \dm
            split := strings.Split(string(buffer), " ")  

            // detect "i hate professor"
            warning := ""
            violation := false
            found := false;
            if strings.Contains(string(bytes.ToUpper(buffer[1:])), "I HATE PROFESSOR") {
                // fmt.Println("i hate professor has detected")
                delete(UserDatabase, string(client.Conn.RemoteAddr().String()));
                client.Conn.Write([]byte("2"))
                warning = "[" + client.nickname + " is disconnected. There are " + strconv.Itoa(len(UserDatabase)) + " users in the chat room.]\n"
                violation = true
            }
            // send dm to the specific user
            for ipport, user := range UserDatabase {
                if ipport != string(client.Conn.RemoteAddr().String()) {
                    if strings.Compare(user.nickname, split[1]) == 0 {
                        m := "from: " + client.nickname + "> " + string(buffer[3 + len(split[1]):n])
                        user.Conn.Write([]byte("3" + m))
                        // fmt.Print(m)
                        found = true
                        break;
                    }
                }
            }

            if violation {
                for ipport, user := range UserDatabase {
                    if ipport != string(client.Conn.RemoteAddr().String()) {
                        if strings.Compare(user.nickname, split[1]) == 0 {
                            user.Conn.Write([]byte(warning))
                        } else {
                            user.Conn.Write([]byte("3" + warning))
                        }
                    }
                }
            }

            if found == false {
                client.Conn.Write([]byte("3There's no such other user, " + split[1] +"\n"))
            }
            if violation {
                fmt.Print(warning)
                client.Conn.Close()
                return;
            }
        default :
            fmt.Printf("invalid command\n")
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

// code 0 : chatting room is full
// code 1 : duplicated nickname
// code 2 : "i hate professor" detected
// code 3 : ordinary message
// code 5 : close connection
// code 7 : \list
// code 6 : \dm
// code 5 : \exit
// code 8 : \ver
// code 9 : \rtt
var UserDatabase = make(map[string]*Client)   // ip,port : user
var Listener net.Conn
var ServerPort = "44089"
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

    for {
        // create client socket
        clientConn, err:= Listener.Accept()
        if err != nil {
            // fmt.Printf("create client socket failed\n")
            log.Fatal(err)
        }
        // check whether chatting room is full, it can handle only 8 users
        if (len(UserDatabase) >= 8) {
            fmt.Printf("chatting room full. cannot connect\n",)
            clientConn.Write([]byte("0"))
            clientConn.Close()
            continue
        }
        
        // save user info at the UserDatabase
        _, exists := UserDatabase[clientConn.RemoteAddr().String()]
        if !exists {
            nickname := make([]byte, 1024)
            n, _ := clientConn.Read(nickname)
            // fmt.Printf("user nickname is  %s\n", string(nickname[:32]))

            // check whether nickname is already taken
            if duplicatedNickname(string(nickname[:n])) {
                fmt.Printf("duplicated nickname\n",)
                clientConn.Write([]byte("1"))
                clientConn.Close()  
                continue
            }

            // add user to the database
            UserDatabase[clientConn.RemoteAddr().String()] = &Client{clientConn, string(nickname[0:n])}

            // send welcome message
            welcomeMessage := "[Welcome "+string(nickname[:n])+" to CAU network class chat room at 165.194.35.202:" + ServerPort + "]\n" +"[There are " + strconv.Itoa(len(UserDatabase)) + " users connected.]\n" 
            // clientConn.Write([]byte(welcomeMessage))
            log := string(nickname[0:n]) + " joined from " + clientConn.RemoteAddr().String() + ". There are " + strconv.Itoa(len(UserDatabase)) + " users connected"
            fmt.Printf("%s\n",log)
            for _, user := range UserDatabase {
                user.Conn.Write([]byte("3"+ welcomeMessage))
            }
        }
        // fmt.Printf("Connection request from %s\n", clientConn.RemoteAddr().String())
        // printUserDatabase()

        // client socket do its work
        go handleConnection(UserDatabase[clientConn.RemoteAddr().String()])
        
    }
}
 
 