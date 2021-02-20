# Linker

Linker is a go-zepto module for those who don't want to lose time with repetitive tasks when building an API.

With few lines of codes you can build a complete CRUD API. Also, Linker is extensible, wich means you can easily apply business rules.


### Supported Datasources

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


## Include Filter

You can include results from related models data when you retrieve objects.

Include works with:
- One-to-One
- One-to-Many
- Many-to-Many


Example in JSON filter:

```json
{
  "include": [
    {
      "relation": "Author"
    }
  ]
}
```

Example in URI filter:

```
GET /books?filter[include][][relation]=Author
```

You can also, perform a where in related model.

Let's assume we have an ecommerce and we need to get a user with all his active orders:


Example:

```json
{
  "include": [
    {
      "relation": "Orders",
      "where": {
        "active": {
          "eq": true
        }
      }
    }
  ]
}
```


# Hooks

Sometimes just a simple CRUD does not meet our needs. Sometimes it is necessary to insert a business rule or have a certain customized control under a Model.

Hooks was created for that. For any kind of request (List, Show, Create, Update or Destroy) it intercepts the request in 4 moments:

**Remote Hooks:**
- BeforeRemote 
- AfterRemote

**Operation Hooks:**
- BeforeOperation
- AfterOperation


## Remote Hooks

Remote hooks is called in web context. It means that you can intercept in two moments:

- BeforeRemote: When Liker receives the HTTP Request from user
- AfterRemote: When Linker will send the response back to the user

If you need to intercept the request in the sitations above, you will need to implement the RemoteHooks interface:

```go
type RemoteHooks interface {
	BeforeRemote(info RemoteHooksInfo) error
	AfterRemote(info RemoteHooksInfo) error
}
```

#### Quick Example RemoteHooks

1 - Create a struct that implements RemoteHooks:

```go
type BookRemoteHooks struct{}

func (h *BookRemoteHooks) BeforeRemote(info hooks.RemoteHooksInfo) error {
	if info.Ctx.Session().Get("user") == nil {
		info.Ctx.SetStatus(401)
		return errors.New("not logged")
	}
	return nil
}

func (h *BookRemoteHooks) AfterRemote(info hooks.RemoteHooksInfo) error {
	d := *info.Data
	d["custom_field"] = "I added this custom field to the current response data"
	if info.Endpoint == "Show" {
		d["only_show_endpoint"] = "This custom field only is available in endpoint of kind 'Show'"
	}
	return nil
}
```

2 - Add hooks to the resource configuration:

```go
...

	lr.AddResource(linker.Resource{
		Name:        "Book",
		Datasource:  gormds.NewGormDatasource(db, &Book{}),
		RemoteHooks: &BookRemoteHooks{},
	})

...
```


In Remote Hooks (Before and After), you have access of a `info` object like below:

```go
type RemoteHooksInfo struct {
	Endpoint string
	Ctx      web.Context
	ID       *string
	Data     *map[string]interface{}
}
```


| Field | Description
| --- | --- |
| info.Endpoint | kind of endpoint: `List`,`Show`,`Create`,`Update` or `Destroy`
| info.Ctx | Zepto Web Context. You have access to Request (Host, Cookies, Session, and every Zepto feature)
| info.ID | The object ID to be requested. It can be nil in some kind of endpoints |
| info.Data | The object Data that will be send to user as json. You can change it in AfterRemote to customize the final response


## Operation Hooks

Operation hooks is called in Repository context. It means that you can intercept in two moments:

- BeforeOperation: When Liker have all parsed data needed to perform operation in Model (Find, FindOne, Update, Destroy)
- AfterOperation: When Linker sucessfully executed the operation and will send the result to Remote (Web)

If you need to intercept the remote operations in the sitations above, you will need to implement the OperationHooks interface:

```go
type OperationHooks interface {
	BeforeOperation(info OperationHooksInfo) error
	AfterOperation(info OperationHooksInfo) error
}
```

In Operation Hooks (Before and After), you have access of a `info` object like below:

```go
type OperationHooksInfo struct {
	Operation    string
	ID           *string
	Data         *map[string]interface{}
	QueryContext *datasource.QueryContext
}
```

| Field | Description
| --- | --- |
| info.Operation | kind of operation: `FindOne`,`Find`,`Update`,`Create` or `Destroy`
| info.ID | The object ID to be requested. It can be nil in some kind of operation |
| info.Data | The object Data received from AfterOperation. You can change it with custom fields or removing some sensitive data, for example
| info.QueryContext | The operation Query context. You can modify the Where object here in BeforeOperation

#### Quick Example OperationHooks

1 - Create a struct that implements OperationHooks:

```go
type BookOperationHooks struct{}

func (h *BookOperationHooks) BeforeOperation(info hooks.OperationHooksInfo) error {
	if info.Operation == "Find" {
		maxLimit := int64(10)
		filter := info.QueryContext.Filter
		if filter.Limit != nil && *filter.Limit > maxLimit {
			fmt.Printf("User requested limit=%d. Changing to %d (max)\n", *filter.Limit, maxLimit)
			filter.Limit = &maxLimit
		}
	}
	return nil
}

func (h *BookOperationHooks) AfterOperation(info hooks.OperationHooksInfo) error {
	return nil
}
```

2 - Add hooks to the resource configuration:

```go
...

	lr.AddResource(linker.Resource{
		Name:        "Book",
		Datasource:  gormds.NewGormDatasource(db, &Book{}),
		OperationHooks: &BookOperationHooks{},
	})

...
```


# Manual Repository Operations

You can access the Linker Repository object and call operations manually.

Example:

```go
lr := linker.NewLinker(api)

lr.AddResource(linker.Resource{
  Name:           "Book",
  Datasource:     gormds.NewGormDatasource(db, &Book{}),
  OperationHooks: &BookOperationHooks{},
})

res, err := lr.Repository("Book").Create(context.Background(), map[string]interface{}{
  "title": "Harry Potter",
})

if err != nil {
  // do something...
}

var createdBook Book
err = res.Decode(&createdBook)
if err != nil {
  // do something
}
```

> In this case res is a `SingleResult` type and you will need to decode using `res.Decode(&out)`


You can use `lr.RepositoryDecoder("Book")` instead `lr.Repository("Book")` and the book is created and decoded in just one call.

Example:

```go
var createdBook Book

err := lr.RepositoryDecoder("Book").Create(context.Background(), map[string]interface{}{"title": "Harry Potter"}, &createdBook)
```

Also, you can create or update using the GORM model object:

```go
var createdBook Book
book := Book {
  Title: "Harry Potter",
}
err := lr.RepositoryDecoder("Book").Create(context.Background(), &book, &createdBook)
```

You can use the same GORM object as input and output. After the creation, the book object will be filled after operation:

```go
book := Book {
  Title: "Harry Potter",
}
err := lr.RepositoryDecoder("Book").Create(context.Background(), &book, &book)
if err != nil {
  // do something
}
fmt.Println(fmt.Sprintf("Book created! ID=%d", book.ID))
```
