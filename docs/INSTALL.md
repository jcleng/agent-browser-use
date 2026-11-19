# Install and configuration

## Docker

```sh
git clone https://github.com/cavaldos/agent-browser-use.git
cd agent-browser-use
make docker          # start
make docker-logs     # logs
make docker-stop     # stop
```

## From source (Go 1.22+ and Chrome/Chromium)

```sh
make run             # starts the server on port 3000
make build           # builds bin/agent-browser-use
make install         # copies it to ~/.local/bin (BINDIR=~/bin to change)
```

If `~/.local/bin` is not on your `PATH`: `echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc && exec zsh`

Pin a release version: `curl -fsSL .../scripts/install.sh | sh -s -- v1.2.3`

## Options

```sh
./run -p 8080 -pw mysecret            # run from source with any flags (make can't take -p itself)
make run p=8080                       # another port (then use http://localhost:8080/mcp)
make run pw=mysecret                  # require a password for the web UI and /mcp
make run engine=brave                 # pick the browser: chromium, chrome, edge, brave, vivaldi, headless-shell
make run -- -p 8080 -pw mysecret      # or pass any app flags after --
make run -- -agent-speed 1            # slow the agent's cursor and typing so people can follow
make run -- -mcp-remote               # allow other machines on your LAN to call /mcp
make run -- -no-mcp                   # disable /mcp, use only as a shared browser room
make run -- -tunnel                   # public https URL for friends outside your network (no router setup)
make run -- -tunnel -subdomain myapp  # choose the name: https://myapp.nport.link
./bin/agent-browser-use -h            # all options
./run -h|--help                       # launcher help
```

Run `./run` in a terminal to see the startup menu: **Restart**, **Run in background**, or **Quit**.

## Browser extensions

Extensions are **off by default**. Turn them on with an environment variable (the flag
`-extensions` does the same and overrides it):

```sh
AGENT_BROWSER_USE_EXTENSIONS=1 agent-browser-use          # load what is already installed in the profile
AGENT_BROWSER_USE_EXTENSIONS=1 \
  AGENT_BROWSER_USE_EXTENSIONS_DIR=/path/to/ext1,/path/to/ext2 agent-browser-use
```

Each folder must be an **unpacked** extension (a directory with `manifest.json`), not a
`.crx` file. Chrome does not load packed extensions from `--load-extension`. When
`AGENT_BROWSER_USE_EXTENSIONS_DIR` is empty, `<folder of the profile>/extensions` is used.

In Docker the default folder is `/data/extensions` (on the data volume), so it survives
restarts:

```sh
docker run -d --name agent-browser-use --init --shm-size=1g \
  -p 127.0.0.1:3000:3000 -v agent-browser-use-data:/data \
  -e AGENT_BROWSER_USE_EXTENSIONS=1 \
  agent-browser-use
# then copy the unpacked extensions in:
docker cp ./my-extension agent-browser-use:/data/extensions/
docker restart agent-browser-use
```

Notes:

- The browser runs headless (`--headless=new`), which loads extensions but throttles MV3
  service workers; an extension that needs a visible window may misbehave.
- The folders are passed to Chromium as `--load-extension` together with
  `--disable-extensions-except`; without the second flag Chromium ignores the first.
- A folder that does not exist is reported in the log and skipped, so a typo shows up
  instead of silently starting with no extensions.

**Browser engine.** Chromium by default (or the first installed Chromium-based browser: Chrome, Edge, Brave, Vivaldi); each keeps its own profile. Firefox/Safari are unsupported (they need CDP). For the lightest setup, install `headless-shell` via `npx @puppeteer/browsers install chrome-headless-shell@stable`.

**Tunnel.** `-tunnel` gives a public `https://….nport.link` URL via embedded `cloudflared`; without `-pw`, a random password is printed.

**Password.** `-pw` guards the UI and `/mcp`; MCP clients send `Authorization: Bearer <password>`.

**Recovery.** A frozen tab reloads (then closes); a crashed browser restarts with the same profile.

