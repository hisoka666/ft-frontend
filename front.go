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

type PerShift struct {
	Pagi  int `json:"pagi"`
	Sore  int `json:"sore"`
	Malam int `json:"malam"`
	Total int `json:"total"`
}
type KursorIGD struct {
	Bulan string `json:"bulan"`
	Point string `json:"point"`
}
type SupervisorList struct {
	StatusServer    string                 `json:"status"`
	ListPasien      []SupervisorListPasien `json:"listpasien"`
	Token           string                 `json:"token"`
	SupervisorName  string                 `json:"user"`
	ListBulan       []string               `json:"listbulan"`
	PerHari         []int                  `json:"perhari"`
	PerDeptPerHari  []Departemen           `json:"perdept"`
	PerShiftPerHari []PerShift             `json:"shift"`
}
type Departemen struct {
	Interna   int `json:"interna"`
	Bedah     int `json:"bedah"`
	Anak      int `json:"anak"`
	Obgyn     int `json:"obgyn"`
	Saraf     int `json:"saraf"`
	Anestesi  int `json:"anes"`
	Psikiatri int `json:"psik"`
	THT       int `json:"tht"`
	Kulit     int `json:"kulit"`
	Kardio    int `json:"jant"`
	Umum      int `json:"umum"`
	Mata      int `json:"mata"`
	MOD       int `json:"mod"`
}
type MainView struct {
	Token      string         `json:"token"`
	User       string         `json:"user"`
	Bulan      []string       `json:"bulan"`
	Pasien     []Pasien       `json:"pasien"`
	IKI        []ListIKI      `json:"list"`
	Admin      Admin          `json:"admin"`
	Supervisor SupervisorList `json:"supervisor"`
	Peran      string         `json:"peran"`
	LinkID     string         `json:"link"`
}
type SupervisorListPasien struct {
	TglKunjungan time.Time `json:"tgl"`
	ATS          string    `json:"ats"`
	Dept         string    `json:"dept"`
	Diagnosis    string    `json:"diag"`
	LinkID       string    `json:"link"`
}

type Admin struct {
	Staff []Staff `json:"list"`
	Token string  `json:"token"`
}
type Staff struct {
	Email, NamaLengkap, LinkID, Peran string
}
type ListIKI struct {
	Tanggal   int      `json:"tgl"`
	SumIKI1   int      `json:"iki1"`
	SumIKI2   int      `json:"iki2"`
	P3K       []string `json:"p3k"`
	Rapat     []string `json:"rapat"`
	Pelatihan []string `json:"pelatihan"`
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
	NamaPasien string    `json:"namapts"`
	NomorCM    string    `json:"nocm"`
	JenKel     string    `json:"jenkel"`
	Alamat     string    `json:"alamat"`
	TglDaftar  time.Time `json:"tgldaf"`
	TglLahir   time.Time `json:"tgllhr"`
}
type PasienResep struct {
	Nama      string `json:"nama"`
	Umur      string `json:"umur"`
	Berat     string `json:"berat"`
	Alamat    string `json:"alamat"`
	Alergi    string `json:"alergi"`
	Diagnosis string `json:"diag"`
	NoCM      string `json:"nocm"`
	LinkID    string `json:"link"`
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

type InputObat struct {
	MerkDagang     string   `json:"merk"`
	Kandungan      string   `json:"kand"`
	MinDose        string   `json:"mindose"`
	MaxDose        string   `json:"maxdose"`
	Tablet         []string `json:"tab"`
	Sirup          []string `json:"syr"`
	Drop           []string `json:"drop"`
	Lainnya        string   `json:"lainnya"`
	SediaanLainnya []string `json:"lainnya_sediaan"`
	Rekomendasi    string   `json:"rekom"`
	Dokter         string   `json:"doc"`
}

type ServerResponse struct {
	Error string `json:"error"`
}

type IndexObat struct {
	MerkDagang string `json:"merk"`
	Kandungan  string `json:"kandungan"`
	Link       string `json:"link"`
}
type ObatView struct {
	Rekomendasi string   `json:"rekom"`
	Kemasan     string   `json:"kemasan"`
	Sediaan     []string `json:"sediaan"`
	Dosis       string   `json:"dosis"`
	Satuan      string   `json:"satuan"`
	Link        string   `json:"link"`
}

type Resep struct {
	Dokter    string      `json:"dokter"`
	Tanggal   string      `json:"tanggal"`
	ListObat  []Obat      `json:"listobat"`
	ListPuyer []Puyer     `json:"listpuyer"`
	Pasien    PasienResep `json:"pasien"`
}

type Obat struct {
	NamaObat   string `json:"obat"`
	Jumlah     string `json:"jumlah"`
	Instruksi  string `json:"instruksi"`
	Keterangan string `json:"keterangan"`
}

type Puyer struct {
	Obat       []SatuObat `json:"satuobat"`
	Racikan    string     `json:"racikan"`
	JmlRacikan string     `json:"jml-racikan"`
	Instruksi  string     `json:"instruksi"`
	Keterangan string     `json:"keterangan"`
}

type SatuObat struct {
	NamaObat string `json:"obat"`
	Takaran  string `json:"takaran"`
}

type KegiatanEkstra struct {
	P3K       []string `json:"p3k"`
	Rapat     []string `json:"rapat"`
	Pelatihan []string `json:"pelatihan"`
}

type DetailPasien struct {
	Pasien    DataPasien        `json:"datapts"`
	Kunjungan []KunjunganPasien `json:"kunjungan"`
	LinkID    string            `json:"link"`
}

type DetailStaf struct {
	NamaLengkap  string    `json:"nama"`
	NIP          string    `json:"nip"`
	NPP          string    `json:"npp"`
	GolonganPNS  string    `json:"golpns"`
	Alamat       string    `json:"alamat"`
	Bagian       string    `json:"bagian"`
	LinkID       string    `json:"link"`
	TanggalLahir time.Time `json:"tgl"`
	Umur         string    `json:"umur"`
}

type SuratSakit struct {
	LinkID        string `json:"link"`
	TglLahir      string `json:"tgl"`
	Pekerjaan     string `json:"pekerjaan"`
	Alamat        string `json:"alamat"`
	LamaIstirahat string `json:"lama"`
	StatusData    string `json:"status"`
	Dokter        string `json:"dokter"`
}
type DataSuratSakit struct {
	TglSurat      time.Time `json:"tglsurat"`
	NomorSurat    int       `json:"nomor"`
	LamaIstirahat string    `json:"lama"`
	Pekerjaan     string    `json:"pekerjaan"`
	LinkSurat     string    `json:"link"`
	Dokter        string    `json:"dokter"`
	NamaPasien    string    `json:"namapts"`
	NoCM          string    `json:"nocm"`
	Umur          string    `json:"umur"`
	Alamat        string    `json:"alamat"`
}

type LembarATS struct {
	LinkID         string    `json:"link"`
	KeluhanUtama   string    `json:"kelut"`
	Subyektif      string    `json:"subyektif"`
	TDSistolik     string    `json:"tdsis"`
	TDDiastolik    string    `json:"tddi"`
	Nadi           string    `json:"nadi"`
	LajuPernafasan string    `json:"rr"`
	SuhuBadan      string    `json:"temp"`
	LokasiNyeri    string    `json:"nyerilok"`
	NRS            string    `json:"nrs"`
	Keterangan     string    `json:"keterangan"`
	GCSE           string    `json:"gcse"`
	GCSV           string    `json:"gcsv"`
	GCSM           string    `json:"gcsm"`
	TglInput       time.Time `json:"input"`
	Dokter         string    `json:"dokter"`
}

type RekamMedis struct {
	Pasien    Pasien    `json:"pasien"`
	LembarATS LembarATS `json:"lembarats"`
}
type ListPasienResiden struct {
	List []Pasien `json:"listpasien"`
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
	http.HandleFunc("/cariobt", cariObat)
	http.HandleFunc("/getobat", getObat)
	http.HandleFunc("/tambahdokter", tambahDokter)
	http.HandleFunc("/hapusdokter", hapusDokter)
	http.HandleFunc("/getobatedit", editObat)
	http.HandleFunc("/inputobatedit", confEditObat)
	http.HandleFunc("/formpuyer", formPuyer)
	http.HandleFunc("/cariobatpuyer", cariObatPuyer)
	http.HandleFunc("/getpuyer", getObat)
	http.HandleFunc("/buatresep", buatResep)
	http.HandleFunc("/supgeteachmonth", supGetMonth)
	http.HandleFunc("/supmonthnow", supGetMonthNow)
	http.HandleFunc("/get-detail-pts", getDetailPasien)
	http.HandleFunc("/input-detail-pts", inputDetailPts)
	http.HandleFunc("/buat-resep-pts", buatResepPts)
	http.HandleFunc("/docpage", docPage)
	http.HandleFunc("/simpandoc", simpanDoc)
	http.HandleFunc("/get-surat-sakit-page", getSuratSakit)
	http.HandleFunc("/pdf-surat-sakit", pdfSuratSakit)
	http.HandleFunc("/simpan-lembar-ats", simpanLembarATS)
	http.HandleFunc("/get-rm-kun", getRMKun)
	http.HandleFunc("/get-residen-pasien-list", getResidenPasienList)
	http.HandleFunc("/get-residen-refresh-pasien-list", getResidenRefresh)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8001", nil))

}

func getResidenRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	kun := &KunjunganPasien{
		Dokter: r.FormValue("email"),
	}
	url := "https://residen-dot-igdsanglah.appspot.com/get-refresh"
	send, err := sendPost(kun, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("terjadi kesalahan: %v", err)
	}
	pts := &ListPasienResiden{}
	json.NewDecoder(send.Body).Decode(pts)
	responseTemplate(w, "", GenTemplate(pts.List, "residenpage-new-content"), "", nil)
}
func getResidenPasienList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	// zone, _ := time.LoadLocation("Asia/Makassar")
	dur, _ := time.ParseDuration("3s")
	s := time.Now()
	kun := &KunjunganPasien{
		Dokter:    r.FormValue("email"),
		JamDatang: s.Add(-dur),
	}
	url := "https://residen-dot-igdsanglah.appspot.com/get-new-pasien"
	send, err := sendPost(kun, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("terjadi kesalahan: %v", err)
	}
	pts := &ListPasienResiden{}
	json.NewDecoder(send.Body).Decode(pts)
	defer send.Body.Close()
	fmt.Printf("body adalah: %v", pts)
	// fmt.Print(GenTemplate(pts, "residenpage-new-content"))
	responseTemplate(w, "", GenTemplate(pts.List, "residenpage-new-content"), "", nil)
}
func getRMKun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	ats := &LembarATS{
		LinkID: r.FormValue("link"),
	}
	url := "https://ats-dot-igdsanglah.appspot.com/get-rm-kun"
	send, err := sendPost(ats, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("terjadi kesalahan: %v", err)
	}
	rm := &RekamMedis{}
	json.NewDecoder(send.Body).Decode(rm)
	send.Body.Close()
	responseTemplate(w, "", GenTemplate(rm, "rekam-medik-page"), "", nil)
}
func simpanLembarATS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	ats := &LembarATS{}
	con := []byte(r.FormValue("ats"))
	err := json.Unmarshal(con, ats)
	if err != nil {
		log.Fatalf("terjadi kesalahan: %v", err)
	}
	fmt.Printf("dokter adalah: %v", r.FormValue("dokter"))
	ats.Dokter = r.FormValue("dokter")
	// fmt.Printf("Subyektif adalah: %v", ats.Subyektif)
	url := "https://ats-dot-igdsanglah.appspot.com/simpan"
	resp, err := sendPost(ats, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("terjadi kesalahan: %v", err)
	}
	json.NewDecoder(resp.Body).Decode(ats)
	responseTemplate(w, "", "ok", "", nil)
}
func pdfSuratSakit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	sur := &SuratSakit{}
	// fmt.Print(r.FormValue("content"))
	surat := []byte(r.FormValue("content"))
	err := json.Unmarshal(surat, sur)
	if err != nil {
		log.Fatalf("Terjadi kesalahan dalam unmarshaling json: %v", err)
		return
	}
	sur.StatusData = r.FormValue("status-data")
	sur.Dokter = r.FormValue("dokter")
	url := "https://get-data-pasien-dot-igdsanglah.appspot.com/add-surat-sakit"
	resp, err := sendPost(sur, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server : %v", err)
		return
	}
	dat := &DataSuratSakit{}
	json.NewDecoder(resp.Body).Decode(dat)
	defer resp.Body.Close()
	buatPDFSuratSakit(dat, w)
}
func ubahBulanKeRomawi(s string) string {
	t := ""
	switch s {
	case "01":
		t = "I"
	case "02":
		t = "II"
	case "03":
		t = "III"
	case "04":
		t = "IV"
	case "05":
		t = "V"
	case "06":
		t = "VI"
	case "07":
		t = "VII"
	case "08":
		t = "VIII"
	case "09":
		t = "IX"
	case "10":
		t = "X"
	case "11":
		t = "XI"
	case "12":
		t = "XII"
	}
	return t
}

func buatPDFSuratSakit(n *DataSuratSakit, w http.ResponseWriter) {
	zone, _ := time.LoadLocation("Asia/Makassar")
	lama, _ := strconv.Atoi(n.LamaIstirahat)
	tglnow := time.Now()
	tglakhir := tglnow.AddDate(0, 0, lama-1)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 16)
	pdf.AddPageFormat("P", gofpdf.SizeType{Wd: 210, Ht: 145})
	pdf.CellFormat(0, 7, "DEPARTEMEN KESEHATAN RI", "", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(0, 7, "RSUP SANGLAH DENPASAR", "B", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	surat := "SURAT KETERANGAN SAKIT"
	pdf.SetFont("Arial", "", 12)
	wd := pdf.GetStringWidth(surat) + 6
	pdf.SetX((210 - wd) / 2)
	pdf.CellFormat((wd + 6), 10, surat, "B", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	fmt.Print(n.NomorSurat)
	nomor := "Nomor: RSUP 01/" + strconv.Itoa(n.NomorSurat) + "/" + ubahBulanKeRomawi(tglnow.In(zone).Format("01")) + "/2017"
	wd = pdf.GetStringWidth(nomor) + 6
	pdf.SetX((210 - wd) / 2)
	pdf.CellFormat((wd + 6), 9, nomor, "", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, "Yang bertandatangan di bawah ini menerangkan bahwa: ")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.Cell(40, 5, "Nama")
	pdf.Cell(70, 5, ": "+ProperCapital(n.NamaPasien))
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.Cell(40, 5, "Umur")
	pdf.Cell(70, 5, ": "+n.Umur+" tahun")
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.Cell(40, 5, "Pekerjaan")
	pdf.Cell(70, 5, ": "+ProperCapital(n.Pekerjaan))
	pdf.Ln(-1)
	pdf.SetX(30)
	pdf.Cell(40, 5, "Alamat")
	pdf.Cell(70, 5, ": "+ProperCapital(n.Alamat))
	pdf.Ln(-1)
	pdf.Cell(0, 5, "oleh karena SAKIT, perlu diberikan ISTIRAHAT")
	pdf.Ln(-1)

	pdf.Cell(0, 5, "selama "+n.LamaIstirahat+" hari, terhitung mulai tanggal "+tglnow.In(zone).Format("02/01/2006")+" s.d. "+tglakhir.In(zone).Format("02/01/2006"))
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Cell(0, 5, "Demikian surat keterangan ini dibuat dengan sebenarnya dan untuk dipergunakan semestinya")
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetX(160)
	pdf.CellFormat(0, 5, "Denpasar, "+tglnow.Format("02/01/2006"), "", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetX(160)
	pdf.CellFormat(0, 5, "("+n.Dokter+")", "", 0, "C", false, 0, "")
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
func getSuratSakit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	pts := &Pasien{
		LinkID: r.FormValue("link"),
	}
	url := "https://get-data-pasien-dot-igdsanglah.appspot.com/get-detail"
	resp, err := sendPost(pts, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}

	det := &DataPasien{}
	json.NewDecoder(resp.Body).Decode(det)
	defer resp.Body.Close()
	// fmt.Print(GenTemplate(nil, "modsuratsakit"))

	responseTemplate(w, "", GenTemplate(det, "modsuratsakit"), "", nil)
}
func simpanDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	var nip string
	var npp string
	i := &nip
	p := &npp
	if r.FormValue("gol") != "" {
		*i = r.FormValue("nopeg")
	} else {
		*p = r.FormValue("nopeg")
	}

	det := &DetailStaf{
		NamaLengkap: r.FormValue("nama"),
		NIP:         nip,
		NPP:         npp,
		GolonganPNS: r.FormValue("gol"),
		Bagian:      r.FormValue("docbag"),
		LinkID:      r.FormValue("link"),
	}
	url := "https://dokter-dot-igdsanglah.appspot.com/simpan-dokter"
	send, err := sendPost(det, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	json.NewDecoder(send.Body).Decode(det)
	defer send.Body.Close()
	responseTemplate(w, "", GenTemplate(det, "docpage"), "", nil)
}
func docPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	staf := &Staff{
		LinkID: r.FormValue("link"),
	}
	url := "https://dokter-dot-igdsanglah.appspot.com"

	send, err := sendPost(staf, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	det := &DetailStaf{}
	json.NewDecoder(send.Body).Decode(det)
	defer send.Body.Close()
	responseTemplate(w, "", GenTemplate(det, "docpage"), "", nil)
}

func buatResepPts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	url := "https://pasien-dot-igdsanglah.appspot.com"
	pts := &Pasien{
		LinkID: r.FormValue("link"),
	}
	resp, err := sendPost(pts, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	det := &DetailPasien{}
	json.NewDecoder(resp.Body).Decode(det)
	defer resp.Body.Close()
	// fmt.Print(GenTemplate(det, "modresepwithdata"))
	responseTemplate(w, "", GenTemplate(det, "modresepwithdata"), "", nil)
}

func inputDetailPts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	url := "https://input-detail-pts-dot-igdsanglah.appspot.com/"
	tgl, err := time.Parse("02-01-2006", r.FormValue("tgl"))
	if err != nil {
		log.Fatalf("Kesalahan parsing tgl lahir: %v", err)
	}
	data := DataPasien{
		NamaPasien: r.FormValue("nama"),
		TglLahir:   tgl,
		Alamat:     r.FormValue("almt"),
		JenKel:     r.FormValue("jenkel"),
	}
	det := &DetailPasien{
		Pasien: data,
		LinkID: r.FormValue("link"),
	}

	resp, err := sendPost(det, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}

	json.NewDecoder(resp.Body).Decode(det)
	defer resp.Body.Close()
	// fmt.Print(GenTemplate(det, "detailpasienonly"))
	responseTemplate(w, "", GenTemplate(det, "detailpasienonly"), "", nil)
}
func getDetailPasien(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	pts := &Pasien{
		LinkID: r.FormValue("link"),
	}
	url := "https://pasien-dot-igdsanglah.appspot.com"
	resp, err := sendPost(pts, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}

	det := &DetailPasien{}
	json.NewDecoder(resp.Body).Decode(det)
	defer resp.Body.Close()
	// fmt.Print(det.Kunjungan[0].JamDatang)
	// fmt.Print(GenTemplate(det, "detailpasien"))
	responseTemplate(w, "OK", GenTemplate(det, "detailpasien"), "", nil)

}
func supGetMonthNow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://igdsanglah.appspot.com/getsupmonthnow"
	resp, err := sendPost(nil, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	list := &MainView{}
	json.NewDecoder(resp.Body).Decode(list)
	defer resp.Body.Close()
	responseTemplate(w, "OK", GenTemplate(list.Supervisor, "contentsupervisor"), "", nil)
}
func supGetMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://igdsanglah.appspot.com/getsupmonth"
	send := &KursorIGD{
		Bulan: r.FormValue("tgl"),
	}

	resp, err := sendPost(send, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	list := &SupervisorList{}
	json.NewDecoder(resp.Body).Decode(list)
	resp.Body.Close()
	responseTemplate(w, "", GenTemplate(list, "contentsupervisor"), "", nil)
}

func buatResep(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	pts := &PasienResep{}
	tab := []Obat{}
	pyr := []Puyer{}
	err := json.Unmarshal([]byte(r.FormValue("pts")), pts)
	if err != nil {
		fmt.Printf("Gagal membaca pts: %v", err)
		return
	}
	err = json.Unmarshal([]byte(r.FormValue("puyer")), &pyr)
	if err != nil {
		fmt.Printf("Gagal membaca pyr: %v", err)
		return
	}
	err = json.Unmarshal([]byte(r.FormValue("tablet")), &tab)
	if err != nil {
		fmt.Printf("Gagal membaca pts: %v", err)
		return
	}
	rec := &Resep{
		Dokter:    r.FormValue("dokter"),
		ListObat:  tab,
		ListPuyer: pyr,
		Pasien:    *pts,
	}
	zone, _ := time.LoadLocation("Asia/Makassar")
	rec.Tanggal = time.Now().In(zone).Format("02/01/2006")
	// fmt.Println(pts)
	fmt.Println(rec.Pasien.Nama)
	fmt.Println(rec.Pasien.NoCM)
	fmt.Println(rec.ListObat[0].NamaObat)
	if rec.Pasien.NoCM == "" {
		pdfResep(w, *rec)
	} else {
		url := "https://add-obat-pasien-dot-igdsanglah.appspot.com"
		_, err := sendPost(rec, r.FormValue("token"), url)
		if err != nil {
			log.Fatalf("Terjadi kesalahan di server: %v", err)
		}
		pdfResep(w, *rec)
	}
}

func pdfResep(w http.ResponseWriter, r Resep) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 9)
	pdf.AddPageFormat("P", gofpdf.SizeType{Wd: 100, Ht: 210})
	pdf.Cell(20, 5, "Nama Dokter")
	//Nama Dokter
	pdf.Cell(60, 5, (": dr. " + r.Dokter))
	pdf.Ln(-1)
	pdf.Cell(20, 5, "Ruangan")
	pdf.Cell(60, 5, ": Fasttrack")
	pdf.Ln(-1)
	pdf.CellFormat(20, 5, "Diagnosis", "B", 0, "L", false, 0, "")
	//Diagnosis
	pdf.CellFormat(60, 5, (": " + ProperCapital(r.Pasien.Diagnosis)), "B", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetX(60)
	//Tanggal
	pdf.Cell(20, 5, ("Tanggal: " + r.Tanggal))
	//Nama Obat
	pdf.Ln(-1)
	for _, v := range r.ListObat {
		pdf.Cell(70, 5, ("Rx. " + ProperCapital(v.NamaObat)))
		//Jumlah obat
		pdf.Cell(10, 5, ("No. " + v.Jumlah))
		pdf.Ln(-1)
		//Aturan minum
		pdf.Cell(20, 5, "")
		pdf.CellFormat(40, 5, ("S " + v.Instruksi), "B", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.Ln(-1)
	}
	if len(r.ListPuyer) > 0 {
		pdf.Cell(5, 5, "Rx.")
		pdf.Ln(-1)
		for _, v := range r.ListPuyer {
			for _, n := range v.Obat {
				pdf.Cell(5, 5, "")
				pdf.Cell(60, 5, ProperCapital(n.NamaObat))
				pdf.Cell(10, 5, n.Takaran)
				pdf.Ln(-1)
			}
			pdf.Cell(20, 5, "")
			pdf.Cell(40, 5, v.Racikan)
			pdf.Cell(10, 5, ("No. " + v.JmlRacikan))
			pdf.Ln(-1)
			pdf.Cell(20, 5, "")
			pdf.CellFormat(40, 5, ("S " + v.Instruksi), "B", 0, "", false, 0, "")
			pdf.Ln(-1)
			pdf.Ln(-1)
		}
	}

	// pdf.Cell(70, 5, "Rx. Paracetamol tab 500 mg ")
	// //Jumlah obat
	// pdf.Cell(10, 5, "No. X")
	// pdf.Ln(-1)
	// //Aturan minum
	// pdf.Cell(20, 5, "")
	// pdf.CellFormat(40, 5, "S 3 dd tab 1", "B", 0, "C", false, 0, "")
	// pdf.Ln(-1)
	//Identitas Pasien
	pdf.SetY(170)
	pdf.Cell(20, 5, "Pro")
	pdf.Cell(60, 5, (": " + ProperCapital(r.Pasien.Nama)))
	pdf.Ln(-1)
	pdf.Cell(20, 5, "No. CM")
	pdf.Cell(20, 5, (": " + r.Pasien.NoCM))
	pdf.Cell(20, 5, "Umur")
	pdf.Cell(20, 5, (": " + r.Pasien.Umur + " th"))
	pdf.Ln(-1)
	pdf.Cell(20, 5, "Alamat")
	pdf.Cell(20, 5, (": " + ProperCapital(r.Pasien.Alamat)))
	pdf.Cell(20, 5, "Berat Badan")
	pdf.Cell(20, 5, (": "+r.Pasien.Berat)+" kg")

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
func formPuyer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	responseTemplate(w, "", GenTemplate(nil, "formpuyer"), "", nil)
}

func cariObatPuyer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://get-obat-dot-igdsanglah.appspot.com"
	obat := r.FormValue("obat")
	pos := &IndexObat{
		MerkDagang: obat,
	}

	resp, err := sendPost(pos, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	listobt := []IndexObat{}
	json.NewDecoder(resp.Body).Decode(&listobt)
	defer resp.Body.Close()
	responseTemplate(w, "OK", GenTemplate(listobt, "listpuyer"), pos.MerkDagang, nil)
}

func confEditObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	var dat InputObat
	m := []byte(r.FormValue("send"))
	err := json.Unmarshal(m, &dat)
	if err != nil {
		fmt.Printf("Gagal mengubah json: %v", err)
	}
	// fmt.Printf("Json adlaha: %v", dat)

	url := "https://input-edit-obat-dot-igdsanglah.appspot.com/" + r.FormValue("link")
	resp, err := sendPost(dat, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""), nil)
		log.Print("Terjadi kesalahan server")
		return
	}
	n := ServerResponse{}
	json.NewDecoder(resp.Body).Decode(&n)
	if n.Error != "" {
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Gagal menyimpan ke datastore. Ulangi lagi menginput data", ""), nil)
	}
	responseTemplate(w, "OK", "", "", nil)
}

func editObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://link-obat-dot-igdsanglah.appspot.com"
	link := r.FormValue("link")
	pos := &IndexObat{
		Link: link,
	}
	resp, err := sendPost(pos, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}

	obt := &InputObat{}
	json.NewDecoder(resp.Body).Decode(obt)
	defer resp.Body.Close()
	responseTemplate(w, "OK", GenTemplate(nil, "modinputobatedit"), link, obt)
}
func hapusDokter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://hapus-dokter-dot-igdsanglah.appspot.com"
	st := &Staff{
		LinkID: r.FormValue("link"),
	}
	resp, err := sendPost(st, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "not-ok", fmt.Sprintf("Kesalahan pada server %v", err), "", nil)
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	json.NewDecoder(resp.Body).Decode(st)
	defer resp.Body.Close()
	if st.Email != "OK" {
		responseTemplate(w, "not-ok", "", GenModal("Gagal", st.Email, ""), nil)
	} else {
		responseTemplate(w, "OK", "", GenModal("Sukses", "Berhasil menghapus data", ""), nil)
	}

}
func tambahDokter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	url := "https://tambah-dokter-dot-igdsanglah.appspot.com"
	st := &Staff{
		NamaLengkap: r.FormValue("nama"),
		Email:       r.FormValue("email"),
		Peran:       r.FormValue("peran"),
	}
	resp, err := sendPost(st, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "not-ok", fmt.Sprintf("Kesalahan pada server %v", err), "", nil)
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	json.NewDecoder(resp.Body).Decode(st)
	defer resp.Body.Close()
	responseTemplate(w, "OK", GenTemplate(st, "dokrow"), GenModal("Sukses", "Berhasil menambahkan data", ""), nil)
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

func getObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://link-obat-dot-igdsanglah.appspot.com"
	link := r.FormValue("link")
	pos := &IndexObat{
		Link: link,
	}
	resp, err := sendPost(pos, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}

	obt := &InputObat{}
	json.NewDecoder(resp.Body).Decode(obt)
	defer resp.Body.Close()
	// fmt.Println(obt.MaxDose)
	// fmt.Println(obt.MinDose)
	maxDo, _ := strconv.ParseFloat(obt.MaxDose, 32)
	minDo, _ := strconv.ParseFloat(obt.MinDose, 32)
	bb, _ := strconv.ParseFloat(r.FormValue("berat"), 32)
	maxD := strconv.FormatFloat((maxDo * bb), 'f', 2, 32)
	minD := strconv.FormatFloat((minDo * bb), 'f', 2, 32)
	view := &ObatView{
		Rekomendasi: obt.Rekomendasi,
		Link:        link,
	}
	if obt.Lainnya != "" {
		view.Sediaan = obt.SediaanLainnya
		view.Kemasan = obt.Lainnya
		view.Dosis = minD + " - " + maxD + " mg tiap kali pemberian /(" + obt.MinDose + " - " + obt.MaxDose + ") perKGBB/kali pemberian"
	} else if obt.Sirup[0] != "" {
		view.Sediaan = obt.Sirup
		view.Kemasan = "sirup"
		view.Dosis = minD + " - " + maxD + " mg tiap kali pemberian /(" + obt.MinDose + " - " + obt.MaxDose + ") perKGBB/kali pemberian"
		view.Satuan = "mg per 5 ml"
	} else if obt.Drop[0] != "" {
		view.Sediaan = obt.Drop
		view.Kemasan = "drop"
		view.Dosis = minD + " - " + maxD + " mg tiap kali pemberian /(" + obt.MinDose + " - " + obt.MaxDose + ") perKGBB/kali pemberian"
		view.Satuan = "mg per 1 ml"
	} else {
		view.Sediaan = obt.Tablet
		view.Kemasan = "tablet"
		view.Dosis = minD + " - " + maxD + " mg tiap kali pemberian /(" + obt.MinDose + " - " + obt.MaxDose + ") perKGBB/kali pemberian"
		view.Satuan = "mg"
	}
	responseTemplate(w, "OK", GenTemplate(view, "viewobatbaru"), obt.MerkDagang, view)

}
func cariObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://get-obat-dot-igdsanglah.appspot.com"
	obat := r.FormValue("obat")
	pos := &IndexObat{
		MerkDagang: obat,
	}

	resp, err := sendPost(pos, r.FormValue("token"), url)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	listobt := []IndexObat{}
	json.NewDecoder(resp.Body).Decode(&listobt)
	defer resp.Body.Close()
	// fmt.Println(len(listobt))
	if len(listobt) == 0 {
		// fmt.Println("This means empty slice")
		responseTemplate(w, "OK", GenTemplate(pos, "listobtnil"), "", nil)
	} else {
		responseTemplate(w, "OK", GenTemplate(listobt, "listobt"), "", nil)
	}

}
func getPDFNow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	url := "https://get-bulan-dot-igdsanglah.appspot.com/bulanini"

	gettgl := r.FormValue("tgl")
	url2 := "https://dokter-dot-igdsanglah.appspot.com"
	staf := &Staff{
		LinkID: r.FormValue("linkpdf"),
	}
	kirim, err := sendPost(staf, r.FormValue("token"), url2)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	det := &DetailStaf{}
	json.NewDecoder(kirim.Body).Decode(det)
	defer kirim.Body.Close()
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
	nopeg := ""
	n := &nopeg
	if det.NIP == "" {
		*n = det.NPP
	} else {
		*n = det.NIP
	}
	log.Printf("NIp adalah : %v", r.FormValue("linkdok"))
	createPDF(w, pts, iki, gettgl, r.FormValue("nama"), nopeg)
}

func getPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	// fmt.Print("Request masuk")
	url := "https://get-bulan-dot-igdsanglah.appspot.com/"

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
	url2 := "https://dokter-dot-igdsanglah.appspot.com"
	staf := &Staff{
		LinkID: r.FormValue("linkpdf"),
	}
	kirim, err := sendPost(staf, r.FormValue("token"), url2)
	if err != nil {
		log.Fatalf("Terjadi kesalahan di server: %v", err)
	}
	det := &DetailStaf{}
	json.NewDecoder(kirim.Body).Decode(det)
	defer kirim.Body.Close()
	nopeg := ""
	n := &nopeg
	if det.NIP == "" {
		*n = det.NPP
	} else {
		*n = det.NIP + " " + det.GolonganPNS
	}
	log.Printf("NIp adalah : %v", r.FormValue("linkdok"))
	createPDF(w, pts, countIKI(pts), gettgl, r.FormValue("nama"), nopeg)
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
func createPDF(w http.ResponseWriter, p []Pasien, l []ListIKI, tgl, email, nip string) {
	a, b, c, d, e, f, g := countTotalIKI(l)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	// Tabel IKI //////////////////////////////////////////////////////////
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 210, Ht: 297})
	pdf.Cell(160, 6, "Bukti Kegiatan Harian")
	pdf.Cell(120, 6, ("Nama Pegawai: " + email))
	pdf.Ln(-1)
	pdf.Cell(160, 6, "Pegawai RSUP Sanglah Denpasar")
	pdf.Cell(120, 6, "NIP/Gol: "+nip)
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
	pdf.Ln(-1)
	linep3k := ""
	p3k := &linep3k
	linerapat := ""
	rapat := &linerapat
	linepelatihan := ""
	pelatihan := &linepelatihan
	// for _, v := range l {

	// }
	p3 := []string{}
	rap := []string{}
	pel := []string{}
	for _, v := range l {
		for _, n := range v.P3K {
			*p3k = ": Tanggal "
			if n == "" {
				*p3k = ""
				continue
			} else {
				p3 = append(p3, n)
				// *p3k = *p3k + n + " "
			}
		}
		for _, n := range v.Pelatihan {
			*pelatihan = ": Tanggal "
			if n == "" {
				*pelatihan = ""
				continue
			} else {
				pel = append(pel, n)
				// *pelatihan = *pelatihan + n + " "
			}
		}
		for _, n := range v.Rapat {
			*rapat = ": Tanggal "
			if n == "" {
				*rapat = ""
				continue
			} else {
				rap = append(rap, n)
				// *rapat = *rapat + n + " "
			}
		}
	}

	for _, v := range p3 {
		*p3k = *p3k + v + ", "
	}
	for _, v := range rap {
		*rapat = *rapat + v + ", "
	}
	for _, v := range pel {
		*pelatihan = *pelatihan + v + ", "
	}

	// fo
	linep3k = strings.TrimSuffix(linep3k, ", ")
	linepelatihan = strings.TrimSuffix(linepelatihan, ", ")
	linerapat = strings.TrimSuffix(linerapat, ", ")
	pdf.Cell(40, 6, "P3K")
	pdf.Cell(30, 6, linep3k)
	pdf.Ln(-1)
	pdf.Cell(40, 6, "Pelatihan/Simulasi/Lainnya")
	pdf.Cell(30, 6, linepelatihan)
	pdf.Ln(-1)
	pdf.Cell(40, 6, "Rapat")
	pdf.Cell(30, 6, linerapat)
	pdf.Ln(-1)
	// for _, v := range l {
	// 	if v.P3K != nil {
	// 		pdf.Cell(20, 6, "P3K")
	// 		pdf.Cell(30, 6, linep3k)
	// 		pdf.Ln(-1)
	// 	}
	// 	if v.Pelatihan != nil {
	// 		pdf.Cell(40, 6, "Pelatihan/Simulasi/Lainnya")
	// 		pdf.Cell(30, 6, linepelatihan)
	// 		pdf.Ln(-1)
	// 	}
	// 	if v.Rapat != nil {
	// 		pdf.Cell(20, 6, "Rapat")
	// 		pdf.Cell(30, 6, linerapat)
	// 		pdf.Ln(-1)
	// 	}
	// }
	// for _, v := range l {
	// 	linerapat = linerapat + v + " "
	// }

	// for _, v := range l {
	// 	linepelatihan = linepelatihan + v + " "
	// }
	// pdf.Cell(20, 6, "P3K")
	// pdf.Cell(30, 6, linep3k)
	// pdf.Ln(-1)
	// pdf.Cell(20, 6, "Rapat")
	// pdf.Cell(30, 6, linerapat)
	// pdf.Ln(-1)
	// pdf.Cell(20, 6, "Pelatihan")
	// pdf.Cell(30, 6, linepelatihan)
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

	url := "https://get-bulan-dot-igdsanglah.appspot.com/"

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
	jaga := dataJaga(perBagian(pts), countIKI(pts), ConvertEkstra(countIKI(pts)))
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

func ConvertEkstra(n []ListIKI) KegiatanEkstra {
	var keg KegiatanEkstra
	for _, v := range n {
		// fmt.Println(v.P3K == nil)
		if v.P3K != nil {
			for _, n := range v.P3K {
				keg.P3K = append(keg.P3K, n)
			}
		}
		if v.Pelatihan != nil {
			for _, n := range v.Pelatihan {
				keg.Pelatihan = append(keg.Pelatihan, n)
			}
		}
		if v.Rapat != nil {
			for _, n := range v.Rapat {
				keg.Rapat = append(keg.Rapat, n)
			}
		}
	}
	// fmt.Println(keg)
	return keg
}
func getMonthly(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}
	url := "https://get-bulan-dot-igdsanglah.appspot.com/bulanini"
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
	jaga := dataJaga(perBagian(pts), countIKI(pts), ConvertEkstra(countIKI(pts)))
	// fmt.Printf("LIst iki adalah: %v", iki)

	responseTemplate(w, "OK", GenTemplate(pts, "contentrefresh"), GenTemplate(iki, "tabeliki"), jaga)

}

func countIKI(n []Pasien) []ListIKI {

	g := []ListIKI{}

	for h := 1; h <= 31; h++ {
		var u1, u2 int
		var p3k []string
		var rapat []string
		var pelatihan []string
		for _, v := range n {
			tgl, _ := strconv.Atoi(v.TglKunjungan[:2])
			if tgl != h {
				continue
			}

			// if v.NamaPasien == "P3K" || v.NamaPasien == "Rapat Rutin" || v.NamaPasien == "Pelatihan" {
			if v.NamaPasien == "P3K" {
				p3k = append(p3k, v.TglKunjungan[:2])
			} else if v.NamaPasien == "Rapat/Pertemuan" {
				rapat = append(rapat, v.TglKunjungan[:2])
			} else if v.NamaPasien == "Pelatihan" {
				pelatihan = append(pelatihan, v.TglKunjungan[:2])
			} else if v.IKI == "1" {
				u1++
			} else {
				u2++
			}
		}
		f := ListIKI{
			Tanggal:   h,
			SumIKI1:   u1,
			SumIKI2:   u2,
			P3K:       p3k,
			Rapat:     rapat,
			Pelatihan: pelatihan,
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

//////////////////////////////////////////////////////////////////////////////
func inputObat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request please", http.StatusMethodNotAllowed)
		return
	}

	var dat InputObat
	m := []byte(r.FormValue("send"))
	err := json.Unmarshal(m, &dat)
	if err != nil {
		fmt.Printf("Gagal mengubah json: %v", err)
	}

	url := "https://input-obat-dot-igdsanglah.appspot.com/"
	resp, err := sendPost(dat, r.FormValue("token"), url)
	if err != nil {
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Terjadi kesalahan server. Hubungi admin", ""), nil)
		log.Print("Terjadi kesalahan server")
		return
	}
	n := ServerResponse{}
	json.NewDecoder(resp.Body).Decode(&n)
	if n.Error != "" {
		responseTemplate(w, "not-OK", "", GenModal("Kesalahan Server", "Gagal menyimpan ke datastore. Ulangi lagi menginput data", ""), nil)
	}
	responseTemplate(w, "OK", "", "", nil)
}

//////////////////////////////////////////////////////////////////////////////
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
	// fmt.Print(GenTemplate(nil, "modlistresepnew"))
	responseTemplate(w, "OK", GenTemplate(nil, "modlistresepnew"), "", nil)
}

func getPtsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get request please", http.StatusMethodNotAllowed)
		return
	}
	responseTemplate(w, "OK", GenTemplate(nil, "modresepnew"), "", nil)
}

func confEditTanggal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post request only", http.StatusMethodNotAllowed)
	}

	url := "https://update-tanggal-dot-igdsanglah.appspot.com/"

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

	url := "https://ubah-tanggal-dot-igdsanglah.appspot.com/"
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

	url := "https://first-entries-dot-igdsanglah.appspot.com/"

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

	url := "https://delete-entri-dot-igdsanglah.appspot.com/"
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

	url := "https://update-entri-dot-igdsanglah.appspot.com/"
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

	url := "https://get-edit-entri-dot-igdsanglah.appspot.com/"

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

	url := "https://add-pasien-dot-igdsanglah.appspot.com/"

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
	// b := new(bytes.Buffer)
	// tmp := template.Must(template.New("baristabel.html").ParseFiles("templates/baristabel.html"))
	// err = tmp.Execute(b, pts)
	// if err != nil {
	// 	responseTemplate(w, "kesalahan-template", "", "", nil)
	// }
	responseTemplate(w, "OK", GenTemplate(pts, "baristabel"), GenTemplate(pts, "modformats"), nil)

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

	url := "https://get-data-pasien-dot-igdsanglah.appspot.com/"

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
		"tgl": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			return t.In(zone).Format("02-01-2006")
		},
		"jam": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			return t.In(zone).Format("15:04")
		},
		"umur": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			now := time.Now().In(zone)
			t = t.In(zone)
			yr := now.Year() - t.Year()
			var har int
			if now.Day() >= t.Day() {
				har = now.Day() - t.Day()
			} else {
				har = (now.Day() + now.AddDate(0, 0, -(now.Day())).Day()) - t.Day()
			}
			if now.YearDay() < t.YearDay() {
				yr--
				return fmt.Sprintf("%v Tahun %v Bulan %v Hari", yr, 12+int(now.Month()-t.Month()), har)
			}

			return fmt.Sprintf("%v Tahun %v Bulan %v Hari", yr, int(now.Month()-t.Month()), har)

		},
		"umurtahun": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			now := time.Now().In(zone)
			t = t.In(zone)
			yr := now.Year() - t.Year()
			if now.YearDay() < t.YearDay() {
				yr--
			}
			return strconv.Itoa(yr)
		},
		"bag": func(n string) string {
			switch n {
			case "1":
				return "Interna"
				break
			case "2":
				return "Bedah"
			case "3":
				return "Anak"
			case "4":
				return "Obgyn"
			case "5":
				return "Neuro"
			case "6":
				return "Anestesi"
			case "7":
				return "Psikiatri"
			case "8":
				return "THT"
			case "9":
				return "Kulit/Kelamin"
			case "10":
				return "Kardio"
			case "11":
				return "Umum"
			case "12":
				return "Mata"
			case "13":
				return "MOD"
			}
			return "Undefined"
		},
		"jenkel": func(n string) string {
			switch n {
			case "1":
				return "Laki-laki"
			case "2":
				return "Perempuan"
			}
			return "Belum dilengkapi"
		},
		"rekmed": func(n string) bool {
			if n == "4" || n == "5" {
				return true
			} else {
				return false
			}
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
	resp, err := http.Get("https://igdsanglah.appspot.com/login?token=" + token)
	if err != nil {
		log.Fatal(err)
		return
	}
	web := &MainView{}
	json.NewDecoder(resp.Body).Decode(web)
	defer resp.Body.Close()
	if web.Peran == "admin" {
		responseTemplate(w, web.Token, GenTemplate(web, "adminpage"), "", nil)
	} else if web.Peran == "supervisor" {
		responseTemplate(w, web.Token, GenTemplate(web, "supervisorpage"), web.Supervisor.SupervisorName, web)
	} else if web.Peran == "staf" {
		responseTemplate(w, web.Token, GenTemplate(web, "main", "input", "content"), web.User, nil)
	} else {
		responseTemplate(w, web.Token, GenTemplate(web, "residenpage", "residenpage-content"), web.User, nil)
	}
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
