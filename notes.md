Runs migrations. Basically interacts with the DB
- goose postgres postgres://postgres:12345678@localhost:5432/rssagg up

Generates GO code from SQL queries
Ensures that generated code is type-sage
- sqlc generate