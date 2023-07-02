# Tutorial: Accessing a relational database
Tutorial from [go dev tutorial database access mysql](https://go.dev/doc/tutorial/database-access).

## To create database, table and run mysql container:
`NOTE`: See the [Dockerfile](./Dockerfile) and [database-script](./create-tables.sql) for this tutorial.

- Build the mysql image:
```bash
$ docker build -t go-tutorial-mysql .
```

- Run the container:
```bash
$ docker run --name go-tutorial-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d go-tutorial-mysql
```

- Enter inside instance container:
```bash
$ docker exec -it go-tutorial-mysql bash
```

- Run `mysql` client:
```bash
$ mysql -u root -p
```

- List databases:
```mysql
mysql> show databases;
```

- Use database:
```mysql
mysql> use recordings;
```

- List tables:
```mysql
mysql> show tables;
```

- Select items:
```mysql
mysql> select * from album;
```

- To close `mysql` client:
```mysql
mysql> ctrl + \
```

- To exit from `container`:
```bash
$ exit
```

- To stop the `container`:
```bash
$ docker container stop go-tutorial-mysql
```

- To start the `container`:
```bash
$ docker container start go-tutorial-mysql
```

## To run application:
```bash
$ DBUSER=root DBPASS=123456 go run main.go
```