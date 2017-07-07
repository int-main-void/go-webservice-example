package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func authMiddleware(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.Header.Get("IMPORTANTSTUFF") == "" {
		logrus.Error("no auth creds")
		http.Error(resp, "you can't do that on television", http.StatusUnauthorized)
		return
	}
	next(resp, req)
}
