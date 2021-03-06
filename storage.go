package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"time"
	"unicode"
)

type noteStorage interface {
	// Insert stores a new note.
	Insert(ctx context.Context, markdown string) (id int64, err error)
	// Update overwrites the existing markdown content for the note, renders it to HTML and sets the updated time.
	Update(ctx context.Context, id int64, markdown string) error
	// Search returns all notes for which their markdown content matches the given pattern.
	Search(ctx context.Context, pattern string) (notes []Note, err error)
	// Delete removes a note.
	Delete(ctx context.Context, id int64) error
	// Note returns the note having the given id.
	Note(ctx context.Context, id int64) (note *Note, err error)
	// Notes returns all notes.
	// TODO: add from/to time.Time arguments for pagination.
	Notes(ctx context.Context) (notes []Note, err error)
	// Renew renders the HTML output of all notes again.
	// This is useful if settings of the markdown parser changed.
	Renew(ctx context.Context) error
}

type sqlCipherNotes struct {
	db             *sql.DB
	markdownToHTML func(string) (string, error)
}

func newSQLCipherNotes(db *sql.DB, markdownToHTML func(string) (string, error)) (*sqlCipherNotes, error) {
	return &sqlCipherNotes{db, markdownToHTML}, nil
}

// Insert implements NoteStorage.
func (s *sqlCipherNotes) Insert(ctx context.Context, markdown string) (id int64, err error) {
	html, err := s.markdownToHTML(markdown)
	if err != nil {
		return
	}

	res, err := s.db.ExecContext(ctx, `INSERT INTO note(date_created, markdown, html) VALUES(?,?,?)`, time.Now().Format(time.RFC3339), markdown, html)
	if err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

// Update implements NoteStorage.
func (s *sqlCipherNotes) Update(ctx context.Context, id int64, markdown string) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			if errors.Is(rollbackErr, sql.ErrTxDone) {
				return
			}
			if err != nil {
				err = fmt.Errorf("%s: %w", rollbackErr.Error(), err)
			}
			err = rollbackErr
		}
	}()
	err = s.updateTx(ctx, tx, id, markdown)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (s *sqlCipherNotes) updateTx(ctx context.Context, tx *sql.Tx, id int64, markdown string) error {
	html, err := s.markdownToHTML(markdown)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE note SET markdown=?, html=?, date_updated=? WHERE id=?;`, markdown, html, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("expected one row to change but was %d", n)
	}
	return nil
}

// Search implements NoteStorage.
func (s *sqlCipherNotes) Search(ctx context.Context, pattern string) (notes []Note, err error) {
	var containsOtherUnicode bool
	for _, c := range pattern {
		if unicode.In(c, unicode.Lo) {
			containsOtherUnicode = true
			break
		}
	}

	var rows *sql.Rows
	if containsOtherUnicode {
		// Use a LIKE query because the fts4 index has problems matching languages with implicit whitespace, e.g. japanese
		rows, err = s.db.QueryContext(ctx, `SELECT id, date_created, date_updated, markdown, html FROM note WHERE markdown LIKE ?;`, "%"+pattern+"%")
	} else {
		rows, err = s.db.QueryContext(ctx, `SELECT id, date_created, date_updated, markdown, html FROM note WHERE id IN (SELECT id FROM note_fts WHERE markdown MATCH ?);`, pattern)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		var (
			rawDateCreated string
			rawDateUpdated = new(string)
		)
		err = rows.Scan(&note.ID, &rawDateCreated, &rawDateUpdated, &note.Markdown, &note.HTML)
		if err != nil {
			return
		}

		var dateCreated time.Time
		dateCreated, err = time.Parse(time.RFC3339, rawDateCreated)
		if err != nil {
			return
		}
		note.DateCreated = dateCreated

		if rawDateUpdated != nil {
			var dateUpdated time.Time
			dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
			if err != nil {
				return
			}
			note.DateUpdated = dateUpdated
		}

		notes = append(notes, note)
	}
	err = rows.Err()
	return
}

// Delete implents NoteStorage.
func (s *sqlCipherNotes) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM note WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("expected to delete one row but was %d", n)
	}
	return nil
}

// Note implements NoteStorage.
func (s *sqlCipherNotes) Note(ctx context.Context, id int64) (note *Note, err error) {
	var (
		rawDateCreated, markdown, noteHTML string
		rawDateUpdated                     = new(string)
		dateCreated                        time.Time
	)
	err = s.db.QueryRowContext(ctx, `SELECT date_created, date_updated, markdown, html FROM note WHERE id = ?`, id).Scan(&rawDateCreated, &rawDateUpdated, &markdown, &noteHTML)
	if err != nil {
		return
	}

	dateCreated, err = time.Parse(time.RFC3339, rawDateCreated)
	if err != nil {
		return
	}

	dateUpdated := time.Time{}
	if rawDateUpdated != nil {
		dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
		if err != nil {
			return
		}
	}

	note = &Note{ID: id, Markdown: markdown, HTML: template.HTML(noteHTML), DateCreated: dateCreated, DateUpdated: dateUpdated}
	return
}

// Notes implements NoteStorage.
func (s *sqlCipherNotes) Notes(ctx context.Context) (notes []Note, err error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, date_created, date_updated, markdown, html FROM note ORDER BY date_created DESC;`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Note{}, nil
		}
		return
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%s: %w", closeErr.Error(), err)
			}
			err = closeErr
		}
	}()

	for rows.Next() {
		var note Note
		var (
			rawDateCreated string
			rawDateUpdated = new(string)
		)
		err = rows.Scan(&note.ID, &rawDateCreated, &rawDateUpdated, &note.Markdown, &note.HTML)
		if err != nil {
			return
		}

		var dateCreated time.Time
		dateCreated, err = time.Parse(time.RFC3339, rawDateCreated)
		if err != nil {
			return
		}
		note.DateCreated = dateCreated

		if rawDateUpdated != nil {
			var dateUpdated time.Time
			dateUpdated, err = time.Parse(time.RFC3339, *rawDateUpdated)
			if err != nil {
				return
			}
			note.DateUpdated = dateUpdated
		}

		notes = append(notes, note)
	}
	err = rows.Err()
	return
}

// Renew implements NoteStorage.
func (s *sqlCipherNotes) Renew(ctx context.Context) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			if errors.Is(rollbackErr, sql.ErrTxDone) {
				return
			}
			if err != nil {
				err = fmt.Errorf("%s: %w", rollbackErr.Error(), err)
			}
			err = rollbackErr
		}
	}()

	rows, err := tx.QueryContext(ctx, `SELECT id, markdown FROM note;`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var markdown string
		err = rows.Scan(&id, &markdown)
		if err != nil {
			return err
		}

		err = s.updateTx(ctx, tx, id, markdown)
		if err != nil {
			return
		}
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return tx.Commit()
}
