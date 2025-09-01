package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/google/uuid"
)

func GetToken(ctx context.Context, s Server, arg *client_to_server.GetToken) *server_to_client.Message {
	slog.Info("getting token", "data", arg)
	var token uuid.UUID
	if t, err := uuid.Parse(arg.Token); err == nil {
		token = t
	} else {
		if arg.Token != "" {
			slog.Error("unable to parse token", "data", arg, "err", err)
			return makeError("Invalid token")
		}
		token, err = uuidUtils.UuidFromMnemonic(arg.TwelveWords)
		if err != nil {
			slog.Error("Could not get uuid from twelve words", "words", arg.TwelveWords, "err", err)
			return makeError("Invalid token")
		}
	}

	users, err := s.Db.FindUserByLogin(ctx, token.String())
	if err != nil {
		slog.Error("Could not find user", "token-uuid", token.String(), "err", err)
		return makeError("Invalid User")
	}
	if len(users) == 0 {
		slog.Error("no users", "token-uuid", token.String())
		return makeError("Internal Server Error")
	} else if len(users) > 1 {
		slog.Error("There should be only one user with a given userid", "token-uuid", token.String())
		return makeError("Internal Server Error")
	}

	slog.Info("Registering", "params", sql_queries.RegisterLoginParams{
		LoginUuid:   token.String(),
		CurrentTime: time.Now().UTC().Unix(),
	})
	err = s.Db.RegisterLogin(ctx, sql_queries.RegisterLoginParams{
		LoginUuid:   token.String(),
		CurrentTime: time.Now().UTC().Unix(),
	})

	if err != nil {
		slog.Error("Unable to register login", "token", token, "err", err)
		return makeError("Internal Server Error")
	}

	slog.Info("about to look for the old key")
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		slog.Info("deleting users[0]", "users", users[0])
		delete(server.UsersWaiting, users[0])
		slog.Info("done deleting")

	}()
	slog.Info("Done")
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GetToken{GetToken: &server_to_client.GetToken{Token: token.String()}}})

}
