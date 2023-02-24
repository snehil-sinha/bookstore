
# BOOKSTORE CRUD API

This is a RESTful CRUD API for managing books in a bookstore, built in Golang using the Gin web framework and integrated with MongoDB using the mgm ODM package.


## Features

- Create, Read, Update, and Delete books using RESTful API endpoints
- MongoDB integration using the mgm ODM package
- Hooks for Creating and Updating books
- Validation for the book model
- Logging using the Zap package
- Configurations using YAML files
- Support for environment variables



## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_NAME`

`MONGODB_URI`

`LOG_PATH`


## Run Locally

Clone the project

```bash
  git clone https://github.com/snehil-sinha/bookstore.git
```

Go to the project directory

```bash
  cd bookstore
```

Install dependencies

```bash
  go mod download
```

The service configuration needs to be specified as a YAML file

```bash
  env: development  # allowed values: development production
  port: 8080
  bind: 0.0.0.0
```

Start the server

```bash
  go run main.go
```

Note: You can alternatively build the project to generate a binary executable file called goBookStore and use it to start the server

```bash
  go build
```




## API Reference

Link to the postman API doc: https://documenter.getpostman.com/view/25819993/2s93CLrtDN
## Usage

The following endpoints are available:

- `GET /health`: Health check endpoint.
- `GET /api/v1/books`: Get all books.
- `GET /api/v1/books/:id`: Get a specific book by ID.
- `POST /api/v1/books`: Create a new book.
- `PUT /api/v1/books/:id`: Update an existing book by ID.
- `DELETE /api/v1/books/:id`: Delete an existing book by ID.


## Presentation

[PPT](https://docs.google.com/presentation/d/1wmuzLwG5qvKy8W1nL5c17geIJTDKmRldkksrYsrwUfE/edit?usp=sharing)


## Authors

- [@snehil-sinha](https://www.github.com/snehil-sinha)


## 🚀 About Me
Exploring backend with Golang...

