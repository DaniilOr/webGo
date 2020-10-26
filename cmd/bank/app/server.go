package app

import (
	"encoding/json"
	"github.com/DaniilOr/webGo/cmd/bank/app/dto"
	"github.com/DaniilOr/webGo/pkg/CardGiverService"
	"log"
	"net/http"
	"strconv"

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
		log.Println(cards)
		found := false
		dtos := make([]*dto.CardDTO, len(cards))
		for i, c := range cards {
			if c.OwnerId == uid {
				found = true
				dtos[i] = &dto.CardDTO{
					Id:      c.Id,
					Number:  c.Number,
					Issuer:  c.Issuer,
					Type: c.Type,
				}
			}
		}
		if !found {
			response := dto.Result{Result: "No cards"}
			respBody, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			makeResponse(respBody, w, r)
			return
		}
		respBody, err := json.Marshal(dtos)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		makeResponse(respBody, w, r)
	}
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	issuer := r.PostForm.Get("issuer")
	cardType := r.PostForm.Get("type")
	suid := r.PostForm.Get("uid")
	if suid != "" {
		uid, err := strconv.ParseInt(suid, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		err = s.cardSvc.IsueCard(uid, issuer, cardType, r.Context())
		var result dto.Result
		if err != nil{
			result = dto.Result{Result: "Error"}


		} else {
			result = dto.Result{Result: "Ok"}
			makeResponse([]byte("Ok"), w, r)
		}
		response, err := json.Marshal(result)
		if err != nil{
			log.Println(err)
			return
		}
		makeResponse(response, w, r)
	}

}

func makeResponse(respBody []byte, w http.ResponseWriter, r*http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}