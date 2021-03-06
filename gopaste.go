package main

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz"
const idLength = 6
const dataDir = "data"
const pasteMaxSize = 1048576
const pasteMaxSizeDesc = "1MB"

var validPath = regexp.MustCompile(
	fmt.Sprintf("^/([a-zA-Z0-9]{%d})$", idLength))
var filenameReg = regexp.MustCompile(
	fmt.Sprintf("^%s/([a-zA-Z0-9]{%d})-([a-z0-9]+)-([0-9]+)$",
		dataDir,
		idLength))

var languages = map[string]map[string]string{
	"text": {
		"name":  "Plain/Text",
		"brush": "shBrushPlain.js",
	},
	"bash": {
		"name":  "Bash/Shell",
		"brush": "shBrushBash.js",
	},
	"c": {
		"name":  "C/C++",
		"brush": "shBrushCpp.js",
	},
	"diff": {
		"name":  "Diff",
		"brush": "shBrushDiff.js",
	},
	"erlang": {
		"name":  "Erlang",
		"brush": "shBrushErlang.js",
	},
	"html": {
		"name":  "HTML",
		"brush": "shBrushXml.js",
	},
	"js": {
		"name":  "JavaScript",
		"brush": "shBrushJScript.js",
	},
	"perl": {
		"name":  "Perl",
		"brush": "shBrushPerl.js",
	},
	"php": {
		"name":  "PHP",
		"brush": "shBrushPhp.js",
	},
	"python": {
		"name":  "Python",
		"brush": "shBrushPython.js",
	},
	"ruby": {
		"name":  "Ruby",
		"brush": "shBrushRuby.js",
	},
	"sql": {
		"name":  "SQL",
		"brush": "shBrushSql.js",
	},
	"xml": {
		"name":  "XML",
		"brush": "shBrushXml.js",
	},
}

var templates = template.Must(template.ParseFiles("templates/add.html",
	"templates/share.html",
	"templates/view.html"))

type Paste struct {
	Name       string
	Language   string
	Expiration int64
	Content    template.HTML
}

func (p *Paste) save() error {
	filename := fmt.Sprintf("%s/%s-%s-%d",
		dataDir, p.Name, p.Language, p.Expiration)
	return ioutil.WriteFile(filename, []byte(p.Content), 0600)
}

func loadPaste(name string) (*Paste, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/%s-*", dataDir, name))
	if err != nil {
		return nil, err
	}
	if files == nil {
		return nil, errors.New("Not found.")
	}
	matches := filenameReg.FindSubmatch([]byte(files[0]))
	if matches == nil {
		return nil, errors.New("Pattern not found.")
	}
	Content, err := ioutil.ReadFile(files[0])
	if err != nil {
		return nil, err
	}
	return &Paste{
		Name:     name,
		Language: string(matches[2]),
		Content:  template.HTML(Content)}, nil
}

func GenerateID() string {
	buf := make([]byte, idLength)
	for i := 0; i < idLength; i++ {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func showHandler(w http.ResponseWriter, r *http.Request, name string) {
	p, err := loadPaste(name)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, "view", struct {
		Host  string
		P     *Paste
		Brush string
	}{r.Host, p, languages[p.Language]["brush"]})
}

//r.URL.Host
func pasteitHandler(w http.ResponseWriter, r *http.Request) {
	name := GenerateID()
	pasteLanguage := "text"
	if _, exist := languages[r.FormValue("lang")]; exist {
		pasteLanguage = r.FormValue("lang")
	}

	pasteExpiration := time.Now().Unix()

	if formExpiration, err := strconv.ParseInt(r.FormValue("expire"), 10, 64); err != nil &&
		formExpiration > 3600 && formExpiration < 604800 {
		pasteExpiration = pasteExpiration + formExpiration
	} else {
		pasteExpiration = pasteExpiration + 3600
	}

	p := &Paste{
		Name:       name,
		Language:   pasteLanguage,
		Expiration: pasteExpiration,
		Content:    template.HTML(html.EscapeString(r.FormValue("paste"))),
	}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "share", struct{ Host, Name string }{r.Host, p.Name})
}

func frontHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		showHandler(w, r, m[1])
	} else if len(r.URL.Path) > 1 {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		pasteLen := len(r.FormValue("paste"))
		errorMsg := ""
		if pasteLen > 0 && pasteLen <= pasteMaxSize {
			pasteitHandler(w, r)
			return
		} else if pasteLen > pasteMaxSize {
			errorMsg = "Paste refused, too big. Size limit is " +
				pasteMaxSizeDesc
		}
		renderTemplate(w, "add", struct {
			Langs    map[string]map[string]string
			ErrorMsg string
		}{
			languages,
			errorMsg})
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", frontHandler)
	http.ListenAndServe(":8080", nil)
}

// vim: noexpandtab
