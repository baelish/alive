package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

const header=`
<head>
  <link rel='stylesheet' type='text/css' href='/static/standard.css'/>
  <script src='/static/scripts.js'></script>
</head>
<body onresize='rightSizeBigBox()' onload='rightSizeBigBox(); keepalive()'>
<div id='big-box' class='big-box'>
  <div id='status-bar' class='status fullwidth box'>
    <p class='title'>Status</p>
	<p class='message'></p>
	<p class='lastUpdated'></p>
  </div>
`

const footer = `
  </div>
</html>`

var templates *template.Template

func handleRoot(w http.ResponseWriter, _ *http.Request) {
  _, err := fmt.Fprintf(w, header)
	if err != nil {log.Println(err)}


	for i := 0; i < len(boxes); i++ {
		err := templates.ExecuteTemplate(w, "box", boxes[i])
		if err !=nil {log.Println(err)}
	}

	_, err = fmt.Fprintf(w, footer)
	if err != nil {log.Print(err)}
}

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, `{"status":"ok"}`)
	if err != nil {log.Println(err)}
}


func handleBox(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, header)
	if err != nil {log.Print(err)}
	vars := mux.Vars(r)
	i, _ := findBoxByID(vars["id"])
	err = templates.ExecuteTemplate(w, "infoBox", boxes[i])
	if err !=nil {log.Println(err)}
	_, err = fmt.Fprintf(w, footer)
	if err != nil {log.Println(err)}
}

func loadTemplates() (err error){
	boxTemplate := `
		<div onclick='boxClick(this.id)' id='{{.ID}}' class='{{.Status}} {{.Size}} box'>
		  <p class='title'>{{.Name}}</p>
		  <p class='message'>{{.LastMessage}}</p>
		  <p class='lastUpdated'>{{.LastUpdate}}</p>
		  <p class='maxTBU'>{{.MaxTBU}}</p>
		</div>
	`
	templates, err = template.New("box").Parse(boxTemplate)
	if err !=nil {return err}

	infoBoxTemplate := `
		<div id="{{.ID}}" class="{{.Status}} fullwidth info box">
		  <h2>{{.Name}}</h2>
		  {{if .Links}}{{range .Links}}<a href="{{.URL}}" target="_blank" rel="noopener noreferrer">{{.Name}}</a><br />{{end}}{{end}}

		  <table>
			<tr><th>Last message :</th><td class="message">{{.LastMessage}}</td></tr>
			<tr><th>Last updated :</th><td class="lastUpdated">{{.LastUpdate}}</td></tr>
			{{if ne .MaxTBU ""}}<tr><th>Max TBU :</th><td>{{.MaxTBU}}</td></tr>{{end}}
			{{if ne .ExpireAfter ""}}<tr><th>Expires after :</th><td>{{.ExpireAfter}}</td></tr>{{end}}
		</div>
	`
	templates, err = templates.New("infoBox").Parse(infoBoxTemplate)
	if err !=nil {return err}

	return nil
}

func runPages() {
	err := loadTemplates()
	if err != nil {log.Printf("Unable to load templates: %v", err)}
	r := mux.NewRouter()
	r.HandleFunc("/box/{id}", handleBox)
	http.Handle("/box/", r)
	http.HandleFunc("/", handleRoot)
	http.HandleFunc ("/health", handleStatus)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.staticFilePath))))
}
