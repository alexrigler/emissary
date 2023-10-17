package model

import (
	"github.com/benpate/rosetta/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OAuthUserTokenSchema() schema.Element {
	return schema.Object{
		Properties: schema.ElementMap{
			"oauthUserTokenId":   schema.String{Format: "objectId", Required: true},
			"userId":             schema.String{Format: "objectId", Required: true},
			"oauthApplicationId": schema.String{Format: "objectId", Required: true},
			"token":              schema.String{},
			"clientSecret":       schema.String{},
			"scopes":             schema.Array{Items: schema.String{}},
		},
	}
}

func (userToken *OAuthUserToken) GetPointer(name string) (any, bool) {
	switch name {

	case "token":
		return &userToken.Token, true

	case "clientSecret":
		return &userToken.ClientSecret, true

	case "scopes":
		return &userToken.Scopes, true
	}

	return nil, false
}

func (userToken *OAuthUserToken) GetStringOK(name string) (string, bool) {

	switch name {

	case "oauthUserTokenId":
		return userToken.OAuthUserTokenID.Hex(), true

	case "oauthApplicationId":
		return userToken.OAuthApplicationID.Hex(), true

	case "userId":
		return userToken.UserID.Hex(), true
	}

	return "", false
}

func (userToken *OAuthUserToken) SetString(name string, value string) bool {

	switch name {

	case "oauthUserTokenId":
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			userToken.OAuthUserTokenID = objectID
			return true
		}

	case "oauthApplicationId":
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			userToken.OAuthApplicationID = objectID
			return true
		}

	case "userId":
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			userToken.UserID = objectID
			return true
		}
	}

	return false
}