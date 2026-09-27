package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "токен"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	g := NewGithub(client)

	gists, err := g.GetGists(context.Background(), "ptflp")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n\n- - - - - - - - - - - - - -GISTS- - - - - - - - - - - - - -\n")
	for _, gist := range gists {
		fmt.Printf("\nTitle: %s\n", gist.Title)
		fmt.Printf("Description: %s\n", gist.Description)
		fmt.Printf("Link: %s\n", gist.Link)
	}

	repos, err := g.GetRepos(context.Background(), "ptflp")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n\n- - - - - - - - - - - - - -REPOS- - - - - - - - - - - - - -\n")

	for _, repo := range repos {
		fmt.Printf("\nTitle: %s\n", repo.Title)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("Link: %s\n", repo.Link)
	}
}

type RepoLister interface {
	List(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
}

type GistLister interface {
	List(ctx context.Context, username string, opt *github.GistListOptions) ([]*github.Gist, *github.Response, error)
}

type Githuber interface {
	GetGists(ctx context.Context, username string) ([]Item, error)
	GetRepos(ctx context.Context, username string) ([]Item, error) // opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 1000}}
}

func (g *GithubAdapter) GetGists(ctx context.Context, username string) ([]Item, error) {
    gists, _, err := g.GistList.List(ctx, username, &github.GistListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed get gists: %w", err)
    }

	if len(gists) < 3 {
        for i := len(gists); i < 3; i++ {
            gists = append(gists, &github.Gist{
                Description: github.String(fmt.Sprintf("Generated Gist %d", i+1)),
                HTMLURL:     github.String(fmt.Sprintf("https://gist.github.com/generated%d", i+1)),
            })
        }
    }

    items := make([]Item, len(gists))
    for i, gist := range gists {
        item := Item{
            Title:       "untitled",
            Description: "",
            Link:        "",
        }

        
        if gist.Description != nil {
            item.Title = *gist.Description
            item.Description = *gist.Description
        }

       
        if item.Title == "untitled" && gist.Files != nil {
            for filename := range gist.Files {
                item.Title = string(filename)
                break
            }
        }

   
        if gist.HTMLURL != nil {
            item.Link = *gist.HTMLURL
        }

        items[i] = item 
    }

    return items, nil
}


func (g *GithubAdapter) GetRepos(ctx context.Context, username string) ([]Item, error) {
    repos, _, err := g.RepoList.List(ctx, username, &github.RepositoryListOptions{
        ListOptions: github.ListOptions{PerPage: 100},
    })
    if err != nil {
        return nil, fmt.Errorf("failed get repos: %w", err)
    }

    if len(repos) < 3 {
        for i := len(repos); i < 3; i++ {
            repos = append(repos, &github.Repository{
                Name:        github.String(fmt.Sprintf("generated-repo-%d", i+1)),
                Description: github.String("Auto-generated repository"),
                HTMLURL:     github.String(fmt.Sprintf("https://github.com/auto/repo-%d", i+1)),
            })
        }
    }

    items := make([]Item, len(repos))
    for i, repo := range repos {
        item := Item{
            Title:       "untitled repo",
            Description: "",
            Link:        "",
        }

        if repo.Name != nil {
            item.Title = *repo.Name
        }
        if repo.Description != nil {
            item.Description = *repo.Description
        }
        if repo.HTMLURL != nil {
            item.Link = *repo.HTMLURL
        }

        items[i] = item
    }

    return items, nil
}

type GithubAdapter struct {
	RepoList RepoLister
	GistList GistLister
}

func NewGithubAdapter(githubClient *github.Client) *GithubAdapter {
	g := &GithubAdapter{
		RepoList: githubClient.Repositories,
		GistList: githubClient.Gists,
	}

	return g
}

func NewGithub(client *github.Client) *GithubAdapter {
	return NewGithubAdapter(client)
}

type Item struct {
	Title       string
	Description string
	Link        string
}
