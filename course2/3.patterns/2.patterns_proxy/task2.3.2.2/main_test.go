package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/v53/github"
)

type MockRepoLister struct {
	Repos []*github.Repository
	Err   error
}

func (m *MockRepoLister) List(ctx context.Context, username string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return m.Repos, nil, m.Err
}

type MockGistLister struct {
	Gists []*github.Gist
	Err   error
}

func (m *MockGistLister) List(ctx context.Context, username string, opt *github.GistListOptions) ([]*github.Gist, *github.Response, error) {
	return m.Gists, nil, m.Err
}

func TestGetRepos(t *testing.T) {
	tests := []struct {
		name     string
		username string
		repos    []*github.Repository
		want     []Item
		wantErr  bool
	}{
		{
			name:     "success case",
			username: "testuser",
			repos: []*github.Repository{
				{
					Name:        github.String("repo1"),
					Description: github.String("test repo"),
					HTMLURL:     github.String("http://github.com/testuser/repo1"),
				},
				{
					Name:        github.String("generated-repo-2"),
					Description: github.String("Auto-generated repository"),
					HTMLURL:     github.String("https://github.com/auto/repo-2"),
				},
				{
					Name:        github.String("generated-repo-3"),
					Description: github.String("Auto-generated repository"),
					HTMLURL:     github.String("https://github.com/auto/repo-3"),
				},
			},
			want: []Item{
				{
					Title:       "repo1",
					Description: "test repo",
					Link:        "http://github.com/testuser/repo1",
				},
				{
					Title:       "generated-repo-2",
					Description: "Auto-generated repository",
					Link:        "https://github.com/auto/repo-2",
				},
				{
					Title:       "generated-repo-3",
					Description: "Auto-generated repository",
					Link:        "https://github.com/auto/repo-3",
				},
			},
			wantErr: false,
		},
		{
			name:     "empty description",
			username: "testuser",
			repos: []*github.Repository{
				{
					Name:    github.String("repo1"),
					HTMLURL: github.String("http://github.com/testuser/repo1"),
				},
				{
					Name:        github.String("generated-repo-2"),
					Description: github.String("Auto-generated repository"),
					HTMLURL:     github.String("https://github.com/auto/repo-2"),
				},
				{
					Name:        github.String("generated-repo-3"),
					Description: github.String("Auto-generated repository"),
					HTMLURL:     github.String("https://github.com/auto/repo-3"),
				},
			},
			want: []Item{
				{
					Title: "repo1",
					Link:  "http://github.com/testuser/repo1",
				},
				{
					Title:       "generated-repo-2",
					Description: "Auto-generated repository",
					Link:        "https://github.com/auto/repo-2",
				},
				{
					Title:       "generated-repo-3",
					Description: "Auto-generated repository",
					Link:        "https://github.com/auto/repo-3",
				},
			},
			wantErr: false,
		},
		{
			name:     "error case",
			username: "testuser",
			repos:    nil,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoLister := &MockRepoLister{
				Repos: tt.repos,
			}
			if tt.wantErr {
				mockRepoLister.Err = fmt.Errorf("test error")
			}

			adapter := &GithubAdapter{
				RepoList: mockRepoLister,
			}

			got, err := adapter.GetRepos(context.Background(), tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !compareItems(got, tt.want) {
				t.Errorf("GetRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGists(t *testing.T) {
    tests := []struct {
        name     string
        username string
        gists    []*github.Gist
        want     []Item
        wantErr  bool
    }{
        {
            name:     "success case",
            username: "testuser",
            gists: []*github.Gist{
                {
                    Description: github.String("test gist"),
                    HTMLURL:     github.String("http://gist.github.com/gist1"),
                },
            },
            want: []Item{
                {
                    Title:       "test gist",
                    Description: "test gist",
                    Link:        "http://gist.github.com/gist1",
                },
                {
                    Title:       "Generated Gist 2",
                    Description: "Generated Gist 2",
                    Link:        "https://gist.github.com/generated2",
                },
                {
                    Title:       "Generated Gist 3",
                    Description: "Generated Gist 3",
                    Link:        "https://gist.github.com/generated3",
                },
            },
            wantErr: false,
        },
        {
            name:     "empty description",
            username: "testuser",
            gists: []*github.Gist{
                {
                    Files: map[github.GistFilename]github.GistFile{
                        "file1.txt": {},
                    },
                    HTMLURL: github.String("http://gist.github.com/gist1"),
                },
            },
            want: []Item{
                {
                    Title:       "file1.txt",
                    Description: "",
                    Link:        "http://gist.github.com/gist1",
                },
                {
                    Title:       "Generated Gist 2",
                    Description: "Generated Gist 2",
                    Link:        "https://gist.github.com/generated2",
                },
                {
                    Title:       "Generated Gist 3",
                    Description: "Generated Gist 3",
                    Link:        "https://gist.github.com/generated3",
                },
            },
            wantErr: false,
        },
        {
            name:     "fully empty gist",
            username: "testuser",
            gists: []*github.Gist{
                {},
            },
            want: []Item{
                {
                    Title:       "untitled",
                    Description: "",
                    Link:        "",
                },
                {
                    Title:       "Generated Gist 2",
                    Description: "Generated Gist 2",
                    Link:        "https://gist.github.com/generated2",
                },
                {
                    Title:       "Generated Gist 3",
                    Description: "Generated Gist 3",
                    Link:        "https://gist.github.com/generated3",
                },
            },
            wantErr: false,
        },
        {
            name:     "error case",
            username: "testuser",
            gists:    nil,
            want:     nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockGistLister := &MockGistLister{
                Gists: tt.gists,
            }
            if tt.wantErr {
                mockGistLister.Err = fmt.Errorf("test error")
            }

            adapter := &GithubAdapter{
                GistList: mockGistLister,
            }

            got, err := adapter.GetGists(context.Background(), tt.username)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetGists() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && !compareItems(got, tt.want) {
                t.Errorf("GetGists() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestProxyCache(t *testing.T) {
	mockRepoLister := &MockRepoLister{
		Repos: []*github.Repository{
			{
				Name:    github.String("repo1"),
				HTMLURL: github.String("http://github.com/repo1"),
			},
		},
	}

	mockGistLister := &MockGistLister{
		Gists: []*github.Gist{
			{
				ID:      github.String("gist1"),
				HTMLURL: github.String("http://gist.github.com/gist1"),
			},
		},
	}

	adapter := &GithubAdapter{
		RepoList: mockRepoLister,
		GistList: mockGistLister,
	}

	proxy := &GithubProxy{
		github: adapter,
		cache:  make(map[string][]Item),
	}

	_, err := proxy.GetRepos(context.Background(), "user1")
	if err != nil {
		t.Fatalf("GetRepos failed: %v", err)
	}

	_, err = proxy.GetRepos(context.Background(), "user1")
	if err != nil {
		t.Fatalf("GetRepos failed: %v", err)
	}

	if len(mockRepoLister.Repos) != 1 {
		t.Errorf("Expected cache to work, but got %d calls", len(mockRepoLister.Repos))
	}

	_, err = proxy.GetGists(context.Background(), "user1")
	if err != nil {
		t.Fatalf("GetGists failed: %v", err)
	}

	_, err = proxy.GetGists(context.Background(), "user1")
	if err != nil {
		t.Fatalf("GetGists failed: %v", err)
	}

	if len(mockGistLister.Gists) != 1 {
		t.Errorf("Expected cache to work, but got %d calls", len(mockGistLister.Gists))
	}
}

func compareItems(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Title != b[i].Title ||
			a[i].Description != b[i].Description ||
			a[i].Link != b[i].Link {
			return false
		}
	}
	return true
}