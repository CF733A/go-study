package main

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v53/github"
)

type MockGistLister struct {
	gists []*github.Gist
	err   error
}

func (mgl *MockGistLister) List(ctx context.Context, username string, opt *github.GistListOptions) ([]*github.Gist, *github.Response, error) {
	return mgl.gists, nil, mgl.err
}

type MockRepoLister struct {
	repos []*github.Repository
	err   error
}

func (mrl *MockRepoLister) List(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return mrl.repos, nil, mrl.err
}

func TestGetRepos(t *testing.T) {
	testRepos := []*github.Repository{
		{
			Name:        github.String("repo1"),
			Description: github.String("First repository"),
			HTMLURL:     github.String("https://github.com/user/repo1"),
		},
		{
			Name:        github.String("repo2"),
			Description: github.String("Second repository"),
			HTMLURL:     github.String("https://github.com/user/repo2"),
		},
		{
			Name:        github.String("repo3"),
			Description: github.String("Third repository"),
			HTMLURL:     github.String("https://github.com/user/repo3"),
		},
	}

	mock := &MockRepoLister{
		repos: testRepos,
	}

	adapter := GithubAdapter{
		RepoList: mock,
	}

	result, err := adapter.GetRepos(context.Background(), "testuser")

	if err != nil {
		t.Fatalf("Ожидалось отсутствие ошибки, но получили: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("Ожидалось 3 репозитория, но получили %d", len(result))
	}

	if result[0].Title != "repo1" {
		t.Errorf("Ожидалось название 'repo1', но получили '%s'", result[0].Title)
	}
}

func TestGetRepos_Empty(t *testing.T) {
    mock := &MockRepoLister{
        repos: []*github.Repository{
            {Name: github.String("generated-repo-1")},
            {Name: github.String("generated-repo-2")},
            {Name: github.String("generated-repo-3")},
        },
    }

    adapter := &GithubAdapter{
        RepoList: mock,
    }

    result, err := adapter.GetRepos(context.Background(), "emptyuser")

    if err != nil {
        t.Fatal("Не ожидалась ошибка")
    }

    if len(result) != 3 {  // Теперь ожидаем 3 элемента, так как код их добавляет
        t.Fatalf("Ожидалось 3 элемента, но получили %d", len(result))
    }
}

func TestGetGists_WithDescription(t *testing.T) {
	mock := &MockGistLister{
		gists: []*github.Gist{
			{
				Description: github.String("My first gist"),
				HTMLURL:     github.String("https://gist.github.com/1"),
				Files: map[github.GistFilename]github.GistFile{
					"main.go": {},
				},
			},
			{
				Description: github.String("My second gist"),
				HTMLURL:     github.String("https://gist.github.com/2"),
				Files: map[github.GistFilename]github.GistFile{
					"main.go": {},
				},
			},
			{
				Description: github.String("My third gist"),
				HTMLURL:     github.String("https://gist.github.com/3"),
				Files: map[github.GistFilename]github.GistFile{
					"main.go": {},
				},
			},
		},
	}

	adapter := &GithubAdapter{
		GistList: mock,
	}

	result, err := adapter.GetGists(context.Background(), "user")

	if err != nil {
		t.Fatal("Не ожидалась ошибка")
	}

	if result[0].Title != "My first gist" {
		t.Errorf("Ожидалось название из описания, но получили '%s'", result[0].Title)
	}

	if result[1].Title != "My second gist" {
		t.Errorf("Ожидалось название из описания, но получили '%s'", result[1].Title)
	}
}

func TestGetGists_WithoutDescription(t *testing.T) {
	mock := &MockGistLister{
		gists: []*github.Gist{
			{
				Description: nil,
				HTMLURL:     github.String("https://gist.github.com/1"),
				Files: map[github.GistFilename]github.GistFile{
					"main.go":   {},
					"readme.md": {},
				},
			},
			{
				Description: nil,
				HTMLURL:     github.String("https://gist.github.com/2"),
				Files: map[github.GistFilename]github.GistFile{
					"client.go": {},
					"readme.md": {},
				},
			},
			{
				Description: nil,
				HTMLURL:     github.String("https://gist.github.com/3"),
				Files: map[github.GistFilename]github.GistFile{
					"readme.txt": {},
				},
			},
		},
	}

	adapter := &GithubAdapter{
		GistList: mock,
	}

	result, err := adapter.GetGists(context.Background(), "user")

	if err != nil {
		t.Fatal("Не ожидалась ошибка")
	}

	if result[0].Title != "main.go" {
		t.Errorf("Ожидалось имя файла 'main.go', но получили '%s'", result[0].Title)
	}
}

func TestGetRepos_Error(t *testing.T) {
    expectedErr := errors.New("API rate limit exceeded")

    mock := &MockRepoLister{
        err: expectedErr,
    }

    adapter := &GithubAdapter{
        RepoList: mock,
    }

    _, err := adapter.GetRepos(context.Background(), "user")

    if err == nil {
        t.Fatal("Ожидалась ошибка, но получили nil")
    }

    // Ожидаем точное сообщение об ошибке, которое формируется в коде
    expectedMsg := "failed get repos: API rate limit exceeded"
    if err.Error() != expectedMsg {
        t.Errorf("Ожидалась ошибка '%s', но получили '%v'", expectedMsg, err)
    }
}

func TestNewGithubAdapter(t *testing.T) {
	client := github.NewClient(nil)

	adapter := NewGithubAdapter(client)

	if adapter.RepoList != client.Repositories {
		t.Error("RepoList не соответствует client.Repositories")
	}

	if adapter.GistList != client.Gists {
		t.Error("GistList не соответствует client.Gists")
	}
}
