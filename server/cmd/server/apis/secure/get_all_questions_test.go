package secure_test

import (
	"cmp"
	"context"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
)

func TestGetAllQuestions(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addUser(ctx, t, q, 1230)
	addQuestion(ctx, t, q, 123, 321)
	addTag(ctx, t, q, 321, []string{"a", "b"})

	addQuestion(ctx, t, q, 123, 322)
	addTag(ctx, t, q, 322, []string{"b", "c"})

	addQuestion(ctx, t, q, 1230, 3210)
	addTag(ctx, t, q, 3210, []string{"d", "e", "f"})

	addQuestion(ctx, t, q, 1230, 1230)
	addTag(ctx, t, q, 1230, []string{"d", "e", "f"})
	server := secure.Server{Db: q, OrigDB: db}
	res, err := secure.GetAllQuestions(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("get all questions failed", err)
	}
	flashcards, valid := res["flashcards"].([]secure.Flashcard)
	if !valid {
		t.Fatal("flashcards is not array of secure.Flashcard", reflect.TypeOf(res["flashcards"]))
	}
	flashcardIds := []int64{}
	for _, id := range flashcards {
		flashcardIds = append(flashcardIds, id.Id)
		switch id.Id {
		case 321:
			assertEqualArray(t, id.Tags, []string{"a", "b"}, "321 not equal")
		case 322:
			assertEqualArray(t, id.Tags, []string{"b", "c"}, "321 not equal")
		default:
			t.Fatal("wrong id", id.Id)
		}
	}
	slices.Sort(flashcardIds)
	if !slices.Equal(flashcardIds, []int64{321, 322}) {
		t.Fatal("wrong ids", flashcardIds)
	}
}

func assertEqualArray[A cmp.Ordered](t *testing.T, a []A, b []A, message string) {
	slices.Sort(a)
	slices.Sort(b)
	if !slices.Equal(a, b) {
		t.Fatal(message, a, b)
	}
}
