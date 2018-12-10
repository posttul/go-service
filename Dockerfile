FROM golang as build
WORKDIR /go/src/app
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure 
RUN go install -v ./...

CMD ["app"]
