# Copy definitions file to container
docker cp cartrmq.json rabbitmq:/etc/rabbitmq/rabbitmq-definitions.json

# Import definitions
docker exec rabbitmq rabbitmqctl import_definitions /etc/rabbitmq/rabbitmq-definitions.json