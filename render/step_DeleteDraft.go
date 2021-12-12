package render

import (
	"io"

	"github.com/benpate/datatype"
	"github.com/benpate/derp"
	"github.com/benpate/ghost/service"
)

// StepStreamDraftDelete represents an action-step that can delete a StreamDraft from the Domain
type StepStreamDraftDelete struct {
	draftService *service.StreamDraft
}

func NewStepStreamDraftDelete(draftService *service.StreamDraft, config datatype.Map) StepStreamDraftDelete {
	return StepStreamDraftDelete{
		draftService: draftService,
	}
}

func (step StepStreamDraftDelete) Get(buffer io.Writer, renderer *Stream) error {
	return nil
}

func (step StepStreamDraftDelete) Post(buffer io.Writer, renderer *Stream) error {

	if err := step.draftService.Delete(renderer.stream, "Deleted"); err != nil {
		return derp.Wrap(err, "ghost.render.StepStreamDraftDelete.Post", "Error deleting stream draft")
	}

	renderer.ctx.Response().Header().Set("hx-redirect", "/"+renderer.stream.Token)
	return nil
}