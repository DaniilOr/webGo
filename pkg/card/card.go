package card

type Card struct {
	Id int64 `json:"id"`
	Issuer string `json:"issuer"`
	Type string `json:"type"`
	OwnerId int64 `json:"owner_id"`
	Number string `json:"number"`
	Balance int64 `json:"balance"`
}
