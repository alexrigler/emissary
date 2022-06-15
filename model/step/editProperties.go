package step

import (
	"github.com/benpate/datatype"
)

// EditProperties contains the configuration data for a modal that lets users edit the features attached to a stream.
type EditProperties struct {
	Paths []string
}

func NewEditProperties(stepInfo datatype.Map) (EditProperties, error) {
	paths := stepInfo.GetSliceOfString("paths")

	if len(paths) == 0 {
		paths = []string{"token", "label", "description"}
	}

	return EditProperties{
		Paths: paths,
	}, nil
}

// AmStep is here only to verify that this struct is a render pipeline step
func (step EditProperties) AmStep() {}