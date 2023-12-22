# videos-microservice
This is the videos microservice application.

The videos microservice is responsible for fetching the videos data of each playlist.

In order to use the videos microservice container, we will need to configure the following environment variables:
- ENVIRONMENT = DEBUG (For now the only option).
- REDIS_HOST = IP address or domain name of Redis database. In our case, we will use Azure Cache for Redis.
- REDIS_PORT = In our case, we will use 6379 (NON TLS PORT for Redis).
- PASSWORD = Redis database password.
- JAEGER_ENDPOINT = JAEGER endpoint (For traces).
- FLAKY = False.

This microservice will expose metrics at the port 8000.

We expose the port 10010 for serving our web application. 

We can get videos by browsing /{id} route.

Additionally, we can take a look at the health of the application by checking /healthz route.