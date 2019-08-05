package user

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
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
			got := New()
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
	newuser := New()
	newuser.Email = "oleh@example.com"
	newuser.Password = "oleh12345"
	clm := Claim{}
	clm.Resource = "accounts.example.com"
	clm.Asserts = make(AssertsMap)
	clm.Asserts["Account"] = "12345"
	newuser.Claims = append(newuser.Claims, clm)
	tests := []struct {
		name     string
		testuser *User
		wantErr  bool
	}{
		{
			"userOleh",
			newuser,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.testuser
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

func TestNewClaim(t *testing.T) {
	type args struct {
		resource string
		asserts  AssertsMap
	}
	tests := []struct {
		name string
		args args
		want *Claim
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClaim(tt.args.resource, tt.args.asserts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}
