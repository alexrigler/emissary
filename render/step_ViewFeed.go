package render

import (
	"io"
	"net/http"
	"time"

	"github.com/EmissarySocial/emissary/model"
	"github.com/EmissarySocial/emissary/tools/convert"
	"github.com/EmissarySocial/emissary/tools/iterators"
	"github.com/benpate/data"
	"github.com/benpate/derp"
	"github.com/gorilla/feeds"
	"github.com/kr/jsonfeed"
	accept "github.com/timewasted/go-accept-headers"
)

// StepViewFeed represents an action-step that can render a Stream into HTML
type StepViewFeed struct{}

// Get renders the Stream HTML to the context
func (step StepViewFeed) Get(renderer Renderer, buffer io.Writer) error {

	const location = "render.StepViewFeed.Get"

	factory := renderer.factory()

	// Get all child streams from the database
	children, err := factory.Stream().ListByParent(renderer.objectID())

	if err != nil {
		return derp.Wrap(err, location, "Error querying child streams")
	}

	mimeType := step.detectMimeType(renderer)

	// Special case for JSONFeed
	if mimeType == model.MimeTypeJSONFeed {
		return step.asJSONFeed(renderer, buffer, children)
	}

	// Initialize the result RSS feed
	result := feeds.Feed{
		Title:       renderer.PageTitle(),
		Description: renderer.Summary(),
		Link:        &feeds.Link{Href: renderer.Permalink()},
		Author:      &feeds.Author{Name: ""},
		Created:     time.Now(),
	}

	result.Items = iterators.Map(children, model.NewStream, convert.StreamToGorillaFeed)

	// Now write the feed into the requested format
	{
		var xml string
		var err error

		// Thank you gorilla/feeds for this awesome API.
		switch mimeType {

		case model.MimeTypeAtom:
			mimeType = "application/atom+xml; charset=UTF=8"
			xml, err = result.ToAtom()

		case model.MimeTypeRSS:
			mimeType = "application/rss+xml; charset=UTF=8"
			xml, err = result.ToRss()
		}

		if err != nil {
			return derp.Wrap(err, location, "Error generating feed. This should never happen")
		}

		// Write the result to the buffer and then success.
		header := renderer.context().Response().Header()
		header.Add("Content-Type", mimeType)
		buffer.Write([]byte(xml))
		return nil
	}
}

func (step StepViewFeed) UseGlobalWrapper() bool {
	return false
}

func (step StepViewFeed) Post(renderer Renderer) error {
	return nil
}

func (step StepViewFeed) detectMimeType(renderer Renderer) string {

	context := renderer.context()

	// First, try to get the format from the query string
	switch context.QueryParam("format") {
	case "json":
		return model.MimeTypeJSONFeed
	case "atom":
		return model.MimeTypeAtom
	case "rss":
		return model.MimeTypeRSS
	}

	// Otherwise, get the format from the "Accept" header
	header := context.Request().Header

	if result, err := accept.Negotiate(header.Get("Accept"), model.MimeTypeJSONFeed, model.MimeTypeAtom, model.MimeTypeRSS); err == nil {
		return result
	}

	// Finally, use JSONFeed as the default
	return model.MimeTypeJSONFeed
}

func (step StepViewFeed) asJSONFeed(renderer Renderer, buffer io.Writer, children data.Iterator) error {

	context := renderer.context()

	feed := jsonfeed.Feed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       renderer.PageTitle(),
		HomePageURL: renderer.Permalink(),
		FeedURL:     renderer.Permalink() + "/feed?format=json",
		Description: renderer.Summary(),
		Hubs: []jsonfeed.Hub{
			{
				Type: "WebSub",
				URL:  renderer.Permalink() + "/websub",
			},
		},
	}

	feed.Items = iterators.Map(children, model.NewStream, convert.StreamToJsonFeed)

	context.Response().Header().Add("Content-Type", model.MimeTypeJSONFeed)

	return context.JSON(http.StatusOK, feed)
}