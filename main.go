package main

import (
	"log"
	"net/http"
	"time"

	coremiddleware "golang-xlsx/middleware"

	responsebuilder "github.com/featriadi/golang-libs/response_builder"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.With(coremiddleware.RequestXLSXSheetRows(5), coremiddleware.RequestXLSXConvertDataStruct[BulkAll]()).Post("/bulk", BulkXlsx)

	http.ListenAndServe(":8080", r)
}

type Song struct {
	Code             string     `json:"code" validate:"required"`
	Type             string     `json:"type"`
	ParentCode       string     `json:"parent_code"`
	Title            string     `json:"title" validate:"required"`
	Composers        string     `json:"composers"`
	Notes            string     `json:"notes"`
	Language         string     `json:"language"`
	CopyrightDate    *time.Time `json:"copyright_date"`
	CopyrightNo      string     `json:"copyright_no"`
	Duration         string     `json:"duration"`
	Recorded         bool       `json:"recorded"`
	Vocal            string     `json:"vocal"`
	ISWC             string     `json:"iswc"`
	ApprovalRequired string     `json:"approval_required"`
	ApprovalNotes    string     `json:"approval_notes"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	CopyrightLine    string     `json:"copyright_line"`
	ISANNumber       string     `json:"isan_number"`
	ISRCNumber       string     `json:"isrc_number"`
	SiteID           string     `json:"site_id"`
}

type SongRecording struct {
	Code           string `json:"code"`
	SongCode       string `json:"song_code"`
	MainRecording  bool   `json:"main_recording"`
	ExternalID     string `json:"external_id"`
	Title          string `json:"title"`
	RecordingDate  string `json:"recording_date"`
	CatalogNumber  string `json:"catalog_number"`
	ISRC           string `json:"isrc"`
	Recorded       string `json:"recorded"`
	AudioFilePath  string `json:"audio_file_path"`
	SampleFilePath string `json:"sample_file_path"`
	Duration       string `json:"duration"`
	MasterOwner    string `json:"master_owner"`
	Label          string `json:"label"`
	MusicianFee    string `json:"musician_fee"`
	Note           string `json:"note"`
	Release        string `json:"release"`
	AlbumTitle     string `json:"album_title"`
	UPC            string `json:"upc"`
	Artist         string `json:"artist"`
}

type BulkAll struct {
	Song          Song          `json:"song"`
	SongRecording SongRecording `json:"song_recording"`
}

func BulkXlsx(w http.ResponseWriter, r *http.Request) {
	songs, err := coremiddleware.GetRequestXLSXData[Song](r)
	if err != nil {
		responsebuilder.NewResponseBuilder[any](&w).
			Status(err.GetStatus()).
			Message(err.GetMessage()).
			Build()
		return
	}

	for _, song := range *songs {
		log.Println(song)
		// log.Println(song.CopyrightDate)
	}

	responsebuilder.NewResponseBuilder[any](&w).
		Status(http.StatusOK).
		Message("Ok!").
		Build()
}
