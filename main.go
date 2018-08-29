package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/verify"
)

// Global, because we need to share this with the handler functions
var (
	client       *messagebird.Client
	clientVerify *verify.Verify
)

// RenderDefaultTemplate takes:
// - a http.ResponseWriter
// - an array of strings to contain a list of template files to render
// - data to render to the template. If no data, should enter 'nil'
func RenderDefaultTemplate(w http.ResponseWriter, thisView string, data interface{}) {
	renderthis := []string{thisView, "views/layouts/default.gohtml"}
	t, err := template.ParseFiles(renderthis...)
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(w, "default", data)
	if err != nil {
		log.Fatal(err)
	}
}

// Routes
func step1(w http.ResponseWriter, r *http.Request) {
	RenderDefaultTemplate(w, "views/step1.gohtml", nil)
}

func step2(w http.ResponseWriter, r *http.Request) {
	var err error

	r.ParseForm()
	num := r.FormValue("number")
	clientVerify, err = verify.Create(client, num, nil)
	if err != nil {
		log.Println(err)
	}
	RenderDefaultTemplate(w, "views/step2.gohtml", nil)
}

func step3(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.FormValue("token")
	clientVerifyDone, err := verify.VerifyToken(client, clientVerify.ID, token)
	if err != nil {
		log.Println(err)
	}

	type successState struct {
		Success bool
	}
	var s successState

	if clientVerifyDone.Status == "verified" {
		s = successState{Success: true}
	} else {
		s = successState{Success: false}
	}
	// Execute template and pass verify.Status as a variable into the step3.gohtml template.
	RenderDefaultTemplate(w, "views/step3.gohtml", s)
}

func main() {
	client = messagebird.New("<enter-your-api-key>")

	// Routes
	http.HandleFunc("/", step1)
	http.HandleFunc("/step2", step2)
	http.HandleFunc("/step3", step3)

	// Serve
	port := ":8080"
	log.Println("Serving application on", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println(err)
	}
}
