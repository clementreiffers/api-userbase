FROM golang:1.20-alpine as builder

RUN apk add make

WORKDIR /workspace

COPY ./ ./

RUN make build

FROM alpine as runner

WORKDIR /

COPY --from=builder /workspace/build/api-userbase ./api-userbase

CMD ./api-userbase