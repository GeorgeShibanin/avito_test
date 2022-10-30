FROM golang:latest


WORKDIR /go/src/avito_test
COPY . /go/src/avito_test

#RUN go mod tidy
RUN go build -o ./bin/avito_test ./cmd/avito_test/
#RUN go build -o app
# Для возможности запуска скрипта
RUN chmod +x /go/src/avito_test/scripts/*


CMD ["/go/src/avito_test/bin/avito_test"]