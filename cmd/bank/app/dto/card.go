package dto

type CardDTO struct {
	Id int64 `json:"id"`
	Issuer string `json:"issuer"`
	Type string `json:"type"`
	Number string `json:"number"`
}
 type Result struct{
 	Result string `json:"result"`
 }