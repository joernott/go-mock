package server

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/joernott/mock/rules"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var ruleset *rules.Ruleset

func Router(Port int, Ruleset *rules.Ruleset) {
	ruleset = Ruleset
	router := httprouter.New()
	router.GET("/*path", ApiGET)
	router.PUT("/*path", ApiPUT)
	router.POST("/*path", ApiPOST)
	router.DELETE("/*path", ApiDELETE)
	log.Debug("Router initialized")
	p := ":" + strconv.Itoa(Port)
	if _, err := os.Stat("server.crt"); err == nil {
		if _, err := os.Stat("server.key"); os.IsNotExist(err) {
			log.Error("Found server.crt but missing server.key")
			fmt.Println("Found server.crt but missing server.key")
			os.Exit(10)
		}
		log.Fatal(http.ListenAndServeTLS(p, "server.crt", "server.key", router))

	} else {
		log.Fatal(http.ListenAndServe(p, router))
	}
}

func ApiGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	method := "GET"
	logger := log.WithFields(log.Fields{
		"Path":   path,
		"Query":  query,
		"Method": method,
	})

	logger.Debug("Call ApiGET")
	respond(method, path, query, w, logger)
}

func ApiPOST(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	method := "POST"
	logger := log.WithFields(log.Fields{
		"Path":   path,
		"Query":  query,
		"Method": method,
	})

	logger.Debug("Call ApiPOST")
	respond(method, path, query, w, logger)
}

func ApiPUT(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	method := "PUT"
	logger := log.WithFields(log.Fields{
		"Path":   path,
		"Query":  query,
		"Method": method,
	})

	logger.Debug("Call ApiPUT")
	respond(method, path, query, w, logger)
}

func ApiDELETE(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	method := "DELETE"
	logger := log.WithFields(log.Fields{
		"Path":   path,
		"Query":  query,
		"Method": method,
	})

	logger.Debug("Call ApiDELETE")
	respond(method, path, query, w, logger)
}

func respond(method string, path string, query string, w http.ResponseWriter, logger *log.Entry) error {
	methods := ruleset.Methods[method]
	for _, m := range methods {
		pr := regexp.MustCompile(m.Path)
		if pr.MatchString(path) {
			logger.WithField("Path Pattern", m.Path).Debug("Path matches pattern")
			qr := regexp.MustCompile(m.Query)
			if qr.MatchString(query) {
				logger.WithField("Query Pattern", m.Query).Debug("Query matches pattern")
				for hdr, val := range m.ResponseHeaders {
					logger.WithFields(log.Fields{
						"Header": hdr,
						"Value":  val,
					}).Debug("Setting header")
					w.Header().Add(hdr, val)
				}
				w.WriteHeader(m.ResponseCode)
				w.Write([]byte(m.ResponseBody))
				return nil
			} else {
				logger.WithField("Query Pattern", m.Query).Debug("Query does not match pattern")
			}
		} else {
			logger.WithField("Path Pattern", m.Path).Debug("Path does not match pattern")
		}
	}
	logger.Debug("No matching path or query found, sending 404")
	w.WriteHeader(404)
	w.Write([]byte(path + "?" + query + " not found"))
	return nil
}
