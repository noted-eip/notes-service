# Notes service
___
Service responsible for managing notes and recommandation.


## Build

To run the service you only need to have [golang](https://go.dev) and [docker](https://docs.docker.com/get-docker/) installed.

After cloning the repository run:

```
make update-submodules
```

You can then build and it's dependancies with :

```
make re
```

## Configuration

| Env Name                           | Flag Name           | Default                     | Description                               |
|------------------------------------|---------------------|-----------------------------|-------------------------------------------|
| `ACCOUNTS_SERVICE_PORT`            | `--port`            | `3000`                      | The port the application shall listen on. |
| `ACCOUNTS_SERVICE_ENV`             | `--env`             | `production`                | Either `production` or `development`.     |
| `ACCOUNTS_SERVICE_MONGO_URI`       | `--mongo-uri`       | `mongodb://localhost:27017` | Address of the MongoDB server.            |
| `ACCOUNTS_SERVICE_MONGO_DB_NAME`   | `--mongo-db-name`   | `notes-service`          | Name of the Mongo database.               |


