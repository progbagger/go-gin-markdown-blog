package main

import (
	"blog/pkg/db"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DefaultHostname string = "localhost"
	DefaultUsername string = "postgres"
	DefaultPassword
	DefaultDatabase
	DefaultPort int = 5432

	MaxPostsPerPage = 3

	MaxClientsPerSecond = 100
)

func main() {
	router := gin.Default()

	// get command-line arguments
	serverHost := flag.String("host", "0.0.0.0", "Server address")
	serverPort := flag.Int("port", 8888, "Server port")
	flag.Parse()

	// bind paths for css and assets
	router.Static("/static", "static")

	// bind paths for templates
	router.LoadHTMLGlob("templates/*")

	// create server backend
	server := Server{Authorized: gin.Accounts{
		"root": "root",
	}}

	authorizedHandler := gin.BasicAuth(server.Authorized)
	tooManyRequestsHandler := server.TooManyRequestsMiddleware()

	// setup handlers
	router.GET("/", tooManyRequestsHandler, func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, "/articles") })
	router.GET("/articles", tooManyRequestsHandler, server.GetArticles)
	router.GET("/admin", tooManyRequestsHandler, authorizedHandler, server.GetAdmin)
	router.POST("/admin", tooManyRequestsHandler, authorizedHandler, server.PostAdmin)

	// setup database
	credentials := GetPostgresCredentialsFromEnv()
	conn, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		credentials.Hostname,
		credentials.Username,
		credentials.Password,
		credentials.Database,
		credentials.Port,
	)))
	if err != nil {
		log.Fatalln(err)
	}
	err = conn.AutoMigrate(&db.Post{})
	if err != nil {
		log.Fatalln(err)
	}
	server.Db = conn
	log.Printf(
		"connected to database \"%s\" at %s:%d as %s",
		credentials.Database,
		credentials.Hostname,
		credentials.Port,
		credentials.Username,
	)

	// run server
	router.Run(fmt.Sprintf("%s:%d", *serverHost, *serverPort))
}

type Server struct {
	Db *gorm.DB

	Authorized gin.Accounts
	Clients    atomic.Uint32
}

func (s *Server) GetArticles(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	if pageStr == "" {
		ctx.Redirect(http.StatusFound, "/articles/?page=1")
		return
	}
	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		ctx.Error(err)
		ctx.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"errorCode": http.StatusBadRequest,
		})
		return
	}

	// handling bad paging
	var totalPostsCount int64
	s.Db.Model(&db.Post{}).Count(&totalPostsCount)
	maxPages := totalPostsCount / MaxPostsPerPage
	if maxPages == 0 || totalPostsCount%MaxPostsPerPage != 0 {
		maxPages++
	}
	if page > maxPages {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/articles/?page=%d", maxPages))
		return
	}
	if page < 1 {
		ctx.Redirect(http.StatusFound, "/articles/?page=1")
		return
	}

	// enabling or disabling buttons
	isPrevButtonDisabled := (page == 1)
	isNextButtonDisabled := (page == maxPages)

	// retrieving needed posts
	offset := (page - 1) * MaxPostsPerPage
	records := s.Db.Model(&db.Post{}).Order("timestamp DESC").Offset(int(offset)).Limit(MaxPostsPerPage)
	if records.Error != nil {
		ctx.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"errorCode": http.StatusInternalServerError,
		})
		return
	}
	var posts []db.Post
	records.Find(&posts)

	// rendering Makrdown in posts
	for i := range posts {
		posts[i].Text = template.HTML(blackfriday.MarkdownBasic([]byte(posts[i].Text)))
	}

	// rendering HTML page
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"posts":                posts,
		"isPrevButtonDisabled": isPrevButtonDisabled,
		"isNextButtonDisabled": isNextButtonDisabled,
		"page":                 page,
		"totalPages":           maxPages,
		"nextPage":             page + 1,
		"prevPage":             page - 1,
	})
}

func (s *Server) GetAdmin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "create.tmpl", gin.H{})
}

func (s *Server) PostAdmin(ctx *gin.Context) {
	if err := s.Db.Create(&db.Post{
		Title:     ctx.PostForm("title"),
		Text:      template.HTML(ctx.PostForm("text")),
		Timestamp: time.Now(),
	}).Error; err != nil {
		ctx.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"errorCode": http.StatusInternalServerError,
		})
	} else {
		ctx.Redirect(http.StatusFound, "/articles")
	}
}

func (s *Server) TooManyRequestsMiddleware() gin.HandlerFunc {
	go func() {
		for {
			time.Sleep(time.Second)
			s.Clients.Store(0)
		}
	}()

	return func(ctx *gin.Context) {
		if v := s.Clients.Load(); v >= MaxClientsPerSecond {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			log.Println("too many clients")
		} else {
			s.Clients.Store(v + 1)
			log.Printf("%d clients connected in that second\n", v+1)
		}
	}
}

type Credentials struct {
	Hostname string
	Username string
	Password string
	Database string
	Port     int
}

func GetPostgresCredentialsFromEnv() Credentials {
	result := Credentials{
		Hostname: DefaultHostname,
		Username: DefaultUsername,
		Password: DefaultPassword,
		Database: DefaultDatabase,
		Port:     DefaultPort,
	}

	if hostname, isPresent := os.LookupEnv("POSTGRES_HOST"); isPresent {
		result.Hostname = hostname
	}
	if username, isPresent := os.LookupEnv("POSTGRES_USER"); isPresent {
		result.Username = username
	}
	if password, isPresent := os.LookupEnv("POSTGRES_PASSWORD"); isPresent {
		result.Password = password
	}
	if database, isPresent := os.LookupEnv("POSTGRES_DATABASE"); isPresent {
		result.Database = database
	}
	if port, isPresent := os.LookupEnv("POSTGRES_PORT"); isPresent {
		intPort, err := strconv.ParseInt(port, 10, 32)
		if err == nil {
			result.Port = int(intPort)
		}
	}

	return result
}
