package multimc

import (
	"bytes"
	"text/template"
)

var instanceTemplate = `
InstanceType=OneSix
iconKey=flame
name={{.Name}}
`[1:]

type InstanceConfigData struct {
	Name string
}

func GenerateInstanceConfig(data *InstanceConfigData) (string, error) {
	t, err := template.New("").Parse(instanceTemplate)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if err := t.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}
