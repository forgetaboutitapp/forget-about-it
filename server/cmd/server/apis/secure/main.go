package secure

import (
	"log/slog"
	"net/http"
	"strings"

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
		return
	}
	users, err := s.Db.FindUserByLogin(r.Context(), token)
	if err != nil {
		slog.Error("cannot find user", "err", err)

		return
	}
	if len(users) != 1 {
		slog.Error("wrong amount of users", "len", len(users), "token", token)
		return
	}
	s.Next(uuid.MustParse(users[0]), s, w, r)
}
