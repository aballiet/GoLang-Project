package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

var tabUser = make(map[string]net.Conn)
var myChan = make(chan string)

const servername = "ChatTC"

func main() {

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	//on ne quitte pas le main tant que la go routine n'a pas terminé
	//var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		go sendMessage()
	}

	for {
		// accept connection on port
		conn, error := ln.Accept()
		if error == nil {
			go read(conn)
		} else {
			fmt.Println("Problème de connexion avec le client !")
		}
	}
}

func read(conn net.Conn) {
	//on envoie un message de welcome
	welcome_message := "TCCHAT_WELCOME\t" + servername + "\n"
	conn.Write([]byte(welcome_message))

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, error := bufio.NewReader(conn).ReadString('\n')

		if error != nil {
			managemessage("TCCHAT_DISCONNECT\t", conn)
			break
		}

		if len(message) > 0 {
			message = message[0 : len(message)-1]
		}

		managemessage(message, conn)
	}
}

//Cette fonction choisit l'action à affectuer en fonction du type de message reçu
func managemessage(message string, conn net.Conn) {
	const Welcome = "TCCHAT_WELCOME"
	const UserIn = "TCCHAT_USERIN"
	const UserOut = "TCCHAT_USEROUT"
	const BCast = "TCCHAT_BCAST"
	const Register = "TCCHAT_REGISTER"
	const Message = "TCCHAT_MESSAGE"
	const Disconnect = "TCCHAT_DISCONNECT"

	tabMessage := strings.Split(message, "\t")

	switch tabMessage[0] {

	case Register:
		//il faut envoyer le userId à tous les autres utilisateurs avec USERIN
		for user := range tabUser {
			myChan <- UserIn + "\t" + tabMessage[1] + "@" + user
			//on envoie tous les noms d'utilisateur au client qui vient de se connecter
			myChan <- UserIn + "\t" + user + "@" + tabMessage[1]
		}

		//on ajoute au tableau le nouveau utilisateur
		tabUser[tabMessage[1]] = conn

		fmt.Println("Un nouvel utilisateur s'est connecté : " + tabMessage[1])

		//enregistrement du nickname dans les utilisateurs actifs
		//If second message with same nickname -> terminate connection with client

	case Message:
		username := ""
		//retrieve username from conn
		for user, connTab := range tabUser {
			if conn == connTab {
				username = user
			}
		}
		for user := range tabUser {
			if user != username {
				myChan <- BCast + "\t" + username + "\t" + tabMessage[1] + "@" + user
			}
		}

		//140 characters (verifier length payload)
		//on broadcast by server -> all clients  attention elle peut contenir "\t" !!!!

	case Disconnect:
		username := ""
		//retrieve username from conn
		for user, connTab := range tabUser {
			if conn == connTab {
				username = user
			}
		}

		fmt.Println("L'utilisateur : " + string(username) + " s'est déconnecté ! ")

		delete(tabUser, username)
		for user := range tabUser {
			if user != username {
				myChan <- UserOut + "\t" + username + "@" + user
			}
		}
	}

}

func sendMessage() {
	for {
		// on lit dans le channel : nom utilisateur + message
		command := <-myChan
		tab := strings.Split(command, "@")
		conn := tabUser[tab[1]]
		conn.Write([]byte(tab[0] + "\n"))

	}

}
