package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strings" 

const UserId = "client1"
func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "127.0.0.1:8081")
  for { 
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
    // send to socket
    fmt.Fprintf(conn, text + "\n")
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message from server: "+message)

    managemessage(message, conn)
  }
}
func managemessage(message string,conn net.Conn ){
	const Register = "TCCHAT_REGISTER"
	const Message = "TCCHAT_MESSAGE"
	const Disconnect = "TCCHAT_DISCONNECT"
	const Welcome = "TCCHAT_WELCOME"
	const UserIn = "TCCHAT_USERIN"
	const UserOut = "TCCHAT_USEROUT"
	const BCast = "TCCHAT_BCAST"

	tabMessage := strings.Split(message, "\t")

	switch tabMessage[0] {

	case Welcome:
		fmt.Print("Welcome reçu dans le chat" + tabMessage[1])
    conn.Write([]byte("TCCHAT_REGISTER\t"+ UserId+ "\n"))


	case Message:
		//140 characters (verifier length payload)
		//on broadcast by server -> all clients  attention elle peut contenir "\t" !!!!

	case Disconnect:
		//client send disconnect
		//server should close corresponding TCP connection and send UserOut message to connected clients

	}
}

