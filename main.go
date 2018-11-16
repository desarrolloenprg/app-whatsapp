// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"

// 	"github.com/kenshindeveloper/app-whatsapp/src/libs"

// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/sheets/v4"
// )

// func main() {
// 	b, err := ioutil.ReadFile("credentials.json")
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// If modifying these scopes, delete your previously saved token.json.
// 	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	client := libs.GetClient(config)
// 	srv, err := sheets.New(client)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Sheets client: %v", err)
// 	}

// 	// Prints the names and majors of students in a sample spreadsheet:
// 	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
// 	// spreadsheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
// 	spreadsheetId := "1Znpq-5v4DegW_dVfvKdLZPYp8Zt7KJrJ2ToA7uzb5xM"
// 	readRange := "curso!A1:F8"
// 	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve data from sheet: %v", err)
// 	}

// 	if len(resp.Values) == 0 {
// 		fmt.Println("No data found.")
// 	} else {
// 		fmt.Println("Name, Major:")
// 		for _, row := range resp.Values {
// 			// Print columns A and E, which correspond to indices 0 and 4.
// 			// fmt.Printf("%s, %s\n", row[0], row[4])
// 			if len(row) > 0 {
// 				fmt.Printf("%s \n", row[0])
// 			}
// 		}
// 	}
// }

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
	Text  string
	Error string
}

var qrString = ""
var flagConn = false
var connErr = 0
var wacConn *whatsapp.Conn

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusFound)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./src/templates/base.html", "./src/templates/home.html"))

	data := Data{}
	switch connErr {
	case 1:
		data.Error = "Error, no se ha establecido la conexión con whatsapp."
	case 2:
		data.Error = "Error, no es posible enviar msg con el tipo actual."
	case 10:
		data.Error = "Mensaje enviado."
	case 11:
		data.Error = "Mensaje NO enviado."
	}

	t.ExecuteTemplate(w, "base", data)
}

func handleQrCodePage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./src/templates/base.html", "./src/templates/qrcode.html"))
	data := Data{Text: qrString}
	t.ExecuteTemplate(w, "base", data)
}

func handleSendData(w http.ResponseWriter, r *http.Request) {
	dataType := r.FormValue("send-type")
	dataPhone := r.FormValue("send-phone")
	dataMsg := r.FormValue("send-msg")

	fmt.Printf("type: %s\n", dataType)
	fmt.Printf("phone: %s\n", dataPhone)
	fmt.Printf("msg: %s\n", dataMsg)

	if flagConn && dataType == "0" {
		sendMessage(wacConn, dataPhone, dataMsg)
	} else if !flagConn {
		connErr = 1
	} else if dataType != "0" {
		connErr = 2
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

func main() {
	imgFolder := http.FileServer(http.Dir("src/img"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFolder))

	vendorFolder := http.FileServer(http.Dir("src/vendor"))
	http.Handle("/vendor/", http.StripPrefix("/vendor/", vendorFolder))

	fmt.Printf("Escuchando...\n")
	http.HandleFunc("/", handleIndexPage)
	http.HandleFunc("/home", handleHomePage)
	http.HandleFunc("/sendata", handleSendData)
	http.HandleFunc("/qrcode", handleQrCodePage)
	go connWhatsapp()
	http.ListenAndServe(":8080", nil)
}

//--------------------------------------------------
//--------------------------------------------------
//--------------------------------------------------

func sendMessage(wac *whatsapp.Conn, number string, msg string) {
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Text: msg,
	}
	err := wac.Send(text)
	if err != nil {
		connErr = 11
		fmt.Printf("error al enviar el msg a '%s'.\n", number)
	} else {
		connErr = 10
		fmt.Printf("mensaje enviado a '%s'.\n", number)
	}
}

func connWhatsapp() {
	wac, err := whatsapp.NewConn(120 * time.Second)
	if err != nil {
		log.Fatalf("Error al establecer la conexion...")
	}
	qrChan := make(chan string)
	go func() {
		qrString = <-qrChan
		// fmt.Printf("qrString: %s\n\n", qrString)
		fmt.Printf("abriendo pestaña...\n")
		err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:8080/qrcode").Start()
		if err != nil {
			log.Fatalf("Error al abrir el navegador.\n")
		}
	}()

	sess, err := wac.Login(qrChan)
	if err != nil {
		log.Fatalf("Error al autentificarse...")
	}
	fmt.Printf("abriendo pestaña...\n")
	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:8080/home").Start()
	if err != nil {
		log.Fatalf("Error al abrir el navegador.\n")
	}
	fmt.Printf("id: %s\n", sess.ClientId)
	flagConn = true
	wacConn = wac
	// sendMessage(wac, "584122106942", "prueba de software para envio masivo por whatsapp") //mio
	// sendMessage(wac, "584241951497", "prueba de software para envio masivo por whatsapp") //benjamin
	// fmt.Printf("\nTodos los mensajes han sidos enviados.\n")
}
