
# Supply Guardian

Supply Guardian is a lightweight caching proxy for package registries (such as npm). It improves security, reliability, and performance by sitting between developers and upstream registries.

Supply Guardian caches downloaded packages locally for faster installs, reduces external dependency risks, and provides a central point for enforcing policies. Designed for simplicity and speed, it helps teams ensure consistent, reproducible builds while minimizing exposure to supply chain attacks.

## Features

- 🔒 Supply chain protection – reduces exposure to upstream registry risks.
- ⚡ Local caching – speeds up installs by serving previously fetched packages.
- 🌐 Multiple registries support – configurable upstream sources (default, private).
- 📦 Transparent proxy – acts as a drop-in replacement for public registries.
- 📑 Policy enforcement – filter out potentially dangerous scripts and version that are too young.
- 🛠️ Configurable in TOML – easy setup for cache directory, ports, and registries.


## Installation

Install with git

```bash
  git clone https://github.com/AaronFryer/ChainGuardian.git
  cd ChainGuardian
```
    
## Run
```bash
  go run .
```

## Run with Docker Compose
```yaml
services:
  chainguardian:
    image: aaronfryer/chainguardian:latest
    container_name: chainguardian
    restart: unless-stopped
    ports:
      - "8080:8080"
```

## Build

```bash
  go build .
```

## Running Tests

To run tests, run the following command

```bash
  go test .
```


## Potential Roadmap

- Authentication
- release.sh
  - checksums
  - artifact signing

## License

[GNU GPLv3](https://choosealicense.com/licenses/gpl-3.0/)

