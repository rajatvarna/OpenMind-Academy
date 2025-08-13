package api

import (
	"net/http"
	"strconv"

	"github.com/free-education/forum-service/model"
	"github.com/free-education/forum-service/storage"
	"github.com/gin-gonic/gin"
)

type API struct {
	store *storage.ForumStore
}

func NewAPI(store *storage.ForumStore) *API {
	return &API{store: store}
}

func (a *API) CreateThreadHandler(c *gin.Context) {
	var req model.CreateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// In real app, UserID would come from JWT
	thread, err := a.store.CreateThread(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thread"})
		return
	}
	c.JSON(http.StatusCreated, thread)
}

func (a *API) CreatePostHandler(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// In real app, UserID would come from JWT
	post, err := a.store.CreatePost(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (a *API) GetThreadsForCourseHandler(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	threads, err := a.store.GetThreadsForCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get threads"})
		return
	}
	c.JSON(http.StatusOK, threads)
}

func (a *API) GetPostsForThreadHandler(c *gin.Context) {
	threadID, err := strconv.ParseInt(c.Param("threadId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread ID"})
		return
	}
	posts, err := a.store.GetPostsForThread(c.Request.Context(), threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}
