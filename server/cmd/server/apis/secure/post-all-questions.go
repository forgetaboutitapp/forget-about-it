package secure

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"

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

func PostAllQuestions(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Could not read body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ParseFlashcard(d)
	if err != nil {
		fmt.Println(string(d))
		slog.Error("Could not unmarshal data", "data", data, "d", d, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Info("got questions list", "data", data)
	err = func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		tx, err := s.OrigDB.BeginTx(r.Context(), &sql.TxOptions{})
		if err != nil {
			slog.Error("can't initialize transaction", "err", err)
			return err
		}
		defer tx.Rollback()
		qtx := s.Db.WithTx(tx)

		// get old questions
		oldQuestions, err := qtx.GetAllQuestions(r.Context(), userid)
		if err != nil {
			slog.Error("can't get old questions", "err", err)
			return err
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
			err = qtx.UpdateQuestion(r.Context(), sql_queries.UpdateQuestionParams{Question: card.Question, Answer: card.Answer, QuestionID: card.Id, Enabled: 0})
			if err != nil {
				slog.Error("can't remove question", "id", card.Id, "err", err)
				return err
			}
		}
		slog.Info("finished deleting")
		// to update

		for _, card := range cardsToUpdate {
			err = qtx.UpdateQuestion(r.Context(), sql_queries.UpdateQuestionParams{Question: card.Question, Answer: card.Answer, QuestionID: card.Id, Enabled: 1})
			if err != nil {
				slog.Error("can't update question", "id", card.Id, "err", err)
				return err
			}
		}
		slog.Info("finished updating")
		// to add

		for _, card := range cardsToAdd {
			err = qtx.AddNewQuestion(r.Context(), sql_queries.AddNewQuestionParams{UserID: userid, Question: card.Question, Answer: card.Answer, QuestionID: card.Id, Enabled: 1})
			if err != nil {
				slog.Error("can't add question", "id", card.Id, "err", err)
				return err
			}
		}
		slog.Info("finished adding")

		err = qtx.DeleteAllTags(r.Context(), userid)
		if err != nil {
			slog.Error("can't delete tags", "err", err)
			return err
		}
		slog.Info("deleted all taggs")
		for _, val := range data {
			for _, tag := range val.Tags {
				err = qtx.AddNewTag(r.Context(), sql_queries.AddNewTagParams{
					QuestionID: int64(val.Id),
					Tag:        tag,
				})
				if err != nil {
					slog.Error("can't add new tag", "err", err)
					return err
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			slog.Error("unable to commit", "err", err)
			return err
		}
		return nil
	}()
	if err != nil {
		return
	}
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

func ParseFlashcard(str []byte) ([]Flashcard, error) {
	var data []Flashcard
	err := json.Unmarshal(str, &data)
	if err != nil {
		return []Flashcard{}, fmt.Errorf("cannot parse json (%w)", err)
	}
	return data, nil
}
