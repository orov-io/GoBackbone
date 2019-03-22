# Web service golang backbone

A template to develop API's and Web services. Use it if you see it confortable.

## General points

This software relies on some other packages and makes the following assumptions:

- Use of postgres as your database. (If you use another DB, plase, modify main app database methods).

- To use postgres, this app requires a POSTGRES_CONNECTION string at the .env file. If it not exist, app will assume that no database is used.

- This app uses the [gin](https://github.com/gin-gonic/gin) http framework.

- This app uses the [jmoiron/sqlx](https://github.com/jmoiron/sqlx). However, the official sql package is used for migrations.

## About the app

### Configuring the app

You can find a .env file example. Please, use it if you add configuration values as keys or sensible info.
The app need 3 variables to run:

- __POSTGRES_CONNECTION__: The string to connect to database. Leave it empty if your app doesn't need a database conection. See the `example.env` file.
- __PORT__: application port. If you are running the app on an google app engine, this variable is not needed.
- __ENV__: The environment. Generally, one of [local, dev, prod]. You can add as environments as you want.

### Using firebase

If you want to use the firebase authentication system, you have two ways to use it:

- Put the firebase JSON credential in the App root, and name it 'firebase_credential.json'

- If you use gcloud, leave the 'firebase_credential.json' file on one gcloud storage bucket and provide the bucked name in the '.env' file. Use the __GCLOUD_STORAGE_BUCKET__ variable to do this.

### Serving statics

#### Assets

To serve static assets, you must ensure that a folder called `assets` exits at root level and is not empty.
App will autoload and serve this assets for you.

#### Templates

To load templates to use later with gin, you must ensure that a folder callet `templates` exist at root level
and is not empty. All templates expect have the `.gohtml` extension.
Templates will be available trought gin-gonic.

### Serving endpoints

You need to load your endpoints in `service/routes.go` file.
A _non secure_ ping route example is provided. To maintain code cleanliness, use a function for each logical
group of your API endpoints (see the addRoutes auxiliar function and ad here your routes groups).

### Migrations

By default, migrations are automatically executed when the app starts.
To add a migration:

```bash

$ cd migrations
goose create <migration-name> sql

```

### Dependencies

This projects uses [dep](https://github.com/golang/dep) to manage third party dependencies.

#### Init the project

If no Gopkg.toml is present:

```bash

dep ensure init

```

If project has Gopkg.toml:

```bash

dep ensure --update

```

#### Add a dependency

```bash

dep ensure -add <"repo1"> [<"repo2">...]

```

#### Checking status of dependencies

```bash

dep status

```

#### Updating dependencies

```bash

dep ensure --update

```

## Run locally

There are two ways

### Docker way

- [ ] Use docker compose to attach a postgres container.
- [ ] Force postgres container to use local persistent storage.
- [ ] Add capability to use the gcloud sql docker.
- [ ] Use two step dockerfile proccess.

### Traditional way

First, you need to load the .env file:

```bash

source .env && fresh go run main.go

```

If you are using [autoenv](https://github.com/kennethreitz/autoenv), you must put at the end of your .bashrc file:

```bash

source ~/.autoenv/activate.sh

```

Each time you change your env variables, you must change to the root directori (like ```$cd .```)

## GCloud

### Change to your project

```bash

gcloud config get-value project # actual project
gcloud projects list # list your projects
gcloud config set project <your-project> # change project

```

### Conecting to your database

```bash

gcloud sql connect <dbname> --user=postgres --quiet

```

### See app engine console logs

```bash

gcloud app logs tail -s <service>

```

## Deploys

This project is configured to be deployed with bitbucket pipelines. Steps can be shown on `bitbucket-pipelines.yml`.
You must provide the required variables to your bitbucket project.
Pairing between branches and environments are:
| Branch   |       gcloud project      |
| -------- | ------------------------- |
| Develop  | _dev environmnet_         |
| Master   | _production environment_  |

Commits pushed on one of this branches will trigger a deploy to the paired project.

- [ ] Provide documentation about pipelines variables.
- [ ] Provide circleCI (or another CI github compatible tool) integration.