package secure

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/google/uuid"
)

type Server struct {
	Db   *sql_queries.Queries
	Next func(token uuid.UUID, s Server, w http.ResponseWriter, r *http.Request)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("in middleware")
	token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		slog.Error("No authorization", "header", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	users, err := s.Db.FindUserByLogin(r.Context(), token)
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
	err = s.Db.RegisterLogin(timeout, sql_queries.RegisterLoginParams{
		LoginUuid:   token,
		CurrentTime: time.Now().Unix(),
	})
	if err != nil {
		slog.Error("Unable to register login", "token", token, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.Next(uuid.MustParse(users[0]), s, w, r)
}
