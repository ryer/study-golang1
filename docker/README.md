# docker

DB接続とかgolangアプリサーバー（？）の練習用

 * これを使っています
   * https://github.com/frol/docker-alpine-glibc
   * Dockerfile に apk --no-cache add ca-certificates を足した


```sh
$ docker-compose up -d

$ docker-compose exec app /app/target/debug/study-golang1-linux-amd64 -echo-port 9595

$ docker-compose exec app /app/target/debug/study-golang1-linux-amd64 -url-list /app/image_counter/testdata/data1.json
```
