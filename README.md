# Geppetto - NeuroSC & PromoterDB

The Yale NeuroSC &amp; Promoter DB Project

This repository contains the frontend and backend code supporting the NeuroSC and PromoterDB websites. The frontend is built on React, while the backend is build using Golang.

The backend service contains two primary commands. THe first is an ingestion command that's walks over a specific directory structure of `.gltf` and `.csv` files to populate or update the database. The second starts a REST API server that serves the data to the frontend.

## Folder Structure

The root directory serves as the main entrypoint for the backend, specifically the `cmd/main.go` file. The web frontend is located in the `frontend` directory.

## Building the frontend

The frontend needs to be built for the backend to be able to serve it. The code is rather messy and the build process is more involved than it needs to be so tread carefully. Reference the readme located at `frontend/README.md` for instructions on how to build. Ignore the code in the `frontend/vendor` directory, it is not used in the build process currently.

## Running the backend

### Setup

Copy the .env.example file to .env and update the values as needed. The .env file is used to configure the database connection and other environment variables. Make sure you have a PostgreSQL database running and the connection details are correctly set in the `.env` file. The project uses [Goose](https://github.com/pressly/goose) for database migrations, so you can run `goose up` to apply the migrations to your database.

To run the backend, you need to have at least Go 1.24.3 installed on your machine. You can download it from [the official Go website](https://go.dev/dl/). You can confirm your Go version by running `go version` in the terminal.

### Dependencies

To install the necessary dependencies, run:

```bash
go mod tidy
```

### Migrations

To create a new migration, you can use the Goose CLI. First, install Goose by running:

```bash
go tool goose create some_migration_name sql
```

You can migrate the database with `go tool goose up`.

## Ingestion

Ingesting files is one of the core functionalities of the backend. It is very crucial that the directory structure is followed correctly as it outlines fundamental relationships between the files. For NeuroSC ingestion, we look for the following directory structure:

```bash
<DEVELOPMENTAL_STAGE>/<TIMEPOINT>/<CELL_TYPE>/<FILENAME>.gltf
```

Once you have the files in the correct structure, you can run the ingestion command from the root of the directory:

```bash
go run cmd/main.go ingest -d path/to/neaurosc/files --clean

# for additional options, run:
# go run cmd/main.go ingest --help
```

This will output ingestion progress to the console, it will skip files that are not relevant and the --clean flag will remove any existing data in the database before ingesting the new files.

## Running the API Server

To run the API server, you can use the following command:

```bash
go run cmd/main.go web
```

## TODO

- [ ] Set up a CI/CD pipeline to automate the build and deployment process.
- [ ] Add unit tests for the backend code.
- [ ] Set up docker for easier deployment and development.

```

```

```

```
