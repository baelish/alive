package server

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
	<body onresize='rightSizeBigBox("dashboard");' onload='rightSizeBigBox("dashboard"); keepalive(); initDashboard();'>
		<input type="hidden" id="refreshed" value="no">
		<div id='big-box' class='big-box'>
			{{ template "statusBar" . }}
			<div class="dashboard-viewport">
				<div id="boxes-grid" class="boxes-grid">{{ template "boxGrid" . }}</div>
				<div id="drawer-backdrop" class="drawer-backdrop" onclick="closeDrawer()"></div>
				<aside id="side-drawer-panel" class="side-drawer-panel" aria-hidden="true">
						<div id="drawer-resize-handle" class="drawer-resize-handle" title="Resize details panel"></div>
					<div id="drawer-content"></div>
				</aside>
			</div>
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
	<link rel='icon' type='image/svg+xml' href='/static/favicon.svg'/>
  <script src='/static/scripts.js'></script>
</head>
{{ end }}

{{ define "statusBar" }}
<div id='status-bar' class='status fullwidth box'>
  <p class='title'>Status</p>
	<p class='tooltip' id='tooltip' display="none"></p>
	<span class='status-center'><p class='message' display="none"></p></span>
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
<div onclick='boxClick(this.id)' onmouseover='boxHover("{{ .Name }}")' onmouseout='boxOut()' id='{{ .ID }}' data-sort-name='{{ .Name }}' class='{{ .Status }} {{ .Size }} box'>
    <p class='title'>{{ if .DisplayName }}{{ .DisplayName }}{{ else }}{{ .Name }}{{ end }}</p>
	<p class='box-id'>{{ .ID }}</p>
    <p class='maxTBU'>{{ .MaxTBU }}</p>
    <p class='expireAfter'>{{ .ExpireAfter }}</p>
</div>
<template id='box-detail-{{ .ID }}'>{{ template "drawerInfo" . }}</template>
{{ end }}`

const drawerInfo = `
{{ define "drawerInfo" }}
<div class="drawer-status-strip {{ .Status }}"></div>
<div class="drawer-header">
	<h2 class="drawer-title">{{ if .DisplayName }}{{ .DisplayName }}{{ else }}{{ .Name }}{{ end }}</h2>
	<div class="drawer-actions"><a class="drawer-open-page" href="/box/{{ .ID }}" target="_blank" rel="noopener noreferrer" aria-label="Open full box page" title="Open full box page"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg></a><button class="drawer-close" type="button" aria-label="Close details" onclick="closeDrawer()"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></svg></button></div>
</div>
<div class="drawer-section drawer-metadata">
	{{ if .Description }}<p>{{ .Description }}</p>{{ end }}
	<dl class="drawer-meta">
		<dt>ID</dt><dd><span class="drawer-metadata-value">{{ .ID }}</span></dd>
			<dt>Size</dt><dd><span class="drawer-metadata-value">{{ .Size }}</span></dd>
		{{ if .MaxTBU }}<dt class="drawer-max-tbu-label">Max TBU</dt><dd class="drawer-max-tbu"><span class="drawer-metadata-value">{{ .MaxTBU }}</span><svg class="metadata-icon" viewBox="0 0 24 24" aria-label="Max TBU"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg></dd>{{ end }}
		{{ if .ExpireAfter }}<dt class="drawer-expire-after-label">Expires after</dt><dd class="drawer-expire-after"><span class="drawer-metadata-value">{{ .ExpireAfter }}</span><svg class="metadata-icon" viewBox="0 0 24 24" aria-label="Expires after"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg></dd>{{ end }}
	</dl>
{{ if .Info }}<h3>Metadata</h3><dl class="drawer-meta">{{ range $key, $value := .Info }}<dt>{{ $key }}</dt><dd><span class="drawer-metadata-value">{{ $value }}</span></dd>{{ end }}</dl>{{ end }}
</div>
{{ if .Links }}<div class="drawer-section drawer-links"><h3>Links</h3>{{ range .Links }}<a href="{{ .URL }}" target="_blank" rel="noopener noreferrer"><svg class="external-link-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg><span>{{ .Name }}</span></a><br />{{ end }}</div>{{ end }}
<div class="drawer-section drawer-latest-update"><h3>Latest update</h3><p class="drawer-message">{{ .LastMessage }}</p><time class="drawer-last-updated">{{ .LastUpdate.Format "2006-01-02T15:04:05.000Z07:00" }}</time></div>
<div class="drawer-section drawer-events"><h3>Event log</h3><div class="event-log">{{ range .Messages }}<div class="event-log-card {{ .Status | EventStatus }}"><div class="event-top"><strong class="event-status">{{ .Status | EventStatus }}</strong><div class="event-leading"></div></div><div class="event-message">{{ .Message }}</div><small>{{ .TimeStamp.Format "2006-01-02T15:04:05.000Z07:00" }}</small></div>{{ end }}</div></div>
{{ end }}`

const boxInfo = `
{{ define "boxInfo" }}
<div id="{{ .ID }}" class="{{ .Status }} fullwidth info box">
  <h2>{{ .Name }}</h2>
  {{ if .Links }}{{ range .Links }}<a href="{{ .URL }}" target="_blank" rel="noopener noreferrer">{{ .Name }}</a><br />{{ end }}{{ end }}

  <table>
  {{ if .DisplayName }}<tr><th>Display name:</th><td>{{ .DisplayName }}</td></tr>{{ end }}
  {{ if .Description }}<tr><th>Description:</th><td>{{ .Description }}</td></tr>{{ end }}
  <tr><th>ID:</th><td>{{ .ID }}</td></tr>
  <tr><th>Status:</th><td>{{ .Status }}</td></tr>
  <tr><th>Size:</th><td>{{ .Size }}</td></tr>
  {{ if .Info }}<tr><th>Info:</th><td><table>{{ range $key, $value := .Info }}<tr><th>{{ $key }}:</th><td>{{ $value }}</td></tr>{{ end }}</table></td></tr>{{ end }}
  <tr><th>Last message:</th><td class="message">{{ .LastMessage }}</td></tr>
  <tr><th>Last updated:</th><td class="lastUpdated">{{ .LastUpdate.Format "2006-01-02T15:04:05.000Z07:00" }}</td></tr>
  <tr class="maxTBU" {{ if not .MaxTBU }}style="display: none;"{{ end }}><th>Max TBU:</th><td>{{ .MaxTBU }}</td></tr>
  <tr class="expireAfter" {{ if not .ExpireAfter }}style="display: none;"{{ end }}><th>Expires after:</th><td>{{ .ExpireAfter }}</td></tr>
  <tr><th>Previous Messages:</th><td><ul class="previousMessages">{{ range $m := .Messages }}<li>{{ $m.TimeStamp.Format "2006-01-02T15:04:05.000Z07:00" }}: {{ $m.Status | ToUpper }} ({{ $m.Message }})</li>{{ end }}</ul></td></tr>

</div>
{{ end }}`

var templates *template.Template

func loadTemplates() (err error) {
	funcMap := template.FuncMap{
		"EventStatus": func(status string) string {
			if status == "noUpdate" {
				return "red"
			}
			if status == "" {
				return "grey"
			}
			return status
		},
		"ToUpper": strings.ToUpper,
	}

	// Start with base template and func map
	root := template.New("root").Funcs(funcMap)

	// Parse all template strings into a single tree
	templates, err = root.Parse(generic + boxGrid + drawerInfo + boxInfo + dashboard + infoPage)
	return err
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	// Get all boxes from store (thread-safe)
	boxes := boxStore.GetAll()
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
	id := chi.URLParam(r, "id")

	// Get box from store (thread-safe)
	box, err := boxStore.GetByID(id)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = templates.ExecuteTemplate(w, "infoPage", box)
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
