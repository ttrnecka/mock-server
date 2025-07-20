FROM golang:1.23-alpine AS build

#COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN update-ca-certificates

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /mock_server

# final stage
FROM alpine 
WORKDIR /
COPY --from=build /mock_server /mock_server
COPY --from=build /app/modules /modules
EXPOSE 2222
ENTRYPOINT ["/mock_server"]