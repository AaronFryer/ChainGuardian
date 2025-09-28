
# Supply Guardian

Supply Guardian is a lightweight caching proxy for package registries (such as npm). It improves security, reliability, and performance by sitting between developers and upstream registries.

The proxy caches downloaded packages locally for faster installs, reduces external dependency risks, and provides a central point for enforcing policies. Designed for simplicity and speed, it helps teams ensure consistent, reproducible builds while minimizing exposure to supply chain attacks.

## Features

- 🔒 Supply chain protection – reduces exposure to upstream registry risks.
- ⚡ Local caching – speeds up installs by serving previously fetched packages.
- 🌐 Multiple registries support – configurable upstream sources (default, private).
- 📦 Transparent proxy – acts as a drop-in replacement for public registries.
- 📑 Policy enforcement (optional) – block packages with postinstall scripts or version that are too young.
- 🛠️ Configurable in TOML – easy setup for cache directory, ports, and registries.


## Installation

Install my-project with npm

```bash
  npm install my-project
  cd my-project
```
    
## Run Locally

Clone the project

```bash
  git clone https://link-to-project
```

Go to the project directory

```bash
  cd my-project
```

Install dependencies

```bash
  npm install
```

Start the server

```bash
  npm run start
```


## Running Tests

To run tests, run the following command

```bash
  npm run test
```


## Roadmap

- Additional browser support

- Add more integrations


## License

[MIT](https://choosealicense.com/licenses/mit/)

