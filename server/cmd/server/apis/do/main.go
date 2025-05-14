package do

import (
	"context"
	"database/sql"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"io"
	"log/slog"
	"net/http"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	OrigDB *sql.DB
	Db     *sql_queries.Queries
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Unable to read body", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var message client_to_server.Message
	err = proto.Unmarshal(body, &message)
	if err != nil {
		slog.Error("Unable to parse proto", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if msg := message.GetSecureMessage(); msg != nil {
		processSecureMessage(r.Context(), s, msg, w)
		return
	} else if msg := message.GetInsecureMessage(); msg != nil {
		processInsecureMessage(r.Context(), s, msg, w)
		return
	}
}
func processInsecureMessage(ctx context.Context, s Server, msg *client_to_server.InsecureMessage, w http.ResponseWriter) {
	if msgInner := msg.GetGetToken(); msgInner != nil {
		write(GetToken(ctx, s, msgInner), w)
		return
	} else {
		slog.Error("Wrong type sent", "type", msg.String())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func processSecureMessage(ctx context.Context, s Server, msg *client_to_server.SecureMessage, w http.ResponseWriter) {
	token := msg.GetToken()
	if token == "" {
		write(&server_to_client.Message{ReturnMessage: &server_to_client.Message_ErrorMessage{ErrorMessage: &server_to_client.ErrorMessage{Error: "Please log in", ShouldLogOut: true}}}, w)
		return
	}
	users, err := s.Db.FindUserByLogin(ctx, token)
	if err != nil {
		slog.Error("Unable to find user by login", "token", token, "err", err)
		write(&server_to_client.Message{ReturnMessage: &server_to_client.Message_ErrorMessage{ErrorMessage: &server_to_client.ErrorMessage{Error: "Token is invalid", ShouldLogOut: true}}}, w)
		return
	}
	if len(users) != 1 {
		slog.Error("Unable to find user by login", "token", token, "err", err)
		write(&server_to_client.Message{ReturnMessage: &server_to_client.Message_ErrorMessage{ErrorMessage: &server_to_client.ErrorMessage{Error: "Token is invalid", ShouldLogOut: true}}}, w)
		return
	}
	if msgInner := msg.GetPostAllQuestions(); msgInner != nil {
		write(PostAllQuestions(ctx, users[0], s, msgInner), w)
		return
	} else if msgInner := msg.GetGetAllQuestions(); msgInner != nil {
		write(GetAllQuestions(ctx, users[0], s, msgInner), w)
		return
	} else if msgInner := msg.GetDeleteNewToken(); msgInner != nil {
		write(DeleteNewToken(ctx, users[0], s, msgInner), w)
		return
	} else if msgInner := msg.GetGetAllTags(); msgInner != nil {
		write(GetAllTags(ctx, users[0], s, msgInner), w)
		return
	} else if msgInner := msg.GetGetNextQuestion(); msgInner != nil {
		write(GetNextQuestion(ctx, users[0], s, msgInner), w)
		return
	} else if msgInner := msg.GetGetRemoteSettings(); msgInner != nil {
		write(GetRemoteSettings(ctx, users[0], token, s, msgInner), w)
		return
	} else if msgInner := msg.GetCheckNewToken(); msgInner != nil {
		write(CheckNewToken(ctx, users[0], s, msgInner), w)
		return
	} else if msgInner := msg.GetUploadAlgorithm(); msgInner != nil {
		write(UploadAlgorithm(ctx, users[0], token, s, msgInner), w)
		return
	} else if msgInner := msg.GetUploadAlgorithm(); msgInner != nil {
		write(UploadAlgorithm(ctx, users[0], token, s, msgInner), w)
		return
	} else if msgInner := msg.GetSetDefaultAlgorithm(); msgInner != nil {
		write(SetDefaultAlgorithm(ctx, users[0], token, s, msgInner), w)
		return
	} else if msgInner := msg.GetGradeQuestion(); msgInner != nil {
		write(GradeQuestion(ctx, users[0], token, s, msgInner), w)
		return
	} else if msgInner := msg.GetGetStats(); msgInner != nil {
		write(GetStats(ctx, users[0], token, s, msgInner), w)
		return
	} else {
		slog.Error("Wrong type sent", "type", msg.String())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func write(token *server_to_client.Message, w http.ResponseWriter) {
	slog.Info("writing token", "token", token)
	v, err := proto.Marshal(token)
	if err != nil {
		slog.Error("Unable to serialize message", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(v)
	if err != nil {
		slog.Error("Unable to send message", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func makeError(err string) *server_to_client.Message {
	return &server_to_client.Message{ReturnMessage: &server_to_client.Message_ErrorMessage{ErrorMessage: &server_to_client.ErrorMessage{Error: err}}}
}

func makeOk(t *server_to_client.OkMessage) *server_to_client.Message {
	return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{OkMessage: t}}
}
