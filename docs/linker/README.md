# Linker

Linker is a go-zepto module for those who don't want to lose time with repetitive tasks when building an API.

With few lines of codes you can build a complete CRUD API. Also, Linker is extensible, wich means you can easily apply business rules.


# Supported Datasources

 - Gorm

 > Gorm is currently the first and only datasource that exists. However, there is planning to develop others. Contributions are welcome

# Quickstart

Your First CRUD in 5 minutes. See how easy it is to create a CRUD from a Model:

```go
package main

import (
	"time"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/linker"
	gormds "github.com/go-zepto/zepto/linker/datasource/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	Name        string     `json:"name"`
	ID          uint       `gorm:"primaryKey" json:"id"`
	Edition     uint       `json:"edition"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func SetupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&Book{},
	)
	return db
}

func main() {
	db := SetupDB()
	z := zepto.NewZepto(
		zepto.Name("books-api"),
		zepto.Version("0.0.1"),
	)

	app := z.NewWeb()
	api := app.Router("/api")

	lr := linker.NewLinker(api)
	lr.AddResource(linker.Resource{
		Name:       "Book",
		Datasource: gormds.NewGormDatasource(db, &Book{}),
	})

	z.SetupHTTP("0.0.0.0:8000", app)
	z.Start()
}
```

Now we have a full CRUD in http://localhost:8000/api/books endpoint:

| Method | Endpoint | Description | Query Args
| --- | --- | --- | --- |
| **POST** | /api/books | Create a book |
| **GET** | /api/books | Retrieve a list of books | <code>filter</code>
| **GET** | /api/books/{id} | Retrieve a book by id
| **PUT** | /api/books/{id} | Update a book by id
| **DELETE** | /api/books/{id} | delete a book by id

> Note: Linker currently only supports `application/json` body format


# Filter

You can get a list of records using filters. They help you to make the exact query of what you need.

Example in JSON format:

```json
{
  "where": {
    "published_at": {
      "gt": "1997-06-26T00:00:00.000Z"
    }
  }
}
```

```
GET /api/books?filter=<JSON_ENCODED_FILTER>
```

Linker has another simpler format:

```
GET /api/books?filter[where][published_at][gt]=1995-06-26T00:00:00.000Z
```

## Where Filter

| Method | Description | |
| --- | --- | --- |
| **eq** | Equivalence (Equals a) |
| **and** | Logical AND operator. | [Example](#example-logical-operators)
| **or** | Logical OR operator. | [Example](#example-logical-operators)
| **gt** | Greater than (>) | Example
| **gte** | Greater than or Equal (>=) | Example
| **lt** | Less than (<) | Example
| **lte** | Less than or Equal (<=) | Example
| **between** | Between a range os values | Example
| **in** | In array of values | Example
| **nin** | NOT in array of values | Example
| **like** | Matches a specified pattern | Example
| **nlike** | Does not match a specified pattern | Example


### Example - Logical Operators

```json
{
  "where": {
    "and": [
      {
        "or": [
          {
            "title": {
              "eq": "Harry Potter"
            }
          },
          {
            "title": {
              "eq": "Game of Thrones"
            }
          }
        ]
      },
      {
        "published_at": {
          "gt": "1997-06-26T00:00:00.000Z"
        }
      }
    ],
  }
}
```

> Searching for (`Harry Potter` **OR** `Game of Thrones` titles) **AND** published_date > 1997-06-26
