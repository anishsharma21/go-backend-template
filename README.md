# Go Backend Template

Boiler plate for a simple Golang backend. It includes the following technologies:

- **PostgreSQL Database**: Integrated with a PostgreSQL database.
- **Docker**: Containerised with a Dockerfile.
- **Goose**: Uses Goose for database migration handling.
- **Docker Compose**: Uses Docker Compose for local development setup.
- **Air**: Supports hot module reloading with Air.
- **Clerk**: Middleware authentication integrates with FE.

## Getting Started

To get started with this project, clone the repository and follow the instructions below (more instructions will be added soon).

```bash
git clone https://github.com/anishsharma21/go-backend-template.git
cd go-backend-template
```

Then, run the following command to install all the go dependencies:

```bash
go mod download
```

## Developing locally

Begin by setting the following environment variables:

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=gobe password=gobesecret dbname=gobedb sslmode=disable"
export DATABASE_URL="postgresql://gobe:gobesecret@localhost:5432/gobedb?sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations
export JWT_SECRET_KEY=jwtsecret
export CLERK_SECRET_KEY=sk_test_(get this from Clerk dashboard)
```

### Docker + postgres

To run the code locally you'll need to spin up a local `postgres` database instance. If you don't have `docker` install it using [this link](https://docs.docker.com/desktop/). You can check if its installed by running `docker version` and `docker compose version`. Then, you can run the following command to start the database on its own with persistent data which will remain even after you close it:

```bash
docker compose up -d postgres
```

The `-d` flag is to run it in detached mode - without it, all the logs will appear in your terminal and you will have start a new terminal session to run further commands. It's useful to learn about `docker` and `docker compose` so you understand how to build images and manage containers locally. You can leave this postgres database running, but if you ever want to stop it, you can run `docker compose down`. For reference, the data is persisted because a `docker volume` is created for it on your disk.

You can run the backend separately with the following command:

```bash
docker compose up -d backend --build
```

Use the `--build` flag to rebuild changes in the backend. But, its better to use `air` (see below) to run the backend because it will rebuild automatically after code changes occur.

To check logs when the containers are run in detached mode, you can use one of the following commands:

```bash
docker logs gobe_template_backend # faster in CLI
docker compose logs backend
```

### Hot Module Reloading (`air`)

This tool isn't required but its a quality of life / dx booster. `air` is used for Hot Module Reloading (HMR), which enables your code to automatically recompiled and re-run when changes are made:

```bash
go install github.com/air-verse/air@latest
```

The configuration for `air` is already present in the `.air.toml` file so you can simply run the command `air` on its own from the root of the project, and your server will be started up with HMR:

```bash
air
```

### Local Database migrations (`goose`)

Use the following command to install `goose` locally as it will not be included in the project as a dependency:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

With your database running in the background from the previous `docker compose` instructions, check that `goose` is correctly connected to your database by running the following command:

```bash
goose status
```

Ensure your database is running, then, run the following command to run the migration up:

```bash
goose up
```

If the migration went well, you should see `OK` messages next to each applied sql file, and the final line should say `successfully migrated database to version: ...`. You can check the status again to confirm the migrations occurred successfully. Further migration files can be created using the following command:

```bash
goose create {name of migration} sql
```

With the database running, run the following command to run the migration down:

```bash
goose down
```

### Postman API Testing

Protected routes are authenticated using `Clerk`. To call these endpoints in Postman, you need to create a session and generate a JWT token to authorise requests. You can do this by running the `Clerk Create Session` request, and then the `Create Clerk Session Token` request right after. These will automatically generate the session data and set the JWT token for all future requests in the `Go Backend Template / Local` folder in Postman. The tokens expire after 5min.

### Testing

Tests run locally use the local postgres database. To replicate the CICD environment, you can clear your database before running the tests. Use the following command to run tests locally:

```bash
go test ./tests -v
```

## Production

When deploying to production, you'll need to set all the above environment variables with their production variations. Assuming you're deploying to `Railway`, you can spin-up a `Postgres` database and set some of the database related environment variables to those provided by that db instance. You will also need to set the `RUN_MIGRATION` env variable to `true` in production:

```bash
RUN_MIGRATION=true
```

For `clerk`, you will need to head to the clerk dashboard and ensure that OAuth has been configured with your own credentials. This is more important for the frontend, for the backend though you will need to replace the development api key with the production variation.

## Todos

TODO: use dotenv package to get env variables from .env file instead

## License

This project is licensed under the MIT License.
Feel free to customise the content further as needed!
