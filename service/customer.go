package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/db/mongo"
	"github.com/nairufan/yh-weixin/apperror"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionCustomer = "customer"
)

func AddCustomer(customer *models.Customer) *models.Customer {
	customer.MetaFields = models.NewMetaFields()
	if customer.Tel == "" {
		panic(apperror.NewInvalidParameterError("tel"))
	}
	if customer.UserId == "" {
		panic(apperror.NewInvalidParameterError("userId"))
	}
	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionCustomer, customer)

	return customer
}

func UpdateCustomer(customer *models.Customer) *models.Customer {
	if customer.Id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}
	c := GetCustomerById(customer.Id)
	c.Tel = customer.Tel
	c.Name = customer.Name
	c.Address = customer.Address
	c.Note = customer.Note

	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionCustomer, customer.Id, c)
	return c
}

func GetCustomerById(id string) *models.Customer {
	session := mongo.Get()
	defer session.Close()
	document := &models.Customer{}
	session.MustFindId(collectionCustomer, id, document)
	return document
}

func RemoveCustomerById(id string) {
	if id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}

	session := mongo.Get()
	defer session.Close()
	session.RemoveId(collectionCustomer, id)
}

func GetCustomers(userId string, offset int, limit int) []*models.Customer {
	session := mongo.Get()
	defer session.Close()
	customers := []*models.Customer{}

	option := mongo.Option{
		Sort: []string{"+name"},
		Limit: &limit,
		Offset: &offset,
	}
	session.MustFindWithOptions(collectionCustomer, bson.M{"userId": userId}, option, &customers)
	return customers
}

func GetCustomerByTel(tel string) *models.Customer {
	session := mongo.Get()
	defer session.Close()
	customers := []*models.Customer{}
	session.MustFind(collectionCustomer, bson.M{"tel": tel}, &customers)
	if len(customers) > 0 {
		return customers[0]
	}
	return nil
}

func GetCustomerByIds(ids []string) []*models.Customer {
	session := mongo.Get()
	defer session.Close()
	customers := []*models.Customer{}

	session.MustFind(collectionCustomer, bson.M{"_id": bson.M{"$in": ids}}, &customers)
	return customers
}