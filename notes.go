package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark"
	goldmarkEmoji "github.com/yuin/goldmark-emoji"
	goldmarkExtension "github.com/yuin/goldmark/extension"
	goldmarkHTML "github.com/yuin/goldmark/renderer/html"
)

// AppName is the name of the application.
const AppName = "notes"

var (
	// Version is the build version set by make.
	Version = "unset"

	indexTemplate,
	deleteTemplate,
	errorTemplate *template.Template
)

// Note contains information you want to remember in markdown as well as HTML format and addiitonal metadata.
type Note struct {
	ID int64
	DateCreated,
	DateUpdated time.Time
	Markdown string
	HTML     template.HTML
}

// TemplateData contains all information required to render a template.
type TemplateData struct {
	Title  string
	Header TemplateHeaderData
	Main   TemplateMainData
	Footer TemplateFooterData
}

// TemplateHeaderData contains data required by the header template.
type TemplateHeaderData struct {
	AppName,
	Title string
}

// TemplateMainData contains data required by the main template.
type TemplateMainData struct {
	Heading string
	Content interface{}
}

// TemplateIndexContent contains data required by the index template.
type TemplateIndexContent struct {
	NotesByDay map[time.Time][]Note
	Days       []time.Time
	EditText,
	SubmitAction string
}

// TemplateErrorContent contains data required by the error template.
type TemplateErrorContent struct {
	ErrorMessage string
}

// TemplateFooterData contains data required by the footer template.
type TemplateFooterData struct {
	AppName,
	Version string
	RenderDate time.Time
}

func respondWithTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data interface{}) {
	buf := bytes.NewBuffer(nil)
	err := tmpl.ExecuteTemplate(buf, "", data)
	if err != nil {
		respondWithErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	_, err = io.CopyN(w, buf, int64(buf.Len()))
	if err != nil {
		log.Println(err)
	}
}

func respondWithErrorPage(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err.Error())
	td := TemplateData{
		Title:  "notes",
		Header: TemplateHeaderData{AppName: AppName, Title: "notes"},
		Main:   TemplateMainData{Heading: "something went wrong 😿", Content: TemplateErrorContent{ErrorMessage: err.Error()}},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	buf := bytes.NewBuffer(nil)
	err = errorTemplate.ExecuteTemplate(buf, "", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, err = io.CopyN(w, buf, int64(buf.Len()))
	if err != nil {
		log.Println(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, noteStor NoteStorage) {
	notes, err := noteStor.Notes(r.Context())
	if err != nil {
		respondWithErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	notesByDay := make(map[time.Time][]Note)
	for _, note := range notes {
		date, _ := time.Parse("2006-01-02", note.DateCreated.Format("2006-01-02"))
		notesByDay[date] = append(notesByDay[date], note)
	}
	var days []time.Time
	for day := range notesByDay {
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].After(days[j]) })

	td := TemplateData{
		Title:  "notes",
		Header: TemplateHeaderData{AppName: AppName, Title: "notes"},
		Main: TemplateMainData{Heading: "What do you want to remember?", Content: TemplateIndexContent{
			NotesByDay:   notesByDay,
			Days:         days,
			SubmitAction: "/submit",
		}},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	respondWithTemplate(w, r, indexTemplate, td)
}

func noteSubmitHandler(w http.ResponseWriter, r *http.Request, noteStor NoteStorage) {
	err := r.ParseForm()
	if err != nil {
		respondWithErrorPage(w, err, http.StatusBadRequest)
		return
	}
	markdown := r.FormValue("note")
	if strings.TrimSpace(markdown) == "" {
		respondWithErrorPage(w, fmt.Errorf("note is empty"), http.StatusBadRequest)
		return
	}

	_, err = noteStor.Insert(r.Context(), markdown)
	if err != nil {
		respondWithErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func noteEditHandler(w http.ResponseWriter, r *http.Request, noteStor NoteStorage) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		respondWithErrorPage(w, err, http.StatusBadRequest)
		return
	}

	note, err := noteStor.Note(r.Context(), int64(noteID))
	if err != nil {
		respondWithErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	td := TemplateData{
		Title:  "notes",
		Header: TemplateHeaderData{AppName: AppName, Title: "notes"},
		Main: TemplateMainData{Heading: "notes", Content: TemplateIndexContent{
			SubmitAction: fmt.Sprintf("/note/%d/edit", noteID),
			EditText:     note.Markdown,
		}},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	respondWithTemplate(w, r, indexTemplate, td)
}

func noteUpdateHandler(w http.ResponseWriter, r *http.Request, noteStor NoteStorage) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		respondWithErrorPage(w, err, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		respondWithErrorPage(w, err, http.StatusBadRequest)
		return
	}
	markdown := r.FormValue("note")
	if strings.TrimSpace(markdown) == "" {
		respondWithErrorPage(w, fmt.Errorf("note is empty"), http.StatusBadRequest)
		return
	}

	err = noteStor.Update(r.Context(), int64(noteID), markdown)
	if err != nil {
		respondWithErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func noteDeleteHandler(w http.ResponseWriter, r *http.Request, noteStor NoteStorage) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		respondWithErrorPage(w, err, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		td := TemplateData{
			Title:  "notes",
			Header: TemplateHeaderData{AppName: AppName, Title: "notes"},
			Main: TemplateMainData{Heading: "Confirm Delete", Content: struct {
				NoteID    int
				DeleteURL string
			}{noteID, r.URL.Path}},
			Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
		}

		respondWithTemplate(w, r, deleteTemplate, td)
		return
	case "POST":
		err = r.ParseForm()
		if err != nil {
			respondWithErrorPage(w, err, http.StatusBadRequest)
			return
		}

		switch r.FormValue("submit") {
		case "cancel":
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		case "delete":
			err = noteStor.Delete(r.Context(), int64(noteID))
			if err != nil {
				respondWithErrorPage(w, err, http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		respondWithErrorPage(w, fmt.Errorf("unknown submit value %q", r.FormValue("submit")), http.StatusBadRequest)
		return
	}
}

func noteSearchHandler(w http.ResponseWriter, r *http.Request, noteStor NoteStorage) {
	err := r.ParseForm()
	if err != nil {
		respondWithErrorPage(w, err, http.StatusBadRequest)
		return
	}
	pattern := r.FormValue("search-pattern")
	if strings.TrimSpace(pattern) == "" {
		respondWithErrorPage(w, fmt.Errorf("search-pattern is missing"), http.StatusBadRequest)
		return
	}

	notes, err := noteStor.Search(r.Context(), pattern)
	if err != nil {
		respondWithErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	notesByDay := make(map[time.Time][]Note)
	for idx := range notes {
		note := notes[len(notes)-1-idx]
		date, _ := time.Parse("2006-01-02", note.DateCreated.Format("2006-01-02"))
		notesByDay[date] = append(notesByDay[date], note)
	}
	var days []time.Time
	for day := range notesByDay {
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].After(days[j]) })

	td := TemplateData{
		Title:  "notes",
		Header: TemplateHeaderData{AppName: AppName, Title: "notes"},
		Main: TemplateMainData{Heading: fmt.Sprintf("Search Results for %q", pattern), Content: TemplateIndexContent{
			NotesByDay:   notesByDay,
			Days:         days,
			SubmitAction: "/submit",
		}},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	respondWithTemplate(w, r, indexTemplate, td)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/")
	data := Embeds.File(file)
	if data == nil {
		http.Error(w, "asset not found", http.StatusNotFound)
		return
	}
	contentType := mime.TypeByExtension(filepath.Ext(file))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	w.Header().Add("content-type", contentType)
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func runAction(c *cli.Context) error {
	dbPath := strings.TrimSpace(c.String("database-path"))
	if dbPath == "" {
		return fmt.Errorf("required database path is empty")
	}

	dbPassphrase := strings.TrimSpace(c.String("database-passphrase"))
	if dbPassphrase == "" {
		return fmt.Errorf("required database passphrase is empty")
	}

	indexTemplate = parseTemplate("views/layouts/base.gohtml", "views/index.gohtml")
	deleteTemplate = parseTemplate("views/layouts/base.gohtml", "views/delete.gohtml")
	errorTemplate = parseTemplate("views/layouts/base.gohtml", "views/error.gohtml")

	return run(c.Context, dbPath, dbPassphrase, c.String("listen-addr"))
}

func parseTemplate(layout, content string) *template.Template {
	t := template.New("")
	t = template.Must(t.Parse(Embeds.FileString(content)))
	return template.Must(t.Parse(Embeds.FileString(layout)))
}

func run(ctx context.Context, dbPath, dbPassphrase, httpAddr string) error {
	dbURI := fmt.Sprintf("file:%s?_pragma_key=%s&_pragma_cipher_page_size=4096&_foreign_keys=1", dbPath, url.QueryEscape(dbPassphrase))
	db, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return err
	}

	mdParser := goldmark.New(
		goldmark.WithExtensions(goldmarkExtension.GFM, goldmarkEmoji.Emoji),
		goldmark.WithRendererOptions(goldmarkHTML.WithUnsafe()),
	)
	markdownToHTML := func(markdown string) (string, error) {
		buf := bytes.NewBuffer(nil)
		err = mdParser.Convert([]byte(markdown), buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	noteStor, err := NewSQLCipherNotes(db, markdownToHTML)
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer, middleware.Logger, middleware.Compress(5))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, noteStor)
	})
	r.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		noteSubmitHandler(w, r, noteStor)
	})
	r.Route("/note/{noteID}/edit", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { noteEditHandler(w, r, noteStor) })
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { noteUpdateHandler(w, r, noteStor) })
	})
	r.Route("/note/{noteID}/delete", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { noteDeleteHandler(w, r, noteStor) })
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { noteDeleteHandler(w, r, noteStor) })
	})
	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		noteSearchHandler(w, r, noteStor)
	})
	r.Get("/assets/*", assetHandler)

	log.Printf("listening on http://%s", httpAddr)
	return http.ListenAndServe(httpAddr, r)
}

func renewAction(c *cli.Context) error {
	dbPath := strings.TrimSpace(c.String("database-path"))
	if dbPath == "" {
		return fmt.Errorf("required database path is empty")
	}

	dbPassphrase := strings.TrimSpace(c.String("database-passphrase"))
	if dbPassphrase == "" {
		return fmt.Errorf("required database passphrase is empty")
	}

	dbURI := fmt.Sprintf("file:%s?_pragma_key=%s&_pragma_cipher_page_size=4096&_foreign_keys=1", dbPath, url.QueryEscape(dbPassphrase))
	db, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return err
	}
	mdParser := goldmark.New(
		goldmark.WithExtensions(goldmarkExtension.GFM, goldmarkEmoji.Emoji),
		goldmark.WithRendererOptions(goldmarkHTML.WithUnsafe()),
	)
	markdownToHTML := func(markdown string) (string, error) {
		buf := bytes.NewBuffer(nil)
		err = mdParser.Convert([]byte(markdown), buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	noteStor, err := NewSQLCipherNotes(db, markdownToHTML)
	if err != nil {
		return err
	}

	return noteStor.Renew(c.Context)
}

func main() {
	app := cli.App{
		Name:    AppName,
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run the note server",
				Action: runAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "database-passphrase",
						Usage:    "SQLcipher database passphrase",
						EnvVars:  []string{"DATABASE_PASSPHRASE"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "database-path",
						Usage:    "path to the database file",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "listen-addr",
						Usage: "HTTP listen address",
						Value: "localhost:13333",
					},
				},
			},
			{
				Name:   "renew",
				Usage:  "renew refreshes all existing notes by rendering them again",
				Action: renewAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "database-passphrase",
						Usage:    "SQLcipher database passphrase",
						EnvVars:  []string{"DATABASE_PASSPHRASE"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "database-path",
						Usage:    "path to the database file",
						Required: true,
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
