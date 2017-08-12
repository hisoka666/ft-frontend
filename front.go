package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// type ReturnError struct {
// 	Error string `json:"error"`
// }

type MainView struct {
	Token  string    `json:"token"`
	User   string    `json:"user"`
	Bulan  []string  `json:"bulan"`
	Pasien []Pasien  `json:"pasien"`
	IKI    []ListIKI `json:"list"`
}

type ListIKI struct {
	Tanggal int `json:"tgl"`
	SumIKI1 int `json:"iki1"`
	SumIKI2 int `json:"iki2"`
}

type Pasien struct {
	StatusServer string    `json:"stat"`
	TglKunjungan string    `json:"tgl"`
	ShiftJaga    string    `json:"shift"`
	ATS          string    `json:"ats"`
	Dept         string    `json:"dept"`
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
	Token  string      `json:"token"`
	Script string      `json:"script"`
	Modal  string      `json:"modal"`
	Data   interface{} `json:"data"`
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
	http.HandleFunc("/getmonthly", getMonthly)
	http.HandleFunc("/getbcpmonth", getBCPMonth)
	http.HandleFunc("/getpdf", getPDF)
	http.HandleFunc("/getpdfnow", getPDFNow)
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
		Dept:       r.FormValue("bagian"),
		LinkID:     r.FormValue("link"),
		ShiftJaga:  r.FormValue("shift"),
	}

	return n
}

func getPDFNow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "http://2.igdsanglah.appspot.com/getbulanini"

	gettgl := r.FormValue("tgl")

	send := &MainView{
		User:  r.FormValue("email"),
		Bulan: []string{gettgl},
	}
	resp, err := sendPost(send, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	pts := []Pasien{}
	json.NewDecoder(resp.Body).Decode(&pts)
	defer resp.Body.Close()
	iki := countIKI(pts)
	// jaga := dataJaga(perBagian(pts), countIKI(pts))

	createPDF(w, pts, iki, gettgl, r.FormValue("nama"))
}

func getPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	// fmt.Print("Request masuk")
	url := "http://2.igdsanglah.appspot.com/getbulan"

	gettgl := r.FormValue("tgl")
	fmt.Print(gettgl)
	send := &MainView{
		User:  r.FormValue("email"),
		Bulan: []string{gettgl},
	}
	resp, err := sendPost(send, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	pts := []Pasien{}
	json.NewDecoder(resp.Body).Decode(&pts)
	defer resp.Body.Close()

	createPDF(w, pts, countIKI(pts), gettgl, r.FormValue("nama"))
}
func countTotalIKI(l []ListIKI) (int, int, int, int, float32, float32, float32) {
	var a, b, c, d int

	for k, v := range l {
		switch {
		case k < 16:
			a = a + v.SumIKI1
			b = b + v.SumIKI2
		case k >= 16:
			c = c + v.SumIKI1
			d = d + v.SumIKI2
		}
	}

	e := float32(a+c)*0.0032 + float32(b+d)*0.01
	// m := a.(float32)
	// n := c.(float32)
	// // n := b.(float32) + d.(float32)
	// // e = (a.(float32)+c.(float32))*g + (b.(float32)+d.(float32))*h
	return a, b, a + c, b + d, float32(a+c) * 0.0032, float32(b+d) * 0.01, e
}
func createPDF(w http.ResponseWriter, p []Pasien, l []ListIKI, tgl, email string) {
	a, b, c, d, e, f, g := countTotalIKI(l)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	// Tabel IKI //////////////////////////////////////////////////////////
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 210, Ht: 297})
	pdf.Cell(160, 6, "Bukti Kegiatan Harian")
	pdf.Cell(120, 6, ("Nama Pegawai: " + email))
	pdf.Ln(-1)
	pdf.Cell(160, 6, "Pegawai RSUP Sanglah Denpasar")
	pdf.Cell(120, 6, "NIP/Gol: ")
	pdf.Ln(-1)
	pdf.Cell(160, 6, ("Bulan: " + tgl))
	pdf.Cell(120, 6, "Tempat Tugas: IGD RSUP Sanglah")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	// range list iki

	for i := 1; i < 17; i++ {
		pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.SetXY(266, 28)
	pdf.CellFormat(25, 20, "Jumlah Poin", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 48)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for k, v := range l {
		if k < 16 {
			pdf.CellFormat(11, 24, strconv.Itoa(v.SumIKI1), "1", 0, "C", false, 0, "")
		}
	}
	// for i := 1; i < 17; i++ {
	// 	pdf.CellFormat(11, 24, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(25, 24, strconv.Itoa(a), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 72)
	pdf.CellFormat(20, 12, "0,01", "1", 0, "C", false, 0, "")
	for k, v := range l {
		if k < 16 {
			pdf.CellFormat(11, 12, strconv.Itoa(v.SumIKI2), "1", 0, "C", false, 0, "")
		}
	}
	// for i := 1; i < 17; i++ {
	// 	pdf.CellFormat(11, 12, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(25, 12, strconv.Itoa(b), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	// Baris ke dua
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	for i := 17; i < 32; i++ {
		pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.SetFont("Arial", "B", 7)
	pdf.MultiCell(11, 5, "Jumlah Poin", "1", "C", false)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(266, 96)
	pdf.MultiCell(25, 20, "Jumlah X Poin", "1", "C", false)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 116)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for k, v := range l {
		if k >= 16 {
			pdf.CellFormat(11, 24, strconv.Itoa(v.SumIKI1), "1", 0, "C", false, 0, "")
		}
	}
	// for i := 17; i <= 32; i++ {
	// 	pdf.CellFormat(11, 24, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(11, 24, strconv.Itoa(c), "1", 0, "C", false, 0, "")

	pdf.CellFormat(25, 24, fmt.Sprintf("%.4f", e), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 140)
	pdf.CellFormat(20, 12, "0,01", "1", 0, "C", false, 0, "")
	for k, v := range l {
		if k >= 16 {
			pdf.CellFormat(11, 12, strconv.Itoa(v.SumIKI2), "1", 0, "C", false, 0, "")
		}
	}
	// for i := 17; i <= 32; i++ {
	// 	pdf.CellFormat(11, 12, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.CellFormat(11, 12, strconv.Itoa(d), "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 12, fmt.Sprintf("%.4f", f), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(256, 6, "Jumlah Point X Volume kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, fmt.Sprintf("%.4f", g), "1", 1, "C", false, 0, "")
	pdf.CellFormat(256, 6, "Target Point kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, "1,111", "1", 1, "C", false, 0, "")
	////////////////// Buku Catatan Pasien ///////////////////////////////
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	wd := pdf.GetStringWidth("Buku Catatan Pribadi")
	pdf.SetX((210 - wd) / 2)
	pdf.Cell(wd, 9, "Buku Catatan Pribadi")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(20, 5, "Nama")
	pdf.Cell(105, 5, (": " + email))
	pdf.Ln(-1)
	pdf.Cell(20, 5, "Bulan")
	pdf.Cell(105, 5, (": " + tgl))
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(9, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(18, 20, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(17, 20, "No. CM", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 20, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 20, "Diagnosis", "1", 0, "C", false, 0, "")

	pdf.MultiCell(20, 5, "Melakukan pelayanan medik umum", "1", "C", false)

	pdf.SetXY(174, 35)
	pdf.MultiCell(25, 4, "Melakukan tindakan medik umum tingkat sederhana", "1", "C", false)

	for k, v := range p {
		pdf.SetFont("Arial", "", 8)
		diag := ProperCapital(v.Diagnosis)
		if len(diag) > 20 {
			diag = diag[:21]
		}
		// 11/02/1987
		tang := v.TglKunjungan[:10]
		num := strconv.Itoa(k + 1)
		nocm := v.NoCM
		nam := ProperCapital(v.NamaPasien)
		if len(nam) > 25 {
			nam = nam[:26]
		}
		pdf.CellFormat(9, 7, num, "1", 0, "C", false, 0, "")
		pdf.CellFormat(18, 7, tang, "1", 0, "C", false, 0, "")
		pdf.CellFormat(17, 7, nocm, "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 7, nam, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, diag, "1", 0, "L", false, 0, "")
		pdf.SetFont("ZapfDingbats", "", 8)
		if v.IKI == "1" {
			pdf.CellFormat(20, 7, "4", "1", 0, "C", false, 0, "")
			pdf.CellFormat(25, 7, "", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		} else {
			pdf.CellFormat(20, 7, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(25, 7, "4", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}
	}

	t := new(bytes.Buffer)
	err := pdf.Output(t)
	if err != nil {
		log.Fatalf("Error reading pdf %v", err)
	}

	w.Header().Set("Content-type", "application/pdf")
	if _, err := t.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

}

func getBCPMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "http://2.igdsanglah.appspot.com/getbulan"

	gettgl := r.FormValue("tgl")

	send := &MainView{
		User:  r.FormValue("email"),
		Bulan: []string{gettgl},
	}
	resp, err := sendPost(send, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	pts := []Pasien{}
	json.NewDecoder(resp.Body).Decode(&pts)
	defer resp.Body.Close()
	jaga := dataJaga(perBagian(pts), countIKI(pts))
	// bag := perBagian(pts)
	// bl := countDaysOfMonth(r.FormValue("year"), r.FormValue("month"))
	iki := countIKI(pts)
	// fmt.Printf("LIst iki adalah: %v", iki)

	responseTemplate(w, "OK", GenTemplate(pts, "contentrefresh"), GenTemplate(iki, "tabeliki"), jaga)

}

func dataJaga(m ...interface{}) interface{} {
	j := make(map[string]interface{})
	for k, v := range m {
		keymap := "data" + strconv.Itoa(k)
		j[keymap] = v
	}

	return j
}
func getMonthly(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	url := "http://2.igdsanglah.appspot.com/getbulanini"
	// month di sini dikirim dalam bentuk int satu digit
	month, _ := strconv.Atoi(r.FormValue("month"))
	// diubah menjadi 2 digit untuk bisa mengambil kursor
	strmon := fmt.Sprintf("%02d", month)
	gettgl := r.FormValue("year") + "/" + strmon

	send := &MainView{
		User:  r.FormValue("email"),
		Bulan: []string{gettgl},
	}
	resp, err := sendPost(send, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	pts := []Pasien{}
	json.NewDecoder(resp.Body).Decode(&pts)
	defer resp.Body.Close()
	iki := countIKI(pts)
	jaga := dataJaga(perBagian(pts), countIKI(pts))
	// fmt.Printf("LIst iki adalah: %v", iki)

	responseTemplate(w, "OK", GenTemplate(pts, "contentrefresh"), GenTemplate(iki, "tabeliki"), jaga)

}

func countIKI(n []Pasien) []ListIKI {

	g := []ListIKI{}

	for h := 1; h <= 31; h++ {
		var u1, u2 int
		for _, v := range n {
			tgl, _ := strconv.Atoi(v.TglKunjungan[:2])
			if tgl != h {
				continue
			}
			if v.IKI == "1" {
				u1++
			} else {
				u2++
			}
		}

		f := ListIKI{
			Tanggal: h,
			SumIKI1: u1,
			SumIKI2: u2,
		}

		g = append(g, f)

	}
	return g
}

func countDaysOfMonth(y, m string) int {
	yr, _ := strconv.Atoi(y)
	mo, _ := strconv.Atoi(m)

	return time.Date(yr, time.Month(mo), 0, 0, 0, 0, 0, time.UTC).Day()
}
func perBagian(n []Pasien) map[string]int {
	var interna, bedah, anak, obgyn, saraf, anes, psik, tht, kulit, jant, um, mata, mod int
	for _, v := range n {
		switch v.Dept {
		case "1":
			interna++
		case "2":
			bedah++
		case "3":
			anak++
		case "4":
			obgyn++
		case "5":
			saraf++
		case "6":
			anes++
		case "7":
			psik++
		case "8":
			tht++
		case "9":
			kulit++
		case "10":
			jant++
		case "11":
			um++
		case "12":
			mata++
		case "13":
			mod++
		}
	}

	m := make(map[string]int)
	m["interna"] = interna
	m["bedah"] = bedah
	m["anak"] = anak
	m["obgyn"] = obgyn
	m["saraf"] = saraf
	m["anes"] = anes
	m["psik"] = psik
	m["tht"] = tht
	m["kulit"] = kulit
	m["jant"] = jant
	m["umum"] = um
	m["mata"] = mata
	m["mod"] = mod
	return m
}

type InputObat struct {
	MerkDagang     string            `json:"merk`
	Kandungan      string            `json:"kand"`
	MinDose        string            `json:"mindose"`
	MaxDose        string            `json:"maxdose"`
	Tablet         map[string]string `json:"tab"`
	Sirup          map[string]string `json:"syr"`
	Drop           map[string]string `json:"drop"`
	Lainnya        string            `json:"lainnya"`
	SediaanLainnya map[string]string `json:"lainnya_sediaan"`
	Rekomendasi    string            `json:"rekom"`
	Dokter         string            `json:"doc"`
}

func inputObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	var dat map[string]interface{}
	// fmt.Println(r.FormValue("merk"))
	// fmt.Println(r.FormValue("kand"))
	// fmt.Println(r.FormValue("mindose"))
	// fmt.Println(r.FormValue("maxdose"))
	// // x := r.FormValue("tab")
	// // fmt.Println(x[1:(len(x) - 1)])
	// // m := []byte("{" + x[1:(len(x)-1)] + "}")
	// fmt.Println(r.FormValue("send"))
	m := []byte(r.FormValue("send"))
	// // fmt.Printf("Hasil ubah byte: %v", m)
	err := json.Unmarshal(m, &dat)
	if err != nil {
		fmt.Printf("Gagal mengubah json: %v", err)
	}

	// var blet []interface{}
	tab := dat["tab"].([]interface{})
	fmt.Println(tab)

	// for k, v := range tab {
	// 	// fmt.Println(k)
	// 	// fmt.Println(v)
	// 	c := v.(map[string]interface{})
	// 	l := strconv.Itoa(k)
	// 	fmt.Println(c[l])
	// 	// b := v.(map[string]interface{})
	// 	// blet[k] = b["0"]
	// }

	// fmt.Println(blet)
	// fmt.Println(r.FormValue("tab"))
	// fmt.Println(r.FormValue("syr"))
	// fmt.Println(r.FormValue("drop"))
	// fmt.Println(r.FormValue("lainnya_sediaan"))

	// if err := json.Unmarshal(r.FormValue("tab"), &dat); err == nil {
	// 	for _, v := range dat {
	// 		fmt.Println(v)
	// 	}
	// }
	// if err := json.Unmarshal(r.FormValue("syr"), &dat); err == nil {
	// 	for _, v := range dat {
	// 		fmt.Println(v)
	// 	}
	// }
	// if err := json.Unmarshal(r.FormValue("drop"), &dat); err == nil {
	// 	for _, v := range dat {
	// 		fmt.Println(v)
	// 	}
	// }
	// if err := json.Unmarshal(r.FormValue("lainnya_sediaan"), &dat); err == nil {
	// 	for _, v := range dat {
	// 		fmt.Println(v)
	// 	}
	// }
	// fmt.Println(r.FormValue("lainnya"))
	// fmt.Println(r.FormValue("rekom"))
	// fmt.Println(r.FormValue("doc"))
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

func getInputObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}

	responseTemplate(w, "OK", GenTemplate(nil, "modinputobatbaru"), "", nil)
}

func getPresPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}

	responseTemplate(w, "OK", GenTemplate(nil, "modlistresep"), "", nil)
}

func getPtsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}
	responseTemplate(w, "OK", GenTemplate(nil, "modresep"), "", nil)
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
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""), nil)
		log.Fatalf("Terjadi kesalahan pengiriman ke server")
	}

	list := MainView{}
	json.NewDecoder(resp.Body).Decode(&list)
	defer resp.Body.Close()
	fmt.Printf("Isi dari token adalah: %v", list.Token)
	if list.Token != "OK" {
		log.Fatalf("Terjadi kesalahan server")
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""), nil)
		return
	}

	responseTemplate(w, "OK", GenTemplate(list.Pasien, "contentrefresh"), "", nil)

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
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""), nil)
		log.Print("Terjadi kesalahan server")
		return
	}

	json.NewDecoder(resp.Body).Decode(pts)
	pts.TglKunjungan = pts.TglAsli.Format("Mon 02/01/2006 15:04:05")
	pts.LinkID = link

	script := GenTemplate(pts, "modubahtgl")
	responseTemplate(w, "OK", script, "", nil)
	defer resp.Body.Close()
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
		responseTemplate(w, "kesalahan-client", "", "", nil)
		return
	}

	json.NewDecoder(resp.Body).Decode(send)
	responseTemplate(w, "OK", GenTemplate(send.Pasien, "contentrefresh"), "", nil)
	defer resp.Body.Close()
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
		responseTemplate(w, "kesalahan-client", "", "", nil)
		return
	}

	json.NewDecoder(resp.Body).Decode(del)

	if del.StatusServer != "OK" {
		responseTemplate(w, "kesalahan-server", "", "", nil)
		return
	}
	defer resp.Body.Close()
	responseTemplate(w, "OK", "", "", nil)
}

func confEditEntri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "http://2.igdsanglah.appspot.com/entri/confirmedit"
	ubah := ConvertToUbah(r)
	resp, err := sendPost(ubah, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "kesalahan-server", "", "", nil)
		return
	}
	res := &Pasien{}
	json.NewDecoder(resp.Body).Decode(res)

	if res.StatusServer != "OK" {
		responseTemplate(w, res.StatusServer, "", GenModal("Peringatan", res.NoCM, ""), nil)
		return
	}
	defer resp.Body.Close()
	responseTemplate(w, "OK", GenTemplate(res, "baristabel"), GenModal("Sukses", "Data berhasil diubah", ""), nil)
}

func deleteEntri(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}
	fmt.Print(GenModal("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus"))
	responseTemplate(w, "OK", "", GenModal("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus"), nil)

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
		responseTemplate(w, "kesalahan-server", "", "", nil)
	}
	kun := &Pasien{}
	json.NewDecoder(resp.Body).Decode(kun)
	defer resp.Body.Close()
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("modedit.html").ParseFiles("templates/modedit.html"))
	err = tmp.Execute(b, nil)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "", "", nil)
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
		responseTemplate(w, "kesalahan-server", "", "", nil)
	}

	pts := &Pasien{}

	json.NewDecoder(resp.Body).Decode(pts)
	defer resp.Body.Close()
	if pts.NoCM == "kesalahan-database" {
		responseTemplate(w, "kesalahan-database", "", "", nil)
	}
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("baristabel.html").ParseFiles("templates/baristabel.html"))
	err = tmp.Execute(b, pts)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "", "", nil)
	}
	responseTemplate(w, "OK", b.String(), "", nil)

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
		responseTemplate(w, "kesalahan-server", "", "", nil)
	}

	json.NewDecoder(resp.Body).Decode(pts)
	defer resp.Body.Close()
	b := new(bytes.Buffer)
	tmp := template.Must(template.New("inputpts.html").ParseFiles("templates/inputpts.html"))
	err = tmp.Execute(b, pts)
	if err != nil {
		responseTemplate(w, "kesalahan-template", "", "", nil)
	}
	responseTemplate(w, "OK", b.String(), "", nil)

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

func responseTemplate(w http.ResponseWriter, token, script, modal string, data interface{}) {
	res := &Response{
		Token:  token,
		Script: script,
		Modal:  modal,
		Data:   data,
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
	defer resp.Body.Close()

	responseTemplate(w, web.Token, GenTemplate(web, "main", "input", "content"), "", nil)
}

func ProperCapital(input string) string {
	words := strings.Fields(input)
	smallwords := " dan atau dr. "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}
