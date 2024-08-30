To get the packages:
- go install github.com/pressly/goose/v3/cmd/goose@latest
- go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

Ensure both your global go bin, and user-level go binaries are exposed on path. Check /home/{user}/go/bin is exposed.

Once you have the packages move to the next steps.

Runs migrations. Basically interacts with the DB
- goose postgres postgres://postgres:12345678@localhost:5432/rssagg up

Generates GO code from SQL queries
Ensures that generated code is type-sage
- sqlc generate