package do

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log/slog"
	"math/big"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
	"github.com/hashicorp/go-set"
)

type LocalFlashcard struct {
	Id          uint32
	Question    string
	Answer      string
	Explanation string
	MemoHint    string
	Tags        []string
}

// Hash id is unique
func (i LocalFlashcard) Hash() uint32 {
	return i.Id
}

func PostAllQuestions(ctx context.Context, user sql_queries.User, s *Server, req *v1.PostAllQuestionsRequest) *v1.PostAllQuestionsResponse {
	var data []LocalFlashcard
	for _, v := range req.Flashcards {
		slog.Info("v:", "v", v)

		data = append(data, LocalFlashcard{
			Id:          v.Id,
			Question:    v.Question,
			Answer:      v.Answer,
			Tags:        v.Tags,
			MemoHint:    v.MemoHint,
			Explanation: v.Explanation,
		})
	}
	slog.Info("got questions list", "data", data)

	tx, err := s.OrigDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		slog.Error("can't initialize transaction", "err", err)
		return &v1.PostAllQuestionsResponse{
			Result: &v1.PostAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	defer func(tx *sql.Tx) {
		tx.Rollback()
	}(tx)
	qtx := s.Db.WithTx(tx)

	// get old questions
	oldQuestions, err := qtx.GetAllQuestions(ctx, user.UserID)
	if err != nil {
		slog.Error("can't get old questions", "err", err)
		return &v1.PostAllQuestionsResponse{
			Result: &v1.PostAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	fillData(data)
	var oldQuestionsAsFlashcard []LocalFlashcard
	for _, v := range oldQuestions {
		oldQuestionsAsFlashcard = append(oldQuestionsAsFlashcard, LocalFlashcard{
			Id:          uint32(v.QuestionID),
			Question:    v.Question,
			Answer:      v.Answer,
			MemoHint:    v.MemoHint,
			Explanation: v.Explanation,
		})
	}
	slog.Info("old questions", "oldQuestionsAsFlashcard", oldQuestionsAsFlashcard)
	cardsToDelete, cardsToAdd, cardsToUpdate := UpdateCards(data, oldQuestionsAsFlashcard)
	slog.Info("Updating questions", "cardsToDelete", cardsToDelete, "cardsToAdd", cardsToAdd, "cardsToUpdate", cardsToUpdate)

	// to delete
	for _, card := range cardsToDelete {
		err = qtx.UpdateQuestion(ctx, sql_queries.UpdateQuestionParams{Question: card.Question, Answer: card.Answer, QuestionID: int64(card.Id), Enabled: 0})
		if err != nil {
			slog.Error("can't remove question", "id", card.Id, "err", err)
			return &v1.PostAllQuestionsResponse{
				Result: &v1.PostAllQuestionsResponse_Error{
					Error: &v1.ErrorMessage{Error: "Internal Server Error"},
				},
			}
		}
	}
	slog.Info("finished deleting")

	// to update
	for _, card := range cardsToUpdate {
		err = qtx.UpdateQuestion(ctx, sql_queries.UpdateQuestionParams{Question: card.Question, Answer: card.Answer, QuestionID: int64(card.Id), Enabled: 1, Explanation: card.Explanation, MemoHint: card.MemoHint})
		if err != nil {
			slog.Error("can't update question", "id", card.Id, "err", err)
			return &v1.PostAllQuestionsResponse{
				Result: &v1.PostAllQuestionsResponse_Error{
					Error: &v1.ErrorMessage{Error: "Internal Server Error"},
				},
			}
		}
	}
	slog.Info("finished updating")

	// to add
	for _, card := range cardsToAdd {
		err = qtx.AddNewQuestion(ctx, sql_queries.AddNewQuestionParams{UserID: user.UserID, Question: card.Question, Answer: card.Answer, QuestionID: int64(card.Id), Enabled: 1, Explanation: card.Explanation, MemoHint: card.MemoHint})
		if err != nil {
			slog.Error("can't add question", "userid", user.UserID, "id", card.Id, "err", err)
			return &v1.PostAllQuestionsResponse{
				Result: &v1.PostAllQuestionsResponse_Error{
					Error: &v1.ErrorMessage{Error: "Internal Server Error"},
				},
			}
		}
	}
	slog.Info("finished adding")

	err = qtx.DeleteAllTags(ctx, user.UserID)
	if err != nil {
		slog.Error("can't delete all tags", "err", err)
		return &v1.PostAllQuestionsResponse{
			Result: &v1.PostAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	slog.Info("deleted all taggs")
	for _, val := range data {
		for _, tag := range val.Tags {
			err = qtx.AddNewTag(ctx, sql_queries.AddNewTagParams{
				QuestionID: int64(val.Id),
				Tag:        tag,
			})
			if err != nil {
				slog.Error("can't add new tag", "err", err)
				return &v1.PostAllQuestionsResponse{
					Result: &v1.PostAllQuestionsResponse_Error{
						Error: &v1.ErrorMessage{Error: "Internal Server Error"},
					},
				}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("unable to commit", "err", err)
		return &v1.PostAllQuestionsResponse{
			Result: &v1.PostAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	return &v1.PostAllQuestionsResponse{
		Result: &v1.PostAllQuestionsResponse_Ok{
			Ok: &v1.PostAllQuestions{},
		},
	}
}

func fillData(data []LocalFlashcard) {
	for i := range data {
		var questionId uint32 = 0
		if data[i].Id == 0 {
			bigUserid, err := rand.Int(rand.Reader, big.NewInt(IntPow(2, 52)))
			if err != nil {
				slog.Error("can't read random values", "err", err)
				return
			}
			questionId = uint32(bigUserid.Int64())
		} else {
			questionId = data[i].Id
		}
		data[i].Id = questionId
	}
}

// https://stackoverflow.com/questions/64108933/how-to-use-math-pow-with-integers-in-go
func IntPow(n, m int64) int64 {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	result := n
	for i := int64(2); i <= m; i++ {
		result *= n
	}
	return result
}

// Returns cardsToDelete, cardsToAdd, cardsToUpdate
func UpdateCards(submittedQuestions []LocalFlashcard, oldQuestions []LocalFlashcard) ([]LocalFlashcard, []LocalFlashcard, []LocalFlashcard) {
	submittedQuestionSets := set.HashSetFrom(submittedQuestions)
	oldQuestionSets := set.HashSetFrom(oldQuestions)

	submittedQuestionMap := make(map[uint32]LocalFlashcard)
	oldQuestionMap := make(map[uint32]LocalFlashcard)

	for _, r := range submittedQuestions {
		submittedQuestionMap[r.Id] = r
	}

	for _, r := range oldQuestions {
		oldQuestionMap[r.Id] = r
	}
	slog.Info("diff", "submitted", submittedQuestionSets, "oldQ", oldQuestionSets)
	cardsToDelete := oldQuestionSets.Difference(submittedQuestionSets)
	cardsToAdd := submittedQuestionSets.Difference(oldQuestionSets)
	potentialToUpdate := submittedQuestionSets.Intersect(oldQuestionSets)
	cardsToUpdate := set.NewHashSet[LocalFlashcard](0)
	for _, item := range potentialToUpdate.Slice() {
		oldQuestion := oldQuestionMap[item.Id]
		newQuestion := submittedQuestionMap[item.Id]
		if oldQuestion.Question != newQuestion.Question || oldQuestion.Answer != newQuestion.Answer || oldQuestion.Explanation != newQuestion.Explanation || oldQuestion.MemoHint != newQuestion.MemoHint {
			cardsToUpdate.Insert(newQuestion)
		}
	}
	return cardsToDelete.Slice(), cardsToAdd.Slice(), cardsToUpdate.Slice()
}
