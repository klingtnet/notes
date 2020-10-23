package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	dbPassphrase := t.Name()
	dbURI := fmt.Sprintf(":memory:?_pragma_key=%s&_pragma_cipher_page_size=4096&_foreign_keys=1", url.QueryEscape(dbPassphrase))
	db, err := sql.Open("sqlite3", dbURI)
	require.NoError(t, err)
	return db
}

func testConvert(t *testing.T, prefix string) func(string) (string, error) {
	return func(s string) (string, error) { return prefix + s, nil }
}

func TestInsert(t *testing.T) {
	db := setupTestDB(t)

	ns, err := newSQLCipherNotes(db, testConvert(t, t.Name()))
	require.NoError(t, err)

	notes, err := ns.Notes(context.Background())
	require.NoError(t, err)
	require.Empty(t, notes)

	id, err := ns.Insert(context.Background(), "A bit of ~~asciidoc~~ _markdown_ content. 🎉")
	require.NoError(t, err)

	note, err := ns.Note(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, note)

	require.Equal(t, id, note.ID)
	require.Equal(t, "A bit of ~~asciidoc~~ _markdown_ content. 🎉", note.Markdown)
	require.Equal(t, template.HTML(t.Name()+"A bit of ~~asciidoc~~ _markdown_ content. 🎉"), note.HTML)
	require.True(t, note.DateUpdated.IsZero())
	require.WithinDuration(t, time.Now(), note.DateCreated, 1*time.Second)
}

func TestEdit(t *testing.T) {
	db := setupTestDB(t)

	ns, err := newSQLCipherNotes(db, testConvert(t, t.Name()))
	require.NoError(t, err)

	id, err := ns.Insert(context.Background(), "Does not matter much.")
	require.NoError(t, err)

	tCases := []struct {
		name string
		id   int64
		markdown,
		errMsg string
	}{
		{"existing-note", id, "An updated note.", ""},
		{"invalid-id", -123, "Does not matter much.", "expected one row to change but was 0"},
	}

	parentTestName := t.Name()
	for _, tCase := range tCases {
		t.Run(tCase.name, func(t *testing.T) {
			err = ns.Update(context.Background(), tCase.id, tCase.markdown)
			if tCase.errMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tCase.errMsg)
				return
			}

			note, err := ns.Note(context.Background(), id)
			require.NoError(t, err)
			require.NotNil(t, note)

			require.Equal(t, tCase.markdown, note.Markdown)
			require.Equal(t, template.HTML(parentTestName+tCase.markdown), note.HTML)
			require.WithinDuration(t, time.Now(), note.DateCreated, 1*time.Second)
			require.WithinDuration(t, time.Now(), note.DateCreated, 1*time.Second)
		})
	}
}

func TestSearch(t *testing.T) {
	db := setupTestDB(t)
	ns, err := newSQLCipherNotes(db, testConvert(t, t.Name()))
	require.NoError(t, err)

	mustRead := func(p string) string {
		data, err := ioutil.ReadFile(p)
		require.NoError(t, err, "mustRead")
		return string(data)
	}

	tCases := []struct {
		name,
		markdown,
		searchPattern string
	}{
		{"english with symbol", "Remember the milk™", "milk"},
		{"japanese", "おすし大好き！", "大好き"},
		{"russian", "Здравствуй, мир", "мир"},
		{"inside markdown link text", mustRead("README.md"), "Releases"},
		{"link inside text file", mustRead("LICENSE"), "fsf.org"},
	}

	for _, tCase := range tCases {
		t.Run(tCase.name, func(t *testing.T) {
			id, err := ns.Insert(context.Background(), tCase.markdown)
			require.NoError(t, err)

			notes, err := ns.Search(context.Background(), tCase.searchPattern)
			require.NoError(t, err)
			require.NotNil(t, notes)
			require.Len(t, notes, 1)
			require.Equal(t, tCase.markdown, notes[0].Markdown)
			require.Equal(t, int64(id), notes[0].ID)
		})
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)

	ns, err := newSQLCipherNotes(db, testConvert(t, t.Name()))
	require.NoError(t, err)

	id, err := ns.Insert(context.Background(), "I will be deleted 😥")
	require.NoError(t, err)

	tCases := []struct {
		name   string
		id     int64
		errMsg string
	}{
		{"existing-note", id, ""},
		{"already-deleted", id, "expected to delete one row but was 0"},
		{"invalid-id", -123, "expected to delete one row but was 0"},
	}

	for _, tCase := range tCases {
		t.Run(tCase.name, func(t *testing.T) {
			err = ns.Delete(context.Background(), tCase.id)
			if tCase.errMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tCase.errMsg)
				return
			}
		})
	}
}

func TestRenew(t *testing.T) {
	db := setupTestDB(t)

	ns, err := newSQLCipherNotes(db, testConvert(t, t.Name()))
	require.NoError(t, err)

	// prepare some notes
	for i := 0; i < 10; i++ {
		_, err = ns.Insert(context.Background(), "I will be renewed ✔️")
		require.NoError(t, err)
	}

	// create a new instance to simulate a changed markdown configuration
	ns, err = newSQLCipherNotes(db, func(s string) (string, error) { return "this simulates a changed markdown parser congifuration", nil })
	require.NoError(t, err)

	err = ns.Renew(context.Background())
	require.NoError(t, err)
	notes, err := ns.Notes(context.Background())
	require.NoError(t, err)
	require.Len(t, notes, 10)
	for _, note := range notes {
		require.Equal(t, template.HTML("this simulates a changed markdown parser congifuration"), note.HTML)
	}
}
