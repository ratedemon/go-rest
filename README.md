# go-rest

go-rest is application, that allows to create users and watch at them accounts.
Application is served by

## Installation

Previously, you need to have [docker](https://www.docker.com/get-started) and [docker-compose](https://docs.docker.com/compose/install/) on your machine.
And, that's it.

## Run

you can run:
```
docker-compose up -d --build
```

`go-rest` will be running in detached mode.


## Test

***You must have installed Golang on your machine before running these tests***

To run golang tests, just go to [/server](/server) directory and run there `make test`.