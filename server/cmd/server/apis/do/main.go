package do

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

type Server struct {
	OrigDB *sql.DB
	Db     *sql_queries.Queries
}

func (s *Server) getUser(ctx context.Context, token string) (sql_queries.User, error) {
	if token == "" {
		return sql_queries.User{}, errors.New("please log in")
	}
	users, err := s.Db.FindUserByLogin(ctx, token)
	if err != nil {
		slog.Error("Unable to find user by login", "token", token, "err", err)
		return sql_queries.User{}, errors.New("token is invalid")
	}
	if len(users) != 1 {
		slog.Error("Unable to find user by login (multiple or zero users)", "token", token, "count", len(users))
		return sql_queries.User{}, errors.New("token is invalid")
	}
	// Note: FindUserByLogin returns []int64, but we need User struct.
	// Actually, FindUserByLogin returns []int64 in query.sql.go.
	// Let's check getUser result from FindUserByLogin.
	// Ah, I see. FindUserByLogin returns []int64.
	// So we need to fetch the full user.
	// Actually, most handlers only need the UserID (int64).
	// I'll change getUser to return (int64, error).

	return sql_queries.User{UserID: users[0]}, nil
}

func (s *Server) GetAllQuestions(ctx context.Context, req *v1.GetAllQuestionsRequest) (*v1.GetAllQuestionsResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GetAllQuestionsResponse{
			Result: &v1.GetAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GetAllQuestions(ctx, user, s, req), nil
}

func (s *Server) PostAllQuestions(ctx context.Context, req *v1.PostAllQuestionsRequest) (*v1.PostAllQuestionsResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.PostAllQuestionsResponse{
			Result: &v1.PostAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return PostAllQuestions(ctx, user, s, req), nil
}

func (s *Server) GenerateNewToken(ctx context.Context, req *v1.GenerateNewTokenRequest) (*v1.GenerateNewTokenResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GenerateNewTokenResponse{
			Result: &v1.GenerateNewTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GenerateNewToken(ctx, user, s, req), nil
}

func (s *Server) CheckNewToken(ctx context.Context, req *v1.CheckNewTokenRequest) (*v1.CheckNewTokenResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.CheckNewTokenResponse{
			Result: &v1.CheckNewTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return CheckNewToken(ctx, user, s, req), nil
}

func (s *Server) DeleteNewToken(ctx context.Context, req *v1.DeleteNewTokenRequest) (*v1.DeleteNewTokenResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.DeleteNewTokenResponse{
			Result: &v1.DeleteNewTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return DeleteNewToken(ctx, user, s, req), nil
}

func (s *Server) GetRemoteSettings(ctx context.Context, req *v1.GetRemoteSettingsRequest) (*v1.GetRemoteSettingsResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GetRemoteSettingsResponse{
			Result: &v1.GetRemoteSettingsResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GetRemoteSettings(ctx, user, req.Token, s, req), nil
}

func (s *Server) GetAllTags(ctx context.Context, req *v1.GetAllTagsRequest) (*v1.GetAllTagsResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GetAllTagsResponse{
			Result: &v1.GetAllTagsResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GetAllTags(ctx, user, s, req), nil
}

func (s *Server) GetNextQuestion(ctx context.Context, req *v1.GetNextQuestionRequest) (*v1.GetNextQuestionResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GetNextQuestionResponse{
			Result: &v1.GetNextQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GetNextQuestion(ctx, user, s, req), nil
}

func (s *Server) GradeQuestion(ctx context.Context, req *v1.GradeQuestionRequest) (*v1.GradeQuestionResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GradeQuestion(ctx, user, req.Token, s, req), nil
}

func (s *Server) UploadAlgorithm(ctx context.Context, req *v1.UploadAlgorithmRequest) (*v1.UploadAlgorithmResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.UploadAlgorithmResponse{
			Result: &v1.UploadAlgorithmResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return UploadAlgorithm(ctx, user, req.Token, s, req), nil
}

func (s *Server) SetDefaultAlgorithm(ctx context.Context, req *v1.SetDefaultAlgorithmRequest) (*v1.SetDefaultAlgorithmResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.SetDefaultAlgorithmResponse{
			Result: &v1.SetDefaultAlgorithmResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return SetDefaultAlgorithm(ctx, user, req.Token, s, req), nil
}

func (s *Server) RemoveLogin(ctx context.Context, req *v1.RemoveLoginRequest) (*v1.RemoveLoginResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return RemoveLogin(ctx, user, req.Token, s, req), nil
}

func (s *Server) RemoveAlgorithm(ctx context.Context, req *v1.RemoveAlgorithmRequest) (*v1.RemoveAlgorithmResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.RemoveAlgorithmResponse{
			Result: &v1.RemoveAlgorithmResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return RemoveAlgorithm(ctx, user, req.Token, s, req), nil
}

func (s *Server) GetStats(ctx context.Context, req *v1.GetStatsRequest) (*v1.GetStatsResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.GetStatsResponse{
			Result: &v1.GetStatsResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return GetStats(ctx, user, req.Token, s, req), nil
}

func (s *Server) LogData(ctx context.Context, req *v1.LogDataRequest) (*v1.LogDataResponse, error) {
	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		return &v1.LogDataResponse{
			Result: &v1.LogDataResponse_Error{
				Error: &v1.ErrorMessage{Error: err.Error(), ShouldLogOut: true},
			},
		}, nil
	}
	return LogData(ctx, user, s, req), nil
}

func (s *Server) GetToken(ctx context.Context, req *v1.GetTokenRequest) (*v1.GetTokenResponse, error) {
	// Insecure, no getUser call
	return GetToken(ctx, s, req), nil
}
