## Shakespearean Pokemon Service

### Requirements

To build from source:
* go1.14 ([install](https://golang.org/doc/install))

To run static code checks:
* golangci-lint ([install](https://github.com/golangci/golangci-lint#install))

To start service using docker:
* docker ([install](https://docs.docker.com/install/))

### Usage

The easiest way to run the service (and avoid having to install go) is by
loading the provided `image.tar` from disk and running it using the following commands:
```
make load run
```

The service will start serving on port 3000 and you should see the following log line in your
current terminal session.
```
2020/mm/dd hh:mm:ss starting server on: 0.0.0.0:3000
```

To check if the service works:
```
make example
```

You should see the following output on your terminal:
```
{"name":"charizard","description":"Charizard flies 'round the sky in search of powerful opponents. 't breathes fire of such most wondrous heat yond 't melts aught. However,  't nev'r turns its fiery breath on any opponent weaker than itself."}
```

To test, build and run the service, use the following commands:
* `make lint`: run linter
* `make test`: run unit tests
* `make build`: build go binary and docker image
* `make run`: run docker image as container, requires the image to built (run `make build`)
* `make save`: save docker image to .tar file
* `make load`: load docker image from .tar file
* `make example`: run an example command to see if the service works (fetches description for Charizard)

### Notes
* I've added the saved docker image to ease testing the service. Normally, I would suggest uploading it to a shared container registry as part of some automated build process, without checking in any build artifacts to the repo.
* There are no structured error messages returned. (However, I've attempted to returned the most helpful http error codes to the user where possible. Errors should also appear in the logs.)
