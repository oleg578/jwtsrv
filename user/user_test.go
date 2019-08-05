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
	pool := newPool(":6379", 2)
	defer pool.Close()
	id := uuid.New().String()
	newuser := New(id, "oleh@example.com", "oleh12345")
	clm := Claim{}
	clm.AppID = uuid.New().String()
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
			if err := u.Save(pool); (err != nil) != tt.wantErr {
				t.Errorf("User.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClaim(t *testing.T) {
	newid := uuid.New().String()
	rsc := "specialresource.com"
	asserts := make(AssertsMap, 3)
	asserts["role"] = "admin"
	asserts["account"] = "12846978"
	asserts["another"] = "something"
	clm := &Claim{
		AppID:    newid,
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
				newid,
				rsc,
				asserts,
			},
			clm,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClaim(tt.args.appid, tt.args.resource, tt.args.asserts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	id := uuid.New().String()
	email := "oleh@example.com"
	pswd := "oleh12345"
	userT := &User{
		ID:       id,
		Email:    email,
		Password: pswd,
	}
	type args struct {
		id    string
		email string
		pswd  string
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
				pswd,
			},
			userT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.id, tt.args.email, tt.args.pswd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSave(b *testing.B) {
	pool := newPool(":6379", 2)
	defer pool.Close()
	id := uuid.New().String()
	newuser := New(id, "oleh@example.com", "oleh12345")
	clm := Claim{}
	clm.AppID = uuid.New().String()
	clm.Resource = "accounts.example.com"
	clm.Asserts = make(AssertsMap)
	clm.Asserts["Account"] = "12345"
	newuser.Claims = append(newuser.Claims, clm)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newuser.Save(pool)
	}
}

func TestGetByEmail(t *testing.T) {
	pool := newPool(":6379", 2)
	defer pool.Close()
	user1Email := "oleg.nagornij@gmail.com"
	user1Pswd := "corner578"
	user1 := New("", user1Email, user1Pswd)
	asserts := make(AssertsMap)
	asserts["role"] = "admin"
	asserts["account"] = "*"
	claim := NewClaim(
		uuid.New().String(),
		"sales.bw-api.com",
		asserts)
	user1.Claims = append(user1.Claims, *claim)
	type args struct {
		email string
		pool  *redis.Pool
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
				pool,
			},
			*user1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, err := GetByEmail(tt.args.email, tt.args.pool)
			tt.wantU.ID = gotU.ID
			tt.wantU.Claims[0].AppID = gotU.Claims[0].AppID
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
