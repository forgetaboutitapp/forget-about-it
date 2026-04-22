package do

import (
	"context"
	"log/slog"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
	"github.com/google/uuid"
)

func GetToken(ctx context.Context, s *Server, req *v1.GetTokenRequest) *v1.GetTokenResponse {
	slog.Info("getting token", "req", req)
	var token uuid.UUID
	if t, err := uuid.Parse(req.Token); err == nil {
		token = t
	} else {
		if req.Token != "" {
			slog.Error("unable to parse token", "data", req.Token, "err", err)
			return &v1.GetTokenResponse{
				Result: &v1.GetTokenResponse_Error{
					Error: &v1.ErrorMessage{Error: "Invalid token"},
				},
			}
		}
		token, err = uuidUtils.UuidFromMnemonic(req.TwelveWords)
		if err != nil {
			slog.Error("Could not get uuid from twelve words", "words", req.TwelveWords, "err", err)
			return &v1.GetTokenResponse{
				Result: &v1.GetTokenResponse_Error{
					Error: &v1.ErrorMessage{Error: "Invalid token"},
				},
			}
		}
	}

	users, err := s.Db.FindUserByLogin(ctx, token.String())
	if err != nil {
		slog.Error("Could not find user", "token-uuid", token.String(), "err", err)
		return &v1.GetTokenResponse{
			Result: &v1.GetTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: "Invalid User"},
			},
		}
	}
	if len(users) == 0 {
		slog.Error("no users", "token-uuid", token.String())
		return &v1.GetTokenResponse{
			Result: &v1.GetTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	} else if len(users) > 1 {
		slog.Error("There should be only one user with a given token", "token-uuid", token.String(), "count", len(users))
		return &v1.GetTokenResponse{
			Result: &v1.GetTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
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
		return &v1.GetTokenResponse{
			Result: &v1.GetTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}

	slog.Info("about to look for the old key")
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		slog.Info("deleting user from waiting list", "userid", users[0])
		delete(server.UsersWaiting, users[0])
		slog.Info("done deleting")
	}()
	slog.Info("Done")
	return &v1.GetTokenResponse{
		Result: &v1.GetTokenResponse_Ok{
			Ok: &v1.GetToken{
				Token: token.String(),
			},
		},
	}
}
