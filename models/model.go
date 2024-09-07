package models

type Item struct {
	Id    int
	Name  string
	Price int
}

var id = 0

var Items []Item

func NewItem(name string, price int) Item {
	id++
	return Item{
		Id:    id,
		Name:  name,
		Price: price,
	}
}
