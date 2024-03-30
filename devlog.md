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