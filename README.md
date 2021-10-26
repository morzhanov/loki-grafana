# loki-grafana

An example of using Loki with Grafana for an Golang app monitoring.

## Running
```bash

cd ./deploy
docker-compose up -d

cd ..
go run ./cmd/main.go

```

After that you could go to `http://127.0.0.1:3000` and setup Loki data source in the Grafana UI.
