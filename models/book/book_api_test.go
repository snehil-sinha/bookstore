package book_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/snehil-sinha/goBookStore/db"
	"github.com/snehil-sinha/goBookStore/models/book"
	"github.com/snehil-sinha/goBookStore/test"
)

var baseUrl string

func TestMain(m *testing.M) {
	t := &testing.T{}
	ctx := context.TODO()
	cfg, err := test.LoadTestConfig()
	if err != nil {
		t.Fatalf(err.Error())
	}
	s, err := test.GetMockAppInstance(cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// start the server before running tests
	ts, err := test.StartTestSever(s)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// set the base url
	baseUrl = "http://" + ts.Addr
	// run tests
	exitCode := m.Run()

	// perform teardown

	// start a goroutine to send the shutdown signal
	go func() {
		time.Sleep(1 * time.Second)
		test.SignalShutDown(t)
	}()

	// shut down server after running tests
	test.SetupServerShutdown(s, ts)
	// clear the database to make sure there's no existing test data
	test.ClearDB(ctx)
	// close DB connection
	test.CloseDBConnection(db.Client.Client, ctx)

	// exit test
	os.Exit(exitCode)
}

func TestHealth(t *testing.T) {
	Convey("Given the /health endpoint", t, func() {

		url := baseUrl + "/health"

		Convey("When called with a GET request to /health", func() {
			client := resty.New()
			resp, err := client.R().Get(url)
			So(err, ShouldBeNil)

			Convey("Then the response status code should be 200", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusOK)

				Convey("Then the response body should be Pong!", func() {
					So(string(resp.Body()), ShouldEqual, "Pong!")

				})
			})
		})
	})
}

func TestFindBooks(t *testing.T) {
	Convey("Given the /books endpoint", t, func() {

		url := baseUrl + "/api/v1/books"

		Convey("When called with a GET request to /books", func() {
			resp, err := resty.New().R().
				Get(url)
			So(err, ShouldBeNil)

			Convey("Then the response should have a 200 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusOK)

				Convey("Then the response should return an empty list of books", func() {
					var response struct {
						Data []book.Book `json:"data"`
					}
					err := json.Unmarshal(resp.Body(), &response)
					So(err, ShouldBeNil)

					books := response.Data
					So(len(books), ShouldEqual, 0)
				})
			})
		})
	})
}

func TestFindBook(t *testing.T) {
	Convey("Given a book successfully created using the create endpoint", t, func() {

		url := baseUrl + "/api/v1/books"

		book := map[string]interface{}{
			"title": "The Chronicles of Narnia",
			"pages": 222,
		}

		resp, err := resty.New().R().
			SetBody(book).
			Post(url)
		So(err, ShouldBeNil)

		var response struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Title     string    `json:"title"`
			Pages     int       `json:"pages"`
		}
		err = json.Unmarshal(resp.Body(), &response)
		So(err, ShouldBeNil)

		Reset(func() {
			test.ClearDB(context.TODO())
		})

		Convey("When called with a GET request to /books with the existing book ID", func() {

			resp, err := resty.New().R().
				Get(url + "/" + response.ID)
			So(err, ShouldBeNil)

			Convey("Then the response should have a 200 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusOK)

				Convey("Then the response body should return the found book details", func() {

					err = json.Unmarshal(resp.Body(), &response)
					So(err, ShouldBeNil)

					So(response.CreatedAt, ShouldNotBeNil)
					So(response.UpdatedAt, ShouldNotBeNil)
					So(response.Title, ShouldEqual, "The Chronicles of Narnia")
					So(response.Pages, ShouldEqual, 222)
				})
			})
		})

		Convey("When called with a GET request to /books endpoint with non-existing book ID", func() {
			id := "63edfd5fa73563fc9a6dcde5"
			resp, err := resty.New().R().
				Get(url + "/" + id)
			So(err, ShouldBeNil)

			var response struct {
				Error string `json:"error"`
			}
			err = json.Unmarshal(resp.Body(), &response)
			So(err, ShouldBeNil)

			Convey("Then the response should have 404 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusNotFound)

				Convey("Then the response body should contain the appropriate error message", func() {

					So(string(response.Error), ShouldEqual, "mongo: no documents in result")
				})
			})

		})
	})
}

func TestCreateBook(t *testing.T) {
	Convey("Given the /books endpoint", t, func() {
		url := baseUrl + "/api/v1/books"

		Reset(func() {
			test.ClearDB(context.TODO())
		})

		Convey("When called with a POST request to /books with the given book data", func() {
			book := map[string]interface{}{
				"title": "The Chronicles of Narnia",
				"pages": 222,
			}

			resp, err := resty.New().R().
				SetBody(book).
				Post(url)
			So(err, ShouldBeNil)

			Convey("Then the response should have a 201 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusCreated)

			})

			Convey("Then the response should return the created book object with the created_at and updated_at timestamp attached", func() {
				var response struct {
					ID        string    `json:"id"`
					CreatedAt time.Time `json:"created_at"`
					UpdatedAt time.Time `json:"updated_at"`
					Title     string    `json:"title"`
					Pages     int       `json:"pages"`
				}
				err := json.Unmarshal(resp.Body(), &response)
				So(err, ShouldBeNil)

				So(response.CreatedAt, ShouldNotBeNil)
				So(response.UpdatedAt, ShouldNotBeNil)
				So(response.Title, ShouldEqual, "The Chronicles of Narnia")
				So(response.Pages, ShouldEqual, 222)

			})

		})

		Convey("When called with a POST request to /books with already existing book data", func() {
			book := map[string]interface{}{
				"title": "The Chronicles of Narnia",
				"pages": 222,
			}

			_, err := resty.New().R().
				SetBody(book).
				Post(url)
			So(err, ShouldBeNil)

			resp, err := resty.New().R().
				SetBody(book).
				Post(url)
			So(err, ShouldBeNil)

			Convey("Then the response code should have a 400 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then the response body should have an appropriate error message", func() {
				var response struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(resp.Body(), &response)
				So(err, ShouldBeNil)

				So(response.Error, ShouldNotBeBlank)
			})
		})

		Convey("When called with a POST request to /books with incomplete book data", func() {
			book := map[string]interface{}{
				"title": "test title",
			}

			resp, err := resty.New().R().
				SetBody(book).
				Post(url)
			So(err, ShouldBeNil)

			Convey("Then the response code should have a 400 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then the response body should have an appropriate error message", func() {
				var response struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(resp.Body(), &response)
				So(err, ShouldBeNil)

				So(response.Error, ShouldEqual, "validation error: book validation failed, err: Key: 'Book.Pages' Error:Field validation for 'Pages' failed on the 'required' tag")
			})
		})

		Convey("When called with a POST request to /books with invalid book data", func() {
			book := map[string]interface{}{
				"title": "test title",
				"pages": -1,
			}

			resp, err := resty.New().R().
				SetBody(book).
				Post(url)
			So(err, ShouldBeNil)

			Convey("Then the response code should have a 400 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then the response body should have an appropriate error message", func() {
				var response struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(resp.Body(), &response)
				So(err, ShouldBeNil)

				So(response.Error, ShouldEqual, "validation error: book validation failed, err: Key: 'Book.Pages' Error:Field validation for 'Pages' failed on the 'gte' tag")
			})
		})
	})
}

func TestUpdateBook(t *testing.T) {
	Convey("Given a book successfully created using the create endpoint", t, func() {

		url := baseUrl + "/api/v1/books"

		Reset(func() {
			test.ClearDB(context.TODO())
		})

		book := map[string]interface{}{
			"title": "The Chronicles of Narnia",
			"pages": 222,
		}

		resp, err := resty.New().R().
			SetBody(book).
			Post(url)
		So(err, ShouldBeNil)

		var response struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Title     string    `json:"title"`
			Pages     int       `json:"pages"`
		}
		err = json.Unmarshal(resp.Body(), &response)
		So(err, ShouldBeNil)

		Convey("When called a PUT request to /books with the existing book ID and new book details", func() {

			newBook := map[string]interface{}{
				"title": "The Chronicles of Narnia",
				"pages": 223,
			}

			resp, err := resty.New().R().
				SetBody(newBook).
				Put(url + "/" + response.ID)
			So(err, ShouldBeNil)

			Convey("Then the response should have a 200 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusOK)

				Convey("Then the response body should return the updated book object with the created_at and updated_at timestamp attached", func() {
					var response struct {
						ID        string    `json:"id"`
						CreatedAt time.Time `json:"created_at"`
						UpdatedAt time.Time `json:"updated_at"`
						Title     string    `json:"title"`
						Pages     int       `json:"pages"`
					}
					err = json.Unmarshal(resp.Body(), &response)
					So(err, ShouldBeNil)

					So(response.CreatedAt, ShouldNotBeNil)
					So(response.UpdatedAt, ShouldNotBeNil)
					So(response.Title, ShouldEqual, "The Chronicles of Narnia")
					So(response.Pages, ShouldEqual, 223)
				})
			})
		})
	})
}

func TestDeleteBook(t *testing.T) {
	Convey("Given a book successfully created using the create endpoint", t, func() {
		url := baseUrl + "/api/v1/books"

		book := map[string]interface{}{
			"title": "The Chronicles of Narnia",
			"pages": 222,
		}

		resp, err := resty.New().R().
			SetBody(book).
			Post(url)
		So(err, ShouldBeNil)

		var response struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Title     string    `json:"title"`
			Pages     int       `json:"pages"`
		}
		err = json.Unmarshal(resp.Body(), &response)
		So(err, ShouldBeNil)

		Convey("When called with a DELETE request to /books endpoint with the existing book ID", func() {

			resp, err := resty.New().R().
				Delete(url + "/" + response.ID)
			So(err, ShouldBeNil)

			Convey("Then the response should have a 200 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusOK)

				Convey("Then the response body should be yes bro.", func() {
					So(string(resp.Body()), ShouldEqual, "\"yes bro.\"")

				})
			})
		})

		Convey("When called with a DELETE request to /books endpoint with an invalid book ID", func() {
			id := "abcd"
			resp, err := resty.New().R().
				Delete(url + "/" + id)
			So(err, ShouldBeNil)

			var response struct {
				Error string `json:"error"`
			}
			err = json.Unmarshal(resp.Body(), &response)
			So(err, ShouldBeNil)

			Convey("Then the response should have 400 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusBadRequest)

				Convey("Then the response body should contain the appropriate error message", func() {

					So(string(response.Error), ShouldEqual, "the provided hex string is not a valid ObjectID")
				})
			})

		})

		Convey("When called with a DELETE request to /books endpoint with non-existing book ID", func() {
			id := "63edfd5fa73563fc9a6dcde5"
			resp, err := resty.New().R().
				Delete(url + "/" + id)
			So(err, ShouldBeNil)

			var response struct {
				Error string `json:"error"`
			}
			err = json.Unmarshal(resp.Body(), &response)
			So(err, ShouldBeNil)

			Convey("Then the response should have 404 status code", func() {
				So(resp.StatusCode(), ShouldEqual, http.StatusNotFound)

				Convey("Then the response body should contain the appropriate error message", func() {

					So(string(response.Error), ShouldEqual, "mongo: no documents in result")
				})
			})

		})
	})
}
