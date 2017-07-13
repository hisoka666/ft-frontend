package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
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
	Token  string   `json:"token"`
	User   string   `json:"user"`
	Bulan  []string `json:"bulan"`
	Pasien []Pasien `json:"pasien"`
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

type KunjunganPasien struct {
	Diagnosis, LinkID      string
	GolIKI, ATS, ShiftJaga string
	JamDatang              time.Time
	Dokter                 string
	Hide                   bool
	JamDatangRiil          time.Time
	Bagian                 string
}

type DataPasien struct {
	NamaPasien, NomorCM, JenKel, Alamat string
	TglDaftar, Umur                     time.Time
}

type Response struct {
	Token  string `json:"token"`
	Script string `json:"script"`
	Modal  string `json:"modal"`
}

type InputPts struct {
	*DataPasien      `json:"datapts"`
	*KunjunganPasien `json:"kunjungan"`
}

type UbahPasien struct {
	NoCM       string    `json:"nocm"`
	NamaPasien string    `json:"namapts"`
	Diagnosis  string    `json:"diag"`
	ATS        string    `json:"ats"`
	Shift      string    `json:"shift"`
	Bagian     string    `json:"bagian"`
	IKI        string    `json:"iki"`
	LinkID     string    `json:"link"`
	TglAsli    time.Time `json:"tgl"`
}

type ModalTemplate struct {
	Script  string      `json:"script"`
	Content *UbahPasien `json:"content"`
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
	http.HandleFunc("/getcm", getCM)
	http.HandleFunc("/inputdata", inputData)
	http.HandleFunc("/editentri", editEntri)
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

func editEntri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/edit"

	pts := Pasien{
		LinkID: r.FormValue("link"),
	}

	resp, err := sendPost(pts, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "kesalahan-server", "")
	}
	kun := &UbahPasien{}
	json.NewDecoder(resp.Body).Decode(kun)

	b := new(bytes.Buffer)
	tmp := template.Must(template.New("modedit.html").ParseFiles("templates/modedit.html"))
	err = tmp.Execute(b, nil)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "")
	}

	mod := &ModalTemplate{
		Script:  b.String(),
		Content: kun,
	}
	fmt.Print(mod.Content)
	json.NewEncoder(w).Encode(mod)

}
func inputData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	data := &DataPasien{}
	if r.FormValue("baru") == "true" {
		data.NamaPasien = r.FormValue("namapts")
		data.NomorCM = r.FormValue("nocm")
		data.TglDaftar = CreateTime()
	} else {
		data.NamaPasien = r.FormValue("namapts")
		data.NomorCM = r.FormValue("nocm")
	}

	kun := &KunjunganPasien{
		Diagnosis:     r.FormValue("diag"),
		GolIKI:        r.FormValue("iki"),
		ATS:           r.FormValue("ats"),
		ShiftJaga:     r.FormValue("shift"),
		JamDatang:     CreateTime(),
		JamDatangRiil: CreateTime(),
		Dokter:        r.FormValue("dok"),
		Bagian:        r.FormValue("bagian"),
	}

	input := InputPts{data, kun}

	url := "http://2.igdsanglah.appspot.com/inputpts"

	resp, err := sendPost(input, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "kesalahan-server", "")
	}

	pts := &Pasien{}

	json.NewDecoder(resp.Body).Decode(pts)
	if pts.NoCM == "kesalahan-database" {
		responseTemplate(w, "kesalahan-database", "")
	}
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("baristabel.html").ParseFiles("templates/baristabel.html"))
	err = tmp.Execute(b, pts)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "")
	}
	responseTemplate(w, "OK", b.String())

}

func CreateTime() time.Time {
	t := time.Now()
	zone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		fmt.Println("Err: ", err.Error())
	}
	jam := t.In(zone)
	return jam
}

func getCM(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintln(w, "Akses ditolak")
	}

	token := r.FormValue("token")
	// fmt.Print(token)
	// fmt.Println(r.FormValue("nocm"))
	pts := &Pasien{
		NoCM: r.FormValue("nocm"),
	}

	url := "http://2.igdsanglah.appspot.com/getcm"

	resp, err := sendPost(pts, token, url)
	if err != nil {
		// fmt.Print(err)
		responseTemplate(w, "kesalahan-server", "")
	}

	json.NewDecoder(resp.Body).Decode(pts)
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("inputpts.html").ParseFiles("templates/inputpts.html"))
	err = tmp.Execute(b, pts)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "")
	}
	responseTemplate(w, "OK", b.String())

}

func sendPost(n interface{}, token, url string) (*http.Response, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(n)
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
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
