package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

var tabUser = make(map[string]net.Conn)
var myChan = make(chan string)

const servername = "ChatTC"

func main() {

	fmt.Println(tabUser)
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	//on ne quitte pas le main tant que la go routine n'a pas terminé
	//var wg sync.WaitGroup

	for i:=0; i<10; i++{
		go sendMessage()
	}

	for{
		// accept connection on port
		conn, _ := ln.Accept()
		go read(conn)
	}
}

func read(conn net.Conn) {
	//on envoie un message de welcome
	welcome_message := "TCCHAT_WELCOME\t" + servername + "\n"
	conn.Write([]byte(welcome_message))

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		managemessage(message, conn)
	}
}

func managemessage(message string, conn net.Conn) {
	const Register = "TCCHAT_REGISTER"
	const Message = "TCCHAT_MESSAGE"
	const Disconnect = "TCCHAT_DISCONNECT"
	const Welcome = "TCCHAT_WELCOME"
	const UserIn = "TCCHAT_USERIN"
	const UserOut = "TCCHAT_USEROUT"
	const BCast = "TCCHAT_BCAST"

	tabMessage := strings.Split(message, "\t")

	switch tabMessage[0] {

	
	case Register:
		//il faut envoyer le userId à tous les autres utilisateurs avec USERIN
		for user,_ := range(tabUser){
			myChan <- UserIn+"\t"+tabMessage[1]+"@"+user
		}

		//on ajoute au tableau le nouveau utilisateur
		tabUser[tabMessage[1]] = conn
		fmt.Println("Un nouvel utilisateur connecté : " + tabMessage[1])

		//enregistrement du nickname dans les utilisateurs actifs
		//If second message with same nickname -> terminate connection with client

	case Message:
		for user,_ := range(tabUser){
			myChan <- Message+"\t"+tabMessage[1]+"@"+user
		}
		//140 characters (verifier length payload)
		//on broadcast by server -> all clients  attention elle peut contenir "\t" !!!!

	case Disconnect:
		//client send disconnect
		//server should close corresponding TCP connection and send UserOut message to connected clients

	}
}

func sendMessage(){
	for {
		// on lit dans le channel : nom utilisateur + message
		command := <- myChan
		fmt.Println(command)
		tab := strings.Split(command, "@")
		conn := tabUser[tab[1]]
		fmt.Println(tab[0])
		conn.Write([]byte(tab[0]))
		fmt.Println("sent !")

	}

}
