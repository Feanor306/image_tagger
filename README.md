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

# Build templates and project
make build

# Start up development database
make db-up

# shut down development database
make db-down

# run the app
make run
```