FROM golang:1.25.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o application cmd/api/main.go 

CMD [ "./application" ]