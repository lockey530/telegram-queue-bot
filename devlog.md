Creating this so that I will have an easier time developing Go stuff in the future

- [General Go Stuff](#general-go-stuff)
  - [Core commands](#core-commands)
  - [Getting dependencies](#getting-dependencies)
  - [Setting up your own packages](#setting-up-your-own-packages)
  - [Sub-packages](#sub-packages)
  - [Capitalization](#capitalization)
- [Telegram-specific](#telegram-specific)
  - [How to set up communications to your bot](#how-to-set-up-communications-to-your-bot)
  - [How to Hello World - listen for updates](#how-to-hello-world---listen-for-updates)
  - [Setting up Telegram menu selections](#setting-up-telegram-menu-selections)
  - [How to](#how-to)
- [Postgres](#postgres)
  - [How to set-up](#how-to-set-up)
- [Remote deployment - Docker + Railway](#remote-deployment---docker--railway)
  - [Docker Install and intro](#docker-install-and-intro)
  - [Docker setup issues](#docker-setup-issues)
    - [1: Volumes](#1-volumes)
    - [2: env files](#2-env-files)
    - [3: Postgres connection within Docker](#3-postgres-connection-within-docker)
    - [4: Nvm, ditch Docker Compose entirely](#4-nvm-ditch-docker-compose-entirely)
    - [5: Secrets and env variables](#5-secrets-and-env-variables)


# General Go Stuff

## Core commands
- `go mod init XYZ` - run in your root folder of your Go project to set up dependencies and package names (more on this [below](#setting-up-your-own-packages))

## Getting dependencies
- `go mod install` or `go mod tidy` if the Go repository already has a go.mod or go.sum file.
- Otherwise, do `go get -u ...`, getting it through the URL of the required package (typically on GitHub).
- If you are scared that the remote repository may disappear, you can consider using `go mod vendor` to essentially get some local copies of the dependencies you need.

## Setting up your own packages

Assuming you are hosting your code on github:
- Upload some skeleton code onto github
- Check the link of your repository
- Within your root folder, run the command `go mod init github.com/josh1248/nusc-queue-bot`, as an example. This allows people to `go get -u` under the same URL as your package name, which is the convention - this is because it ensures that your packages are uniquely named, and helps others see your public repo code easily. This is not strictly needed - you can name your module funny things like `go mod init hehehe` - however, your code would be easily published as Go tools use the name of your package to download the packages themselves, and may clash with other package names. See: https://go.dev/doc/tutorial/create-module#start 
- Afterwards, install the required dependencies for your code using `go mod tidy`.
  
## Sub-packages

If you are a sane programmer, you would have split your code into multiple sub-folders, such as under `internal/xx` and `internal/yy`, whereby you would need to export code from one sub-folder to the other.
Within the go files in these folders, it is convention to write `import xx` or `import yy` corresponding to the folder name as convention. Again, this is not needed, but is good practice to track your sup-packages easily.

Unlike other languages like Python or JavaScript where relative imports are useable, Golang discourages relative imports, and instead has their packages use an absolute path. This is done in 2 ways: GOPATH (older way) or `go mod init <module name>` (newer, better way). We shall focus on the 2nd way.

- Suppose that you have previously run `go mod init XYZ`. Your package's "root path" is henceforth `XYZ`.
- You have the folders `internal/xx` and `internal/yy`, and .go files within these folders have `package xx` and `package yy` respectively.
- You can then use these folders in any other folder using `XYZ/internal/xx` or `XYZ/internal/yy`.

## Capitalization

https://go.dev/tour/basics/3 
- Golang has a peculiar way of marking things as shareable to external parties beyond your package - if your function or type stats with a lower letter case, it is private. If it starts with an upper case, it is public.
- This is the reason why all functions and types that you use from external packages all start with an upper case.

# Telegram-specific
Uses github.com/go-telegram-bot-api/telegram-bot-api - module name of tgbotapi typically given.

## How to set up communications to your bot

[Link to my README file](./README.md#register-your-bot-on-telegram)

## How to Hello World - listen for updates

## Setting up Telegram menu selections

## How to 

# Postgres

## How to set-up
Find any appropriate download method in https://www.postgresql.org/download/ to get Postgres up and running. I used the Postgres.app method as explained in https://www.youtube.com/watch?v=wTqosS71Dc4 for ease of setup.

My version of the app used is the following:

![alt text](images/postgres_setup.png)

Then, configure $PATH so that u can use psql, a CLI to interact with Postgres:

```
sudo mkdir -p /etc/paths.d &&
echo /Applications/Postgres.app/Contents/Versions/latest/bin | sudo tee /etc/paths.d/postgresapp
```

After you are done, you should be able to connect to the database:

![alt text](images/postgres_setupdone.png)

`psql postgres://<system username>@localhost:5432/<system username>`
system username refers to the username you use to log in to your actual file system.

# Remote deployment - Docker + Railway

I decided to use Railway due to its free $5 credits monthly, which is perfect for a small scale, hobby project like this.

I could have used the Railway CLI tool (https://docs.railway.app/guides/cli), but I might as well start learning Docker since it is more applicable and prevents lock-in into the Railway deployment service in case the free credits disappear - plus it helps in making the application runnable in multiple places.

## Docker Install and intro
I got docker here https://docs.docker.com/get-docker/, utilizing Docker Desktop within Visual Studio Code.

As a dummy, I accepted the recommended default settings for the desktop app.

I methodically went through the tutorials(within the desktop app) in https://docs.docker.com/guides/get-started/ to get an idea of what I am trying to do.
   1. Setting up a single-container app that is built based on the Dockerfile: `docker build -t welcome-to-docker .`
   2. Example of multi-container app that uses compose.yaml as well as the Dockerfile: `docker compose up`
   3. To stop, do `docker compose down --volume`
   4. Monitor changes on save for container: `docker compose watch`

## Docker setup issues
I used `docker init` to use the defaults provided by Docker (these are the files you see as Dockerfile and compose.yaml within this directory.)

I used the Docker Compose example repository that they have provided, which is extremely useful: https://github.com/docker/awesome-compose/blob/master/nginx-golang-postgres/backend/main.go

### 1: Volumes
However, I kept facing this extremely annoying issue:

```
PostgreSQL Database directory appears to contain a database; Skipping initialization
...
FATAL: role 'postgres' does not exist
```

The solution was to clear out my volumes (which is Docker's term for persistent data tied to local files) using the `docker compose down --volumes` to clear off my persistent data. I will probably need to investigate a better solution that can avoid the need to wipe the volumes every time. (Interestingly, after a few runs of this, `docker compose down` now allows my container to work fine without the error.)

### 2: env files

My go files read off a `.env` file, which I had to copy into the container with `COPY .env .env` within my Go Dockerfile. Additionally, I had to comment out `

:warning: Note: This is a small-scale project - using secrets stored in .env files is not the recommended way to keep your secrets since `.env` is in plaintext. If you happen to use this code for more high-value applications, you should definitely use a secrets manager. (but a lot of repos still use the .env file :eyes:)

### 3: Postgres connection within Docker

First of all, always make sure that your db container is working first! When testing changes locally, make sure things are running in port 5432 (or the Postgres port of your choice).

Another issue I faced was the inability to connect to the Postgres DB within Go. The issue I got was `127.0.0.1:5432: connect: connection refused`. This issue was echoed in https://stackoverflow.com/questions/57696882/docker-golang-lib-pq-dial-tcp-127-0-0-15432-connect-connection-refused 

The solution was to switch from using the following connection string:
```
	db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("user=%s dbname=%s port=%s sslmode=disable",
			user, dbname, port))
```

to the following instead, e.g. `"postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"`, as stated in the documentation in https://pkg.go.dev/github.com/lib/pq#section-readme:

This solution was mentioned by user nsandalov in the github issue https://github.com/quay/clair/issues/134: 

>Faced the same problem. Figured out that connection URL to the database should not be localhost or 127.0.0.1. It should be URL to your container with the Postgres. So the easiest way to do it is to define something like container_name: clair_postgres in the docker-compose.yml and use it as a connection string postgresql://postgres:password@clair_postgres:5432?sslmode=disable

The main issue is that the 1st connection method used the default `localhost` into its connection URL. In actuality, you need to connect to the name of the container itself. Docker auto-generates this container name if it is not specified. Therefore, you will need to specify either `--name=XXXX` if running a single container, or `container_name:XXX` under the db services in the compose yaml.

Additionally, make sure that your Postgres instance at the port is running if you are testing this change! Faced a few minutes accidentally stuck in this.

The solution is an unholy mess of `postgres`es spammed in compose.yaml (including within the passwords file):

```Dockerfile
 db:
    image: postgres
    restart: always
    user: postgres
    container_name: postgres
    secrets:
      - db-password
    # volumes:
    #   - db-data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=postgres
      # - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    expose:
      - 5432
    healthcheck:
      test: [ "CMD", "pg_isready" ]
      interval: 10s
      timeout: 5s
      retries: 5
# volumes:
#   db-data:

secrets:
  db-password:
    file: db/password.txt
```
### 4: Nvm, ditch Docker Compose entirely

Except - most cloud services for dynamic things do not support docker compose! There have been plans to consider Docker Compose in Railway, but it has not been implemented. I tried looking for free options around, but things like Google Cloud Run do not support Docker compose as well :( Looks like I spent a lot of time finding an un-useable solution!

I had to scale back and throw away `compose.yaml` entirely. Thankfully, the solution in [#2](#2-env-files) still works. This solution also helped to avoid some of the pitfalls

### 5: Secrets and env variables



Then, further configure the compose.yaml and Dockerfiles based on documentation if required: https://docs.docker.com/reference/dockerfile/ https://docs.docker.com/compose/compose-file/?uuid=14e6d05b-8c4a-4389-b002-3e27079fd972%0A
Lastly, store persistent data and stuff with volumes: https://docs.docker.com/storage/volumes/