package models

type Customer struct {
	MetaFields             `bson:",inline"`
	UserId  string         `bson:"userId"  json:"userId"`
	Name    string         `bson:"name"  json:"name"`
	NamePY  string         `bson:"name_py"  json:"name_py"`
	Tel     string         `bson:"tel"  json:"tel"`
	Address string         `bson:"address"  json:"address"`
	Note    string         `bson:"note"  json:"note"`
}
