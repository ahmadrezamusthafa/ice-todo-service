# Todo Service

A simple todo service with file upload capabilities using S3, MySQL, and Redis.

## Prerequisites

- Docker and Docker Compose
- Go 1.24 or later (for local development)

## Running the Project

1. Clone the repository:
   ```
   git clone https://github.com/ahmadrezamusthafa/ice-todo-service.git
   cd ice-todo-service
   ```

2. Create a `.env` file (or use the existing one) with the following content:
   ```
   # Server Configuration
   PORT=8080

   # MySQL Configuration
   MYSQL_HOST=localhost
   MYSQL_PORT=3306
   MYSQL_USER=root
   MYSQL_PASSWORD=password
   MYSQL_DATABASE=todo_db

   # Redis Configuration
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=

   # AWS Configuration
   AWS_ENDPOINT=http://localhost:4566
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY=test
   AWS_SECRET_KEY=test
   S3_BUCKET=todo-files
   ```

3. Start the application using Docker Compose:
   ```
   make run
   ```
   This will start all required services (MySQL, Redis, LocalStack for S3) and the application.

4. The API will be available at http://localhost:8080

## Applying Migrations

Migrations are automatically applied when starting the application with Docker Compose. The migration service runs before the application starts.

If you need to apply migrations manually:

```
docker-compose up mysql -d
docker-compose run --rm migrate -path=/migrations -database="mysql://root:password@tcp(mysql:3306)/todo_db" up
```

## Executing Tests

To run all tests:

```
make test
```

To run benchmark tests:

```
make benchmark
```

To run specific benchmark tests (e.g., for S3 file uploads):

```
go test -bench=BenchmarkFileRepository -benchmem ./adapter/persistence/s3/
```

## Cleaning Up

To stop all services and remove volumes:

```
make clean
```