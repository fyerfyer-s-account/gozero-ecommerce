docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=pa55word -d mysql:latest
docker run --name my-redis -d redis
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management