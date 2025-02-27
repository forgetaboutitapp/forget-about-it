package login

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
)

type Server struct {
	Db *sql_queries.Queries
}

type postData struct {
	TwelveWordsData []string `json:"twelve-words"`
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		slog.Error("Method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var data postData
	d, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Could not read body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &data)
	if err != nil {
		slog.Error("Could not unmarshal request", "d", string(d), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uuid, err := uuidUtils.UuidFromMnemonic(data.TwelveWordsData)
	if err != nil {
		slog.Error("Could not get uuid from twelve words", "words", data.TwelveWordsData, "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(5*time.Second))
	defer cancel()
	users, err := s.Db.FindUserByLogin(ctx, uuid.String())
	if err != nil {
		slog.Error("Could not find user", "token-uuid", uuid.String(), "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if len(users) == 0 {
		slog.Error("Could not find user", "token-uuid", uuid.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if len(users) > 1 {
		slog.Error("There should be only one user with a given userid", "token-uuid", uuid.String())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(map[string]string{"token": uuid.String()})
	if err != nil {
		slog.Error("Could not encode token", "val", map[string]string{"token": uuid.String()}, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}
