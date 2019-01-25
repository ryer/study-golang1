# docker

DB接続とかgolangアプリサーバー（？）の練習用

 * これを使っています
   * https://github.com/frol/docker-alpine-glibc
   * Dockerfile に apk --no-cache add ca-certificates を足した


```sh
$ ./app.sh -echo-port 9595

$ ./app.sh -url-list /app/image_counter/testdata/data1.json
```
