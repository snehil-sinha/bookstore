package book

import (
	"fmt"

	"github.com/snehil-sinha/goBookStore/service/validators"
)

// Validate the fields of the book
func (b *Book) Validate() (err error) {

	v := validators.New()
	err = v.Struct(b)
	if err != nil {
		err = fmt.Errorf("book validation failed, err: %s", err.Error())
		return
	}
	return
}
