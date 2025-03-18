package secure

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

type Server struct {
	OrigDB *sql.DB
	Db     *sql_queries.Queries
	Next   func(ctx context.Context, token int64, s Server, body map[string]any) (map[string]any, error)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("in middleware")
	token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		slog.Error("No authorization", "header", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := func() ([]int64, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		users, err := s.Db.FindUserByLogin(r.Context(), token)
		return users, err
	}()
	if err != nil {
		slog.Error("cannot find user", "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if len(users) != 1 {
		slog.Error("wrong amount of users", "len", len(users), "token", token)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	timeout, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	slog.Info("Registering", "params", sql_queries.RegisterLoginParams{
		LoginUuid:   token,
		CurrentTime: time.Now().Unix(),
	})

	err = func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		err := s.Db.RegisterLogin(timeout, sql_queries.RegisterLoginParams{
			LoginUuid:   token,
			CurrentTime: time.Now().Unix(),
		})
		return err
	}()
	if err != nil {
		slog.Error("Unable to register login", "token", token, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Unable to read body", "token", token, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	m := map[string]any{}
	if len(body) != 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			slog.Error("Unable to unmarsal body", "body", body, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	returnMap, err := s.Next(r.Context(), users[0], s, m)
	if err != nil {
		slog.Error("Writing server error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if returnMap != nil {
		b, err := json.Marshal(returnMap)
		if err != nil {
			slog.Error("cannot marshal returnMap", "returnMap", returnMap, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(b))
		if err != nil {
			panic(err)
		}
	}
}
