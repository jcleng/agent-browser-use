# agent-browser-use in a container: the Go server plus headless Chromium.
#   make docker      build the image and run it on http://localhost:3000
# Build context is the repo root: docker build -f docker/Dockerfile .

# ---- Build ----
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY app/go.mod app/go.sum ./
RUN go mod download
COPY app/ .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/agent-browser-use .

# ---- Run ----
FROM debian:bookworm-slim
# Chromium from Debian works on both amd64 and arm64 (Apple Silicon).
# Fonts so pages in Vietnamese, CJK and emoji render properly.
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
      chromium ca-certificates tzdata \
      fonts-liberation fonts-dejavu-core fonts-noto-color-emoji fonts-noto-cjk \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/agent-browser-use /usr/local/bin/agent-browser-use
COPY docker/chromium /usr/local/bin/chromium-docker
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/chromium-docker /usr/local/bin/entrypoint.sh \
 && useradd --create-home --uid 1000 app \
 && mkdir -p /data && chown app:app /data

# Unpacked extensions live on the volume so they survive restarts. With
# AGENT_BROWSER_USE_EXTENSIONS=1 and no AGENT_BROWSER_USE_EXTENSIONS_DIR, the
# default folders are <profile dir>/../extensions and /data/extensions.
ENV AGENT_BROWSER_USE_EXTENSIONS_DIR=/data/extensions
RUN mkdir -p /data/extensions && chown app:app /data/extensions

USER app
# Chrome profile (cookies, logins) and chat history live here: mount a volume to keep them.
# Unpacked browser extensions go in /data/extensions; enable them with
# AGENT_BROWSER_USE_EXTENSIONS=1 (off by default).
VOLUME /data
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
