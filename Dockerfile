
FROM node:22-alpine3.19 AS builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable


WORKDIR /build

COPY --link ui/package.json ui/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY --link ui .

RUN pnpm build && pnpm prune

FROM golang:1.23.2-alpine3.20 AS go-builder

ENV CGO_ENABLED=1

WORKDIR /build

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN rm -rf ui

COPY ./ui/embed.go ./ui/embed.go

COPY --from=builder /build/build ./ui/build

RUN go build -o ./app ./cmd/chmoly

FROM golang:1.23.2-alpine3.20 AS final

WORKDIR /app
COPY --from=go-builder /build/app /app/app

EXPOSE 80

CMD ["/bin/sh", "-c", "./app serve"]