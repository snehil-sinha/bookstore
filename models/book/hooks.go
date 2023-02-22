package book

import (
	"context"
	"fmt"
)

// Callbacks for book
func (b *Book) Saving(ctx context.Context) (err error) {

	// Call the DefaultModel Saving hook to update created_at and updated_at timestamps
	if err = b.DefaultModel.Saving(); err != nil {
		return
	}

	// Validate the model fields
	err = b.Validate()
	if err != nil {
		err = fmt.Errorf("book must have a non-empty title and at least 1 page: %s", err)
		return
	}

	// Check if book is already present
	if isBookAlreadyPresent(b.Title, b.Pages) {
		err = fmt.Errorf("book with title %s and pages %d already exists", b.Title, b.Pages)
		return
	}

	return
}
