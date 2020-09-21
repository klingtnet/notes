package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli"
)

func pathToVar(path string) string {
	return fmt.Sprintf("file%x", []byte(path))
}

func encodeFile(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
}

var (
	fileTemplate = template.Must(template.New("").Funcs(template.FuncMap{"pathToVar": pathToVar, "encode": encodeFile}).Parse(`package {{ .Package }}

import (
	"encoding/base64"
	"sort"
)

const (
{{- range $path, $data := .Files }}
	{{ pathToVar $path }} = "{{ encode $data }}"
{{- end }}
)

var	embedMap = map[string]string{
{{- range $path, $_ := .Files }}
	"{{ $path }}": {{ pathToVar $path }},
{{- end }}
}

// EmbeddedFiles returns an alphabetically sorted list of the embedded files.
func EmbeddedFiles() []string {
	var fs []string
	for f := range embedMap {
		fs = append(fs,f)
	}
	sort.Strings(fs)
	return fs
}

// EmbeddedFile returns the content of the file embedded as path.
// The function will panic if the content is not properly encoded.
func EmbeddedFile(path string) []byte {
	e, ok := embedMap[path]
	if !ok {
		return nil
	}
	d, err := base64.RawStdEncoding.DecodeString(e)
	if err != nil {
		panic(err)
	}
	return d
}

// EmbeddedFileString is a convenience function
// that works like EmbeddedFile but returns a string
// instead of a byte slice.
func EmbeddedFileString(path string) string {
	return string(EmbeddedFile(path))
}
`))
)

func readFile(path string) (data []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	return
}

func embed(c *cli.Context) error {
	files := make(map[string][]byte)

	for _, includePath := range c.StringSlice("include") {
		info, err := os.Stat(includePath)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
		if info.IsDir() {
			walkFn := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				data, err := readFile(path)
				if err != nil {
					return fmt.Errorf("readFile: %w", err)
				}
				files[path] = data

				return nil
			}
			err = filepath.Walk(includePath, walkFn)
			if err != nil {
				return fmt.Errorf("filepath.Walk: %w", err)
			}
		} else {
			data, err := readFile(includePath)
			if err != nil {
				return fmt.Errorf("readFile: %w", err)
			}
			files[includePath] = data
		}
	}

	templateData := struct {
		Package string
		Files   map[string][]byte
	}{
		Package: c.String("package"),
		Files:   files,
	}

	buf := bytes.NewBuffer(nil)
	err := fileTemplate.Execute(buf, templateData)
	if err != nil {
		return fmt.Errorf("fileTemplate.Execute: %w", err)
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format.Source: %w", err)
	}
	dest, err := os.OpenFile(c.String("destination"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer dest.Close()
	_, err = dest.Write(source)
	if err != nil {
		return fmt.Errorf("dest.Write: %w", err)
	}
	return nil
}

func main() {
	app := cli.App{
		Name: "embed",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "package",
				Usage: "name of the package the generated Go file is associated to",
				Value: "main",
			},
			&cli.StringFlag{
				Name:  "destination",
				Usage: "where to store the generated Go file",
				Value: "embeds.go",
			},
			&cli.StringSliceFlag{
				Name:     "include",
				Usage:    "paths to embed, directories are stored recursively (can be used multiple times)",
				Required: true,
			},
		},
		Action: embed,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
