package service

import (
	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/rosetta/list"
	"github.com/whisperverse/mediaserver"
	"github.com/whisperverse/whisperverse/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Attachment manages all interactions with the Attachment collection
type Attachment struct {
	collection  data.Collection
	mediaServer mediaserver.MediaServer
}

// NewAttachment returns a fully populated Attachment service
func NewAttachment(collection data.Collection, mediaServer mediaserver.MediaServer) Attachment {
	return Attachment{
		collection:  collection,
		mediaServer: mediaServer,
	}
}

// New creates a newly initialized Attachment that is ready to use
func (service Attachment) New() model.Attachment {
	return model.Attachment{
		AttachmentID: primitive.NewObjectID(),
	}
}

// List returns an iterator containing all of the Attachments who match the provided criteria
func (service Attachment) List(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {
	return service.collection.List(notDeleted(criteria), options...)
}

// Load retrieves an Attachment from the database
func (service Attachment) Load(criteria exp.Expression, result *model.Attachment) error {

	if err := service.collection.Load(notDeleted(criteria), result); err != nil {
		return derp.Wrap(err, "service.Attachment", "Error loading Attachment", criteria)
	}

	return nil
}

// Save adds/updates an Attachment in the database
func (service Attachment) Save(attachment *model.Attachment, note string) error {

	if err := service.collection.Save(attachment, note); err != nil {
		return derp.Wrap(err, "service.Attachment", "Error saving Attachment", attachment, note)
	}

	return nil
}

// Delete removes an Attachment from the database (virtual delete)
func (service Attachment) Delete(attachment *model.Attachment, note string) error {

	// Delete uploaded files from MediaServer
	if err := service.mediaServer.Delete(attachment.AttachmentID.Hex()); err != nil {
		derp.Report(derp.Wrap(err, "service.Attachment", "Error deleting attached files", attachment))
		// Fail loudly, but do not stop.
	}

	// Delete Attachment record last.
	if err := service.collection.Delete(attachment, note); err != nil {
		return derp.Wrap(err, "service.Attachment", "Error deleting Attachment", attachment, note)
	}

	return nil
}

// DeleteByStream removes all attachments from the provided stream (virtual delete)
func (service Attachment) DeleteAllFromStream(streamID primitive.ObjectID, note string) error {

	var attachment model.Attachment
	it, err := service.ListByObjectID(streamID)

	if err != nil {
		return derp.Wrap(err, "service.Attachment.DeleteByStream", "Error listing attachments", streamID)
	}

	for it.Next(&attachment) {
		if err := service.Delete(&attachment, note); err != nil {
			derp.Report(derp.Wrap(err, "service.Attachment.DeleteByStream", "Error deleting child stream", attachment))
			// Fail loudly, but do not stop.
		}
	}

	return nil
}

/*******************************************
 * CUSTOM QUERIES
 *******************************************/

func (service Attachment) ListByObjectID(objectID primitive.ObjectID) (data.Iterator, error) {
	return service.List(
		exp.Equal("streamId", objectID).
			AndEqual("journal.deleteDate", 0),
		option.SortAsc("rank"))
}

func (service Attachment) ListFirstByObjectID(objectID primitive.ObjectID) (data.Iterator, error) {
	return service.List(
		exp.Equal("streamId", objectID).
			AndEqual("journal.deleteDate", 0),
		option.SortAsc("rank"), option.FirstRow())
}

func (service Attachment) LoadByStreamID(streamID primitive.ObjectID, attachmentID primitive.ObjectID) (model.Attachment, error) {
	var result model.Attachment
	criteria := exp.Equal("streamId", streamID).AndEqual("_id", attachmentID).AndEqual("journal.deleteDate", 0)
	err := service.Load(criteria, &result)
	return result, err
}

func (service Attachment) LoadByToken(streamID primitive.ObjectID, token string) (model.Attachment, error) {
	var result model.Attachment
	criteria := exp.Equal("streamId", streamID).
		AndEqual("filename", list.Dot(token).Head()).
		AndEqual("journal.deleteDate", 0)
	err := service.Load(criteria, &result)
	return result, err
}
