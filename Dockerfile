FROM golang:1.17-alpine as builder
RUN apk update && apk add --no-cache --update git build-base
RUN addgroup -S scratchgroup && adduser -S scratchuser -G scratchgroup

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGOENABLED=0 go build -a -tags netgo -o RussianWarshipGoFuckYourself ./cmd/RusskijKorablIdiNaxuj-cli && mkdir build_app && cp RussianWarshipGoFuckYourself build_app/ && cp targets/* build_app/



FROM scratch

COPY --from=builder app/build_app /app
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser

WORKDIR /app
ENTRYPOINT [ "/app/RussianWarshipGoFuckYourself", "-i=all.json", "-N=5000", "-s" ]
