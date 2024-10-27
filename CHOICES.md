# Design and Technology Choices
We use golang because it has great tooling, strong typing, superb concurrency tools, is very powerful and is a pleasure to write

## Web framework
(Echo)[github.com/labstack/echo] because I wanted to learn how to use it

## Database
Postgres because I have used it for many of my projects  
(pgx)[github.com/jackc/pgx/v5] is the standard postgres driver  
(squirrel)[github.com/Masterminds/squirrel] because it is a simple sql builder that fits small project requirements  
In retrospect, using a more powerful ORM like GORM would have been easier to manage relations  

## Env
(godotenv)[github.com/joho/godotenv] to read ENV variables from .env file and keep our secrets from prying eyes

## Templating engine
(Templ)[github.com/a-h/templ] because I wanted to learn how to use it

## Testing
(testify)[github.com/stretchr/testify] because assert and require simplify writing tests

## Other
- (uuid)[github.com/google/uuid] for id-s to ensure uniqueness  
- **Makefile** because it allows for 'aliasing' or batching of all the great docker and go tools we have, so we don't have to remember hundreds of individual commands  
- **docker-compose** for db and test-db to make db setup easy and repeatable, keep files in their own dir  
- keeping DB files and other assets in ./data allows for easy backup
- binaries in ./bin in case we need to build several
- entrypoint in ./cmd to make scripting easy and allow us to decouple some dependencies with smart interface use
- most of source in ./src to keep project root clean


# Possible future improvements

1. Pagination for tags and media
2. Safeguard against media.tags that don't exist
3. Full dockerization of app as well as db
4. gorm instead of squirrel for DB access
5. full CRUD functionality for both tag and media

