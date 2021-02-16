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
| **eq** | Equivalence (Equals a) | [Example](#example-equal-operator)
| **and** | Logical AND operator. | [Example](#example-logical-operators)
| **or** | Logical OR operator. | [Example](#example-logical-operators)
| **gt** | Greater than (>) | [Example](#example-greater-than-less-than)
| **gte** | Greater than or Equal (>=) | [Example](#example-greater-than-or-equal-to-less-than-or-equal-to)
| **lt** | Less than (<) |  [Example](#example-greater-than-less-than)
| **lte** | Less than or Equal (<=) | [Example](#example-greater-than-or-equal-to-less-than-or-equal-to)
| **between** | Between a range os values | [Example](#example-between)
| **in** | In array of values | [Example](#example-in)
| **nin** | NOT in array of values | [Example](#example-not-in)
| **like** | Matches a specified pattern | [Example](#example-like) 
| **nlike** | Does not match a specified pattern | [Example](#example-not-like)


### Example - Equal Operator

```json
{
  "where": {
    "title": {
      "eq": "Harry Potter"
    }
  }
}
```

> Searching for books with title = `Harry Potter`


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
    ]
  }
}
```

> Searching for (`Harry Potter` **OR** `Game of Thrones` titles) **AND** published_date > `1997-06-26`


### Example - Greater than / Less than


```json
{
  "where": {
    "edition": {
      "gt": 1
    },
    "id": {
      "lt": 5
    }
  }
}
```

> Searching for any book with edition > 1 **AND** id < 5


### Example - Greater than or Equal to / Less than or Equal to


```json
{
  "where": {
    "edition": {
      "gte": 1
    },
    "id": {
      "lte": 3
    }
  }
}
```

> Searching for any book with edition >= 1 **AND** id <= 3

### Example - Between


```json
{
  "where": {
    "edition": {
      "between": [2, 3]
    }
  }
}
```

> Searching for any book with edition **between** 2 and 3. (Including 2 and 3).


### Example - In


```json
{
  "where": {
    "id": {
      "in": [1, 2, 3]
    }
  }
}
```

> Searching for any book that id is present in [1, 2, 3]


### Example - Not In


```json
{
  "where": {
    "id": {
      "nin": [1, 2, 3]
    }
  }
}
```

> Searching for any book that id is not in [1, 2, 3]

### Example - Like

```json
{
  "where": {
    "title": {
      "like": "%Potter%"
    }
  }
}
```

> Searching for any book that title contains `Potter`. The `%` symbol may change depending on Dataprovider. [Read more about SQL Like](https://www.w3schools.com/sql/sql_like.asp)


### Example - Not Like


```json
{
  "where": {
    "title": {
      "nlike": "%Potter%"
    }
  }
}
```

> Searching for any book that title **NOT** contains `Potter`.


## Skip and Limit Filter

You can easily paginate the list result with `skip` and `limit` as below:

```json
{
  "skip": 10,
  "limit": 10
}
```

Also, as any filter config, you can use URI filter format:

```
GET /api/books?filter[skip]=10&filter[limit]=10
```
