package CardGiverService

import (
	"context"
	"errors"
	"log"
	"sync"
	"webService/pkg/card"
)
var (
	errWrongType = errors.New("Wrong card type")
	errNoKey = errors.New("No key in request")
)
type Service struct{
	Cards []*card.Card
	MaxId int64
	mu sync.RWMutex
}

func CreateService() *Service{
	return &Service{MaxId: 0, mu: sync.RWMutex{}}
}

func (s *Service) IsueCard (uid int64, issuer string, cardType string, ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if(cardType != "Virtual" && cardType != "Real"){
		log.Println(errWrongType)
		return
	}
	number := string(s.MaxId)
	cards := s.GetAll(ctx)
	found := false
	for _, c := range cards{
		found = found || (c.OwnerId == uid)
	}
	if found {
		c := card.Card{Id: s.MaxId, Type: cardType, Issuer: issuer, OwnerId: uid, Number: "000" + number}
		s.Cards = append(s.Cards, &c)
		s.MaxId = s.MaxId + 1
	}
	return
}

func (s * Service) GetAll(ctx context.Context) []*card.Card{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Cards
}