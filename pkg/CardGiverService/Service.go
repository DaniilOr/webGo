package CardGiverService

import (
	"context"
	"errors"
	"fmt"
	"github.com/DaniilOr/webGo/pkg/card"
	"log"
	"strconv"
	"sync"
)
var (
	errWrongType = errors.New("Wrong card type")
	errNoKey = errors.New("No key in request")
	errNoSuchUser = errors.New("No such user")
)
type Service struct{
	cards []*card.Card
	maxId int64
	mu sync.RWMutex
}

func CreateService() *Service{
	s := Service{maxId: 1, mu: sync.RWMutex{}, cards: []*card.Card{&card.Card{OwnerId: 5, Number: "789"}}}
	return &s
}

func (s *Service) IsueCard (uid int64, issuer string, cardType string, ctx context.Context) error {
	if(cardType != "Virtual" && cardType != "Real"){
		log.Println(errWrongType)
		return errWrongType
	}
	number := strconv.Itoa(int(s.maxId))
	cards := s.GetAll(ctx)
	found := false
	for _, c := range cards{
		found = found || (c.OwnerId == uid)
	}
	if found {
		c := card.Card{Id: s.maxId, Type: cardType, Issuer: issuer, OwnerId: uid, Number: fmt.Sprintf("%s%s", "000", number)}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.cards = append(s.cards, &c)
		s.maxId = s.maxId + 1
		return nil
	}
	return errNoSuchUser
}

func (s * Service) GetAll(ctx context.Context) []*card.Card{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cards
}