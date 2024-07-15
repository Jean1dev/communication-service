package application_test

import (
	"testing"

	"github.com/Jean1dev/communication-service/internal/application"
)

func TestDeveCriarPostComSucesso(t *testing.T) {
	authorName := "jean"
	authorAvatar := "avatar"
	message := "message"
	media := "media"

	post := application.NewPostEntityFromInputBody(authorName, authorAvatar, message, media)

	if post.IsLiked == true {
		t.Errorf("post.IsLiked should be false")
	}

	if post.Likes > 0 {
		t.Errorf("post.likes should be 0")
	}

	if media != post.Media {
		t.Errorf("post.Media should be equals a media var")
	}

	if authorName != post.Author.Name {
		t.Errorf("post.Author.Name should be equals a authorName var")
	}

	if authorAvatar != post.Author.Avatar {
		t.Errorf("post.Author.Avatar should be equals a authorAvatar var")
	}
}

func TestDeveCriarUmNovoPostViaApp(t *testing.T) {
	authorName := "jean"
	authorAvatar := "avatar"
	message := "message"
	media := "media"

	err := application.InsertNewPost(authorName, authorAvatar, message, media)

	if err != nil {
		t.Errorf("err ocurred when try to save new post")
	}
}

func TestDeveOcurrerUmErroAoSalvarPostInvalidoViaApp(t *testing.T) {
	authorName := ""
	authorAvatar := "avatar"
	message := "message"
	media := "media"

	err := application.InsertNewPost(authorName, authorAvatar, message, media)

	if err.Error() != "Author name is required" {
		t.Errorf("expect AuthorErrorName received other")
	}
}
