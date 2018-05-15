GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME= foodCrawler
BINARY_UNIX=$(BINARY_NAME)_unix
DIR = src/
MAIN = main

all: clean build run

clean: # remove data and binary folder
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
				rm -f $(BINARY_UNIX)
build: 
				$(GOBUILD)  -o $(BINARY_NAME) -v $(DIR)*.go

run: 
				$(GOBUILD) -o $(BINARY_NAME) -v  $(DIR)*.go
				./$(BINARY_NAME)