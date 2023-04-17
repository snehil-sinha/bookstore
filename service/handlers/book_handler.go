package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/models/book"
)

// Used to ping
func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "Pong!")
	}
}

// FindBooksHandler fetches all books
func FindBooksHandler(svc book.BookService, s *common.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			resp []*book.Book
		)

		resp, err = svc.ReadAll(s.Log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "server encountered an unknown error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// FindBookHandler fetches the book (if present) by the specified ID
func FindBookHandler(svc book.BookService, s *common.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		queryParam := c.Param("id")

		var (
			err  error
			resp *book.Book
		)

		resp, err = svc.ReadById(s.Log, queryParam)
		if err != nil {
			if strings.EqualFold("the provided hex string is not a valid ObjectID", err.Error()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if strings.EqualFold("mongo: no documents in result", err.Error()) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

// CreateBookHandler creates a new book
func CreateBookHandler(s *common.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			err error
			req book.Book
		)

		if err = c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		book, err := book.Create(s.Log, &req)
		if err != nil {
			if strings.Contains(err.Error(), "validation") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, book)
	}
}

// UpdateBookHandler updates a book (if present) by the specified ID
func UpdateBookHandler(s *common.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var (
			err error
			req book.Book
			rsp *book.Book
		)

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error parsing the request body: " + err.Error(),
			})
			return
		}

		rsp, err = book.Update(s.Log, id, &req)
		if err != nil {
			if strings.EqualFold("the provided hex string is not a valid ObjectID", err.Error()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if strings.EqualFold("mongo: no documents in result", err.Error()) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, rsp)
	}
}

// DeleteBookHandler deletes a book (if present) by its specified id
func DeleteBookHandler(s *common.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if err := book.Delete(s.Log, id); err != nil {
			if strings.EqualFold("the provided hex string is not a valid ObjectID", err.Error()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if strings.EqualFold("mongo: no documents in result", err.Error()) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, "yes bro.")
	}
}
