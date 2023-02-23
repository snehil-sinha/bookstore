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
		err = fmt.Errorf("validation error: %s", err)
		return
	}

	return
}
