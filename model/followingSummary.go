package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type FollowingSummary struct {
	FollowingID primitive.ObjectID `bson:"_id"`
	URL         string             `bson:"url"`
	Label       string             `bson:"label"`
	FolderID    primitive.ObjectID `bson:"folderId"`
	ImageURL    string             `bson:"imageUrl"`
	Method      string             `bson:"method"`
	Status      string             `bson:"status"`
	LastPolled  int64              `bson:"lastPolled"`
}

// FollowingSummaryFields returns a slice of all BSON field names for a FollowingSummary
func FollowingSummaryFields() []string {
	return []string{"_id", "url", "label", "folderId", "imageUrl", "method", "status", "lastPolled"}
}

func (summary FollowingSummary) Fields() []string {
	return FollowingSummaryFields()
}

func (summary FollowingSummary) StatusIcon() string {

	var icon string

	switch summary.Method {
	case FollowMethodActivityPub:
		icon = "activitypub"
	case FollowMethodPoll:
		icon = "rss"
	case FollowMethodWebSub:
		icon = "websub"
	}

	switch summary.Status {
	case FollowingStatusLoading:
		return "loading"
	case FollowingStatusSuccess:
		return icon + "-fill"
	default:
		return icon
	}
}

func (summary FollowingSummary) StatusClass() string {

	switch summary.Status {
	case FollowingStatusLoading:
		return "spin"
	case FollowingStatusFailure:
		return "red"
	case FollowingStatusSuccess:
		return "green"
	default:
		return ""
	}
}
