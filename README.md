[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/sourcesoft/ssql) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/sourcesoft/ssql/main/LICENSE)

**This library is still in development and the API may change**

### Table of Contents

- [What it is and what it's not](#what-it-is-and-what-its-not)
  * [Goal](#goal)
  * [Features](#features)
  * [Limitations](#limitations)
- [Getting Started](#getting-started)
- [APIs](#apis)
  * [Insert](#insert)
  * [UpdateOne](#updateone)
  * [Update](#update)
  * [DeleteOne](#deleteone)
  * [Delete](#delete)
  * [FindOne](#findone)
  * [Find](#find)
    + [Find: Minimal](#find-minimal)
    + [Find: Offset Pagination](#find-offset-pagination)
    + [Find: Cursor Pagination](#find-cursor-pagination)
    + [Find: Sorting](#find-sorting)
    + [Find: Conditions](#find-conditions)
    + [Find: GraphQL](#find-graphql)
  * [Raw](#raw)
  * [Helpers working with structs](#helpers-working-with-structs)
  * [Helpers scanning query results](#helpers-scanning-query-results)

## What it is and what it's not

This is not an ORM. The client is just a tiny simple wrapper around `database/sql`
that provides support for simple querying pattern. It supports and provides
extra utilities that can be used with that makes it actually useful.

If you need anything more than what the API provides, you can use the `Raw` method.

[API Reference](https://pkg.go.dev/github.com/sourcesoft/ssql)

[Examples](https://github.com/sourcesoft/ssql/tree/main/_examples)

### Goal

- No unnecessary extra abstraction, should be compatible with standard `database/sql`
- Opt-in for features that make common complex query patterns simple
- Be opinionated and enforce some usage patterns best practices
- Minimum use of `reflect`
- Some common utilities for everyday usage like `sqlx` scan
while still being compatible with standard `sql` lib

### Features

- Super simple, only a few query patterns are supported
- Extra utils to scan rows
- GraphQL cursor pagination
- Utils for that are Relay spec compatible with Relay connections (first, last, before, after, totalCount, cursor, ...)
- Limit and offset pagination

### Limitations

If your DML is not a simple query that is not supported, **just use the `Raw` method instead**. 
We'll keep it intentionally simple:

- No JOINS
- No ORM features
- No transaction support to build complex queries
- Honesty most non-trival query patterns are not added 

If you think you need more patterns/utilities/methods/helpers, and it's actually useful that is hard to do
without a wrapper, feel free to open a PR.

[↩](#table-of-contents)

## Getting Started

Get the library with:

```bash
go get "github.com/sourcesoft/ssql"
```

[API Reference](https://pkg.go.dev/github.com/sourcesoft/ssql)

[Examples](https://github.com/sourcesoft/ssql/tree/main/_examples)

First create the client by connecting to a database of your choice.

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

See queries like [FindOne](#findone) or other APIs to see how to use the client to execute queries.

[↩](#table-of-contents)

## APIs

### Insert

When using insert, you will be passing the variable of a struct, which `ssql` uses `reflect` package to
extract the `sql` tag by default (you can customize it in the client options).

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

[↩](#table-of-contents)

### UpdateOne

Having the ID of the record you can simply update it by passing the struct variable.

```golang
// Update row by ID.
fEmail = "new@test.com"
newUser.Email = &fEmail
res, err := client.UpdateOne(ctx, "user", "id", fID, newUser)
if err != nil {
  log.Error().Err(err).Msg("Postgres update user error")
  panic(err)
}
if count, err := (*res).RowsAffected(); count < 1 {
  log.Error().Err(err).Msg("Postgres update user error, or not found")
  panic(err)
}
```

[↩](#table-of-contents)

### Update

You can also create a condition array to update all the matching fields.

```golang
...
// Add some custom conditions.
conds := []*ssql.ConditionPair{{
  Field: "active",
  Value: true,
  Op:    "=",
}}
fFalse = false
newUser.Active = &fFalse
res, err := client.Update(ctx, "user", conds, updatedUser)
...
```

[↩](#table-of-contents)

### DeleteOne

```golang
res, err = client.DeleteOne(ctx, "user", "id", fID)
if err != nil {
  log.Error().Err(err).Msg("Cannot delete user by ID from Postgres")
  panic(err)
}
if count, err := (*res).RowsAffected(); count < 1 {
  log.Error().Err(err).Msg("User not found")
  panic(err)
}
```

[↩](#table-of-contents)

### Delete

You can also create a condition array to delete all the matching fields.

```golang
...
// Add some custom conditions.
conds := []*ssql.ConditionPair{{
  Field: "active",
  Value: false,
  Op:    "=",
}}
res, err = client.Delete(ctx, "user", conds)
if err != nil {
  log.Error().Err(err).Msg("Cannot delete user")
  panic(err)
}
...
```

[↩](#table-of-contents)

### FindOne

Having a primary key and finding your record using that is a common use case.
You can also pass the key (`id` in the following example) to look up.

```golang
rows, err := client.FindOne(ctx, "user", "id", "7f8d1637-ca82-4b1b-91dc-0828c98ebb34")
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


[↩](#table-of-contents)

### Find

#### Find: Minimal

Let's see how a minimal simple `Find` query looks like.

First build the query options

```golang
// setting up pagination.
limit := 10
params := ssql.Params{
  OffsetParams: &ssql.OffsetParams{
    Limit: &limit,
    // There's also offset available.
  },
  // There's also 'order' available.
  // There's also 'cursor' pagination available.
}
opts := ssql.SQLQueryOptions{
  Table:          "user",
  WithTotalCount: shouldReturnTotalCount,
  Params:         &params,
}
```

Run the query by using the `Find` method.

```golang
// Executing the query.
rows, total, err := client.Find(ctx, &opts)
if err != nil {
  log.Error().Err(err).Msg("Cannot find users")
  panic(err)
}
```

Scan the rows into your arrays of custom structs.

```golang
// Reading through the results.
var users []User
for rows.Next() {
  var user User
  if err := ssql.ScanRow(&user, rows); err != nil {
    log.Error().Err(err).Msg("Cannot scan users")
  }
  users = append(users, user)
}
```

[↩](#table-of-contents)

#### Find: Offset Pagination

You can use offset and limit pagination.

```golang
limit := 10
offset := 3
params := ssql.Params{
  OffsetParams: &ssql.OffsetParams{
    Limit: &limit,
    Offset: &offset
  },
}
// Same as before use the params in query options argument.
opts := ssql.SQLQueryOptions{
  Table:          "user",
  Params:         &params,
}
rows, total, err := client.Find(ctx, &opts)
```

[↩](#table-of-contents)

#### Find: Cursor Pagination

```golang
params := ssql.Params{
	CursorParams: &ssql.CursorParams{
    After:  previousCursor, // If you have the previous cursor, you can pass it here to continue the pagination.
    Before: beforeCursor, // To go back in pagination.
    First:  10, // Get first 20 results (works like LIMIT).
    Last:   nil, // Work same as first (LIMIT) but reverses the order of querying.
  }
}
// Same as before use the params in query options argument.
opts := ssql.SQLQueryOptions{
  Table:          "user",
  Fields:         dbFields,
  WithTotalCount: shouldReturnTotalCount,
  Params:         &params,
  Conditions:     conds,
}
rows, total, err := client.Find(ctx, &opts)
```

[↩](#table-of-contents)

#### Find: Sorting

You can have one or many sorting configs. The order matters.

```golang
params := ssql.Params{
  SortParams = []*ssql.SortParams{{
    Direction: "asc",
    Field:     "created_at",
  }}
}
// Same as before use the params in query options argument.
opts := ssql.SQLQueryOptions{
  Table:          "user",
  Fields:         dbFields,
  WithTotalCount: shouldReturnTotalCount,
  Params:         &params,
  Conditions:     conds,
}
rows, total, err := client.Find(ctx, &opts)
  
```

[↩](#table-of-contents)

#### Find: Conditions

You can set one or many conditions which will translate to WHERE clause in the final query.

```golang
...
userIDs := []string{...} // some list of user IDs
// Add some custom conditions.
conds := []*ssql.ConditionPair{
  {
    Field: "active",
    Value: true,
    Op:    "=", // equal operator
  },
	{
		Field: "user_id",
		Value: userIDs,
		Op:    "in", // Example of "IN" operator.
	},
}
// Same as before use the params in query options argument.
opts := ssql.SQLQueryOptions{
  Table:          "user",
  Fields:         dbFields,
  WithTotalCount: shouldReturnTotalCount,
  Params:         &params,
  Conditions:     conds,
}
rows, total, err := client.Find(ctx, &opts)

```

[↩](#table-of-contents)

#### Find: GraphQL

You can get a full PageInfo GraphQL Relay type which is:

```golang
type PageInfo struct {
	HasPreviousPage *bool   `json:"hasPreviousPage,omitempty"`
	HasNextPage     *bool   `json:"hasNextPage,omitempty"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
	TotalCount      *int    `json:"-"`
}
```

Simply use the helper function `PrepareGraphQLConnection` that returns `PageInfo` type above. Here's
a example of how to use the helper function on the result of a query.

```golang
// First let's get a list of rows.
rows, total, err := rp.Client.Find(ctx, &opts)
if err != nil {
  log.Logger(ctx).Error().Err(err).Msg("Cannot find users")
  panic(err)
}
var users []model.User
for rows.Next() {
	var user model.User
	if err := ssql.ScanRow(&user, rows); err != nil {
		log.Logger(ctx).Error().Err(err).Msg("Cannot scan users")
	}
	users = append(users, user)
}
// Build the pageInfo connection.
result, pageInfo := ssql.PrepareGraphQLConnection(users, query.Params)
// Set StartCursor and EndCursor.
if pageInfo.HasPreviousPage != nil && *pageInfo.HasPreviousPage && result != nil && len(*result) > 0 {
  startCursor := ssql.StrToCursor((*result)[0].GetPaginationCursor())
  pageInfo.StartCursor = &startCursor
}
if pageInfo.HasNextPage != nil && *pageInfo.HasNextPage && result != nil && len(*result) > 0 {
  startCursor := ssql.StrToCursor((*result)[len(*result)-1].GetPaginationCursor())
  pageInfo.EndCursor = &startCursor
}
if total != nil {
  pageInfo.TotalCount = total
}
```

[↩](#table-of-contents)

### Raw

If the provided API doesn't satisfy the usage you need, feel free to just run a raw custom query.

```golang
...
// Add some custom conditions.
values := []interface{}{true}
raw := "SELECT * FROM \"user\" WHERE active = $1"
res, err = client.Raw(ctx, raw, values)
if err != nil {
  log.Error().Err(err).Msg("Cannot execute raw query")
  panic(err)
}
...
```

[↩](#table-of-contents)

### Helpers working with structs

Much of the headache working with SQL in golang is to have a mapping between your struct fields and SQL columns.

Imagine we have type in Go that describes our user object.

```golang
type User struct {
  ID            *string `json:"id,omitempty" sql:"id"`
  Username      *string `json:"username,omitempty" sql:"username"`
  Email         *string `json:"email,omitempty" sql:"email"`
  EmailVerified *bool   `json:"emailVerified,omitempty" sql:"email_verified"`
  Active        *bool   `json:"active,omitempty" sql:"active"`
  UpdatedAt     *int    `json:"updatedAt,omitempty" sql:"updated_at"`
  CreatedAt     *int    `json:"createdAt,omitempty" sql:"created_at"`
  DeletedAt     *int    `json:"deletedAt,omitempty" sql:"deleted_at"`
}

user := User{
  ID: "...",
  Username: "...",
  ...
}
```

We want to `insert` the above user record to our PostgreSQL database.

```golang
// You can pass any struct as is.
_, err = client.Insert(ctx, "user", newUser)
if err != nil {
  panic(err)
}
```

This is because `ssql` internally uses `reflect` package in this case during runtime to get the tags.

However for find records (using `Find` method) with query options, you will need to pass the SQL fields themselves as a argument to the `find` 
method. This is because `Find` internally avoids using `reflect` and expects the plain text fields to be determined as query options.

To do this you can use `ExtractStructMappings(tags []string, s interface{})` helper function that simply returns type of:

```golang
type TagMappings map[string]map[string]string
```

In the following example the first-level map is the tag name requested, and the second-level is the either by tags or by field.

```golang
type User struct {
  ...
  CreatedAt     *int    `json:"createdAt,omitempty" sql:"created_at"`
  ...
}
// Request tags for mappings for `sql` and `json`.
var userMappingsByTags, userMappingsByFields = ssql.ExtractStructMappings([]string{"sql", "json"}, model.User{})
// Get field name by tag name we know.
userMappingsByTags["sql"]["created_at"] // will return `CreatedAt` (struct field name)
userMappingsByTags["json"]["createdAt"] // will return `CreatedAt` (struct field name)
// Get json/sql tag by field name
userMappingsByField["sql"]["CreatedAt"] // will return `created_at` (sql tag value)
userMappingsByField["json"]["CreatedAt"] // will return `createdAt` (json tag value)
```

Knowing this we can use these helper functions to run the expensive extraction of tags/fields only one time during startup
instead of for each find since it uses `reflect` internally.

To do this call `ExtractStructMappings` outside of your insert function/method in the same file for your type, then use it to
populate the fields array.


```golang

// Calling this one time only outside of our function.
var _, userMappingsByFields = ssql.ExtractStructMappings([]string{"rel", "sql"}, model.User{})

...

func MyInsertRowFunction(user *User) {
  dbUserFields := map[string]bool{}
  // Let's convert our struct type to a map of string that keys are the SQL column names.
  for fieldName := range user {
    dbUserFields[userMappingsByFields["sql"][fieldName]] = true
  }
  ...
  opts := ssql.SQLQueryOptions{
    Table:          "user",
    Fields:         dbFields, // Fields expect a map[string]bool which we now have.
    ...
  }
  // Executing the query.
  rows, total, err := client.Find(ctx, &opts)
  if err != nil {
    log.Error().Err(err).Msg("Cannot find users")
    panic(err)
  }
}
```

[↩](#table-of-contents)

### Helpers scanning query results

Use `ScanOne` if you are selecting/expecting one result. You can pass your struct pointer as is to scan it
without mapping individual fields like the standard `database/sql` library forces you to.

```golang
rows, err := client.FindOne(ctx, "user", "id", "7f8d1637-ca82-4b1b-91dc-0828c98ebb34")
if err != nil {
  log.Error().Err(err).Msg("Cannot select by ID")
  panic(err)
}
var resp User
if err := ssql.ScanOne(&resp, rows); err != nil {
  log.Error().Err(err).Msg("Cannot get resp by ID from Postgres")
  panic(err)
}
```

Use `ScanRow` inside the `rows.Next()` loop to populate your users array.

```golang
rows, total, err := client.Find(ctx, &opts)
if err != nil {
  log.Error().Err(err).Msg("Cannot find users")
  panic(err)
}
// Reading through the results.
var users []User
for rows.Next() {
  var user User
  if err := ssql.ScanRow(&user, rows); err != nil {
    log.Error().Err(err).Msg("Cannot scan users")
  }
  users = append(users, user)
}
```

[↩](#table-of-contents)
