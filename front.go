package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// type ReturnError struct {
// 	Error string `json:"error"`
// }

type MainView struct {
	Token  string   `json:"token"`
	User   string   `json:"user"`
	Bulan  []string `json:"bulan"`
	Pasien []Pasien `json:"pasien"`
	//IKI      []List    `json:"list"`
}
type NavBar struct {
	Token   string   `json:"token"`
	User    string   `json:"user"`
	Bulan   []string `json:"bulan"`
	Pasien  []Pasien `json:"pasien"`
}
type Pasien struct {
	TglKunjungan string `json:"tgl"`
	ShiftJaga    string `json:"shift"`
	NoCM         string `json:"nocm"`
	NamaPasien   string `json:"nama"`
	Diagnosis    string `json:"diag"`
	IKI1         string `json:"iki1"`
	IKI2         string `json:"iki2"`
	LinkID       string `json:"link"`
}

type Response struct {
	Token  string `json:"token"`
	Script string `json:"script"`
}

func main() {
	// variable fs membuat folder "script" menjadi sebuah file server,
	// alamat dari file server ini akan diarahkan oleh http.Handle
	// yang akan mengedit semua alamat URL dengan "/script/" menggunakan
	// StripPrefix
	// StripPrefix akan menghapus semua prefix yang berisi "/script/"
	// dan diarahkan ke fs
	fs := http.FileServer(http.Dir("script"))
	http.Handle("/script/", http.StripPrefix("/script/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", mainContent)
	// http.HandleFunc("/getmain", getMain)
	http.ListenAndServe(":9090", nil)
	log.Println("Listening...")
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to ParseFile: %v", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Failed to Execute template: %v", err)
	}
}

func getMain(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	responseTemplate(w, "", "")
	// }

	// token := js.Global.Get("localStorage").Get("token")
	// js.Global.Call("alert", token)
	// fmt.Print(token)
	// responseTemplate(w, token, "")

}

func responseTemplate(w http.ResponseWriter, token, script string) {
	res := &Response{
		Token:  token,
		Script: script,
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&res)

	if err != nil {
		fmt.Print(err)
	}
}

func mainContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintln(w, "Akses ditolak")
	}

	token := r.FormValue("idtoken")
	resp, err := http.Get("http://2.igdsanglah.appspot.com/login?token=" + token)
	if err != nil {
		log.Fatal(err)
		return
	}

	var web MainView

	json.NewDecoder(resp.Body).Decode(&web)

	var b bytes.Buffer
	tmp := template.Must(template.New("main.html").ParseFiles("templates/main.html", "templates/input.html", "templates/content.html"))
	err = tmp.Execute(&b, web)
	if err != nil {
		fmt.Print(err)
	}

	responseTemplate(w, web.Token, b.String())

	// res := &Response{
	// 	Token:  web.Token,
	// 	Script: b.String(),
	// }

	// fmt.Println(b.String())
	// fmt.Fprintln(w, )
	// enc := json.NewEncoder(w)
	// enc.SetEscapeHTML(false)
	// err = enc.Encode(&res)

	// fmt.Fprintln(w, string(data))

}
