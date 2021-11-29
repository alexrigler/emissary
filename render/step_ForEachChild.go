package render

import (
	"io"

	"github.com/benpate/datatype"
	"github.com/benpate/derp"
	"github.com/benpate/ghost/model"
	"github.com/benpate/ghost/service"
)

// StepForChildren represents an action-step that can update the data.DataMap custom data stored in a Stream
type StepForChildren struct {
	streamService *service.Stream
	steps         []datatype.Map
}

// NewStepForChildren returns a fully initialized StepForChildren object
func NewStepForChildren(streamService *service.Stream, stepInfo datatype.Map) StepForChildren {

	return StepForChildren{
		streamService: streamService,
		steps:         stepInfo.GetSliceOfMap("then"),
	}
}

// Get displays a form where users can update stream data
func (step StepForChildren) Get(buffer io.Writer, renderer *Renderer) error {
	return nil
}

// Post updates the stream with approved data from the request body.
func (step StepForChildren) Post(buffer io.Writer, renderer *Renderer) error {

	children, err := step.streamService.ListByParent(renderer.stream.ParentID)

	if err != nil {
		return derp.Wrap(err, "ghost.render.StepForChildren.Post", "Error listing children")
	}

	child := new(model.Stream)

	for children.Next(child) {

		// Make a renderer with the new child stream
		childRenderer, err := renderer.newRenderer(child, renderer.ActionID())

		if err != nil {
			return derp.Wrap(err, "ghost.render.StepForChildren.Post", "Error creating renderer for child")
		}

		// Execute the POST render pipeline on the child
		if err := DoPipeline(&childRenderer, buffer, step.steps, ActionMethodPost); err != nil {
			return derp.Wrap(err, "ghost.render.StepForChildren.Post", "Error executing steps for child")
		}

		// Reset the child object so that old records don't bleed into new ones.
		child = new(model.Stream)
	}

	return nil
}