# go-blog-aggregator

## Prerequisites

To run this program, you will need to have the following installed:

- [PostgreSQL](https://www.postgresql.org/download/)
- [Go](https://golang.org/dl/)

## Installation

First, install the `gator` CLI using `go install`:

```sh
go install github.com/thisantm/go-blog-aggregator@latest
```

## Configuration

Set up the configuration file `.gatorconfig` in your home directory (`$HOME`). Here is an example configuration:

```json
{
  "db_url": "<postgres_connection_string>",
  "current_user_name": ""
}
```

## Database Migrations

Clone the repository:
```sh
git clone https://github.com/thisantm/go-blog-aggregator.git
cd gator
```

To handle database migrations, you will need to install `goose`:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Run the migrations using the following command:

```sh
goose -dir sql/schema postgres <postgres_connection_string> up
```
## Commands

Here are a few commands you can run with `go-blog-aggregator`:

- `go-blog-aggregator register <user>`: Registers a new user.
- `go-blog-aggregator login <user>`: Logs in a user.
- `go-blog-aggregator addfeed <url>`: Adds a new feed.
- `go-blog-aggregator follow <url>`: Follows a feed.
- `go-blog-aggregator browse <limit>`: Browses available feeds.
