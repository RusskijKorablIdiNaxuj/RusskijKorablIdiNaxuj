FROM golang:1.17-alpine as builder
RUN apk update && apk add --no-cache --update git build-base
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o RussianWarshipGoFuckYourself ./cmd/RusskijKorablIdiNaxuj-cli && mkdir build_app && cp RussianWarshipGoFuckYourself build_app/ && cp targets/targets.txt build_app/


FROM alpine

COPY --from=builder app/build_app /app

WORKDIR /app
ENTRYPOINT [ "./RussianWarshipGoFuckYourself", "-i=targets.txt", "-N=100", "-s" ]
