package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

func newPool(addr string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     300,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func TestUser_Save(t *testing.T) {
	id := uuid.New().String()
	newuserEmail := "oleg.nagornij@gmail.com"
	nick := "Oleh"
	newuserPswd := "corner578"
	secret := "secret"
	newuser := New(id, newuserEmail, nick, newuserPswd, secret)
	claim := NewClaim(
		"*",
		"admin")
	newuser.Claims = append(newuser.Claims, *claim)
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

func TestNew(t *testing.T) {
	id := uuid.New().String()
	email := "oleg.nagornij@gmail.com"
	nick := "Oleh"
	pswd := "corner578"
	secret := "secret"
	userT := &User{
		ID:        id,
		Email:     email,
		Password:  pswd,
		SecretKey: secret,
	}
	type args struct {
		id     string
		email  string
		nick   string
		pswd   string
		secret string
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			"UserTest",
			args{
				id,
				email,
				nick,
				pswd,
				secret,
			},
			userT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.id, tt.args.email, tt.args.nick, tt.args.pswd, tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetByEmail(t *testing.T) {
	user1Email := "oleg.nagornij@gmail.com"
	nick := "Oleh"
	user1Pswd := "corner578"
	secret := "secret"
	user1 := New("", user1Email, nick, user1Pswd, secret)
	claim := NewClaim(
		"*",
		"admin")
	user1.Claims = append(user1.Claims, *claim)
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantU   User
		wantErr bool
	}{
		{
			"userGetByEmail",
			args{
				"oleg.nagornij@gmail.com",
			},
			*user1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, err := GetByEmail(tt.args.email)
			tt.wantU.ID = gotU.ID
			tt.wantU.Claims[0].Resource = gotU.Claims[0].Resource
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("GetByEmail() = %v, want %v", gotU, tt.wantU)
			}
		})
	}
}
