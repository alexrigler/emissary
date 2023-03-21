package model

import (
	"github.com/benpate/rosetta/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BlockSchema() schema.Element {
	return schema.Object{
		Properties: schema.ElementMap{
			"blockId":  schema.String{Required: true, Format: "objectId"},
			"userId":   schema.String{Required: true, Format: "objectId"},
			"type":     schema.String{Required: true, Enum: []string{BlockTypeDomain, BlockTypeActor, BlockTypeContent, BlockTypeExternal}},
			"trigger":  schema.String{Required: true},
			"behavior": schema.String{Enum: []string{BlockBehaviorBlock, BlockBehaviorMute, BlockBehaviorAllow}},
			"comment":  schema.String{},
			"origin":   OriginLinkSchema(),
			"isPublic": schema.Boolean{},
		},
	}
}

/******************************************
 * Getter Interfaces
 ******************************************/

func (block *Block) GetBoolOK(name string) (bool, bool) {

	switch name {

	case "isPublic":
		return block.IsPublic, true
	}

	return false, false
}

func (block *Block) GetStringOK(name string) (string, bool) {

	switch name {

	case "blockId":
		return block.BlockID.Hex(), true

	case "userId":
		return block.UserID.Hex(), true

	case "type":
		return block.Type, true

	case "trigger":
		return block.Trigger, true

	case "behavior":
		return block.Behavior, true

	case "comment":
		return block.Comment, true
	}

	return "", false
}

/******************************************
 * Setter Interfaces
 ******************************************/

func (block *Block) SetBool(name string, value bool) bool {

	switch name {

	case "isPublic":
		block.IsPublic = value
		return true
	}

	return false
}

func (block *Block) SetString(name string, value string) bool {

	switch name {

	case "blockId":
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			block.BlockID = objectID
			return true
		}

	case "userId":
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			block.UserID = objectID
			return true
		}

	case "type":
		block.Type = value
		return true

	case "trigger":
		block.Trigger = value
		return true

	case "behavior":
		block.Behavior = value
		return true

	case "comment":
		block.Comment = value
		return true

	}

	return false
}

func (block *Block) GetObject(name string) (any, bool) {

	switch name {

	case "origin":
		return &block.Origin, true
	}

	return nil, false
}
