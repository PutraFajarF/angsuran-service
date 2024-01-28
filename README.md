# angsuran-service


# command untuk jalankan file exe yg sudah dibuild dengan os Windows
- ./cmd/main.go


# evidence result aplikasi berada pada path berikut :
- ./result


# Dockerize app, sudah ada dalam script Makefile berikut step by stepnya :
- clone repo
- masuk directory project hasil clone repo, example = /e/project/angsuran-service
- ketik command di terminal => make build
- setelah dibuild imagenya dalam path directory repo project berada ketik command berikut => docker-compose up -d
- docker logs -f <container_id> => untuk pengecekan aplikasi berhasil running
- docker-compose down jika ingin stop service running container