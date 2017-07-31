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
	StatusServer string    `json:"stat"`
	TglKunjungan string    `json:"tgl"`
	ShiftJaga    string    `json:"shift"`
	ATS          string    `json:"ats"`
	Bagian       string    `json:"bagian"`
	NoCM         string    `json:"nocm"`
	NamaPasien   string    `json:"nama"`
	Diagnosis    string    `json:"diag"`
	IKI          string    `json:"iki"`
	LinkID       string    `json:"link"`
	TglAsli      time.Time `json:"tglasli"`
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

type Obat struct {
	Merk        string `json:"merk"`
	Kandungan   string `json:"kandungan"`
	PerkiloMin  int    `json:"kgmin"`
	PerkiloMax  int    `json:"kgmax"`
	Tablet      []int  `json:"tablet"`
	Sirop       []int  `json:"sirop"`
	Drop        []int  `drop:"drop"`
	Rekomendasi string `json:"rekom"`
	InputBy     string `json:"inputby"`
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

type ModalTemplate struct {
	Script  string  `json:"script"`
	Content *Pasien `json:"content"`
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
	http.HandleFunc("/confedit", confEditEntri)
	http.HandleFunc("/delentri", deleteEntri)
	http.HandleFunc("/confdel", confDelete)
	http.HandleFunc("/firstentries", firstEntries)
	http.HandleFunc("/edittgl", editTanggal)
	http.HandleFunc("/confedittgl", confEditTanggal)
	http.HandleFunc("/getptspage", getPtsPage)
	http.HandleFunc("/getprespage", getPresPage)
	http.HandleFunc("/getinputobat", getInputObat)
	http.HandleFunc("/inputobat", inputObat)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8001", nil))

}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Failed to Execute template: %v", err)
	}
}
func ConvertToUbah(r *http.Request) *Pasien {
	n := &Pasien{
		NamaPasien: r.FormValue("namapts"),
		Diagnosis:  r.FormValue("diag"),
		ATS:        r.FormValue("ats"),
		IKI:        r.FormValue("iki"),
		Bagian:     r.FormValue("bagian"),
		LinkID:     r.FormValue("link"),
		ShiftJaga:  r.FormValue("shift"),
	}

	return n
}

////////////////////////////////////////////////////////////////////////////////////
func inputObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	// obat := &Obat{
	// 	Merk:        r.FormValue("merk"),
	// 	Kandungan:   r.FormValue("kand"),
	// 	PerkiloMin:  r.FormValue("mindose"),
	// 	PerkiloMax:  r.FormValue("maxdose"),
	// 	Tablet:      r.FormValue("tab"),
	// 	Sirop:       r.FormValue("syr"),
	// 	Drop:        r.FormValue("drop"),
	// 	Rekomendasi: r.FormValue("rekom"),
	// 	InputBy:     r.FormValue("doc"),
	// }

	fmt.Print(r.Body)

}

////////////////////////////////////////////////////////////////////////////////////
func getInputObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}

	responseTemplate(w, "OK", GenTemplate(nil, "modinputobat"), "")
}

func getPresPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}

	responseTemplate(w, "OK", GenTemplate(nil, "modlistresep"), "")
}

func getPtsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}
	responseTemplate(w, "OK", GenTemplate(nil, "modresep"), "")
}

func confEditTanggal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/confubahtanggal"

	pts := &Pasien{
		TglKunjungan: r.FormValue("tanggal"),
		LinkID:       r.FormValue("link"),
	}

	resp, err := sendPost(pts, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""))
		log.Fatalf("Terjadi kesalahan pengiriman ke server")
	}

	list := MainView{}
	json.NewDecoder(resp.Body).Decode(&list)

	fmt.Printf("Isi dari token adalah: %v", list.Token)
	if list.Token != "OK" {
		log.Fatalf("Terjadi kesalahan server")
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""))
		return
	}

	responseTemplate(w, "OK", GenTemplate(list.Pasien, "contentrefresh"), "")

}

func editTanggal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/ubahtanggal"
	link := r.FormValue("link")
	pts := &Pasien{
		LinkID: link,
	}

	resp, err := sendPost(pts, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""))
		log.Print("Terjadi kesalahan server")
		return
	}

	json.NewDecoder(resp.Body).Decode(pts)
	pts.TglKunjungan = pts.TglAsli.Format("Mon 02/01/2006 15:04:05")
	pts.LinkID = link

	script := GenTemplate(pts, "modubahtgl")
	responseTemplate(w, "OK", script, "")

}

func firstEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/firstitems"

	send := &MainView{
		User: r.FormValue("email"),
	}

	resp, err := sendPost(send, r.FormValue("token"), url)

	if err != nil {
		responseTemplate(w, "kesalahan-client", "", "")
		return
	}

	json.NewDecoder(resp.Body).Decode(send)
	responseTemplate(w, "OK", GenTemplate(send.Pasien, "contentrefresh"), "")

}

func confDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/delentri"
	del := &Pasien{
		LinkID: r.FormValue("link"),
	}

	resp, err := sendPost(del, r.FormValue("token"), url)

	if err != nil {
		responseTemplate(w, "kesalahan-client", "", "")
		return
	}

	json.NewDecoder(resp.Body).Decode(del)

	if del.StatusServer != "OK" {
		responseTemplate(w, "kesalahan-server", "", "")
		return
	}

	responseTemplate(w, "OK", "", "")
}

func confEditEntri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/confirmedit"
	ubah := ConvertToUbah(r)
	resp, err := sendPost(ubah, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "kesalahan-server", "", "")
		return
	}
	res := &Pasien{}
	json.NewDecoder(resp.Body).Decode(res)

	if res.StatusServer != "OK" {
		responseTemplate(w, res.StatusServer, "", GenModal("Peringatan", res.NoCM, ""))
		return
	}
	fmt.Print(GenTemplate(res, "baristabel"))
	fmt.Print(res)
	responseTemplate(w, "OK", GenTemplate(res, "baristabel"), GenModal("Sukses", "Data berhasil diubah", ""))
}

func deleteEntri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}
	fmt.Print(GenModal("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus"))
	responseTemplate(w, "OK", "", GenModal("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus"))

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
		responseTemplate(w, "kesalahan-server", "", "")
	}
	kun := &Pasien{}
	json.NewDecoder(resp.Body).Decode(kun)

	b := new(bytes.Buffer)
	tmp := template.Must(template.New("modedit.html").ParseFiles("templates/modedit.html"))
	err = tmp.Execute(b, nil)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "", "")
		fmt.Print(err)
	}

	mod := &ModalTemplate{
		Script:  b.String(),
		Content: kun,
	}
	// fmt.Print(mod.Content)
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
		responseTemplate(w, "kesalahan-server", "", "")
	}

	pts := &Pasien{}

	json.NewDecoder(resp.Body).Decode(pts)
	if pts.NoCM == "kesalahan-database" {
		responseTemplate(w, "kesalahan-database", "", "")
	}
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("baristabel.html").ParseFiles("templates/baristabel.html"))
	err = tmp.Execute(b, pts)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "", "")
	}
	responseTemplate(w, "OK", b.String(), "")

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
		responseTemplate(w, "kesalahan-server", "", "")
	}

	json.NewDecoder(resp.Body).Decode(pts)
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("inputpts.html").ParseFiles("templates/inputpts.html"))
	err = tmp.Execute(b, pts)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "", "")
	}
	responseTemplate(w, "OK", b.String(), "")

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

func responseTemplate(w http.ResponseWriter, token, script, modal string) {
	res := &Response{
		Token:  token,
		Script: script,
		Modal:  modal,
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&res)

	if err != nil {
		log.Print(err)
	}
}

func GenModal(title, msg, butname string) string {
	b := new(bytes.Buffer)
	modal := map[string]string{
		"title":  title,
		"msg":    msg,
		"button": butname,
	}

	tmp := template.Must(template.New("modalpopup.html").ParseFiles("templates/modalpopup.html"))
	err := tmp.Execute(b, modal)
	if err != nil {
		fmt.Print(err)
		return ""
	}

	return b.String()
}

func GenTemplate(n interface{}, temp ...string) string {
	b := new(bytes.Buffer)
	funcs := template.FuncMap{"inc": func(i int) int {
		return i + 1
	},
	}

	tmpl := template.New("")

	for k, v := range temp {
		if k == 0 {
			tmp := template.Must(template.New(v + ".html").Funcs(funcs).ParseFiles("templates/" + v + ".html"))
			tmpl = tmp

		}
	}

	for k, v := range temp {
		if k != 0 {
			temp, err := template.Must(tmpl.Clone()).ParseFiles("templates/" + v + ".html")
			if err != nil {
				fmt.Print(err)
				break
			}
			tmpl = temp
		}
	}
	// template.Must(template.New(temp + ".html").Funcs(funcs).ParseFiles("templates/" + temp + ".html"))
	err := tmpl.Execute(b, n)
	if err != nil {
		fmt.Print(err)
		return ""
	}

	return b.String()
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

	// var b bytes.Buffer
	// tmp := template.Must(template.New("main.html").ParseFiles("templates/main.html", "templates/input.html", "templates/content.html"))
	// err = tmp.Execute(&b, web)
	// if err != nil {
	// 	fmt.Print(err)
	// }

	responseTemplate(w, web.Token, GenTemplate(web, "main", "input", "content"), "")

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
