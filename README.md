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

If your DML is not a simple query that is not supported, just use the `Raw` method. 
We'll keep it intentionally simple:

- No JOINS
- No ORM features
- No complex WHERE conditions
- Honesty most non-trival query patterns are not added 

If you think you need more patterns/utilities/methods/helpers, and it's actually useful that is hard to do
without a wrapper, feel free to open a PR.

## Getting Started

```
go get "github.com/sourcesoft/ssql"
```

[API Reference](https://pkg.go.dev/github.com/sourcesoft/ssql)

## Simple usage

TODO
