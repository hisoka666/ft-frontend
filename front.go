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

type WebView struct {
	Token  string   `json:"token"`
	User   string   `json:"user"`
	Pasien []Pasien `json:"pasien"`
	//IKI      []List    `json:"list"`
}
type NavBar struct {
	Token string   `json:"token"`
	User  string   `json:"user"`
	Bulan []string `json:"bulan"`
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

// func jsonError(m string) []byte {
// 	res := &ReturnError{
// 		Error: m,
// 	}

// 	js, _ := json.Marshal(res)

// 	return js
// }

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

	var web NavBar

	json.NewDecoder(resp.Body).Decode(&web)

	var b bytes.Buffer
	tmp := template.Must(template.New("main.html").ParseFiles("templates/main.html"))
	err = tmp.Execute(&b, web)
	if err != nil {
		fmt.Print(err)
	}

	res := &Response{
		Token:  web.Token,
		Script: b.String(),
	}

	// fmt.Println(b.String())
	// fmt.Fprintln(w, )
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	err = enc.Encode(&res)

	// fmt.Fprintln(w, string(data))

}
