# ── Stage 1: Build ──────────────────────────────────────────
FROM golang:1.25-alpine AS builder
# golang:alpine = Go compiler + Alpine Linux (small ~300MB vs ~800MB for debian)
# AS builder = gives this stage a name so stage 2 can reference it

RUN apk add --no-cache gcc musl-dev
# apk = Alpine's package manager (like apt on Ubuntu, brew on macOS)
# gcc = C compiler (required by CGO / go-sqlite3)
# musl-dev = C standard library headers (Alpine uses musl instead of glibc)
# --no-cache = don't store the package index (keeps the layer smaller)

WORKDIR /app
# sets the working directory inside the container
# all following commands run from /app

COPY go.mod go.sum ./
# copy ONLY dependency files first
# this layer gets cached — if go.mod doesn't change, Docker skips the next step

RUN go mod download
# download all dependencies into the module cache
# cached until go.mod or go.sum changes

COPY . .
# now copy the rest of the source code
# this layer changes often, but the download layer above stays cached

RUN go build -o server .
# compile everything in the current directory
# -o server = name the output binary "server"
# CGO is enabled by default, which is what go-sqlite3 needs

# ── Stage 2: Run ────────────────────────────────────────────
FROM alpine:latest
# start fresh — no Go compiler, no source code, nothing from stage 1
# except what we explicitly copy below

WORKDIR /app

COPY --from=builder /app/server .
# --from=builder = grab the file from the stage named "builder"
# copies only the compiled binary — final image stays tiny (~20MB)

EXPOSE 3000
# documents that the container listens on port 3000
# doesn't actually open the port — that's done in docker-compose

CMD ["./server"]
# the command that runs when the container starts
