package secure_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
)

func TestGetAllTagsNoTags(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	server := secure.Server{Db: q, OrigDB: db}
	tags, err := secure.GetAllTags(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("err should be empty")
	}
	if len(tags["tag-set"].([]secure.TagSet)) != 0 {
		t.Fatal("tag-set should be empty", tags["tag-set"])
	}
}

func TestGetAllTagsWithOneQuestionTag(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	addQuestion(ctx, t, q, 123, 323)
	addTag(ctx, t, q, 321, []string{"a", "b", "c"})

	server := secure.Server{Db: q, OrigDB: db}
	tags, err := secure.GetAllTags(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("err should be empty")
	}

	realTags, valid := tags["tag-set"].([]secure.TagSet)
	if !valid {
		t.Fatal("invalid type", tags["tag-set"])
	}
	if !slices.Equal(realTags, []secure.TagSet{
		{
			Tag:          "a",
			NumQuestions: 1,
		},
		{
			Tag:          "b",
			NumQuestions: 1,
		},
		{
			Tag:          "c",
			NumQuestions: 1,
		},
	}) {
		t.Fatal("tags are not equal", realTags)
	}
}

func TestGetAllTagsWithTwoQuestionTag(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addQuestion(ctx, t, q, 123, 321)
	addQuestion(ctx, t, q, 123, 323)
	addTag(ctx, t, q, 321, []string{"a", "b"})
	addTag(ctx, t, q, 323, []string{"b", "c"})

	server := secure.Server{Db: q, OrigDB: db}
	tags, err := secure.GetAllTags(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("err should be empty")
	}

	realTags, valid := tags["tag-set"].([]secure.TagSet)
	if !valid {
		t.Fatal("invalid type", tags["tag-set"])
	}
	if !slices.Equal(realTags, []secure.TagSet{
		{
			Tag:          "a",
			NumQuestions: 1,
		},
		{
			Tag:          "b",
			NumQuestions: 2,
		},
		{
			Tag:          "c",
			NumQuestions: 1,
		},
	}) {
		t.Fatal("tags are not equal", realTags)
	}
}

func TestGetAllTagsWithTwoQuestionTagMultiUser(t *testing.T) {
	q, db := Init(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addUser(ctx, t, q, 123)
	addUser(ctx, t, q, 1230)
	addQuestion(ctx, t, q, 123, 321)
	addQuestion(ctx, t, q, 123, 323)
	addQuestion(ctx, t, q, 1230, 3210)
	addQuestion(ctx, t, q, 1230, 3230)
	addTag(ctx, t, q, 321, []string{"a", "b"})
	addTag(ctx, t, q, 323, []string{"b", "c"})
	addTag(ctx, t, q, 3210, []string{"a", "b"})
	addTag(ctx, t, q, 3230, []string{"b", "c", "d"})

	server := secure.Server{Db: q, OrigDB: db}
	tags, err := secure.GetAllTags(ctx, 123, server, nil)
	if err != nil {
		t.Fatal("err should be empty")
	}

	realTags, valid := tags["tag-set"].([]secure.TagSet)
	if !valid {
		t.Fatal("invalid type", tags["tag-set"])
	}
	if !slices.Equal(realTags, []secure.TagSet{
		{
			Tag:          "a",
			NumQuestions: 1,
		},
		{
			Tag:          "b",
			NumQuestions: 2,
		},
		{
			Tag:          "c",
			NumQuestions: 1,
		},
	}) {
		t.Fatal("tags are not equal", realTags)
	}
}
