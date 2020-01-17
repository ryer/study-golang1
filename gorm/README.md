# gorm

$ sqlite3 /tmp/dat < gorm/dat.sql

$ sqlite3 /tmp/dat "select * from products"

$ make ARGS="main.go -db-file-gorm /tmp/dat" run

$ sqlite3 /tmp/dat "select * from products"
