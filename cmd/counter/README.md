# This example needs redis-cluster, please run the docker-compose file docker-compose-redis-cluster.yaml at the root directory

## To run an atomic counter
```
docker-compose -f docker-compose-atomic.yaml up
```

## To run a non-atomic counter
```
docker-compose -f docker-compose-non-atomic.yaml up
```