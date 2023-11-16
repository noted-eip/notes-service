# Notes service
___
Service responsible for managing notes and recommandations.


## Build

To run the service you only need to have [golang](https://go.dev) or [docker](https://docs.docker.com/get-docker/) installed.

After cloning the repository run:

```
make update-submodules
```
This will update the git submodules referencing our gRPC models and gRPC API definition.


You can then build the project by running the following command :

```
make re
```

Or by building and running the Dockerfile.

## Tests

In order to test the application, you have to have your MongoDB database on the side.

It is as easy as running `go test`. `-v` flag for verbose is recommended.

## Configuration

### CLI configuration

You can either use the env variable or the flag to set the corresponding values.

| Env Name                           | Flag Name           | Default                     | Description                               |
|------------------------------------|---------------------|-----------------------------|-------------------------------------------|
| `NOTES_SERVICE_PORT`            | `--port`            | `3000`                      | The port the application shall listen on. |
| `NOTES_SERVICE_ENV`             | `--env`             | `production`                | Either `production` or `development`.     |
| `NOTES_SERVICE_MONGO_URI`       | `--mongo-uri`       | `mongodb://localhost:27017` | Address of the MongoDB server.            |
| `NOTES_SERVICE_MONGO_DB_NAME`   | `--mongo-db-name`   | `notes-service`          | Name of the Mongo database.               |
| `NOTES_SERVICE_ACCOUNT_SERVICE_URL`   | `--account-service-url`   | `accounts.noted.koyeb:3000`          | Account service's address               |
| `NOTES_SERVICE_JWT_PRIVATE_KEY`   | `--jwt-private-key`   |           | JWT private key used for authentification               |
| `NOTES_SERVICE_GMAIL_SUPER_SECRET`   | `--gmail-super-secret`   |         | Gmail secret to send emails.               |

### Other env variables

| Env Name                           | Description                               |
|------------------------------------|-------------------------------------------|
| `JSON_GOOGLE_CREDS_B64`            | Google credentials used to connect to Google's natural API, in json, converted in base 64  |
| `GOOGLE_API_KEY`            | Google API key used to connect to Google's knowledge graph|
| `OPENAI_API_KEY`            | OpenAI API key used to connect to OpenAI's GPT|
