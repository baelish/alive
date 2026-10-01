package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Template names
const (
	TemplateDashboard = "dashboard"
	TemplateInfoPage  = "infoPage"
	TemplateHead      = "head"
	TemplateStatusBar = "statusBar"
	TemplateBoxGrid   = "boxGrid"
	TemplateBox       = "box"
	TemplateDrawerInfo = "drawerInfo"
	TemplateBoxInfo   = "boxInfo"
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
<div class="drawer-header-wrapper {{ .Status }}">
	<div class="drawer-header">
		<h2 class="drawer-title">{{ if .DisplayName }}{{ .DisplayName }}{{ else }}{{ .Name }}{{ end }}</h2>
		<div class="drawer-actions"><a class="drawer-open-page" href="/box/{{ .ID }}" target="_blank" rel="noopener noreferrer" aria-label="Open full box page" title="Open full box page"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg></a><button class="drawer-close" type="button" aria-label="Close details" onclick="closeDrawer()"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></svg></button></div>
	</div>
</div>
<div class="drawer-tabs">
	<button class="drawer-tab active" onclick="switchDrawerTab(event, 'overview')">Overview</button>
	<button class="drawer-tab" onclick="switchDrawerTab(event, 'events')">Events</button>
</div>
<div class="drawer-content-area">
	<div id="drawer-tab-overview" class="drawer-tab-content active">
		<div class="drawer-section drawer-status-details">
			<h3>Status Details</h3>
			{{ if .Description }}<p class="drawer-description">{{ .Description }}</p>{{ end }}
			<dl class="drawer-meta">
				<dt>ID</dt><dd><span class="drawer-metadata-value">{{ .ID }}</span></dd>
				<dt>Size</dt><dd><span class="drawer-metadata-value">{{ .Size }}</span></dd>
				{{ if .MaxTBU }}<dt class="drawer-max-tbu-label">Max TBU</dt><dd class="drawer-max-tbu"><svg class="metadata-icon" viewBox="0 0 24 24" aria-label="Max TBU"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span class="drawer-metadata-value">{{ .MaxTBU }}</span></dd>{{ end }}
				{{ if .ExpireAfter }}<dt class="drawer-expire-after-label">Expires after</dt><dd class="drawer-expire-after"><svg class="metadata-icon" viewBox="0 0 24 24" aria-label="Expires after"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span class="drawer-metadata-value">{{ .ExpireAfter }}</span></dd>{{ end }}
			</dl>
		</div>
		{{ if .Info }}<div class="drawer-section drawer-metadata"><h3>Metadata</h3><dl class="drawer-meta">{{ range $key, $value := .Info }}<dt>{{ $key }}</dt><dd><span class="drawer-metadata-value">{{ $value }}</span></dd>{{ end }}</dl></div>{{ end }}
		{{ if .Links }}<div class="drawer-section drawer-links"><h3>Links</h3>{{ range .Links }}<a href="{{ .URL }}" target="_blank" rel="noopener noreferrer"><svg class="external-link-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg><span>{{ .Name }}</span></a><br />{{ end }}</div>{{ end }}
		<div class="drawer-section drawer-latest-update"><h3>Latest update</h3><div class="event-log-card {{ .Status | EventStatus }}"><div class="event-top"><strong class="event-status">{{ .Status | EventStatus }}</strong><div class="event-leading">{{ if .Messages }}{{ with index .Messages 0 }}{{ if .MaxTBU }}<span class="event-indicator" aria-label="Max TBU: {{ .MaxTBU }}" title="Max TBU: {{ .MaxTBU }}"><svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span>{{ .MaxTBU }}</span></span>{{ end }}{{ if .ExpireAfter }}<span class="event-indicator" aria-label="Expires after: {{ .ExpireAfter }}" title="Expires after: {{ .ExpireAfter }}"><svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span>{{ .ExpireAfter }}</span></span>{{ end }}{{ end }}{{ end }}</div></div><div class="event-message">{{ $.LastMessage }}</div><small>{{ $.LastUpdate.Format "2006-01-02T15:04:05.000Z07:00" }}</small></div></div>
	</div>
	<div id="drawer-tab-events" class="drawer-tab-content">
		<div class="drawer-section drawer-events"><h3>Event log</h3><div class="event-log">{{ range .Messages }}<div class="event-log-card {{ .Status | EventStatus }}"><div class="event-top"><strong class="event-status">{{ .Status | EventStatus }}</strong><div class="event-leading">{{ if .MaxTBU }}<span class="event-indicator" aria-label="Max TBU: {{ .MaxTBU }}" title="Max TBU: {{ .MaxTBU }}"><svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span>{{ .MaxTBU }}</span></span>{{ end }}{{ if .ExpireAfter }}<span class="event-indicator" aria-label="Expires after: {{ .ExpireAfter }}" title="Expires after: {{ .ExpireAfter }}"><svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span>{{ .ExpireAfter }}</span></span>{{ end }}</div></div><div class="event-message">{{ .Message }}</div><small>{{ .TimeStamp.Format "2006-01-02T15:04:05.000Z07:00" }}</small></div>{{ end }}</div></div>
	</div>
</div>
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

// templateFuncMap returns the template functions available to templates
func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"EventStatus": eventStatusFunc,
		"ToUpper":     strings.ToUpper,
	}
}

// eventStatusFunc converts status to display string
func eventStatusFunc(status interface{}) string {
	statusStr := fmt.Sprintf("%v", status)
	if statusStr == "noUpdate" {
		return "red"
	}
	if statusStr == "" {
		return "grey"
	}
	return statusStr
}

func loadTemplates() error {
	// Start with base template and func map
	root := template.New("root").Funcs(templateFuncMap())

	// Parse all template strings into a single tree
	var err error
	templates, err = root.Parse(generic + boxGrid + drawerInfo + boxInfo + dashboard + infoPage)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}
	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Get all boxes from store (thread-safe)
	boxes := boxStore.GetAll()

	err := templates.ExecuteTemplate(w, TemplateDashboard, boxes)
	if err != nil {
		logger.Error("failed to execute dashboard template", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, `{"status":"ok"}`)
	if err != nil {
		logger.Error("failed to write status response", zap.Error(err))
	}
}

func handleBox(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Box ID is required", http.StatusBadRequest)
		return
	}

	// Get box from store (thread-safe)
	box, err := boxStore.GetByID(id)
	if err != nil {
		logger.Error("failed to get box", zap.String("id", id), zap.Error(err))
		http.Error(w, "Box not found", http.StatusNotFound)
		return
	}

	err = templates.ExecuteTemplate(w, TemplateInfoPage, box)
	if err != nil {
		logger.Error("failed to execute info page template", zap.String("id", id), zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func runDashboard(ctx context.Context) {
	if options.Debug {
		logger.Info("Starting Dashboard")
	}

	// Load templates
	if err := loadTemplates(); err != nil {
		logger.Fatal("Failed to load templates", zap.Error(err))
	}

	// Setup routes
	r := chi.NewRouter()
	r.Get("/box/{id}", handleBox)
	http.Handle("/box/", r)
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health", handleStatus)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(options.StaticPath))))

	// Start server
	listenAddr := fmt.Sprintf(":%s", options.SitePort)
	logger.Info("HTTP server listening", zap.String("address", listenAddr))

	server := &http.Server{
		Addr:    listenAddr,
		Handler: nil, // Use DefaultServeMux
	}

	// Run server with context cancellation support
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down HTTP server")
		if err := server.Close(); err != nil {
			logger.Error("Error closing server", zap.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}
}
