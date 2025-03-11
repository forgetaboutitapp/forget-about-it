package secure_test

import (
	"fmt"
	"testing"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
)

func TestParse(t *testing.T) {
	d, err := secure.ParseFlashcard([]byte(`[{"id":null, "question": "q", "answer": "a", "tags":["t1", "t2"]}]`))
	if err != nil {
		panic(err)
	}
	if len(d) != 1 {
		panic(fmt.Sprintln("d len is", len(d)))
	}
	if d[0].Id != 0 {
		panic("d[0] id is not zero")
	}
	if d[0].Question != "q" {
		panic("d[0] q is not q")
	}

	if d[0].Answer != "a" {
		panic("d[0] a is not a")
	}
}

// https://stackoverflow.com/a/30716481
func Ptr[T any](v T) *T {
	return &v
}

func TestUpdateCards(t *testing.T) {

	submittedCards := []secure.Flashcard{
		{
			Id:       1,
			Question: "Q1",
			Answer:   "A1",
		},
		{
			Id:       2,
			Question: "Q2",
			Answer:   "A2",
		},
		{
			Id:       3,
			Question: "Q3",
			Answer:   "A3",
		},
	}
	oldQuestions := []secure.Flashcard{
		{
			Id:       2,
			Question: "Q2",
			Answer:   "A2",
		},
		{
			Id:       3,
			Question: "Q3-1",
			Answer:   "A3-1",
		},
		{
			Id:       4,
			Question: "Q4",
			Answer:   "A4",
		},
	}
	toDelete, toAdd, toUpdate := secure.UpdateCards(submittedCards, oldQuestions)
	if toDelete[0].Id != 4 {
		panic("toDelete is invalid")
	}

	if toAdd[0].Id != 1 {
		panic("to add is not correct")
	}

	if toUpdate[0].Id != 3 {
		panic("to update is not valid")

	}
}
