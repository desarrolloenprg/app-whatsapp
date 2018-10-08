package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

type myHandler struct{}

func (myHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}

func (myHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Println(message)
}

func (myHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	fmt.Println(message)
}

func (myHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	fmt.Println(message)
}

func (myHandler) HandleJsonMessage(message string) {
	fmt.Println(message)
}

func saveSession(session whatsapp.Session) {
	file, err := os.Create("info_session.txt")
	if err != nil {
		log.Fatalf("Error al crear el fichero de almacenamiento de la sesion.\n")
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(session.ClientId + "\n")
	writer.Flush()
	writer.WriteString(session.ClientToken + "\n")
	writer.WriteString(string(session.EncKey) + "\n")
	writer.WriteString(string(session.MacKey) + "\n")
	writer.WriteString(session.ServerToken + "\n")
	writer.WriteString(session.Wid + "\n")
}

func sendMessage(wac *whatsapp.Conn, number string, msg string) {
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Text: msg,
	}
	err := wac.Send(text)
	if err != nil {
		fmt.Printf("error al enviar el msg a '%s'.\n", number)
	} else {
		fmt.Printf("mensaje enviado a '%s'.\n", number)
	}
}

func main() {
	wac, err := whatsapp.NewConn(30 * time.Second)
	if err != nil {
		log.Fatalf("Error al establecer la conexion...")
	}

	qrChan := make(chan string)
	go func() {
		fmt.Printf("qr code: %v\n", <-qrChan)
		//show qr code or save it somewhere to scan
	}()
	sess, err := wac.Login(qrChan)
	if err != nil {
		log.Fatalf("Error al autentificarse...")
	}
	saveSession(sess)

	// ca, err := wac.RestoreSession(sess)
	fmt.Printf("id: %d\n", sess.ClientId)
	// wac.AddHandler(myHandler{})

	for i := 0; i < 10; i++ {
		sendMessage(wac, "584122106942", "ULTIMA (por hoy, de pana gracias XD), prueba de software para envio masivo por whatsapp") //mio
		sendMessage(wac, "584241951497", "ULTIMA (por hoy, de pana gracias XD),prueba de software para envio masivo por whatsapp")  //benjamin
		sendMessage(wac, "593961165103", "ULTIMA (por hoy, de pana gracias XD),prueba de software para envio masivo por whatsapp")  //valeria
		sendMessage(wac, "56932741963", "ULTIMA (por hoy, de pana gracias XD),prueba de software para envio masivo por whatsapp")   //rosario

		sendMessage(wac, "34649948216", "ULTIMA (por hoy, de pana gracias XD),prueba de software para envio masivo por whatsapp")  //javier
		sendMessage(wac, "593979816906", "ULTIMA (por hoy, de pana gracias XD),prueba de software para envio masivo por whatsapp") //german
		sendMessage(wac, "770741394", "ULTIMA (por hoy, de pana gracias XD),prueba de software para envio masivo por whatsapp")    //armando
	}

	fmt.Printf("\nTodos los mensajes han sidos enviados.\n")

	// text := whatsapp.TextMessage{
	// 	Info: whatsapp.MessageInfo{
	// 		RemoteJid: "584122106942@s.whatsapp.net",
	// 	},
	// 	Text: "Hello Whatsapp",
	// }

	// err = wac.Send(text)
	// if err != nil {
	// 	log.Fatalf("error al enviar el msg...")
	// }
}
