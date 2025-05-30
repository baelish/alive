package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const dashboard = `
{{ define "dashboard" }}
<!DOCTYPE html>
{{ template "head" . }}
<html>
	<body onresize='rightSizeBigBox("dashboard");' onload='rightSizeBigBox("dashboard"); keepalive();'>
		<input type="hidden" id="refreshed" value="no">
		<div id='big-box' class='big-box'>
			{{ template "statusBar" . }}
			{{ template "boxGrid" . }}
		</div>
	</body>
</html>
{{ end }}`

const infoPage = `
{{ define "infoPage" }}
<!DOCTYPE html>
{{ template "head" . }}
<html>
	<body onresize='rightSizeBigBox();' onload='rightSizeBigBox(); keepalive();'>
	<input type="hidden" id="refreshed" value="no">
		<div id='big-box' class='big-box'>
		    {{ template "statusBar" . }}
		    {{ template "boxInfo" . }}
		</div>
	</body>
</html>
{{ end }}`

const generic = `
{{ define "head" }}
<head>
  <link rel='stylesheet' type='text/css' href='/static/standard.css'/>
  <script src='/static/scripts.js'></script>
</head>
{{ end }}

{{ define "statusBar" }}
<div id='status-bar' class='status fullwidth box'>
  <p class='title'>Status</p>
  <p class='tooltip' id='tooltip' display="none"></p>
<p class='message' display="none"></p>
<p class='lastUpdated'></p>
</div>
{{ end }}`

const boxGrid = `
{{ define "boxGrid" }}
  {{ range . }}
    {{ template "box" . }}
  {{ end }}
{{ end }}

{{ define "box" }}
<div onclick='boxClick(this.id)' onmouseover='boxHover("{{ .Name }}")' onmouseout='boxOut()' id='{{ .ID }}' class='{{ .Status }} {{ .Size }} box'>
    <p class='title'>{{ if .DisplayName }}{{ .DisplayName }}{{ else }}{{ .Name }}{{ end }}</p>
    <p class='message'>{{ .LastMessage }}</p>
    <p class='lastUpdated'>{{ .LastUpdate.Format "2006-01-02T15:04:05.000Z07:00"}}</p>
    <p class='maxTBU'>{{ .MaxTBU }}</p>
    <p class='expireAfter'>{{ .ExpireAfter }}</p>
</div>
{{ end }}`

const boxInfo = `
{{ define "boxInfo" }}
<div id="{{ .ID }}" class="{{ .Status }} fullwidth info box">
  <h2>{{ .Name }}</h2>
  {{ if .Links }}{{ range .Links }}<a href="{{ .URL }}" target="_blank" rel="noopener noreferrer">{{ .Name }}</a><br />{{ end }}{{ end }}

  <table>
  <tr><th>ID:</th><td>{{ .ID }}</td></tr>
  {{ if .DisplayName }}<tr><th>Display name:</th><td>{{ .DisplayName }}</td></tr>{{ end }}
  {{ if .Description }}<tr><th>Description:</th><td>{{ .Description }}</td></tr>{{ end }}
  {{ if .Info }}{{ range $key, $value := .Info }}<tr><th>{{ $key }}:</th><td>{{ $value }}</td></tr>{{ end }}{{ end }}
  <tr><th>Last message:</th><td class="message">{{ .LastMessage }}</td></tr>
  <tr><th>Last updated:</th><td class="lastUpdated">{{ .LastUpdate.Format "2006-01-02T15:04:05.000Z07:00" }}</td></tr>
  <tr class="maxTBU" {{ if eq .MaxTBU.Duration 0}}style="display: none;"{{ end }}><th>Max TBU:</th><td>{{ .MaxTBU }}</td></tr>
  <tr class="expireAfter" {{ if eq .ExpireAfter.Duration 0 }}style="display: none;"{{ end }}><th>Expires after:</th><td>{{ .ExpireAfter }}</td></tr>
  <tr><th>Previous Messages:</th><td><ul class="previousMessages">{{ range $m := .Messages }}<li>{{ $m.TimeStamp.Format "2006-01-02T15:04:05.000Z07:00" }}: {{ $m.Status | ToUpper }} ({{ $m.Message }})</li>{{ end }}</ul></td></tr>

</div>
{{ end }}`
var templates *template.Template

func loadTemplates() (err error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	// Start with base template and func map
	root := template.New("root").Funcs(funcMap)

	// Parse all template strings into a single tree
	templates, err = root.Parse(generic + boxGrid + boxInfo + dashboard + infoPage)
	return err
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	err := templates.ExecuteTemplate(w, "dashboard", boxes)
	if err != nil {
		logger.Error(err.Error())
	}
}

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, `{"status":"ok"}`)
	if err != nil {
		logger.Error(err.Error())
	}
}

func handleBox(w http.ResponseWriter, r *http.Request) {
	i, err := findBoxByID(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "infoPage", boxes[i])
	if err != nil {
		logger.Error(err.Error())
	}
}

func runDashboard(_ context.Context) {
	if options.Debug {
		logger.Info("Starting Dashboard")
	}

	err := loadTemplates()
	if err != nil {
		logger.Fatal("Failed to load templates", zap.Error(err))
	}
	r := chi.NewRouter()
	r.HandleFunc("/box/{id}", handleBox)
	http.Handle("/box/", r)
	http.HandleFunc("/", handleRoot)

	http.HandleFunc("/health", handleStatus)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(options.StaticPath))))

	logger.Info("listening", zap.String("port", options.SitePort))
	listenOn := fmt.Sprintf(":%s", options.SitePort)

	log.Fatal(http.ListenAndServe(listenOn, nil))
}
