package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"time"

	whatsapp "github.com/rhymen/go-whatsapp"
)

type Data struct {
	Text string
}

var qrString = ""

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

// http.Redirect(w, r, "/home/", http.StatusFound)

func handleInitPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./src/templates/base.html", "./src/templates/home.html"))
	data := Data{Text: qrString}
	t.ExecuteTemplate(w, "base", data)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./src/templates/base.html", "./src/templates/home.html"))
	data := Data{Text: qrString}
	t.ExecuteTemplate(w, "base", data)
}

func main() {

	cssFolder := http.FileServer(http.Dir("src/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFolder))

	imgFolder := http.FileServer(http.Dir("src/img"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFolder))

	jsFolder := http.FileServer(http.Dir("src/js"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFolder))

	scssFolder := http.FileServer(http.Dir("src/scss"))
	http.Handle("/scss/", http.StripPrefix("/scss/", scssFolder))

	vendorFolder := http.FileServer(http.Dir("src/vendor"))
	http.Handle("/vendor/", http.StripPrefix("/vendor/", vendorFolder))

	fmt.Printf("Escuchando...")
	http.HandleFunc("/init", handleInitPage)
	http.HandleFunc("/home", handleHomePage)
	// go connWhatsapp()
	http.ListenAndServe(":8080", nil)
}

//--------------------------------------------------
//

func connWhatsapp() {
	wac, err := whatsapp.NewConn(60 * time.Second)
	if err != nil {
		log.Fatalf("Error al establecer la conexion...")
	}
	qrChan := make(chan string)
	go func() {
		// fmt.Printf("qr code: %v\n", <-qrChan)
		qrString = <-qrChan
		fmt.Printf("qrString: %s\n\n", qrString)
		// err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "https://api.qrserver.com/v1/create-qr-code/?size=150x150&data="+qrString).Start()
		err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:8080/init").Start()
		if err != nil {
			log.Fatalf("Error al abrir el navegador.\n")
		}
	}()

	sess, err := wac.Login(qrChan)
	if err != nil {
		log.Fatalf("Error al autentificarse...")
	}
	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:8080/home").Start()
	if err != nil {
		log.Fatalf("Error al abrir el navegador.\n")
	}
	fmt.Printf("id: %d\n", sess.ClientId)

	sendMessage(wac, "584122106942", "prueba de software para envio masivo por whatsapp") //mio
	sendMessage(wac, "584241951497", "prueba de software para envio masivo por whatsapp") //benjamin
	fmt.Printf("\nTodos los mensajes han sidos enviados.\n")
	fmt.Printf("Entro!")
}
