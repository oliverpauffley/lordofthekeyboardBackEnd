package app

import (
	"fmt"
	"keyboard/service/quotes"
	"net/http"
)

type Server struct {
	Router      *http.ServeMux
	QuoteSource quotes.API
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleHome())
	s.Router.HandleFunc("/quote", s.handleQuote())
}

func (s *Server) handleHome() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "hello world")
	}
}

func (s *Server) handleQuote() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		quote, err := s.QuoteSource.NewQuote()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, err)

			return
		}

		fmt.Fprintf(res, quote.Text+quote.CharacterName)
	}
}
