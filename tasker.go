package main

import (
	"bytes"
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
	"regexp"
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
	"github.com/yuin/goldmark"
	goldmarkEmoji "github.com/yuin/goldmark-emoji"
	goldmarkExtension "github.com/yuin/goldmark/extension"
)

var (
	AppName = "tasker"
	Version string
	JiraRe  = regexp.MustCompile(`((\s+)((DEV|INT)-\d+))`)

	indexTemplate *template.Template
)

type TemplateData struct {
	Title  string
	Header TemplateHeaderData
	Main   TemplateMainData
	Footer TemplateFooterData
}

type Note struct {
	ID int64
	DateCreated,
	DateUpdated time.Time
	HTML template.HTML
}

type TemplateHeaderData struct {
	Title string
}

type TemplateMainData struct {
	Heading,
	EditText,
	SubmitAction string
	NotesByDay map[time.Time][]Note
	Days       []time.Time
}

type TemplateFooterData struct {
	AppName,
	Version string
	RenderDate time.Time
}

func indexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.QueryContext(r.Context(), `SELECT id, date_created, date_updated, html FROM note ORDER BY date_created DESC;`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var (
			id                       int64
			rawDateCreated, noteHTML string
			rawDateUpdated           = new(string)
		)
		err = rows.Scan(&id, &rawDateCreated, &rawDateUpdated, &noteHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dateCreated, err := time.Parse(time.RFC3339, rawDateCreated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dateUpdated := time.Time{}
		if rawDateUpdated != nil {
			dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		notes = append(notes, Note{ID: id, HTML: template.HTML(noteHTML), DateCreated: dateCreated, DateUpdated: dateUpdated})
	}
	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		Title:  "tasker",
		Header: TemplateHeaderData{Title: "tasker"},
		Main:   TemplateMainData{NotesByDay: notesByDay, Days: days, SubmitAction: "/submit"},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	buf := bytes.NewBuffer(nil)
	err = indexTemplate.ExecuteTemplate(buf, "", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.CopyN(w, buf, int64(buf.Len()))
	if err != nil {
		log.Println(err)
	}
}

func noteSubmitHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, mdParser goldmark.Markdown, mdFn func(string) string) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mdNote := mdFn(r.FormValue("note"))

	buf := bytes.NewBuffer(nil)
	err = mdParser.Convert([]byte(mdNote), buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.ExecContext(r.Context(), `INSERT INTO note(date_created, markdown, html) VALUES(?,?,?)`, time.Now().Format(time.RFC3339), mdNote, buf.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func noteEditHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var md string
	err = db.QueryRowContext(r.Context(), `SELECT markdown FROM note WHERE id = ?`, noteID).Scan(&md)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	td := TemplateData{
		Title:  "tasker",
		Header: TemplateHeaderData{Title: "tasker"},
		Main:   TemplateMainData{EditText: md, SubmitAction: fmt.Sprintf("/note/%d/update", noteID)},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	buf := bytes.NewBuffer(nil)
	err = indexTemplate.ExecuteTemplate(buf, "", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.CopyN(w, buf, int64(buf.Len()))
	if err != nil {
		log.Println(err)
	}
}

func noteUpdateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, mdParser goldmark.Markdown, mdFn func(string) string) {
	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mdNote := mdFn(r.FormValue("note"))

	buf := bytes.NewBuffer(nil)
	err = mdParser.Convert([]byte(mdNote), buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.ExecContext(r.Context(), `UPDATE note SET markdown=?, html=?, date_updated=? WHERE id=?;`, mdNote, buf.String(), time.Now().Format(time.RFC3339), noteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func noteSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pattern := r.FormValue("search-pattern")
	if strings.TrimSpace(pattern) == "" {
		http.Error(w, "search-pattern is missing", http.StatusBadRequest)
		return
	}

	rows, err := db.QueryContext(r.Context(), `SELECT id, date_created, date_updated, html FROM note WHERE id IN (SELECT id FROM note_fts WHERE markdown MATCH ?);`, pattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var (
			id                       int64
			rawDateCreated, noteHTML string
			rawDateUpdated           = new(string)
		)
		err = rows.Scan(&id, &rawDateCreated, &rawDateUpdated, &noteHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dateCreated, err := time.Parse(time.RFC3339, rawDateCreated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dateUpdated := time.Time{}
		if rawDateUpdated != nil {
			dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		notes = append(notes, Note{ID: id, HTML: template.HTML(noteHTML), DateCreated: dateCreated, DateUpdated: dateUpdated})
	}
	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		Title:  "tasker",
		Header: TemplateHeaderData{Title: "tasker"},
		Main:   TemplateMainData{Heading: fmt.Sprintf("Search Resulst for %q", pattern), NotesByDay: notesByDay, Days: days, SubmitAction: "/submit"},
		Footer: TemplateFooterData{Version: Version, AppName: AppName, RenderDate: time.Now()},
	}

	buf := bytes.NewBuffer(nil)
	err = indexTemplate.ExecuteTemplate(buf, "", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.CopyN(w, buf, int64(buf.Len()))
	if err != nil {
		log.Println(err)
	}
}

func run(dbPassword, httpAddr, jiraRootURL string) error {
	indexTemplate = template.Must(template.New("").Parse(Embeds.FileString("views/index.gohtml")))

	dbURI := fmt.Sprintf("file:tasker.db?_pragma_key=%s&_pragma_cipher_page_size=4096&_foreign_keys=1", url.QueryEscape(dbPassword))
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

	linkifyTickets := func(md string) string {
		// link all tickets
		// TODO: maybe use a Goldmark extension for htis?
		return JiraRe.ReplaceAllString(md, `$2[$3](`+jiraRootURL+`$3)`)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, db)
	})
	r.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		noteSubmitHandler(w, r, db, mdParser, linkifyTickets)
	})
	r.Get("/note/{noteID}/edit", func(w http.ResponseWriter, r *http.Request) {
		noteEditHandler(w, r, db)
	})
	r.Post("/note/{noteID}/update", func(w http.ResponseWriter, r *http.Request) {
		noteUpdateHandler(w, r, db, mdParser, linkifyTickets)
	})
	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		noteSearchHandler(w, r, db)
	})
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		file := strings.TrimPrefix(r.URL.Path, "/")
		data := Embeds.File(file)
		if data == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		contentType := mime.TypeByExtension(filepath.Ext(file))
		if contentType == "" {
			contentType = http.DetectContentType(data)
		}
		w.Header().Add("content-type", contentType)
		w.Write(data)
	})

	log.Printf("listening on %q", httpAddr)
	http.ListenAndServe(httpAddr, r)

	return nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <database-password> <http-listen-addr> <jira-root-url>", os.Args[0])
	}

	err := run(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
}
