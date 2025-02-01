docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=pa55word -d -p 3306:3306 mysql:latest
docker run --name my-redis -d -p 6379:6379 redis
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management