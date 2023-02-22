package book

import (
	"fmt"

	"github.com/kamva/mgm/v3"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string `json:"title" bson:"title" validate:"required,gt=0"`
	Pages            int    `json:"pages" bson:"pages" validate:"required,numeric,gte=1"`
}

func NewBook(name string, pages int) *Book {
	return &Book{
		Title: name,
		Pages: pages,
	}
}

// Get book by id
func ReadById(log *common.Logger, id string) (out *Book, err error) {
	out = &Book{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	err = db.GoBookStore.FirstWithCtx(mgm.Ctx(), bson.M{"_id": objID}, out)
	if err != nil {
		return
	}
	return
}

// Get all books
func ReadAll(log *common.Logger) (out []*Book, err error) {

	filter := bson.M{}

	results := []Book{}

	err = db.GoBookStore.SimpleFindWithCtx(mgm.Ctx(), &results, filter)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for i := range results {
		out = append(out, &results[i])
	}
	return
}

// Create a book
func Create(log *common.Logger, in *Book) (out *Book, err error) {

	err = db.GoBookStore.CreateWithCtx(mgm.Ctx(), in)
	if err != nil {
		log.Error(err.Error())
		return
	}
	out = in
	return
}

// Update a book
func Update(log *common.Logger, id string, data *Book) (out *Book, err error) {
	out, err = ReadById(log, id)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("error fetching book, id: %s err: %s", id, err.Error())
	}

	if data.Title != "" {
		out.Title = data.Title
	}
	if data.Pages != 0 {
		out.Pages = data.Pages
	}

	err = db.GoBookStore.UpdateWithCtx(mgm.Ctx(), out)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

// Delete a book by id
func Delete(log *common.Logger, id string) (err error) {

	out, err := ReadById(log, id)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = db.GoBookStore.DeleteWithCtx(mgm.Ctx(), out)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

// check if book is already present in the collection
func isBookAlreadyPresent(title string, pages int) bool {

	filter := bson.M{"title": title, "pages": pages}

	result := []Book{}

	err := db.GoBookStore.SimpleFindWithCtx(mgm.Ctx(), &result, filter)
	if err != nil {
		return false
	}

	return len(result) != 0
}
