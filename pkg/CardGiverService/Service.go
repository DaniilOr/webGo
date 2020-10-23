package CardGiverService

import (
	"context"
	"errors"
	"github.com/DaniilOr/webGo/pkg/card"
	"log"
	"sync"
)
var (
	errWrongType = errors.New("Wrong card type")
	errNoKey = errors.New("No key in request")
	errNoSuchUser = errors.New("No such user")
)
type Service struct{
	Cards []*card.Card
	MaxId int64
	mu sync.RWMutex
}

func CreateService() *Service{
	s := Service{MaxId: 1, mu: sync.RWMutex{}, Cards: []*card.Card{&card.Card{OwnerId: 5, Number: "789"}}}
	return &s
}

func (s *Service) IsueCard (uid int64, issuer string, cardType string, ctx context.Context) error {
	if(cardType != "Virtual" && cardType != "Real"){
		log.Println(errWrongType)
		return errWrongType
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
		log.Println(c)
		log.Println(s.Cards)
		s.MaxId = s.MaxId + 1
		return nil
	}
	return errNoSuchUser
}

func (s * Service) GetAll(ctx context.Context) []*card.Card{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Cards
}