package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func parseFile(filename string) string {
	b, err := ioutil.ReadFile(filename) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	texte := string(b) // convert content to a 'string'
	tabLigne := strings.Split(texte, `\n`)

	fmt.Println(texte) // print the content as a 'string'

	for _, elem := range tabLigne {
		fmt.Println(elem)
	}

	return tabLigne[0]
}

func manageMessage(message string) {
	const Register = "TCCHAT_REGISTER"
	const Message = "TCCHAT_MESSAGE"
	const Disconnect = "TCCHAT_DISCONNECT"
	const Welcome = "TCCHAT_WELCOME"
	const UserIn = "TCCHAT_USERIN"
	const UserOut = "TCCHAT_USEROUT"
	const BCast = "TCCHAT_BCAST"

	tabMessage := strings.Split(message, `\t`)
	fmt.Print(tabMessage)

	switch tabMessage[0] {

	case Welcome:
		//send welcome to all
	case Register:
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

func main() {
	parseFile("file.txt")

}
