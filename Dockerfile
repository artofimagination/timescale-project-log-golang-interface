FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/timescaledb-project-log-go-interface

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++
RUN go mod tidy
RUN cd $GOPATH/src/timescaledb-project-log-go-interface/ && go build main.go

# This container exposes port 8186 to the outside world
EXPOSE 8186

RUN chmod 0766 $GOPATH/src/timescaledb-project-log-go-interface/scripts/init.sh

# Run the executable
CMD ["./scripts/init.sh"]