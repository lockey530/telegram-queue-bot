Creating this so that I will have an easier time developing Go stuff in the future

- [General Go Stuff](#general-go-stuff)
  - [Core commands](#core-commands)
  - [Getting dependencies](#getting-dependencies)
  - [Setting up your own packages](#setting-up-your-own-packages)
  - [Sub-packages](#sub-packages)
  - [Capitalization](#capitalization)
  - [Postgres / sqlx](#postgres--sqlx)
  - [go build](#go-build)
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
    - [3: Postgres connection issues](#3-postgres-connection-issues)
    - [4: Nvm, ditch Docker Compose entirely](#4-nvm-ditch-docker-compose-entirely)
    - [5: First light (Hacky)](#5-first-light-hacky)
    - [6: Fixes](#6-fixes)
    - [X: Secrets and env variables](#x-secrets-and-env-variables)


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

## Postgres / sqlx

After some research, I decided on using the Postgres driver, `lib/pq`, along with `jmoiron/sqlx` to simplify the marshalling / unmarshalling of data between Go structs and SQL entries. Thankfully, these libraries are well documented, so I managed to get a hold of waht they do by reading their documentation and consulting tutorials as needed.

Be careful and remember capitalization! The sql driver will need the fields in your struct to be public (i.e. capitalized) before it can marshall into them. I spent a thankful short 10 mins on this.

Additionally, as I used a SERIAL PRIMARY KEY for auto-incrementing my queue, I also had to reset my serialization within my WIPE_DATA command with `TRUNCATE TABLE queue RESTART IDENTITY` - solution identified in https://stackoverflow.com/questions/3819292/reset-postgresql-primary-key-to-1. Alternatively, a less destructive method is noted, but requires you to name your serial field as `id` or something with no underscores - `ALTER SEQUENCE queue_id_seq RESTART WITH 1`.

To check how many rows were deleted, I used the Postgres driver to get the number of rows, which helped to check if the deletion was correct.
```Go
	result, err := db.Exec(`DELETE FROM queue WHERE user_handle = $1;`, userHandle)
	if err != nil {
		return fmt.Errorf("failed to leave queue. %v", err)
	}
	log.Println(result.RowsAffected())
```
A more generalizable solution is found in https://stackoverflow.com/questions/2251567/how-to-get-the-number-of-deleted-rows-in-postgresql

## go build

`go build -o <binary-name> ./...` into a binary at the root folder will not work as intended in this repository as it has no Go files within the root repository. Instead, you would need to throw this folder into a separate repository, for example, the folder `server`. The binary would then be save-able as `go build -o server ./...`. Of course, you should stick with the simple `go mod cmd/server/main.go` for development.

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

When containerzing, the environment changes (duh) and the environment variables are no longer as easily read. Looking through the proposed solutions in https://stackoverflow.com/questions/46917831/how-to-load-several-environment-variables-without-cluttering-the-dockerfile, I decided to go with copying the .env file directly into the image files. 
My go files read off a `.env` file, which I had to copy into the container with `COPY .env .env` within my Go Dockerfile. Additionally, I had to comment out `**/.env` within `.dockerignore`. 


:warning: Note: This is a small-scale project - using secrets stored in .env files is not the recommended way to keep your secrets since `.env` is in plaintext. If you happen to use this code for more high-value applications, you should definitely use a secrets manager. (but a lot of repos still use the .env file :eyes:)

### 3: Postgres connection issues

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

I had to scale back and throw away `compose.yaml` entirely. Thankfully, the solution in [#2](#2-env-files) still works. This solution also helped to avoid some of the pitfalls associated with args . (If you are advanced, you can consider the --secret option in https://stackoverflow.com/questions/45405212/safe-way-to-use-build-time-argument-in-docker/51921954#51921954 https://github.com/docker/cli/pull/1288, which I did not pursue.)

### 5: First light (Hacky)

The idea of Railway all this time was to visualize your images individually in a nice graphical format. We just had to set this up, with the Go repo connecting to Github for continuous deployment while the Postgres server was automatically booted up via the `New` button at the top right of the visualization.

![alt text](images/railway_progress1.png)

Next up, I removed the hard-coded `postgres` spam glory, reverting back to the use of environment variables which were used to connect to the Postgres database with the following:
```Go
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)
```

I faced the following 2 issues when my dockerfiles were run by Railway. I got rid of both of these issues within the Dockerfile with bandaid fix (to be fixed).

`Cache mounts MUST be in the format --mount=type=cache,id=<cache-id>` - comment out all cache-based commands in Dockerfile.

`ERROR: failed to solve: failed to compute cache key: failed to calculate checksum of ref (some SHA256 key): "/.env": not found` - Remove .env copying. Instead, test un-containerized code using a simple `go run cmd/server/main.go`, while dumping environment variables directly into Railway with the RAW editor to test the containerized version of the code in remote deployment.

Afterwards, I utilized the fields `PGUSER`, `PGPASSWORD`, `PGHOST`, `PGPORT`, and `PGDATABASE` that Railway provided within the service variables in Railway, linking it up into my Golang code using reference variables (see https://docs.railway.app/guides/variables#reference-variables). For now, I also additionally fed it into my local .env files for ease of use (I will clear up the remote/local deployment process later on.)

![alt text](images/railway_progress2.png)

![alt text](images/railway_progress3.png)

### 6: Fixes


### X: Secrets and env variables



Then, further configure the compose.yaml and Dockerfiles based on documentation if required: https://docs.docker.com/reference/dockerfile/ https://docs.docker.com/compose/compose-file/?uuid=14e6d05b-8c4a-4389-b002-3e27079fd972%0A
Lastly, store persistent data and stuff with volumes: https://docs.docker.com/storage/volumes/