FROM rust AS builder

COPY ./ ./

RUN cargo build

FROM alpine AS runner

COPY --FROM=builder ./target/release/api-userbase ./api-userbase

CMD ./api-userbase
