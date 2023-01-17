[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/sourcesoft/ssql) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/sourcesoft/ssql/main/LICENSE)

# **DO NOT USE - WIP**

This is not an ORM. The client is just a tiny simple wrapper around `database/sql`
that provides support for simple querying pattern. It supports and provides
extra utilities that can be used with that makes it actually useful.

If you need anything more than what the API provides, you can use the `Raw` method.

## Goal

- No unnecessary extra abstraction, should be compatible with standard `database/sql`
- Opt-in for features that make common complex query patterns simple
- Be opinionated and enforce some usage patterns best practices
- Minimum use of `reflect`
- Some common utilities for everyday usage like `sqlx` scan
while still being compatible with standard `sql` lib

## Features

- Super simple, only a few query patterns are supported
- Extra utils to scan rows
- GraphQL cursor pagination
- Utils for that are Relay spec compatible with Relay connections (first, last, before, after, totalCount, cursor, ...)
- Limit and offset pagination

**Limitations**

If your DML is not a simple query that is not supported, **just use the `Raw` method instead**. 
We'll keep it intentionally simple:

- No JOINS
- No ORM features
- No transaction support to build complex queries
- Honesty most non-trival query patterns are not added 

If you think you need more patterns/utilities/methods/helpers, and it's actually useful that is hard to do
without a wrapper, feel free to open a PR.

## Getting Started

```bash
go get "github.com/sourcesoft/ssql"
```

[API Reference](https://pkg.go.dev/github.com/sourcesoft/ssql)

[Examples](https://github.com/sourcesoft/ssql/tree/main/_examples)

## Simple usage

First connect to the database.

```golang
package main

import (
  "context"
  "database/sql"

  _ "github.com/lib/pq"
  "github.com/sourcesoft/ssql"
)

func main() {
  ...
  ctx := context.Background()
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable", ...)
  dbCon, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  // You can also pass nil as options.
  options := ssql.Options{
    Tag:      "sql", // Struct tag used for SQL field name (defaults to 'sql').
    LogLevel: ssql.LevelDebug, // By default loggin is disabled.
  }
  client, err := ssql.NewClient(ctx, dbCon, &options)
  if err != nil {
    panic(err)
  }
  ...
```

Inserting is easy:

```golang
  type User struct {
    ID            *string `json:"id,omitempty" sql:"id" graph:"id" rel:"pk"`
    Username      *string `json:"username,omitempty" sql:"username" graph:"username"`
    Email         *string `json:"email,omitempty" sql:"email" graph:"email"`
    EmailVerified *bool   `json:"emailVerified,omitempty" sql:"email_verified" graph:"emailVerified"`
    Active        *bool   `json:"active,omitempty" sql:"active" graph:"active"`
    UpdatedAt     *int    `json:"updatedAt,omitempty" sql:"updated_at" graph:"updatedAt"`
    CreatedAt     *int    `json:"createdAt,omitempty" sql:"created_at" graph:"createdAt"`
    DeletedAt     *int    `json:"deletedAt,omitempty" sql:"deleted_at" graph:"deletedAt"`
  }  
  // Sample record.
  fID := "7f8d1637-ca82-4b1b-91dc-0828c98ebb34"
  fUsername := "test"
  fEmail := "test@domain.com"
  ts := 1673899847

  // Insert a new row.
  newUser := User{
    ID:        &fID,
    Username:  &fUsername,
    Email:     &fEmail,
    UpdatedAt: &ts,
    CreatedAt: &ts,
  }
  // You can pass any struct as is.
  _, err = client.Insert(ctx, "user", newUser)
  if err != nil {
    panic(err)
  }

```

We can then SELECT the row by ID:

```golang
rows, err := client.SelectByID(ctx, "user", "id", "7f8d1637-ca82-4b1b-91dc-0828c98ebb34")
if err != nil {
	panic(err)
}
// You can scan all the fields to the struct directly.
var resp User
if err := ssql.ScanOne(&resp, rows); err != nil {
	panic(err)
}
logger.Print("user %+v", resp)
```

Check the [examples](https://github.com/sourcesoft/ssql/tree/main/_examples) folder to see more.

## GraphQL

TODO

## Useful Helpers

TODO
