package secure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
)

func TestGradeQuestionBasic(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GradeQuestion(ctx, 1, server, map[string]any{"question-id": any(321.0), "is-correct": any(0.0)})
	if err != nil {
		t.Fatal("grade question failed", err)
	}
}

func TestGradeQuestionBadQuestionID(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GradeQuestion(ctx, 1, server, map[string]any{"question-id": any("what?!"), "is-correct": any(0.0)})
	if !errors.Is(err, secure.ErrMapIsInvalidType) {
		t.Fatal("err is not ErrMapIsInvalidType", err)
	}
}

func TestGradeQuestionBadIsCorrect(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GradeQuestion(ctx, 1, server, map[string]any{"question-id": any(2.0), "is-correct": any("hi")})
	if !errors.Is(err, secure.ErrMapIsInvalidType) {
		t.Fatal("err is not ErrMapIsInvalidType", err)
	}
}

func TestGradeQuestionOORValueForIsCorrect(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GradeQuestion(ctx, 1, server, map[string]any{"question-id": any(2.0), "is-correct": any(2.0)})
	if !errors.Is(err, secure.ErrCorrectValIsNotBoolean) {
		t.Fatal("err is not ErrCorrectValIsNotBoolean", err)
	}
}

func TestGradeQuestionQuestionIDIsInvalid(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}

	_, err := secure.GradeQuestion(ctx, 1, server, map[string]any{"question-id": any(2.0), "is-correct": any(1.0)})
	if !errors.Is(err, secure.ErrCantSaveGrades) {
		t.Fatal("err is not ErrCantSaveGrades", err)
	}
}
