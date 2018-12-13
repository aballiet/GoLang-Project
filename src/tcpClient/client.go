package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strings" 

var UserId = "default"

func main() {

	// read in input from stdin
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Entrer un nom d'utilisateur : ")
	text, _ := reader.ReadString('\n')
	UserId = text

	// connect to this socket + ATTENTIION GESTION ERREUR
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	//on crée une go routine qui gère la réception des messages
	go receiver(conn)

	for { 
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')

		// send to server
		conn.Write([]byte ("TCCHAT_MESSAGE\t"+text+"\n"))
	}
}

func receiver(conn net.Conn){
	for{
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message from server: "+message)
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
	fmt.Println(message)

	switch tabMessage[0] {

	case Welcome:
		fmt.Println("Welcome reçu dans le chat" + tabMessage[1])
    	conn.Write([]byte("TCCHAT_REGISTER\t"+ UserId+ "\n"))

	case UserIn:
		fmt.Println(tabMessage[1]+ "vient de se connecter")

	case Message:
		fmt.Println(tabMessage)
		//140 characters (verifier length payload)
		//on broadcast by server -> all clients  attention elle peut contenir "\t" !!!!

	case Disconnect:
		//client send disconnect
		//server should close corresponding TCP connection and send UserOut message to connected clients

	}
}

