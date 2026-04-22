package do_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/do"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func TestGenerateNewTokenWithMnemonic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q, db := start(t)
	s := &do.Server{Db: q, OrigDB: db}
	user := sql_queries.User{UserID: 123}
	addUser(ctx, t, q, 123)
	res := do.GenerateNewToken(ctx, user, s, &v1.GenerateNewTokenRequest{})
	if errRes := res.GetError(); errRes != nil {
		t.Fatal("cannot generate new token: ", errRes.Error)
	}

	mnemonic := res.GetOk().Mnemonic
	if len(mnemonic) != 12 {
		t.Fatal("mnemonic should have 12 strings: ", mnemonic)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	set := false
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		set = true
		returnVal := do.GetToken(ctx, s, &v1.GetTokenRequest{TwelveWords: mnemonic})
		if msg := returnVal.GetError(); msg != nil {
			log.Fatal("cannot get new token", msg.Error)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			t.Fatal("timed out waiting for token check")
		default:
			res := do.CheckNewToken(ctx, user, s, &v1.CheckNewTokenRequest{})
			if msg := res.GetError(); msg != nil {
				t.Fatal("cannot check new token", msg.Error)
			}

			if done := res.GetOk().GetDone(); done {
				if !set {
					t.Fatal("Wrong state, new Token was not set")
				}
				goto done
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
done:
	wg.Wait()
}

func TestGenerateNewTokenWithUUID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q, db := start(t)
	s := &do.Server{Db: q, OrigDB: db}
	user := sql_queries.User{UserID: 123}
	addUser(ctx, t, q, 123)
	res := do.GenerateNewToken(ctx, user, s, &v1.GenerateNewTokenRequest{})
	if errRes := res.GetError(); errRes != nil {
		t.Fatal("cannot generate new token: ", errRes.Error)
	}

	newUuid := res.GetOk().NewUuid

	var wg sync.WaitGroup
	wg.Add(1)
	set := false
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		set = true
		returnVal := do.GetToken(ctx, s, &v1.GetTokenRequest{Token: newUuid})
		if msg := returnVal.GetError(); msg != nil {
			log.Fatal("cannot get new token", msg.Error)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			t.Fatal("timed out waiting for token check")
		default:
			res := do.CheckNewToken(ctx, user, s, &v1.CheckNewTokenRequest{})
			if msg := res.GetError(); msg != nil {
				t.Fatal("cannot check new token", msg.Error)
			}

			if done := res.GetOk().GetDone(); done {
				if !set {
					t.Fatal("Wrong state, new Token was not set")
				}
				goto done
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
done:
	wg.Wait()
}
