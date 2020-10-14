package multimc

import (
	"bytes"
	"strings"
	"text/template"
)

var mainConfigTemplate = `
Analytics={{.Analytics}}
AnalyticsSeen=2
JavaPath={{.JavaPath}}
`[1:]

type MainConfigData struct {
	JavaPath  string
	Analytics bool
}

func GenerateMainConfig(data *MainConfigData) (string, error) {
	// MultiMC strips all left slashes
	data.JavaPath = strings.Replace(data.JavaPath, "\\", "/", -1)
	t, err := template.New("").Parse(mainConfigTemplate)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if err := t.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}
