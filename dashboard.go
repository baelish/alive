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

const header = `
<!DOCTYPE html>
<html>
  <head>
    <link rel='stylesheet' type='text/css' href='/static/standard.css'/>
    <script src='/static/scripts.js'></script>
  </head>
  <body onresize='rightSizeBigBox();' onload='rightSizeBigBox(); keepalive();'>
  <input type="hidden" id="refreshed" value="no">
  <div id='big-box' class='big-box'>
    <div id='status-bar' class='status fullwidth box'>
      <p class='title'>Status</p>
        <p class='tooltip' id='tooltip' display="none"></p>
    <p class='message' display="none"></p>
    <p class='lastUpdated'></p>
    </div>
`

const footer = `
  </div>
</html>`

var templates *template.Template

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, header)
	if err != nil {
		logger.Error(err.Error())
	}

	for i := 0; i < len(boxes); i++ {
		err := templates.ExecuteTemplate(w, "box", boxes[i])
		if err != nil {
			logger.Error(err.Error())
		}
	}

	_, err = fmt.Fprint(w, footer)
	if err != nil {
		log.Print(err)
	}
}

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, `{"status":"ok"}`)
	if err != nil {
		logger.Error(err.Error())
	}
}

func handleBox(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, header)
	if err != nil {
		log.Print(err)
	}

	i, err := findBoxByID(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "infoBox", boxes[i])
	if err != nil {
		logger.Error(err.Error())
	}

	_, err = fmt.Fprint(w, footer)
	if err != nil {
		logger.Error(err.Error())
	}
}

func loadTemplates() (err error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	boxTemplate := `
    <div onclick='boxClick(this.id)' onmouseover='boxHover("{{.Name}}")' onmouseout='boxOut()' id='{{.ID}}' class='{{.Status}} {{.Size}} box'>
        <p class='title'>{{if .DisplayName}}{{.DisplayName}}{{else}}{{.Name}}{{end}}</p>
        <p class='message'>{{.LastMessage}}</p>
        <p class='lastUpdated'>{{.LastUpdate.Format "2006-01-02T15:04:05.000Z07:00"}}</p>
        <p class='maxTBU'>{{.MaxTBU}}</p>
        <p class='expireAfter'>{{.ExpireAfter}}</p>
    </div>
  `
	templates, err = template.New("box").Parse(boxTemplate)
	if err != nil {
		return err
	}

	infoBoxTemplate := `
    <div id="{{.ID}}" class="{{.Status}} fullwidth info box">
      <h2>{{.Name}}</h2>
      {{if .Links}}{{range .Links}}<a href="{{.URL}}" target="_blank" rel="noopener noreferrer">{{.Name}}</a><br />{{end}}{{end}}

      <table>
      <tr><th>ID :</th><td>{{.ID}}</td></tr>
      {{if .DisplayName}}<tr><th>Display name :</th><td>{{.DisplayName}}</td></tr>{{end}}
      {{if .Description}}<tr><th>Description :</th><td>{{.Description}}</td></tr>{{end}}
      <tr><th>Last message :</th><td class="message">{{.LastMessage}}</td></tr>
      <tr><th>Last updated :</th><td class="lastUpdated">{{.LastUpdate.Format "2006-01-02T15:04:05.000Z07:00" }}</td></tr>
      <tr class="maxTBU" {{if eq .MaxTBU.Duration 0}}style="display: none;"{{end}}><th>Max TBU:</th><td>{{.MaxTBU}}</td></tr>
      <tr class="expireAfter" {{if eq .ExpireAfter.Duration 0 }}style="display: none;"{{end}}><th>Expires after :</th><td>{{.ExpireAfter}}</td></tr>
      <tr><th>Previous Messages:</th><td><ul class="previousMessages">{{range $m := .Messages}}<li>{{$m.TimeStamp.Format "2006-01-02T15:04:05.000Z07:00" }}: {{$m.Status | ToUpper}} ({{$m.Message}})</li>{{end}}</ul></td></tr>

    </div>
  `
	templates, err = templates.New("infoBox").Funcs(funcMap).Parse(infoBoxTemplate)
	if err != nil {
		return err
	}

	return nil
}

func runDashboard(_ context.Context) {
	if options.Debug {
		logger.Info("Starting Dashboard")
	}
	err := loadTemplates()
	if err != nil {
		logger.Error(err.Error())
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
