package chmoufrack

import (
	"bytes"
	"path/filepath"
	"text/template"
)

func htmlContent(ts TemplateStruct, content *bytes.Buffer) (err error) {
	t, err := template.ParseFiles(filepath.Join(StaticDir, "templates", "content.tmpl"))
	if err != nil {
		return
	}
	err = t.Execute(content, ts)
	if err != nil {
		return
	}
	return
}

func htmlMainTemplate(programName, content string, outputWriter *bytes.Buffer) (err error) {
	dico := map[string]string{
		"Content":     content,
		"ProgramName": programName,
	}

	t, err := template.ParseFiles(filepath.Join(StaticDir, "templates", "template.tmpl"))
	if err != nil {
		return
	}
	err = t.Execute(outputWriter, dico)
	if err != nil {
		return
	}
	return
}

// HTMLGenerate Generate the HTML of all the programs and rounds -- Sounds to be
// depreacred.
func HTMLGenerate(programName string, rounds []Workout, outputWriter *bytes.Buffer) (err error) {
	var content bytes.Buffer
	var ts TemplateStruct

	for _, workout := range rounds {
		if ts, err = GenerateProgram(workout, TargetVma); err != nil {
			return
		}
		if err = htmlContent(ts, &content); err != nil {
			return
		}
	}

	err = htmlMainTemplate(programName, content.String(), outputWriter)
	return
}
