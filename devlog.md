Creating this so that I will have an easier time developing Go stuff in the future

## Getting package dependencies
- `go mod install` or `go mod tidy` if the Go repository already has a go.mod or go.sum file.
- Otherwise, do `go get -u ...`, getting it through the URL of the required package (typically on GitHub).

## Setting up packages

Assuming you are hosting your code on github:
- Upload some skeleton code onto github
- Check the link of your repository
- Within your root folder, run the command `go mod init github.com/josh1248/nusc-queue-bot`, as an example. This allows people to `go get -u` under the same URL as your package name, which is very useful and is the convention.
- Afterwards, install the required dependencies for your code using `go mod tidy`.
  
## Sub-packages

If you are a sane programmer, you would have split your code into multiple sub-folders, such as under `internal/xx` and `internal/yy`, whereby you would need to export code from one sub-folder to the other.

Unlike other languages like Python or JavaScript where relative imports are useable, Golang supports package importing and exporting via 2 ways: GOPATH and remote packages. Given that you are most likely going to host your code on GitHub, the remote package method would be more user-friendly.

## Capitalization

https://go.dev/doc/code#ImportingRemote