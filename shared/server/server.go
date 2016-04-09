package server

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Server struct {
	Addr string
	Cert string
	Key  string
}

func Load(addr string, cert string, key string) *Server {
	return &Server{
		Addr: addr,
		Cert: cert,
		Key:  key,
	}
}

func (s *Server) Run(handler http.Handler) {
	if len(s.Cert) != 0 {
		logrus.Fatal(http.ListenAndServeTLS(s.Addr, s.Cert, s.Key, handler))
	} else {
		logrus.Fatal(http.ListenAndServe(s.Addr, handler))
	}
}
