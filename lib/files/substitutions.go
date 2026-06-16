package files 

import (
	"html/template"
	"io"
	"os"
)

// TODO this should be more flexible
type Substitutions struct {
	AppName string
	DbPort  int
}

func (s *Substitutions) SubstituteFile(filename string, writer io.Writer) error {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Parse(string(contents))
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, s)
	if err != nil {
		return err
	}

	return nil
}
