package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing
import "sync"

var tabUser = make(map[string]string)
const servername = "ChatTC"

func main() {

	fmt.Println(tabUser)
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	var wg sync.WaitGroup
	conn, _ := ln.Accept()
	wg.Add(1)
	go read(conn)
	wg.Wait()
	
}


func read(conn net.Conn){
	//on envoie un message de welcome
	welcome_message :="TCCHAT_WELCOME\t"+servername + "\n"
	conn.Write([]byte(welcome_message))
	
	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		managemessage(message, conn)
		// output message received
		fmt.Print("Message Received:", string(message))
	}
}

func managemessage(message string, conn net.Conn ){
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
		fmt.Print("Bienvenue sur le serveur !")
	case Register:
		//il faut envoyer le userId à tous les autres utilisateurs
		//on ajoute au tableau le nouveau utilisateur
		tabUser[string(conn.RemoteAddr())]=tabMessage[1]
		fmt.Print(tabUser)
		//enregistrement du nickname dans les utilisateurs actifs
		//If second message with same nickname -> terminate connection with client

	case Message:
		//140 characters (verifier length payload)
		//on broadcast by server -> all clients  attention elle peut contenir "\t" !!!!

	case Disconnect:
		//client send disconnect
		//server should close corresponding TCP connection and send UserOut message to connected clients

	}
}

