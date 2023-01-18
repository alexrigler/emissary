package model

import (
	"testing"

	"github.com/benpate/rosetta/schema"
)

func TestUserSchema(t *testing.T) {

	s := schema.New(UserSchema())
	user := NewUser()

	tests := []tableTestItem{
		{"userId", "000000000000000000000001", nil},
		{"groupIds.0", "000000000000000000000002", nil},
		{"groupIds.1", "000000000000000000000003", nil},
		{"groupIds.2", "000000000000000000000004", nil},
		{"imageId", "000000000000000000000005", nil},
		{"displayName", "USER", nil},
		{"statusMessage", "STATUS", nil},
		{"location", "LOCATION", nil},
		{"links.0.name", "LINK 1", nil},
		{"links.0.profileUrl", "https://profile.url", nil},
		{"profileUrl", "http://profile.url", nil},
		{"emailAddress", "email@address.url", nil},
		{"username", "USERNAME", nil},
		{"followerCount", "1", 1},
		{"followingCount", "2", 2},
		{"blockCount", "3", 3},
		{"isOwner", "true", true},
	}

	tableTest_Schema(t, &s, &user, tests)

	//TODO: Include DefaultAllow?

}
