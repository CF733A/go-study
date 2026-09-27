package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type Item struct {
	Title    string
	Описание string
	Link     string
}

type GithubLister interface {
	GetItems(ctx context.Context, username string) ([]Item, error)
}

type GeneralGithubLister interface {
	GetItems(ctx context.Context, username string, strategy GithubLister) ([]Item, error)
}

type GeneralGithub struct {
	client *github.Client
}

func NewGeneralGithub(client *github.Client) *GeneralGithub {
	return &GeneralGithub{
		client: client,
	}
}

func (g *GeneralGithub) GetItems(ctx context.Context, username string, strategy GithubLister) ([]Item, error) {
	return strategy.GetItems(ctx, username)
}

type GithubGist struct {
	client *github.Client
}

func NewGithubGist(client *github.Client) *GithubGist {
	return &GithubGist{
		client: client,
	}
}

func (g *GithubGist) GetItems(ctx context.Context, username string) ([]Item, error) {
	gists, _, err := g.client.Gists.List(ctx, username, &github.GistListOptions{ListOptions: github.ListOptions{PerPage: 1000}})
	if err != nil {
		return []Item{}, fmt.Errorf("ошибка получения гистов: %v", err)
	}

	items := make([]Item, 0, len(gists))
	for _, gist := range gists {
		items = append(items, Item{
			Title:    safeDereference(gist.ID),
			Описание: safeDereference(gist.Description),
			Link:     safeDereference(gist.HTMLURL),
		})
	}

	return items, nil
}

type GithubRepo struct {
	client *github.Client
}

func NewGithubRepo(client *github.Client) *GithubRepo {
	return &GithubRepo{
		client: client,
	}
}

func (g *GithubRepo) GetItems(ctx context.Context, username string) ([]Item, error) {
	repos, _, err := g.client.Repositories.List(ctx, username, &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 1000}})
	if err != nil {
		return []Item{}, fmt.Errorf("ошибка получения репозитория: %v", err)
	}

	items := make([]Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, Item{
			Title:    safeDereference(repo.Name),
			Описание: safeDereference(repo.Description),
			Link:     safeDereference(repo.HTMLURL),
		})
	}

	return items, nil
}

func safeDereference(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "токен"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	gist := NewGithubGist(client)
	repo := NewGithubRepo(client)

	gg := NewGeneralGithub(client)

	data, err := gg.GetItems(context.Background(), "ptflp", gist)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(data)

	data, err = gg.GetItems(context.Background(), "ptflp", repo)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(data)
}
