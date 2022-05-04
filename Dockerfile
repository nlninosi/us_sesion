FROM golang:1.18

WORKDIR /usr/src/us_sesion_ms

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/us_sesion_ms
EXPOSE 8081

CMD ["us_sesion_ms"]
