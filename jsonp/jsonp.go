package jsonp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
)

type jsonpWriter struct {
	negroni.ResponseWriter
	callback    string
	wroteHeader bool
}

type jsonpHandler struct {
}

func NewJsonp() *jsonpHandler {
	return &jsonpHandler{}
}

func (jsonp *jsonpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	callback := r.URL.Query().Get("callback")
	if callback != "" {
		nw := negroni.NewResponseWriter(w)
		newWriter := &jsonpWriter{ResponseWriter: nw,
			callback:    callback,
			wroteHeader: false}
		next(newWriter, r)
	} else {
		next(w, r)
	}
}

type jsonpWrap struct {
	Data string `json:"data"`
}

func (jsonp *jsonpWriter) WriteHeader(code int) {
	headers := jsonp.ResponseWriter.Header()
	headers.Set("Content-Type", "application/javascript")
	jsonp.ResponseWriter.WriteHeader(code)
	jsonp.wroteHeader = true
}

func (jsonp *jsonpWriter) Write(b []byte) (int, error) {
	if !jsonp.wroteHeader {
		jsonp.WriteHeader(http.StatusOK)
	}
	var callbackFunc string
	if json.Valid(b) {
		callbackFunc = fmt.Sprintf("%s(%s)", jsonp.callback, string(b))
	} else {
		json, err := json.Marshal(jsonpWrap{string(b)})
		if err != nil {
			return -1, err
		} else {
			callbackFunc = fmt.Sprintf("%s(%s)", jsonp.callback, string(json))
		}
	}
	return jsonp.ResponseWriter.Write([]byte(callbackFunc))
}
