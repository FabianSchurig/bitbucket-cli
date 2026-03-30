FROM dhi.io/golang:1

RUN go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest

ENTRYPOINT ["bb-cli"]
