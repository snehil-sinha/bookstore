package db

import (
	"github.com/kamva/mgm/v3"
	"github.com/snehil-sinha/goBookStore/common"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	GoBookStore *mgm.Collection
)

// Initialise the ODM
func InitDbMapper(log *common.Logger, database, uri string) (err error) {
	err = mgm.SetDefaultConfig(nil, database, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error(err.Error())
		return
	}
	initCollections()
	return
}

// Intialise the collections within the DB
func initCollections() {
	GoBookStore = mgm.CollectionByName(CollectionGoBookStore)
}
