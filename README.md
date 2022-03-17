# fun-flights-flight-api
To run API:
```
FLY_PROVIDERS=http://ase.asmt.live:8000/provider/flights1,http://ase.asmt.live:8000/provider/flights2 JAEGER_URL=http://localhost:14268/api/traces go run .
```

```
FLY_PROVIDERS=http://localhost:1323/provider/flights1,http://localhost:1323/provider/flights2,http://localhost:1323/provider/flights3 JAEGER_URL=http://localhost:14268/api/traces go run .
```

To run liters:
```
golangci-lint run
```

To run jaeger
```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.31
  ```
jaeger UI
  http://localhost:16686