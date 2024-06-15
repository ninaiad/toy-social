package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"gateway/internal/service/posts"

	"github.com/gin-gonic/gin"
)

type postTextMsg struct {
	Text string `json:"text"`
}

type postIdMsg struct {
	PostId int64 `json:"post_id"`
}

type postResponse struct {
	Id          int64  `json:"id"`
	Text        string `json:"text"`
	TimeUpdated string `json:"time_updated"`
}

type postsByPageOutput struct {
	PageNum  int32          `json:"page_num"`
	PageSize int32          `json:"page_size"`
	Posts    []postResponse `json:"posts"`
}

func (h *Handler) createPost(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input postTextMsg
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	resp, err := h.service.CreatePost(
		context.Background(),
		&posts_proto.CreateRequest{AuthorId: userId, Text: input.Text})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("successful createPost request")
	c.JSON(http.StatusOK, postIdMsg{PostId: resp.PostId})
}

func (h *Handler) updatePost(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	postIdS, ok := c.GetQuery("id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "no id parameter for updated post")
		return
	}

	postId, err := strconv.ParseInt(postIdS, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "id parameter not a number")
		return
	}

	var input postTextMsg
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	_, err = h.service.UpdatePost(
		context.Background(),
		&posts_proto.UpdateRequest{
			AuthorId: userId,
			PostId:   postId,
			Text:     input.Text,
		})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("successful updatePost request")
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) deletePost(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	postIdS, ok := c.GetQuery("id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "no id parameter for deleted post")
		return
	}

	postId, err := strconv.ParseInt(postIdS, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "id parameter not a number")
		return
	}

	_, err = h.service.DeletePost(
		context.Background(),
		&posts_proto.PostIdRequest{
			AuthorId: userId,
			PostId:   postId,
		})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("successful deletePost request")
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) getPost(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var authorId int64
	authorIdS, ok := c.GetQuery("author_id")
	if !ok {
		authorId = userId
	} else {
		authorId, err = strconv.ParseInt(authorIdS, 10, 64)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "author_id parameter is not a number")
			return
		}
	}

	postIdS, ok := c.GetQuery("id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "no id parameter for post")
		return
	}

	postId, err := strconv.ParseInt(postIdS, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "id parameter not a number")
		return
	}

	post, err := h.service.GetPost(
		context.Background(),
		&posts_proto.PostIdRequest{
			AuthorId: authorId,
			PostId:   postId,
		})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("successful getPost request")
	c.JSON(http.StatusOK,
		postResponse{
			Id:          post.Id,
			Text:        post.Text,
			TimeUpdated: post.TimeUpdated.AsTime().Format(time.RFC3339),
		})
}

func (h *Handler) getPageOfPosts(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var authorId int64
	authorIdS, ok := c.GetQuery("author_id")
	if !ok {
		authorId = userId
	} else {
		authorId, err = strconv.ParseInt(authorIdS, 10, 64)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "author_id parameter is not a number")
			return
		}
	}

	pageNumS, ok := c.GetQuery("page_num")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "no page_num parameter")
		return
	}

	pageNum, err := strconv.Atoi(pageNumS)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "page_num parameter not a number")
		return
	}

	pageSizeS, ok := c.GetQuery("page_size")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "no page_size parameter")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeS)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "page_size parameter not a number")
		return
	}

	posts, err := h.service.GetPageOfPosts(
		context.Background(),
		&posts_proto.GetPageOfPostsRequest{
			AuthorId: authorId,
			PageNum:  int32(pageNum),
			PageSize: int32(pageSize),
		})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	postsData := make([]postResponse, posts.PageSize)
	for i := range (*posts).Posts {
		postsData[i] = postResponse{
			Id:          (*posts).Posts[i].Id,
			Text:        (*posts).Posts[i].Text,
			TimeUpdated: (*posts).Posts[i].TimeUpdated.AsTime().Format(time.RFC3339),
		}
	}

	log.Println("successful getPageOfPosts request")
	c.JSON(http.StatusOK,
		postsByPageOutput{
			PageNum:  posts.PageNum,
			PageSize: posts.PageSize,
			Posts:    postsData,
		})
}
