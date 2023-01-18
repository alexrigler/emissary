package model

import "go.mongodb.org/mongo-driver/bson/primitive"

func (form SignupForm) GetBoolOK(name string) (bool, bool) {

	switch name {

	case "active":
		return form.Active, true

	}

	return false, false
}

func (form SignupForm) GetStringOK(name string) (string, bool) {

	switch name {

	case "title":
		return form.Title, true

	case "message":
		return form.Message, true

	case "groupId":
		return form.GroupID.Hex(), true

	}

	return "", false
}

func (form *SignupForm) SetBoolOK(name string, value bool) bool {

	switch name {

	case "active":
		form.Active = value
		return true

	}

	return false
}

func (form *SignupForm) SetStringOK(name string, value string) bool {

	switch name {

	case "title":
		form.Title = value
		return true

	case "message":
		form.Message = value
		return true

	case "groupId":
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			form.GroupID = objectID
			return true
		}

	}

	return false
}
