package step

import (
	"github.com/benpate/datatype"
)

// Save represents an action-step that can save changes to any object
type Save struct {
	Comment string
}

// NewSave returns a fully initialized Save object
func NewSave(stepInfo datatype.Map) (Save, error) {

	return Save{
		Comment: stepInfo.GetString("comment"),
	}, nil
}

// AmStep is here only to verify that this struct is a render pipeline step
func (step Save) AmStep() {}