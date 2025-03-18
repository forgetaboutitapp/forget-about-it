package secure

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"log/slog"
	"math/big"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/hashicorp/go-set"
)

type Flashcard struct {
	Id       int64    `json:"id"` // id can be null
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Tags     []string `json:"tags"`
}

// id is unique
func (i Flashcard) Hash() int64 {
	return i.Id
}

var ErrCantInitTransaction = errors.New("can't init transaction")
var ErrCantAddTag = errors.New("can't add tag")
var ErrCantRemoveQuestion = errors.New("can't remove question")
var ErrCantUpateQuestion = errors.New("can't update question")
var ErrUpdateQuestion = errors.New("can't update question")
var ErrCantAddQuestion = errors.New("can't add question")
var ErrCantDeleteTag = errors.New("can't delete tag")
var ErrCantCommit = errors.New("can't commit")

func PostAllQuestions(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	data := m["flashcards"].([]Flashcard)
	slog.Info("got questions list", "data", data)
	err := func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		tx, err := s.OrigDB.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			slog.Error("can't initialize transaction", "err", err)
			return errors.Join(ErrCantInitTransaction, err)
		}
		defer tx.Rollback()
		qtx := s.Db.WithTx(tx)

		// get old questions
		oldQuestions, err := qtx.GetAllQuestions(ctx, userid)
		if err != nil {
			slog.Error("can't get old questions", "err", err)
			return errors.Join(ErrCanGetQuestions, err)
		}
		fillData(data)
		oldQuestionsAsFlashcard := []Flashcard{}
		for _, v := range oldQuestions {
			oldQuestionsAsFlashcard = append(oldQuestionsAsFlashcard, Flashcard{Id: v.QuestionID, Question: v.Question, Answer: v.Answer})
		}
		slog.Info("old questions", "oldQuestionsAsFlashcard", oldQuestionsAsFlashcard)
		cardsToDelete, cardsToAdd, cardsToUpdate := UpdateCards(data, oldQuestionsAsFlashcard)
		slog.Info("Updating questions", "cardsToDelete", cardsToDelete, "cardsToAdd", cardsToAdd, "cardsToUpdate", cardsToUpdate)

		// to delete

		for _, card := range cardsToDelete {
			err = qtx.UpdateQuestion(ctx, sql_queries.UpdateQuestionParams{Question: card.Question, Answer: card.Answer, QuestionID: card.Id, Enabled: 0})
			if err != nil {
				slog.Error("can't remove question", "id", card.Id, "err", err)
				return errors.Join(ErrCantRemoveQuestion, err)
			}
		}
		slog.Info("finished deleting")
		// to update

		for _, card := range cardsToUpdate {
			err = qtx.UpdateQuestion(ctx, sql_queries.UpdateQuestionParams{Question: card.Question, Answer: card.Answer, QuestionID: card.Id, Enabled: 1})
			if err != nil {
				slog.Error("can't update question", "id", card.Id, "err", err)
				return errors.Join(ErrUpdateQuestion, err)
			}
		}
		slog.Info("finished updating")
		// to add

		for _, card := range cardsToAdd {
			err = qtx.AddNewQuestion(ctx, sql_queries.AddNewQuestionParams{UserID: userid, Question: card.Question, Answer: card.Answer, QuestionID: card.Id, Enabled: 1})
			if err != nil {
				slog.Error("can't add question", "userid", userid, "id", card.Id, "err", err)
				return errors.Join(ErrCantAddQuestion, err)
			}
		}
		slog.Info("finished adding")

		err = qtx.DeleteAllTags(ctx, userid)
		if err != nil {
			return errors.Join(ErrCantDeleteTag, err)
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
					return errors.Join(ErrCantAddTag)
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			slog.Error("unable to commit", "err", err)
			return errors.Join(ErrCantCommit, err)
		}
		return nil
	}()
	if err != nil {
		slog.Error("error", "err", err)
		return nil, err
	}
	return nil, nil
}

func fillData(data []Flashcard) {
	for i := range data {
		var questionId int64 = 0
		if data[i].Id == 0 {
			bigUserid, err := rand.Int(rand.Reader, big.NewInt(IntPow(2, 52)))
			if err != nil {
				slog.Error("can't read random values", "err", err)
				return
			}
			questionId = bigUserid.Int64()
		} else {
			questionId = int64(data[i].Id)
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
func UpdateCards(submittedQuestions []Flashcard, oldQuestions []Flashcard) ([]Flashcard, []Flashcard, []Flashcard) {
	submittedQuestionSets := set.HashSetFrom(submittedQuestions)
	oldQuestionSets := set.HashSetFrom(oldQuestions)

	submittedQuestionMap := make(map[int64]Flashcard)
	oldQuestionMap := make(map[int64]Flashcard)

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
	cardsToUpdate := set.NewHashSet[Flashcard](0)
	for _, item := range potentialToUpdate.Slice() {
		oldQuestion := oldQuestionMap[item.Id]
		newQuestion := submittedQuestionMap[item.Id]
		if oldQuestion.Question != newQuestion.Question || oldQuestion.Answer != newQuestion.Answer {
			cardsToUpdate.Insert(newQuestion)
		}
	}
	return cardsToDelete.Slice(), cardsToAdd.Slice(), cardsToUpdate.Slice()

}
