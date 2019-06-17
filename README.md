# Backbone Service | My Suite

> A template to develop API's and Web services. Use it if you see it confortable.

## Api documentation

 You can find functional, real state API documentation (with swagger standard) on the gcloud developer portal. For each environment, you can find this documentation visiting the [Developer Portal](https://endpointsportal.${SERVICE_NAME}.cloud.goog/docs)

## About the service

### Configuring the app

You can find a .env file example. Please, use it if you add configuration values as keys or sensible info.
The app need some variables to run:

- __POSTGRES_HOST__
- __POSTGRES_PASSWORD__
- __POSTGRES_USER__
- __POSTGRES_SSL_MODE__
- __SERVICE_DATABASE_NAME__
- __PORT__: application port. If you are running the app on an google app engine, this variable is not needed. By default, in local environment, this service uses 8081
- __ENV__: The environment. Generally, one of [local, dev, prod]. You can add as environments as you want.
- __GOOGLE_IDENTITY_API_KEY__: Firebase identity key, needed to obtain custom tokens
- __ref__: needed to compile the openapi-appengine.example file.
- __SONAR_TOKEN__: Key to send sonar results to sonarcloud
- __PROJECT_PATH__: Path to project
- __SERVICE_NAME__: The name of the service

There are some optional flags:

- __PUBLIC_API__: If true, enables CORS to make cross domain requests.
- __GCLOUD_STORAGE_BUCKET__: Allow app to connect to gcloud storage.

### Auth credentials

This service is build on top of firebase functionalities. You must provide a firebase credential in order to run the service:

- Put the firebase JSON credential in the App root, and name it 'firebase_credential.json'

- If you use gcloud, leave the 'firebase_credential.json' file on one gcloud storage bucket and provide the bucked name in the '.env' file. Use the __GCLOUD_STORAGE_BUCKET__ variable to do this.

### Dependencies

To run locally this service, you should have installed in your system:

- [goconvey](https://github.com/smartystreets/goconvey)
- [fresh](https://github.com/gravityblast/fresh)
- [docker](https://docs.docker.com/install/overview/) & [docker-compose](https://docs.docker.com/compose/)
- [Go language](https://golang.org/)
- [make](https://es.wikipedia.org/wiki/Make)
- [dep](https://github.com/golang/dep)

## Running locally

A make file for local development is provided. Below commands are provided:

- *default* commad simply builds the Dockerfile for the service.
- *up* command starts up the API and runs it in the background.
- *logs* tails the logs on the Docker container.
- *down* shuts down the server.
- *test* runs any tests in the current directory tree, using goconvey.
- *clean* shuts down the API and then clears out saved Docker images from your computer. This can be useful when running another image like Postgres which writes information to your local machine and doesn’t clean it up when the container shuts down.

## General points

This software relies on some other packages and makes the following assumptions:

- Use of postgres as your database. (If you use another DB, plase, modify main app database methods).

- To use postgres, this app requires a POSTGRES_CONNECTION string at the .env file. If it not exist, app will assume that no database is used.

- This app uses the [gin](https://github.com/gin-gonic/gin) http framework.

- This app uses the [jmoiron/sqlx](https://github.com/jmoiron/sqlx). However, the official sql package is used for migrations.

### Migrations

By default, migrations are automatically executed when the app starts.
To add a migration:

```bash

$ cd migrations
goose create <migration-name> sql

```

### Loading env variables

First, you need to load the .env file:

```bash

source .env && fresh go run main.go

```

If you are using [autoenv](https://github.com/kennethreitz/autoenv), you must put at the end of your .bashrc file:

```bash

source ~/.autoenv/activate.sh

```

Each time you change your env variables, you must change to the root directori (like ```$cd .```)

### Testing

Automatic testing can be made outside the docker container. You only need to do 2 things:

- Assert all vendor packages are install: 

```bash
make update
```

- Make sure that host name "db" is in your /etc/host file and is actually pointing to localhost (127.0.0.1)

Now, you can run convey througt:

```bash
make test
```

## tricks

- Rebuild and run service and see logs with one command:

```bash
make restart logs
```
