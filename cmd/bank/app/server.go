package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"webService/cmd/bank/app/dto"
	"webService/pkg/CardGiverService"
)

type Server struct {
	cardSvc *CardGiverService.Service
	mux *http.ServeMux
}

func NewServer(cardSvc *CardGiverService.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	suid := r.URL.Query().Get("uid")
	if suid != ""{
		uid, err := strconv.ParseInt(suid, 10, 64)
		if err != nil{
			log.Println(err)
			return
		}
		cards := s.cardSvc.GetAll(r.Context())
		dtos := make([]*dto.CardDTO, len(cards))
		for i, c := range cards {
			if c.OwnerId == uid {
				dtos[i] = &dto.CardDTO{
					Id:      c.Id,
					Number:  c.Number,
					Issuer:  c.Issuer,
					Type: c.Type,
				}
			}
		}

		makeResponse(dtos, w, r)
	}
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	issuer := r.PostForm.Get("issuer")
	cardType := r.PostForm.Get("type")
	suid := r.URL.Query().Get("uid")
	if suid != "" {
		uid, err := strconv.ParseInt(suid, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		s.cardSvc.IsueCard(uid, issuer, cardType, r.Context())
	}

}

func makeResponse(dtos []*dto.CardDTO, w http.ResponseWriter, r*http.Request) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}