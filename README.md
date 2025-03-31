# Docker Compose Generator

A command-line tool to generate Docker Compose files for common services like MinIO, PostgreSQL, MSSQL, and MySQL.

## Features

- Interactive service selection
- Configurable usernames and passwords
- Custom port mapping
- Automatic volume mounting
- Supports the following services:
  - MinIO (Object Storage)
  - PostgreSQL
  - Microsoft SQL Server
  - MySQL

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose installed on your system

## Installation

```bash
go mod tidy
go build -o docker-generate
```

## Usage

1. Run the application:

```bash
./docker-generate
```

2. Select the services you want to include using space bar and press enter to confirm.

3. For each selected service, you'll be prompted to enter:

   - Username
   - Password
   - Port number

4. The application will generate a `docker-compose.yml` file in the current directory.

5. Start the services using:

```bash
docker-compose up -d
```

## Default Ports

- MinIO: 9000 (API), 9001 (Console)
- PostgreSQL: 5432
- MSSQL: 1433
- MySQL: 3306

## Data Persistence

All services are configured with local volume mounts:

- MinIO: `./minio/data`
- PostgreSQL: `./postgresql/data`
- MSSQL: `./mssql/data`
- MySQL: `./mysql/data`
