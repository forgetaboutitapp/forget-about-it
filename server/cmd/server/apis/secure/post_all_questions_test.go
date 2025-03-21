package secure_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
)

func TestAddTwoQuestionsWithID(t *testing.T) {
	q, db := Init(t)
	fmt.Println(server.DBFilename)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server := secure.Server{Db: q, OrigDB: db}
	r, err := secure.PostAllQuestions(ctx, 1, server, ToAny(map[string][]secure.Flashcard{
		"flashcards": {
			secure.Flashcard{Id: 1, Question: "q1", Answer: "A1", Tags: []string{"t1", "t2"}},
			secure.Flashcard{Id: 2, Question: "q2", Answer: "A2", Tags: []string{"t1", "t2"}},
		},
	}))
	if err != nil || r != nil {
		t.Fatalf("err is not null: %s", err.Error())
	}
}

func TestAddTwoQuestionsWithoutID(t *testing.T) {
	q, db := Init(t)
	fmt.Println(server.DBFilename)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server := secure.Server{Db: q, OrigDB: db}
	r, err := secure.PostAllQuestions(ctx, 1, server, ToAny(map[string][]secure.Flashcard{
		"flashcards": {
			secure.Flashcard{Question: "q1", Answer: "A1", Tags: []string{"t1", "t2"}},
			secure.Flashcard{Question: "q2", Answer: "A2", Tags: []string{"t1", "t2"}},
		},
	}))
	if err != nil || r != nil {
		t.Fatalf("err is not null: %s", err.Error())
	}
}

func TestAddAndThenUpdateQuestionsWithoutID(t *testing.T) {
	q, db := Init(t)
	fmt.Println(server.DBFilename)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server := secure.Server{Db: q, OrigDB: db}
	r, err := secure.PostAllQuestions(ctx, 1, server, ToAny(map[string][]secure.Flashcard{
		"flashcards": {
			secure.Flashcard{Question: "q1", Answer: "A1", Tags: []string{"t1", "t2"}},
			secure.Flashcard{Question: "q2", Answer: "A2", Tags: []string{"t1", "t2"}},
			secure.Flashcard{Id: 3, Question: "q3", Answer: "A3", Tags: []string{"t1", "t2"}},
		},
	}))
	if err != nil || r != nil {
		t.Fatalf("err is not null: %s", err.Error())
	}

	r, err = secure.PostAllQuestions(ctx, 1, server, ToAny(map[string][]secure.Flashcard{
		"flashcards": {
			secure.Flashcard{Id: 3, Question: "q3_new", Answer: "A3_new", Tags: []string{"t1", "t2"}},
		},
	}))
	if err != nil || r != nil {
		t.Fatalf("err is not null: %s", err.Error())
	}
	res, err := q.GetAllQuestions(ctx, 1)
	if err != nil {
		t.Fatalf("err is not null: %s", err.Error())
	}
	found := false
	for _, r := range res {
		if r.QuestionID == 3 {
			found = true
			if r.Question != "q3_new" || r.Answer != "A3_new" {
				t.Fatalf("found wrong question or answer: %s %s", r.Question, r.Answer)
			}
		}
	}
	if !found {
		t.Fatalf("coudn't find val with id")
	}
}

func ToAny(m any) map[string]any {
	str, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	var g map[string]any
	err = json.Unmarshal(str, &g)
	if err != nil {
		panic(err)
	}
	return g
}
