# Running

Run

```sh
docker compose up -d
```

in this directory. And then run 

```sh
docker exec -it sql-pg-container-1 createdb -U postgres gopgtest
```

to create the necessary database.
