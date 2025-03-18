package login

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/google/uuid"
)

type Server struct {
	Db     *sql_queries.Queries
	OrigDB *sql.DB
}

type PostData struct {
	TwelveWordsData []string `json:"twelve-words"`
	Token           string   `json:"token"`
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("Method not allowed", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var data PostData
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
	b, err := RealGetToken(r.Context(), data, s)
	if err != nil {
		slog.Error("Writing server error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(b))
	if err != nil {
		panic(err)
	}

}

var ErrParseToken = errors.New("cannot parse tokens")
var ErrNoUsers = errors.New("no users")
var ErrUserIdNotUnique = errors.New("there should be only one user with a given userid")
var ErrCantRegisterLogin = errors.New("cannot register logins")
var ErrCantEncodeToken = errors.New("cannot encode token")

func RealGetToken(ctx context.Context, data PostData, s Server) (string, error) {
	slog.Info("getting token", "data", data)
	var token uuid.UUID
	if t, err := uuid.Parse(data.Token); err == nil {
		token = t
	} else {
		if data.Token != "" {
			slog.Error("unable to parse token", "data", data, "err", err)
			return "", errors.Join(ErrParseToken, err)
		}
		token, err = uuidUtils.UuidFromMnemonic(data.TwelveWordsData)
		if err != nil {
			slog.Error("Could not get uuid from twelve words", "words", data.TwelveWordsData, "err", err)
			return "", errors.Join(ErrParseToken, err)
		}
	}

	users, err := func() ([]int64, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		users, err := s.Db.FindUserByLogin(ctx, token.String())
		return users, err
	}()
	if err != nil {
		slog.Error("Could not find user", "token-uuid", token.String(), "err", err)
		return "", err
	}
	if len(users) == 0 {
		slog.Error("no users", "token-uuid", token.String())
		return "", ErrNoUsers
	} else if len(users) > 1 {
		slog.Error("There should be only one user with a given userid", "token-uuid", token.String())
		return "", ErrUserIdNotUnique
	}

	slog.Info("Registering", "params", sql_queries.RegisterLoginParams{
		LoginUuid:   token.String(),
		CurrentTime: time.Now().Unix(),
	})
	err = func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		return s.Db.RegisterLogin(ctx, sql_queries.RegisterLoginParams{
			LoginUuid:   token.String(),
			CurrentTime: time.Now().Unix(),
		})
	}()
	if err != nil {
		slog.Error("Unable to register login", "token", token, "err", err)
		return "", errors.Join(ErrCantRegisterLogin, err)
	}

	b, err := json.Marshal(map[string]string{"token": token.String()})
	if err != nil {
		slog.Error("Could not encode token", "val", map[string]string{"token": token.String()}, "err", err)
		return "", errors.Join(ErrCantEncodeToken, err)
	}
	slog.Info("about to look for the old key")
	func() {
		slog.Info("about to lock mutex")
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		slog.Info("deleting users[0]", "users", users[0])
		delete(server.UsersWaiting, users[0])
		slog.Info("done deleting")

	}()
	slog.Info("Done")
	return string(b), nil

}
