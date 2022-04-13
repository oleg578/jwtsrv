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
	asserts := make(AssertsMap)
	asserts["role"] = "admin"
	asserts["account"] = "*"
	claim := NewClaim(
		uuid.New().String(),
		asserts)
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

func TestNewClaim(t *testing.T) {
	rsc := "specialresource.com"
	asserts := make(AssertsMap, 3)
	asserts["role"] = "admin"
	asserts["account"] = "12846978"
	asserts["another"] = "something"
	clm := &Claim{
		Resource: rsc,
		Asserts:  asserts,
	}
	type args struct {
		appid    string
		resource string
		asserts  AssertsMap
	}
	tests := []struct {
		name string
		args args
		want *Claim
	}{
		{
			"NewClaimTest",
			args{
				rsc,
				rsc,
				asserts,
			},
			clm,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClaim(tt.args.appid, tt.args.asserts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClaim() = %v, want %v", got, tt.want)
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

func BenchmarkSave(b *testing.B) {
	newuserEmail := "oleg.nagornij@gmail.com"
	nick := "Oleh"
	newuserPswd := "corner578"
	secret := "secret"
	newuser := New("", newuserEmail, nick, newuserPswd, secret)
	asserts := make(AssertsMap)
	asserts["role"] = "admin"
	asserts["account"] = "*"
	claim := NewClaim(uuid.New().String(), asserts)
	newuser.Claims = append(newuser.Claims, *claim)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newuser.Save()
	}
}

func TestGetByEmail(t *testing.T) {
	user1Email := "oleg.nagornij@gmail.com"
	nick := "Oleh"
	user1Pswd := "corner578"
	secret := "secret"
	user1 := New("", user1Email, nick, user1Pswd, secret)
	asserts := make(AssertsMap)
	asserts["role"] = "admin"
	asserts["account"] = "*"
	claim := NewClaim(uuid.New().String(), asserts)
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
