# angsuran-service


# command untuk jalankan file exe yg sudah dibuild dengan os Windows
- ./cmd/main.go


# evidence result aplikasi berada pada path berikut :
- ./result


# Dockerize app, sudah ada dalam script Makefile berikut step by stepnya :
- ketik command => make build
- setelah image di build masuk directory project hasil clone repo, example = /e/project/angsuran-service
- dalam path directory repo project berada ketika command => docker-compose up -d
- docker logs -f <container_id> => untuk pengecekan aplikasi berhasil running