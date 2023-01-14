package gofed

import (
	"context"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// GetInbox returns the latest page of the inbox corresponding to the inboxIRI.
//
// At first glance this method seems a little odd. It is fine to return an empty
// vocab.ActivityStreamsOrderedCollectionPage. The library expects the very first
// page, which is the most recent chronologically. Therefore, an empty page is
// always treated as the "first zero" items, and the library does not require having
// any items. If you have a caching layer, it can more easily hide under this method
// with proper pagination and delayed writes to the database. The library is simply
// going to prepend an item in the orderedItems property and then call SetInbox.
func (db Database) GetInbox(c context.Context, inboxIRI *url.URL) (inbox vocab.ActivityStreamsOrderedCollectionPage, err error) {

	// It seems like we may not need to serve the Inbox via ActivityPub.
	// Going to leave this empty for now, and will add if it becomes necessary.

	result := streams.NewActivityStreamsOrderedCollectionPage()
	return result, nil
}