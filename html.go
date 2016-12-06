package chmoufrack

import (
	"bytes"
	"path/filepath"
	"text/template"
)

func html_content(ts TemplateStruct, content *bytes.Buffer) (err error) {
	t, err := template.ParseFiles(filepath.Join(STATIC_DIR, "templates", "content.tmpl"))
	if err != nil {
		return
	}
	err = t.Execute(content, ts)
	if err != nil {
		return
	}
	return
}

func html_main_template(program_name, content string, outputWriter *bytes.Buffer) (err error) {
	dico := map[string]string{
		"Content":     content,
		"ProgramName": program_name,
	}

	t, err := template.ParseFiles(filepath.Join(STATIC_DIR, "templates", "template.tmpl"))
	if err != nil {
		return
	}
	err = t.Execute(outputWriter, dico)
	if err != nil {
		return
	}
	return
}

func HTML_generate(program_name string, rounds []Workout, outputWriter *bytes.Buffer) (err error) {
	var content bytes.Buffer
	var ts TemplateStruct

	for _, workout := range rounds {
		if ts, err = GenerateProgram(workout); err != nil {
			return
		}
		if err = html_content(ts, &content); err != nil {
			return
		}
	}

	err = html_main_template(program_name, content.String(), outputWriter)
	return
}
