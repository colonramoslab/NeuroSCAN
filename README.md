# Geppetto - NeuroSC & PromoterDB

The Yale NeuroSC &amp; Promoter DB Project

This repository contains the frontend and backend code supporting the NeuroSC and PromoterDB websites. The frontend is built on React, while the backend is build using Golang.

The backend service contains two primary commands. THe first is an ingestion command that's walks over a specific directory structure of `.gltf` and `.csv` files to populate or update the database. The second starts a REST API server that serves the data to the frontend.

## Requirements

- Go 1.24.3 or later: You can download it from [the official Go website](https://go.dev/dl/).
- [PostgreSQL](https://www.postgresql.org): The backend uses PostgreSQL as the database. Make sure you have it installed and running on your machine. This can be installed via Homebrew on macOS with `brew install postgresql` or on Ubuntu with `sudo apt install postgresql`. It can also be run on mac using the [Postgres app](https://postgresapp.com/downloads.html).
- Node JS: The frontend is built using react and is a little bit fussy about the version. You need version 14.21.3 installed. It is recommended to us [nvm](https://github.com/nvm-sh/nvm) to install and manage this.

## Folder Structure

The root directory serves as the main entrypoint for the backend, specifically the `cmd/main.go` file. The web frontend is located in the `frontend` directory.

The core backend logic is inside `internal` with files named after the entities they relate to (neurons, contacts, synapses, etc). The structure is as follows:

**domain** - Contains the core domain model along with any pertinent logic relating to the model.
**repository** - Contains the database queries and converting the database model to the domain model.
**service** - Contains the business logic and orchestrates the interaction between the repository and handlers.
**handler** - Contains the HTTP handlers that respond to API requests.

## Building the frontend

The frontend needs to be built for the backend to be able to serve it. The code is rather messy and the build process is more involved than it needs to be so tread carefully. Reference the readme located at `frontend/README.md` for instructions on how to build. Ignore the code in the `frontend/vendor` directory, it is not used in the build process currently.

## Running the backend

### Setup

Copy the .env.example file to .env (`cp .env.example .env`) and update the values as needed. The .env file is used to configure the database connection and other environment variables. Make sure you have a PostgreSQL database running and the connection details are correctly set in the `.env` file. The project uses [Goose](https://github.com/pressly/goose) for database migrations, so you can run `goose up` to apply the migrations to your database.

To run the backend, you need to have at least Go 1.24.3 installed on your machine. You can download it from [the official Go website](https://go.dev/dl/). You can confirm your Go version by running `go version` in the terminal.

### Dependencies

To install the necessary dependencies, run:

```bash
go mod tidy
```

## Configuration

The backend uses a `.env` file for configuration. You can copy the `.env.example` file to `.env` and update the values as needed. The `.env` file contains the database connection details and other environment variables.
The following environment variables are required:

````env
APP_ENV="development"
PORT="8080"
LOG_LEVEL="debug"

# Absolute path to the frontend/build folder
APP_FRONTEND_DIR=

# Absolute path the the top level gltf directory. This is the folder with "neuroscan" and "promoters" in it.
APP_GLTF_DIR=

# Database config
DB_DSN="postgres://postgres:@localhost:5432/neuroscan"

# Goose config
GOOSE_DRIVER=postgres
# Goose DB String
GOOSE_DBSTRING=
GOOSE_MIGRATION_DIR=./migrations
```

## Ingestion

Ingesting files is one of the core functionalities of the backend. It is very crucial that the directory structure is followed correctly as it outlines fundamental relationships between the files. For NeuroSC ingestion, we look for the following directory structure:

```bash
<DEVELOPMENTAL_STAGE>/<TIMEPOINT>/<CELL_TYPE>/<FILENAME>.gltf
````

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
# a port can be specified in the .env file or by using the --port(-p) flag. The flag will override the .env file. The default port is 8080.
```

## TODO

- [ ] Set up a CI/CD pipeline to automate the build and deployment process.
- [ ] Add unit tests for the backend code.
- [ ] Set up docker for easier deployment and development.

## Additional Notes

### Migrations

To create a new migration, you can use the Goose CLI. First, install Goose by running:

```bash
go tool goose create some_migration_name sql
```

You can migrate the database with `go tool goose up`.
