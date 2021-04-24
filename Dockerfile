FROM golang
RUN mkdir -p /club
WORKDIR /club
COPY . .
RUN go mod download
RUN go build -o club
ENTRYPOINT ["./club"]