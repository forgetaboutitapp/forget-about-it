package secure_test

import (
	"cmp"
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/google/uuid"
)

func TestGetRemoteSettingsForNotExistantUser(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GetRemoteSettings(ctx, 1024, server, map[string]any{})
	if !errors.Is(err, secure.ErrCantFindUser) {
		t.Fatal("err is not ErrCantFindUser", err)
	}
}

func TestGetRemoteSettingsWhenEmpty(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GetRemoteSettings(ctx, 123, server, map[string]any{})

	if !errors.Is(err, secure.ErrCantFindUser) {
		t.Fatal("err is not ErrCantFindUser", err)
	}
}

func TestGetRemoteSettingsWhenOneLogIn(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	err := q.AddLogin(ctx, sql_queries.AddLoginParams{LoginUuid: uuid.New().String(), UserID: 123, DeviceDescription: "Description", Created: 123})
	if err != nil {
		t.Fatal("cant create login", err)
	}
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	res, err := secure.GetRemoteSettings(ctx, 123, server, map[string]any{})

	if err != nil {
		t.Fatal("cant get remote settings", err)
	}
	if !slices.Equal(res["settings"].(secure.RemoteSettings).RemoteDevices, []secure.RemoteDevice{
		{
			LastUsed:  nil,
			Title:     "Description",
			DateAdded: 123,
		},
	}) {
		t.Fatal("settings is not valid", res["settings"].(secure.RemoteSettings).RemoteDevices)
	}
}

func TestGetRemoteSettingsWhenMultipleLogIn(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	err := q.AddLogin(ctx, sql_queries.AddLoginParams{LoginUuid: uuid.New().String(), UserID: 123, DeviceDescription: "Description", Created: 123})
	if err != nil {
		t.Fatal("cant create login", err)
	}
	err = q.AddLogin(ctx, sql_queries.AddLoginParams{LoginUuid: uuid.New().String(), UserID: 123, DeviceDescription: "Description 123", Created: 234})
	if err != nil {
		t.Fatal("cant create login", err)
	}

	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	res, err := secure.GetRemoteSettings(ctx, 123, server, map[string]any{})

	if err != nil {
		t.Fatal("cant get remote settings", err)
	}
	slices.SortFunc(res["settings"].(secure.RemoteSettings).RemoteDevices, func(a, b secure.RemoteDevice) int {
		return cmp.Compare(a.DateAdded, b.DateAdded)
	})
	if !slices.Equal(res["settings"].(secure.RemoteSettings).RemoteDevices, []secure.RemoteDevice{
		{
			LastUsed:  nil,
			Title:     "Description",
			DateAdded: 123,
		},
		{
			LastUsed:  nil,
			Title:     "Description 123",
			DateAdded: 234,
		},
	}) {
		t.Fatal("settings is not valid", res["settings"].(secure.RemoteSettings).RemoteDevices)
	}
}

func TestGetRemoteSettingsWhenMultipleLogInMultipleUsers(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addUser(ctx, t, q, 234)
	err := q.AddLogin(ctx, sql_queries.AddLoginParams{LoginUuid: uuid.New().String(), UserID: 123, DeviceDescription: "Description", Created: 123})
	if err != nil {
		t.Fatal("cant create login", err)
	}
	err = q.AddLogin(ctx, sql_queries.AddLoginParams{LoginUuid: uuid.New().String(), UserID: 123, DeviceDescription: "Description 123", Created: 234})
	if err != nil {
		t.Fatal("cant create login", err)
	}
	err = q.AddLogin(ctx, sql_queries.AddLoginParams{LoginUuid: uuid.New().String(), UserID: 234, DeviceDescription: "Invalid!", Created: 234})
	if err != nil {
		t.Fatal("cant create login", err)
	}

	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	res, err := secure.GetRemoteSettings(ctx, 123, server, map[string]any{})

	if err != nil {
		t.Fatal("cant get remote settings", err)
	}
	slices.SortFunc(res["settings"].(secure.RemoteSettings).RemoteDevices, func(a, b secure.RemoteDevice) int {
		return cmp.Compare(a.DateAdded, b.DateAdded)
	})
	if !slices.Equal(res["settings"].(secure.RemoteSettings).RemoteDevices, []secure.RemoteDevice{
		{
			LastUsed:  nil,
			Title:     "Description",
			DateAdded: 123,
		},
		{
			LastUsed:  nil,
			Title:     "Description 123",
			DateAdded: 234,
		},
	}) {
		t.Fatal("settings is not valid", res["settings"].(secure.RemoteSettings).RemoteDevices)
	}
}
