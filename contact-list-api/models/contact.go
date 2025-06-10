package models

type Contact struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"João da Silva"`
	Email string `json:"email" example:"joao@email.com"`
	Phone string `json:"phone" example:"11999998888"`
}
