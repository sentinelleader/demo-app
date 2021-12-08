FROM golang:1.17-alpine AS base
COPY . /src/
WORKDIR /src/
RUN go build -o webapp

FROM alpine
WORKDIR /app
COPY --from=base /src/webapp /app/
CMD /app/webapp
