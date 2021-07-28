# This example needs redis-cluster, please run the docker-compose file docker-compose-redis-cluster.yaml at the root directory

## To run a subscriber
```
docker-compose -f docker-compose-publisher.yaml up
```

## To run a publisher
```
docker-compose -f docker-compose-subscriber.yaml up
```
