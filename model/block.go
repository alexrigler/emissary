package model

import (
	"github.com/benpate/data/journal"
	"github.com/benpate/rosetta/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlockSourceInternal represents a block that was created directly by the owner
const BlockSourceInternal = "INTERNAL"

// BlockSourceActivityPub represents a block that was created by an external ActivityPub server
const BlockSourceActivityPub = "ACTIVITYPUB"

// BlockTypeURL blocks all messages that link to a specific domain or URL prefix
const BlockTypeURL = "URL"

// BlockTypeUser blocks all messages from a specific user
const BlockTypeActor = "ACTOR"

// BlockTypeUser blocks all messages that contain a particular phrase (hashtag)
const BlockTypeContent = "CONTENT"

// BlockTypeExternal passes messages to an external block service (TBD) for analysis.
const BlockTypeExternal = "EXTERNAL"

// Block represents many kinds of filters that are applied to messages before they are added into a User's inbox
type Block struct {
	BlockID  primitive.ObjectID `json:"blockId" bson:"_id"`       // Unique identifier of this Block
	UserID   primitive.ObjectID `json:"userId"  bson:"userId"`    // Unique identifier of the User who owns this Block
	Source   string             `json:"source"  bson:"source"`    // Source of the Block (e.g. "INTERNAL", "ACTIVITYPUB")
	Type     string             `json:"type"    bson:"type"`      // Type of Block (e.g. "ACTOR", "ACTIVITY", "OBJECT")
	Trigger  string             `json:"trigger" bson:"trigger"`   // Parameter for this block type)
	Comment  string             `json:"comment" bson:"comment"`   // Optional comment describing why this block exists
	IsPublic bool               `json:"isPublic" bson:"isPublic"` // If TRUE, this record is visible publicly
	IsActive bool               `json:"isActive" bson:"isActive"` // If TRUE, this record is active

	journal.Journal `json:"-" bson:"journal"`
}

func NewBlock() Block {
	return Block{
		BlockID: primitive.NewObjectID(),
	}
}

func BlockSchema() schema.Element {
	return schema.Object{
		Properties: schema.ElementMap{
			"blockId":  schema.String{Format: "objectId"},
			"userId":   schema.String{Format: "objectId"},
			"source":   schema.String{Enum: []string{BlockSourceInternal, BlockSourceActivityPub}},
			"type":     schema.String{Enum: []string{BlockTypeURL, BlockTypeActor, BlockTypeContent, BlockTypeExternal}},
			"trigger":  schema.String{},
			"comment":  schema.String{},
			"isPublic": schema.Boolean{},
			"isActive": schema.Boolean{},
		},
	}
}

/*******************************************
 * data.Object Interface
 *******************************************/

func (block Block) ID() string {
	return block.BlockID.Hex()
}

/*******************************************
 * RoleStateEnumerator Interface
 *******************************************/

// State returns the current state of this object.
// For users, there is no state, so it returns ""
func (block Block) State() string {
	return ""
}

// Roles returns a list of all roles that match the provided authorization.
// Since Block records should only be accessible by the block owner, this
// function only returns MagicRoleMyself if applicable.  Others (like Anonymous
// and Authenticated) should never be allowed on an Block record, so they
// are not returned.
func (block Block) Roles(authorization *Authorization) []string {

	// Folders are private, so only MagicRoleMyself is allowed
	if authorization.UserID == block.UserID {
		return []string{MagicRoleMyself}
	}

	// Intentionally NOT allowing MagicRoleAnonymous, MagicRoleAuthenticated, or MagicRoleOwner
	return []string{}
}
