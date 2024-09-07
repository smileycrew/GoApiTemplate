package models

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var id = 0

var Items []Item

func NewItem(name string, price int) Item {
	id++
	return Item{
		ID:    id,
		Name:  name,
		Price: price,
	}
}
