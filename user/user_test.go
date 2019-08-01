package user

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name string
		want *User
	}{
		{
			"getnewUser",
			&User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser()
			if reflect.TypeOf(got).String() != "*user.User" {
				t.Errorf("NewUser() create error %v", got)
			}
			if len(got.ID) == 0 {
				t.Errorf("NewUser() = %v, set ID error", got)
			}
		})
	}
}

func TestUser_Save(t *testing.T) {
	type fields struct {
		ID       string
		Email    string
		Password string
		Claims   []Claim
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:       tt.fields.ID,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
				Claims:   tt.fields.Claims,
			}
			if err := u.Save(); (err != nil) != tt.wantErr {
				t.Errorf("User.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisConn(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"InitConnect",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := redisConn()
			if err != nil {
				t.Errorf("redisConn() error = %v", err)
				return
			}
			_, err = c.Do("PING")
			if err != nil {
				t.Errorf("redisConn() error = %v", err)
				return
			}
		})
	}
}
