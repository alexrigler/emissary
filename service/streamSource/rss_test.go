package streamSource

import (
	"testing"

	"github.com/benpate/ghost/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRss(t *testing.T) {

	source, err := New(model.StreamSourceAdapterRSS, primitive.NewObjectID(), model.StreamSourceConfig{
		"url": "https://appleinsider.com/rss/news",
	})

	assert.Nil(t, err)

	streams, err := source.Poll()

	assert.Nil(t, err)
	spew.Dump(streams)

	t.Fail()
}