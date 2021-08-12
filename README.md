# Simple Bank

This is a simple bank application created using golang. The main objective was to learn more about testing, handling DB
transactions, docker and using GitHub actions.

## Requirements

- [Docker](https://www.docker.com/products/docker-desktop)
- Postgres

## Run Locally

Clone the project

```bash
  https://github.com/jwambugu/go-simple-bank-class.git
```

Go to the project directory

```bash
  cd go-simple-bank-class
```

Run the docker container

```bash
    make postgres
```

Run the migrations

```bash
    make migrateup
```

Start the server

```bash
  go run main.go
```

## Running Tests

To run tests, run the following command

```bash
  make test
```

## Acknowledgements

- [Tech School](https://www.youtube.com/playlist?list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE)