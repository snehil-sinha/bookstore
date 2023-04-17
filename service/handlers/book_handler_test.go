package handlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/models/book"
	"github.com/snehil-sinha/goBookStore/models/book/mocks"
	"github.com/snehil-sinha/goBookStore/service/handlers"
	"github.com/snehil-sinha/goBookStore/test"
)

var s *common.App

func TestMain(m *testing.M) {
	// Setup
	t := &testing.T{}

	cfg, err := test.LoadTestConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	s, err = test.GetMockAppInstance(cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	gin.SetMode("test")

	// Tests
	exitCode := m.Run()

	// Teardown

	os.Exit(exitCode)
}

func TestPingHandler(t *testing.T) {
	Convey("Given a Ping Handler", t, func() {
		// Create a new gin context for the test
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		h := handlers.PingHandler()

		Convey("When a GET request is made to the endpoint", func() {
			// Make a GET request to the Ping endpoint
			h(c)

			Convey("Then the response should have a 200 status code", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("Then the response body should contain 'Pong!'", func() {
				So(w.Body.String(), ShouldEqual, "Pong!")
			})
		})
	})
}

func TestFindBooksHandler(t *testing.T) {
	Convey("Given a FindBooksHandler", t, func() {
		ctrl := gomock.NewController(t)
		ctrl.Finish()

		// Create a mock book service
		m := mocks.NewMockBookService(ctrl)

		// Get the handler to be tested
		h := handlers.FindBooksHandler(m, s)

		// Create a new gin context for the test
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Test Case 1
		Convey("When a get request is made to the endpoint", func() {
			// Define expected books response
			expectedBooks := []*book.Book{
				{
					Title: "The Theory of Everything",
					Pages: 177,
				},
				{
					Title: "George and the Big Bang",
					Pages: 560,
				},
			}

			// Set up mock book service to return expected books
			m.EXPECT().ReadAll(gomock.Any()).Return(expectedBooks, nil).Times(1)

			// Make a request to FinBooksHandler
			h(c)

			Convey("Then it should return the list of all books", func() {
				So(w.Code, ShouldEqual, http.StatusOK)

				var response struct {
					Data []*book.Book `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response.Data, ShouldResemble, expectedBooks)
			})

		})

		// Test case 2
		Convey("When an unexpected error occurs", func() {
			// Set up mock book service to return an unexpected error
			m.EXPECT().ReadAll(gomock.Any()).Return(nil, fmt.Errorf("server encountered an unknown error")).Times(1)

			// Make a request to the FindBookHandler
			h(c)

			Convey("Then it should return a 500 status code and an error message", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)

				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response["error"], ShouldEqual, "server encountered an unknown error")
			})
		})

	})
}

func TestFindBookHandler(t *testing.T) {
	Convey("Given a FindBookHandler", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create a mock book service
		m := mocks.NewMockBookService(ctrl)

		// Get the handler to be tested
		h := handlers.FindBookHandler(m, s)

		// Create a new gin context for the test
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Define a test book ID
		bookID := "61733b8e9c483c721f65b21d"

		// Test case 1
		Convey("When a valid book ID is provided", func() {
			// Define expected book response
			expectedBook := &book.Book{
				Title: "The Great Gatsby",
				Pages: 223,
			}

			// Set up mock book service to return expected book
			m.EXPECT().ReadById(gomock.Any(), bookID).Return(expectedBook, nil).Times(1)

			// Make a request to the FindBookHandler
			c.Params = []gin.Param{{Key: "id", Value: bookID}}
			h(c)

			Convey("Then it should return the book details", func() {
				So(w.Code, ShouldEqual, http.StatusOK)

				var response map[string]*book.Book
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response["data"], ShouldResemble, expectedBook)
			})
		})

		// Test case 2
		Convey("When an invalid book ID is provided", func() {
			// Set up mock book service to return error
			m.EXPECT().ReadById(gomock.Any(), bookID).Return(nil, errors.New("the provided hex string is not a valid ObjectID")).Times(1)

			// Make a request to the FindBookHandler
			c.Params = []gin.Param{{Key: "id", Value: bookID}}
			h(c)

			Convey("Then it should return an error message", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response["error"], ShouldEqual, "the provided hex string is not a valid ObjectID")
			})
		})

		// Test case 3
		Convey("When a non-existing book ID is provided", func() {
			// Set up mock book service to return nil and error
			m.EXPECT().ReadById(gomock.Any(), bookID).Return(nil, errors.New("mongo: no documents in result")).Times(1)

			// Make a request to the FindBookHandler
			c.Params = []gin.Param{{Key: "id", Value: bookID}}
			h(c)

			Convey("Then it should return an error message", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response["error"], ShouldEqual, "mongo: no documents in result")
			})
		})

		// Test case 4
		Convey("When an unexpected error occurs", func() {
			// Set up mock book service to return an unexpected error
			m.EXPECT().ReadById(gomock.Any(), bookID).Return(nil, errors.New("unexpected error")).Times(1)

			// Make a request to the FindBookHandler
			c.Params = []gin.Param{{Key: "id", Value: bookID}}
			h(c)

			Convey("Then it should return a 500 status code and an error message", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)

				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response["error"], ShouldEqual, "unexpected error")
			})
		})

	})
}
