# videos-microservice
This is the videos microservice application.

The videos microservice  is responsible of fetching the videos data of each playlist.

In order to use the video-microservice container mainly we will need to configure the folowing environment variables :
- ENVIRONMENT = DEBUG (For now the only option)
- REDIS_HOST = Ip address or domain name of Redis database. In our case we will use Azure cache Redis
- REDIS_PORT = in our case we will use 6379 (NON TLS PORT dor Redis)
- PASSWORD = Redis database password
- JAEGER_ENDPOINT = JAEKER Enpoint (For traces)
This microservice will expose metrics at the port 8000
- FLAKY = False


We expose the port 10010 for serving our web application. You can get videos by browsing /{id} routes. Additionally you can tage a look athe health of the application by check / route   


The following repository is inspired by this amazing work [link](https://github.com/kubees/videos-microservice/) 