package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
)

// необходимые структуры
type User struct {
	ID   int
	Name string
}

type Post struct {
	ID       int
	AuthorID int
	Title    string
	Body     string
}

type Comment struct {
	ID       int
	AuthorID int
	PostID   int
	Text     string
}

// интерфейсы для работы с хранилищем
type UserRepository interface {
	Save(id int, name string)
	FindByID(id int) (string, error)
}

type PostRepository interface {
	Save(post *Post)
	FindByID(id int) (*Post, error)
	FindAll() []*Post
}

type CommentRepository interface {
	Save(comment *Comment)
	FindByPostID(postID int) []*Comment
}

// реализация юзер репо

type UserRepo struct {
	mu    sync.Mutex
	users map[int]string
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[int]string),
	}
}

func (u *UserRepo) Save(id int, name string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users[id] = name
}

func (u *UserRepo) FindByID(id int) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if val, ok := u.users[id]; ok {
		return val, nil
	}
	return "", errors.New("user not found")
}

// реализация пост репо

type PostRepo struct {
	mu    sync.Mutex
	posts map[int]*Post
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		posts: make(map[int]*Post),
	}
}

func (p *PostRepo) Save(post *Post) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.posts[post.ID] = post
}

func (p *PostRepo) FindByID(id int) (*Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if val, ok := p.posts[id]; ok {
		return val, nil
	}

	return nil, errors.New("post not found")
}

func (p *PostRepo) FindAll() []*Post {
	p.mu.Lock()
	defer p.mu.Unlock()

	result := make([]*Post, 0, len(p.posts))
	for _, post := range p.posts {
		result = append(result, post)
	}
	return result
}

// реализиция комент репо

type CommentRepo struct {
	mu       sync.Mutex
	comments map[int]*Comment
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments: make(map[int]*Comment),
	}
}

func (c *CommentRepo) Save(comment *Comment) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.comments[comment.ID] = comment
}

func (c *CommentRepo) FindByPostID(postID int) []*Comment {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make([]*Comment, 0)
	for _, comment := range c.comments {
		if comment.PostID == postID {
			result = append(result, comment)
		}
	}

	return result
}

// далее я делаю сервисную часть
// UserServicer интерфейс для операций с пользователями
type UserServicer interface {
	RegisterUser(id int, name string)
	FindUserByID(id int) (*User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) RegisterUser(id int, name string) {
	u.repo.Save(id, name)
}

func (u *UserService) FindUserByID(id int) (*User, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Name: user}, nil
}

// PostServicer интерфейс для операций с постами
type PostServicer interface {
	CreatePost(id int, authorID int, title, body string) *Post
	FindPostByID(id int) (*Post, error)
	FindAllPosts() []*Post
}

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (p *PostService) CreatePost(id int, authorID int, title, body string) *Post {
	post := &Post{
		ID:       id,
		AuthorID: authorID,
		Title:    title,
		Body:     body,
	}
	p.repo.Save(post)
	return post
}

func (p *PostService) FindPostByID(id int) (*Post, error) {
	post, err := p.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PostService) FindAllPosts() []*Post {
	return p.repo.FindAll()
}

// CommentServicer интерфейс для операций с комментариями
type CommentServicer interface {
	AddComment(id, authorID, postID int, text string) *Comment
	FindCommentsByPostID(postID int) []*Comment
}

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (c *CommentService) AddComment(id, authorID, postID int, text string) *Comment {
	comment := &Comment{
		ID:       id,
		AuthorID: authorID,
		PostID:   postID,
		Text:     text,
	}
	c.repo.Save(comment)
	return comment
}

func (c *CommentService) FindCommentsByPostID(postID int) []*Comment {
	return c.repo.FindByPostID(postID)
}

// BlogService интерфейс для агрегирования бизнес-логики блога
type BlogService interface {
	RegisterUser(id int, name string)
	CreatePost(id, authorID int, title, body string) (*AggregatedPost, error)
	AddComment(id, authorID, postID int, text string) (*AggregatedComment, error)
	GetPostWithComments(postID int) (*AggregatedPost, error)
	GetAllPosts() ([]*ListingPost, error)
}

type Blog struct {
	userService    UserServicer
	postService    PostServicer
	commentService CommentServicer
}

func NewBlogFacade(userService *UserService, postService *PostService, commentService *CommentService) BlogService {
	return &Blog{
		userService:    userService,
		postService:    postService,
		commentService: commentService,
	}
}

func (b *Blog) RegisterUser(id int, name string) {
	b.userService.RegisterUser(id, name)
}

func (b *Blog) CreatePost(id, authorID int, title, body string) (*AggregatedPost, error) {
	author, err := b.userService.FindUserByID(authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found")
	}

	post := b.postService.CreatePost(id, authorID, title, body)
	return &AggregatedPost{
		ID:       post.ID,
		Title:    post.Title,
		Body:     post.Body,
		Author:   author,
		Comments: []*AggregatedComment{},
	}, nil
}

func (b *Blog) AddComment(id, authorID, postID int, text string) (*AggregatedComment, error) {
	author, err := b.userService.FindUserByID(authorID)
	if err != nil {
		return nil, err
	}

	comment := b.commentService.AddComment(id, authorID, postID, text)
	return &AggregatedComment{
		ID:     comment.ID,
		Text:   comment.Text,
		Author: author,
	}, nil
}

func (b *Blog) GetPostWithComments(postID int) (*AggregatedPost, error) {

	post, err := b.postService.FindPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти пост с ID %d: %v", postID, err)
	}

	author, err := b.userService.FindUserByID(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти автора поста: %v", err)
	}

	comments := b.commentService.FindCommentsByPostID(postID)

	aggComments := make([]*AggregatedComment, 0, len(comments))
	for _, c := range comments {
		commentAuthor, err := b.userService.FindUserByID(c.AuthorID)
		if err != nil {
			continue // Пропускаем комментарии, если не удалось найти автора
		}

		aggComments = append(aggComments, &AggregatedComment{
			ID:     c.ID,
			Text:   c.Text,
			Author: commentAuthor,
		})
	}

	return &AggregatedPost{
		ID:       post.ID,
		Title:    post.Title,
		Body:     post.Body,
		Author:   author,
		Comments: aggComments,
	}, nil
}

func (b *Blog) GetAllPosts() ([]*ListingPost, error) {
	posts := b.postService.FindAllPosts()
	lp := make([]*ListingPost, 0, len(posts))
	for _, p := range posts {
		author, err := b.userService.FindUserByID(p.AuthorID)
		if err != nil {
			continue
		}
		comment := b.commentService.FindCommentsByPostID(p.ID)

		lp = append(lp, &ListingPost{
			ID:            p.ID,
			Title:         p.Title,
			Body:          p.Body,
			Author:        author,
			CommentsCount: len(comment),
		})
	}
	return lp, nil
}

// AggregatedPost represents a post with author and comments
type AggregatedPost struct {
	ID       int                  `json:"id"`
	Title    string               `json:"title"`
	Body     string               `json:"body"`
	Author   *User                `json:"author"`
	Comments []*AggregatedComment `json:"comments"`
}

type ListingPost struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	Author        *User  `json:"author"`
	CommentsCount int    `json:"comments_count"`
}

// AggregatedComment represents a comment with author information
type AggregatedComment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Author *User  `json:"author"`
}

// AggregatedUserPosts represents a user and their posts
type AggregatedUserPosts struct {
	User  *User   `json:"user"`
	Posts []*Post `json:"posts"`
}

func main() {
	// Backward compatibility with old golang versions
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Repositories
	userRepo := NewUserRepo()
	postRepo := NewPostRepo()
	commentRepo := NewCommentRepo()

	// Services
	userService := NewUserService(userRepo)
	postService := NewPostService(postRepo)
	commentService := NewCommentService(commentRepo)

	// Facade
	blogFacade := NewBlogFacade(userService, postService, commentService)

	// Generate users
	numUsers := 10
	for i := 1; i <= numUsers; i++ {
		name := faker.Name()
		blogFacade.RegisterUser(i, name)
	}

	// Generate posts
	numPosts := 20
	for i := 1; i <= numPosts; i++ {
		authorID := rand.Intn(numUsers) + 1
		title := faker.Sentence()
		body := faker.Paragraph()
		blogFacade.CreatePost(i, authorID, title, body)
	}

	// Generate comments
	numComments := 50
	for i := 1; i <= numComments; i++ {
		authorID := rand.Intn(numUsers) + 1
		postID := rand.Intn(numPosts) + 1
		text := faker.Sentence()
		blogFacade.AddComment(i, authorID, postID, text)
	}

	// Get all posts with their authors and comments
	allPosts, err := blogFacade.GetAllPosts()
	if err != nil {
		log.Fatal(err)
	}

	// Convert to JSON and print
	jsonAllPosts, err := json.MarshalIndent(allPosts, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonAllPosts))

	// Get post with comments and author
	aggregatedPost, _ := blogFacade.GetPostWithComments(5)
	jsonPost, _ := json.MarshalIndent(aggregatedPost, "", "  ")
	fmt.Println(string(jsonPost))
}
