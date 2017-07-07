package server

import (
	"net/http"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const maxContentLength = 2048
const safeRegex = "^[a-zA-Z0-9_>-,/+]*"

func timerMiddleware(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	startTime := time.Now()
	next(resp, req)

	res := resp.(negroni.ResponseWriter)
	duration := time.Since(startTime)
	logrus.WithFields(logrus.Fields{
		"executionTimeMs":  time.Since(startTime) / time.Millisecond,
		"Hostname":         req.Host,
		"Method":           req.Method,
		"Path":             req.URL.Path,
		"Status":           res.Status(),
		"text_status":      http.StatusText(res.Status()),
		"executionTimeStr": duration.String()}).Info("ExecutionInfo")
}

func requestScrubberMiddleware(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	logrus.Debug("Scrubbin'")

	if req.ContentLength > maxContentLength {
		logrus.Error("request too large")
		http.Error(resp, "Request too large", http.StatusRequestEntityTooLarge)
		return
	}

	req.ParseForm()
	for k, v := range req.Form {
		logrus.Debug("k: ", k, " v: ", v)

		if len(k) > 100 || len(v) > 100 {
			logrus.Error("k/v too large")
			http.Error(resp, "Invalid Request", http.StatusBadRequest)
			return
		}

		match, err := regexp.MatchString(safeRegex, k)
		if err != nil {
			logrus.Error(err)
			http.Error(resp, "", http.StatusInternalServerError)
			return
		}
		if !match {
			logrus.Error("illegal characters in key")
			http.Error(resp, "", http.StatusBadRequest)
			return
		}

		for _, vi := range v {
			match, err = regexp.MatchString(safeRegex, vi)
			if err != nil {
				logrus.Error(err)
				http.Error(resp, "", http.StatusInternalServerError)
				return
			}
			if !match {
				logrus.Error("illegal characters in value")
				http.Error(resp, "", http.StatusBadRequest)
				return
			}
		}
	}

	next(resp, req)
}

func authMiddleware(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.Header.Get("Authorization") == "" {
		http.Error(resp, "", http.StatusUnauthorized)
		return
	}
	// TODO: actually inspect the auth data and act appropriately
	next(resp, req)
}
