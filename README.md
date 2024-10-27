# Image tagger 
A simple app that allows for upload of media (images) and setting tags for them

# Requirements
- go 1.23.2
- docker

# Instructions
```bash
# setup .env variables (make changes if needed)
cp .env.stub .env

# You may need to manually create the data directory in the root of the project 
# or change ownership of ./data to current user after starting up DB
mkdir data

# Start up development database
make db-up

# shut down development database
make db-down

# Build templates and project
make build

# run the app
make run

# run tests
make test

# view docs - run make doc and visit local doc url
make doc
http://localhost:6060
```