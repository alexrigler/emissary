package handler

import (
	"github.com/EmissarySocial/emissary/domain"
	"github.com/EmissarySocial/emissary/model"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/davecgh/go-spew/spew"
)

func init() {

	// This funciton handles ActivityPub "Accept/Follow" activities, meaning that
	// it is called with a remote server accepts our follow request.
	inboxRouter.Add(vocab.Any, vocab.Any, func(factory *domain.Factory, user *model.User, activity streams.Document) error {

		spew.Dump("RECEIVED UNKNOWN ACTIVITY -----------------------------", activity)

		return nil
	})
}