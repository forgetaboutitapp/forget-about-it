package secure_test

import (
	"context"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/login"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
)

func TestGenerateNewTokenWithMnemonic(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	server := secure.Server{Db: q, OrigDB: db}

	res, err := secure.GenerateNewToken(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("cannot generate new token")
	}
	_, valid := res["new-uuid"].(string)
	if !valid {
		t.Fatal("new-uuid is not a string", res, reflect.TypeOf(res["new-uuid"]))
	}

	mnemonic, valid := res["mnemonic"].([]string)
	if !valid {
		t.Fatal("mnemonic is not a []string", res, reflect.TypeOf(res["[]mnemonic"]))
	}

	var wg sync.WaitGroup
	wg.Add(1)
	set := false
	go func() {
		time.Sleep(100 * time.Millisecond)
		set = true
		_, err := login.RealGetToken(ctx, login.PostData{TwelveWordsData: mnemonic}, login.Server{Db: server.Db, OrigDB: server.OrigDB})
		if err != nil {
			log.Fatal("cannot delete new token", err)
		}
		wg.Done()
	}()
	for {
		res, err := secure.CheckNewToken(ctx, 123, server, nil)
		if !valid {
			t.Fatal("cannot check new token", err)
		}
		if !set && res["result"].(string) == "done" {
			t.Fatal("wrong state", set, res["result"].(string))
		}
		if res["result"].(string) == "done" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}

func TestGenerateNewTokenWithUuid(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	server := secure.Server{Db: q, OrigDB: db}

	res, err := secure.GenerateNewToken(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("cannot generate new token")
	}
	uuid, valid := res["new-uuid"].(string)
	if !valid {
		t.Fatal("new-uuid is not a string", res, reflect.TypeOf(res["new-uuid"]))
	}

	_, valid = res["mnemonic"].([]string)
	if !valid {
		t.Fatal("mnemonic is not a []string", res, reflect.TypeOf(res["[]mnemonic"]))
	}

	var wg sync.WaitGroup
	wg.Add(1)
	set := false
	go func() {
		time.Sleep(100 * time.Millisecond)
		set = true
		_, err := login.RealGetToken(ctx, login.PostData{Token: uuid}, login.Server{Db: server.Db, OrigDB: server.OrigDB})
		if err != nil {
			log.Fatal("cannot delete new token", err)
		}
		wg.Done()
	}()
	for {
		res, err := secure.CheckNewToken(ctx, 123, server, nil)
		if !valid {
			t.Fatal("cannot check new token", err)
		}
		if !set && res["result"].(string) == "done" {
			t.Fatal("wrong state", set, res["result"].(string))
		}
		if res["result"].(string) == "done" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}
