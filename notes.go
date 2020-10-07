package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
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
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlcipher"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	embedMigrate "github.com/klingtnet/embed/migrate"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark"
	goldmarkEmoji "github.com/yuin/goldmark-emoji"
	goldmarkExtension "github.com/yuin/goldmark/extension"
)

// AppName is the name of the application.
const AppName = "notes"

var (
	// Version is the build version set by make.
	Version = "unset"

	indexTemplate,
	errorTemplate *template.Template
)

type TemplateData struct {
	Title  string
	Header TemplateHeaderData
	Main   TemplateMainData
	Footer TemplateFooterData
}

type NoteRecord struct {
	ID int64
	DateCreated,
	DateUpdated time.Time
	HTML template.HTML
}

type TemplateHeaderData struct {
	Title string
}

type TemplateMainData struct {
	Heading string
	Content interface{}
}

type TemplateIndexContent struct {
	NotesByDay map[time.Time][]NoteRecord
	Days       []time.Time
	EditText,
	SubmitAction string
}

type TemplateErrorContent struct {
	ErrorMessage string
}

type TemplateFooterData struct {
	AppName,
	Version string
	RenderDate time.Time
}

func respondWithTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data interface{}) {
	buf := bytes.NewBuffer(nil)
	err := tmpl.ExecuteTemplate(buf, "", data)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}
	_, err = io.CopyN(w, buf, int64(buf.Len()))
	if err != nil {
		log.Println(err)
	}
}

func respondWithErrorPage(w http.ResponseWriter, r *http.Request, err error, msg string, statusCode int) {
	log.Println(err.Error())
	if msg == "" {
		msg = err.Error()
	}
	td := TemplateData{
		Title:  "notes",
		Header: TemplateHeaderData{Title: "notes"},
		Main:   TemplateMainData{Heading: "something went wrong 😿", Content: TemplateErrorContent{ErrorMessage: msg}},
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

func indexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.QueryContext(r.Context(), `SELECT id, date_created, date_updated, html FROM note ORDER BY date_created DESC;`)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []NoteRecord
	for rows.Next() {
		var (
			id                       int64
			rawDateCreated, noteHTML string
			rawDateUpdated           = new(string)
		)
		err = rows.Scan(&id, &rawDateCreated, &rawDateUpdated, &noteHTML)
		if err != nil {
			respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
			return
		}

		dateCreated, err := time.Parse(time.RFC3339, rawDateCreated)
		if err != nil {
			respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
			return
		}

		dateUpdated := time.Time{}
		if rawDateUpdated != nil {
			dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
			if err != nil {
				respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
				return
			}
		}

		notes = append(notes, NoteRecord{ID: id, HTML: template.HTML(noteHTML), DateCreated: dateCreated, DateUpdated: dateUpdated})
	}
	err = rows.Err()
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}

	notesByDay := make(map[time.Time][]NoteRecord)
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
		Header: TemplateHeaderData{Title: "notes"},
		Main: TemplateMainData{Heading: "notes", Content: TemplateIndexContent{
			NotesByDay:   notesByDay,
			Days:         days,
			SubmitAction: "/submit",
		}},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	respondWithTemplate(w, r, indexTemplate, td)
}

func noteSubmitHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, mdParser goldmark.Markdown) {
	err := r.ParseForm()
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}
	mdNote := r.FormValue("note")

	buf := bytes.NewBuffer(nil)
	err = mdParser.Convert([]byte(mdNote), buf)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}

	_, err = db.ExecContext(r.Context(), `INSERT INTO note(date_created, markdown, html) VALUES(?,?,?)`, time.Now().Format(time.RFC3339), mdNote, buf.String())
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func noteEditHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}

	var md string
	err = db.QueryRowContext(r.Context(), `SELECT markdown FROM note WHERE id = ?`, noteID).Scan(&md)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}

	td := TemplateData{
		Title:  "notes",
		Header: TemplateHeaderData{Title: "notes"},
		Main: TemplateMainData{Heading: "notes", Content: TemplateIndexContent{
			SubmitAction: fmt.Sprintf("/note/%d/update", noteID),
			EditText:     md,
		}},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	respondWithTemplate(w, r, indexTemplate, td)
}

func noteUpdateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, mdParser goldmark.Markdown) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}
	mdNote := r.FormValue("note")

	buf := bytes.NewBuffer(nil)
	err = mdParser.Convert([]byte(mdNote), buf)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}

	_, err = db.ExecContext(r.Context(), `UPDATE note SET markdown=?, html=?, date_updated=? WHERE id=?;`, mdNote, buf.String(), time.Now().Format(time.RFC3339), noteID)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func noteSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusBadRequest)
		return
	}
	pattern := r.FormValue("search-pattern")
	if strings.TrimSpace(pattern) == "" {
		respondWithErrorPage(w, r, fmt.Errorf("search-pattern is missing"), "", http.StatusBadRequest)
		return
	}

	rows, err := db.QueryContext(r.Context(), `SELECT id, date_created, date_updated, html FROM note WHERE id IN (SELECT id FROM note_fts WHERE markdown MATCH ?);`, pattern)
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []NoteRecord
	for rows.Next() {
		var (
			id                       int64
			rawDateCreated, noteHTML string
			rawDateUpdated           = new(string)
		)
		err = rows.Scan(&id, &rawDateCreated, &rawDateUpdated, &noteHTML)
		if err != nil {
			respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
			return
		}

		dateCreated, err := time.Parse(time.RFC3339, rawDateCreated)
		if err != nil {
			respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
			return
		}

		dateUpdated := time.Time{}
		if rawDateUpdated != nil {
			dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
			if err != nil {
				respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
				return
			}
		}

		notes = append(notes, NoteRecord{ID: id, HTML: template.HTML(noteHTML), DateCreated: dateCreated, DateUpdated: dateUpdated})
	}
	err = rows.Err()
	if err != nil {
		respondWithErrorPage(w, r, err, "", http.StatusInternalServerError)
		return
	}

	notesByDay := make(map[time.Time][]NoteRecord)
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
		Header: TemplateHeaderData{Title: "notes"},
		Main: TemplateMainData{Heading: fmt.Sprintf("Search Resulst for %q", pattern), Content: TemplateIndexContent{
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
	dbPassphrase := strings.TrimSpace(c.String("database-passphrase"))
	if dbPassphrase == "" {
		return fmt.Errorf("required database passphrase is empty")
	}

	indexTemplate = parseTemplate("views/layouts/base.gohtml", "views/index.gohtml")
	errorTemplate = parseTemplate("views/layouts/base.gohtml", "views/error.gohtml")

	return run(c.Context, dbPassphrase, c.String("listen-addr"))
}

func parseTemplate(layout, content string) *template.Template {
	t := template.New("")
	t = template.Must(t.Parse(Embeds.FileString(content)))
	return template.Must(t.Parse(Embeds.FileString(layout)))
}

func run(ctx context.Context, dbPassphrase, httpAddr string) error {
	dbURI := fmt.Sprintf("file:notes.db?_pragma_key=%s&_pragma_cipher_page_size=4096&_foreign_keys=1", url.QueryEscape(dbPassphrase))
	db, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return err
	}
	driver, err := sqlcipher.WithInstance(db, &sqlcipher.Config{})
	if err != nil {
		return err
	}

	sourceDriver, err := embedMigrate.WithInstance(Embeds)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("embed", sourceDriver, "sqlite3", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	mdParser := goldmark.New(goldmark.WithExtensions(goldmarkExtension.GFM, goldmarkEmoji.Emoji))

	r := chi.NewRouter()
	r.Use(middleware.Recoverer, middleware.Logger, middleware.Compress(5))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, db)
	})
	r.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		noteSubmitHandler(w, r, db, mdParser)
	})
	r.Get("/note/{noteID}/edit", func(w http.ResponseWriter, r *http.Request) {
		noteEditHandler(w, r, db)
	})
	r.Post("/note/{noteID}/update", func(w http.ResponseWriter, r *http.Request) {
		noteUpdateHandler(w, r, db, mdParser)
	})
	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		noteSearchHandler(w, r, db)
	})
	r.Get("/assets/*", assetHandler)

	log.Printf("listening on %q", httpAddr)
	return http.ListenAndServe(httpAddr, r)
}

func main() {
	app := cli.App{
		Name:    AppName,
		Version: Version,
		Action:  runAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database-passphrase",
				Usage:    "SQLcipher database passphrase",
				EnvVars:  []string{"DATABASE_PASSPHRASE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "listen-addr",
				Usage: "HTTP listen address",
				Value: "localhost:3333",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
