package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
)

/*tableau contenant tous les messages avec une partie initilisée pour expliquer
fonctionnalités à l'utilisateur*/
var posts = []string{
	"Bienvenue dans le client de messagerie !",
	"Appuyer à tout moment sur Echap pour quitter le programme.",
	"Vous pouvez utilisez les flèches pour vous déplacer dans la liste des messages",
}

//variables pour l'interface
var history = tui.NewVBox()
var historyScroll = tui.NewScrollArea(history)
var ui tui.UI

//Liste des utilisateurs connectés
var userList = []string{}

//bare latérale dans laquelle on affiche les utilisateurs présents
var labelSidebar = tui.NewLabel("")

var UserID = "default"

//structure permettant d'inclure un style à une box
type StyledBox struct {
	Style string
	*tui.Box
}

// Draw le widget en respectant le style
func (s *StyledBox) Draw(p *tui.Painter) {
	p.WithStyle(s.Style, func(p *tui.Painter) {
		s.Box.Draw(p)
	})
}

func main() {

	// on récupère le nom d'utilisateur
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Entrer un nom d'utilisateur : ")
	text, _ := reader.ReadString('\n')
	UserID = strings.TrimSuffix(text, "\n")
	userList = append(userList, "Moi : "+UserID)
	labelSidebar.SetText("Utilisateur présents :\n" + "Moi : " + UserID)

	// connect to this socket + ATTENTIION GESTION ERREUR
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("Impossible de se connecter !")
		os.Exit(1)
	}

	//on crée une go routine qui gère la réception des messages
	go receiver(conn)

	//on crée l' interface graphique
	setUserInterface(conn)
}

func receiver(conn net.Conn) {
	for {
		// listen for reply
		message, error := bufio.NewReader(conn).ReadString('\n')
		if error != nil {
			fmt.Println("Server disconnected !")
			os.Exit(1)
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

	message = strings.TrimSuffix(message, "\n")
	tabMessage := strings.Split(message, "\t")

	switch tabMessage[0] {

	case Welcome:
		displayMessage("Serveur", "Bienvenue dans le "+tabMessage[1])
		conn.Write([]byte("TCCHAT_REGISTER\t" + UserID + "\n"))

	case UserIn:
		displayMessage("Serveur", tabMessage[1]+" est connecté")
		addUserSidebar(tabMessage[1])

	case BCast:
		displayMessage(tabMessage[1], tabMessage[2])
		addUserSidebar(tabMessage[1])

	case UserOut:
		displayMessage("Serveur", tabMessage[1]+" s'est déconnecté")
		removeUserSidebar(tabMessage[1])

	}
}

func setUserInterface(conn net.Conn) {
	theme := tui.NewTheme()

	theme.SetStyle("label.sidebarTheme", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorGreen})
	theme.SetStyle("label.timeTheme", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorYellow})
	theme.SetStyle("label.meSending", tui.Style{Bold: tui.DecorationOn, Bg: tui.ColorDefault, Fg: tui.ColorCyan})
	theme.SetStyle("label.theySending", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorWhite})
	theme.SetStyle("label.message", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorDefault})
	theme.SetStyle("entry", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})

	//on crée la sidebar
	labelSidebar.SetStyleName("sidebarTheme")
	sidebar := &StyledBox{
		Style: "entry",
		Box:   tui.NewVBox(labelSidebar),
	}
	sidebar.SetTitle("Utilisateurs présents")
	sidebar.SetBorder(true)

	input := tui.NewEntry()

	input.OnSubmit(func(e *tui.Entry) {

		//on envoie le message
		conn.Write([]byte("TCCHAT_MESSAGE\t" + e.Text() + "\n"))

		//on descend tout en bas du scrollArea
		historyScroll.ScrollToBottom()

		//on ajoute le message à l'interface
		displayMessage(UserID, e.Text())
		input.SetText("")
	})

	historyScroll.SetAutoscrollToBottom(false)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := &StyledBox{
		Style: "entry",
		Box:   tui.NewHBox(input),
	}

	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := &StyledBox{
		Style: "entry",
		Box:   tui.NewVBox(historyBox, inputBox),
	}
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(sidebar, chat)

	ui, _ = tui.New(root)

	ui.SetKeybinding("Esc", func() {
		//on envoie le message de deconnexion
		displayMessage("bot", "message envoyé")
		conn.Write([]byte("TCCHAT_DISCONNECT\t"))
		ui.Quit()
	})
	ui.SetKeybinding("Right", func() {
		historyScroll.Scroll(2, 0)
		ui.Repaint()
	})

	ui.SetKeybinding("Left", func() {
		historyScroll.Scroll(-2, 0)
		ui.Repaint()
	})

	ui.SetKeybinding("Down", func() {
		historyScroll.Scroll(0, 2)
		ui.Repaint()
	})

	ui.SetKeybinding("Up", func() {
		historyScroll.Scroll(0, -2)
		ui.Repaint()
	})

	//on ajoute les messages de tuto

	for _, post := range posts {
		displayMessage("Bot", post)
	}

	//on applique le theme
	ui.SetTheme(theme)

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

//permet d'afficher un message à l'écran ATTENTION ON ACTUALISE l'UI !!!
//style : 1-> normal (blanc), 2->
func displayMessage(username string, message string) {
	timeLabel := tui.NewLabel(time.Now().Format("15:04"))
	timeLabel.SetStyleName("timeTheme")

	userLabel := tui.NewLabel(fmt.Sprintf("<%s>", username))
	if username == UserID {
		userLabel.SetStyleName("meSending")
	} else {
		userLabel.SetStyleName("theySending")
	}

	padder := tui.NewPadder(1, 0, userLabel)
	messageLabel := tui.NewLabel(message)
	messageLabel.SetStyleName("theySending")

	history.Append(tui.NewHBox(timeLabel, padder, messageLabel, tui.NewSpacer()))

	//on scroll jusqu'en bas pour afficher le message envoyé à l'instant puis on rafraichit
	historyScroll.ScrollToBottom()
	ui.Repaint()
}

//on met à jour la barre latérale
func updateSidebar() {
	textUser := "---------------------\n"
	for _, user := range userList {
		textUser = textUser + user + "\n"
	}
	labelSidebar.SetText(textUser)

	ui.Repaint()
}

//on regarde si l'utilisateur est dans la liste, s'il n'est pas dedans, on l'ajoute
func addUserSidebar(userToAdd string) {
	estPresent := false
	for _, user := range userList {
		if userToAdd == user {
			estPresent = true
		}
	}

	if !estPresent {
		userList = append(userList, userToAdd)
		updateSidebar()
	}
}

//on regarde si l'utilisateur était déjà présent dans la liste et on le supprime
func removeUserSidebar(userToRemove string) {
	for index, user := range userList {
		if userToRemove == user {
			userList = append(userList[:index], userList[index+1:]...)
		}
	}

	updateSidebar()
}
