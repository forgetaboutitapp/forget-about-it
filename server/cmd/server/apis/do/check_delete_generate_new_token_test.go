package do_test

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/do"
)

func TestGenerateNewTokenWithMnemonic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	q, db := start(t)
	server := do.Server{Db: q, OrigDB: db}
	addUser(ctx, t, q, 123)
	res := do.GenerateNewToken(ctx, 123, server, nil)
	if newRes := res.GetErrorMessage(); newRes != nil {
		t.Fatal("cannot generate new token: ", newRes.Error)
	}

	mnemonic := res.GetOkMessage().GetGenerateNewToken().Mnemonic
	if len(mnemonic) != 12 {
		t.Fatal("mnemonic should have 12 strings: ", mnemonic)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	set := false
	go func() {
		time.Sleep(100 * time.Millisecond)
		set = true
		returnVal := do.GetToken(ctx, do.Server{Db: server.Db, OrigDB: server.OrigDB}, &client_to_server.GetToken{TwelveWords: mnemonic})
		if msg := returnVal.GetErrorMessage(); msg != nil {
			log.Fatal("cannot get new token", msg.Error)
		}
		wg.Done()
	}()
	for {
		res := do.CheckNewToken(ctx, 123, server, &client_to_server.CheckNewToken{})
		if msg := res.GetErrorMessage(); msg != nil {
			t.Fatal("cannot check new token", msg.Error)
		}

		if done := res.GetOkMessage().GetCheckNewToken().GetDone(); done {
			if !set {
				t.Fatal("Wrong state, new Token was not set")
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}

func TestGenerateNewTokenWithUUID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	q, db := start(t)
	server := do.Server{Db: q, OrigDB: db}
	addUser(ctx, t, q, 123)
	res := do.GenerateNewToken(ctx, 123, server, nil)
	if newRes := res.GetErrorMessage(); newRes != nil {
		t.Fatal("cannot generate new token: ", newRes.Error)
	}

	uuid := res.GetOkMessage().GetGenerateNewToken().NewUuid

	var wg sync.WaitGroup
	wg.Add(1)
	set := false
	go func() {
		time.Sleep(100 * time.Millisecond)
		set = true
		returnVal := do.GetToken(ctx, do.Server{Db: server.Db, OrigDB: server.OrigDB}, &client_to_server.GetToken{Token: uuid})
		if msg := returnVal.GetErrorMessage(); msg != nil {
			log.Fatal("cannot get new token", msg.Error)
		}
		wg.Done()
	}()
	for {
		res := do.CheckNewToken(ctx, 123, server, &client_to_server.CheckNewToken{})
		if msg := res.GetErrorMessage(); msg != nil {
			t.Fatal("cannot check new token", msg.Error)
		}

		if done := res.GetOkMessage().GetCheckNewToken().GetDone(); done {
			if !set {
				t.Fatal("Wrong state, new Token was not set")
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}
