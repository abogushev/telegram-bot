package model

type Category struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func NewCategory(id int, name string) *Category {
	return &Category{Id: id, Name: name}
}
