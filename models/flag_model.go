package models

type FlagModel struct {
	Id   int64
	Name string
}

func NewFlagModel(id int64, name string) FlagModel {
	return FlagModel{
		Id:   id,
		Name: name,
	}
}
