package web

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zond/god/templates"
	htmlTemplate "html/template"
	"net/http"
	"reflect"
	"strings"
	textTemplate "text/template"
)

var apiMethods string

func getFormat(t reflect.Type) interface{} {
	if t.Kind() == reflect.Struct {
		result := make(map[string]interface{})
		var field reflect.StructField
		for i := 0; i < t.NumField(); i++ {
			field = t.Field(i)
			if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Uint8 {
				result[field.Name] = "[]byte"
			} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
				result[field.Name] = "int"
			} else {
				result[field.Name] = field.Type.Name()
			}
		}
		return result
	}
	return t.Name()
}

func SetApi(t reflect.Type) {
	m := make(map[string]map[string]interface{})
	var meth reflect.Method
	var in reflect.Type
	for i := 0; i < t.NumMethod(); i++ {
		meth = t.Method(i)
		if strings.ToUpper(string(meth.Name[0])) == string(meth.Name[0]) && meth.Type.NumIn() == 3 {
			in = meth.Type.In(1)
			m[meth.Name] = map[string]interface{}{
				"name":      meth.Name,
				"parameter": getFormat(in),
			}
		}
	}
	var bts []byte
	var err error
	if bts, err = json.Marshal(m); err != nil {
		panic(err)
	}
	apiMethods = string(bts)
}

type baseData struct {
	Timestamp int64
}

func getBaseData(w http.ResponseWriter, r *http.Request) baseData {
	return baseData{
		Timestamp: templates.Timestamp,
	}
}
func (self baseData) T() string {
	return fmt.Sprint(self.Timestamp)
}
func (self baseData) ApiMethods() string {
	return apiMethods
}

func allCss(w http.ResponseWriter, r *http.Request) {
	data := getBaseData(w, r)
	w.Header().Set("Cache-Control", "public, max-age=864000")
	w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	renderText(w, r, templates.CSS, "bootstrap.min.css", data)
	renderText(w, r, templates.CSS, "common.css", data)
}

func allJs(w http.ResponseWriter, r *http.Request) {
	data := getBaseData(w, r)
	w.Header().Set("Cache-Control", "public, max-age=864000")
	w.Header().Set("Content-Type", "application/javascript; charset=UTF-8")
	renderText(w, r, templates.JS, "underscore-min.js", data)
	renderText(w, r, templates.JS, "jquery-1.8.3.min.js", data)
	renderText(w, r, templates.JS, "bootstrap.min.js", data)
	renderText(w, r, templates.JS, "easeljs-0.5.0.min.js", data)
	renderText(w, r, templates.JS, "jquery.websocket-0.0.1.js", data)
	renderText(w, r, templates.JS, "big.min.js", data)
	renderText(w, r, templates.JS, "jquery.base64.js", data)
	renderText(w, r, templates.JS, "god.js", data)
}

func renderHtml(w http.ResponseWriter, r *http.Request, templates *htmlTemplate.Template, template string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := templates.ExecuteTemplate(w, template, data); err != nil {
		panic(fmt.Errorf("While rendering HTML: %v", err))
	}
}

func renderText(w http.ResponseWriter, r *http.Request, templates *textTemplate.Template, template string, data interface{}) {
	if err := templates.ExecuteTemplate(w, template, data); err != nil {
		panic(fmt.Errorf("While rendering text: %v", err))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	renderHtml(w, r, templates.HTML, "index.html", getBaseData(w, r))
}

func Route(handler websocket.Handler, router *mux.Router) {
	router.HandleFunc("/js/{ver}/all.js", allJs)
	router.HandleFunc("/css/{ver}/all.css", allCss)
	router.Path("/ws").Handler(handler)
	router.HandleFunc("/", index)
}
