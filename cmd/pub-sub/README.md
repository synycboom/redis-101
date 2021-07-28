# This example needs redis-cluster, please run the docker-compose file docker-compose-redis-cluster.yaml at the root directory

## To run a consumer
```
docker-compose -f docker-compose-publisher.yaml up
```

## To run a producer
```
docker-compose -f docker-compose-subscriber.yaml up
```
