.PHONY: all mockgen-todo mockgen-file mockgen-todo-repo mockgen-file-repo mockgen-stream-repo mockgen-logger mockgen-validator mockgen-mysql mockgen-redis mockgen-s3 mockgen-fiber mockgen-mysql-adapter mockgen-redis-adapter mockgen-s3-adapter mockgen-todo-handler mockgen-file-handler

all: mockgen-todo mockgen-file mockgen-todo-repo mockgen-file-repo mockgen-stream-repo mockgen-logger mockgen-validator mockgen-mysql mockgen-redis mockgen-s3 mockgen-fiber mockgen-mysql-adapter mockgen-redis-adapter mockgen-s3-adapter mockgen-todo-handler mockgen-file-handler

mockgen-todo:
	mkdir -p usecase
	mockgen -package=mock_usecase -source=../usecase/todo_usecase.go -destination=usecase/mock_todo_usecase.go

mockgen-file:
	mkdir -p usecase
	mockgen -package=mock_usecase -source=../usecase/file_usecase.go -destination=usecase/mock_file_usecase.go

mockgen-todo-repo:
	mkdir -p repository
	mockgen -package=mock_repository -source=../domain/repository/todo_repository.go -destination=repository/mock_todo_repository.go

mockgen-file-repo:
	mkdir -p repository
	mockgen -package=mock_repository -source=../domain/repository/file_repository.go -destination=repository/mock_file_repository.go

mockgen-stream-repo:
	mkdir -p repository
	mockgen -package=mock_repository -source=../domain/repository/stream_repository.go -destination=repository/mock_stream_repository.go

mockgen-logger:
	mkdir -p infrastructure
	mockgen -package=mock_infrastructure -source=../infrastructure/logger/logger.go -destination=infrastructure/mock_logger.go

mockgen-validator:
	mkdir -p validator
	mockgen -package=mock_validator -source=../domain/validator/validator.go -destination=validator/mock_validator.go

mockgen-mysql:
	mkdir -p infrastructure/database
	mockgen -package=mock_database -source=../infrastructure/database/mysql.go -destination=infrastructure/database/mock_mysql.go

mockgen-redis:
	mkdir -p infrastructure/database
	mockgen -package=mock_database -source=../infrastructure/database/redis.go -destination=infrastructure/database/mock_redis.go

mockgen-s3:
	mkdir -p infrastructure/storage
	mockgen -package=mock_storage -source=../infrastructure/storage/s3.go -destination=infrastructure/storage/mock_s3.go

mockgen-fiber:
	mkdir -p infrastructure/server
	mockgen -package=mock_server -source=../infrastructure/server/fiber.go -destination=infrastructure/server/mock_fiber.go

mockgen-mysql-adapter:
	mkdir -p adapter/persistence/mysql
	mockgen -package=mock_mysql_adapter -source=../adapter/persistence/mysql/todo_repository.go -destination=adapter/persistence/mysql/mock_todo_repository.go

mockgen-redis-adapter:
	mkdir -p adapter/persistence/redis
	mockgen -package=mock_redis_adapter -source=../adapter/persistence/redis/stream_repository.go -destination=adapter/persistence/redis/mock_stream_repository.go

mockgen-s3-adapter:
	mkdir -p adapter/persistence/s3
	mockgen -package=mock_s3_adapter -source=../adapter/persistence/s3/file_repository.go -destination=adapter/persistence/s3/mock_file_repository.go

mockgen-todo-handler:
	mkdir -p adapter/api
	mockgen -package=mock_api -source=../adapter/api/todo_handler.go -destination=adapter/api/mock_todo_handler.go

mockgen-file-handler:
	mkdir -p adapter/api
	mockgen -package=mock_api -source=../adapter/api/file_handler.go -destination=adapter/api/mock_file_handler.go