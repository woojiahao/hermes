# hermes

Simple web forum built to encourage users to start discussions.

## Local Deployment

There are to two ways to deploy hermes locally.

First, go ahead and clone the repository.

```bash
git clone http://github.com/woojiahao/hermes.git
cd hermes/
```

In either case, please create a `.env` in `hermes/hermes-backend` with the following configurations:

```env
DATABASE_HOST=localhost
DATABASE_NAME=hermes
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=root
DATABASE_PORT=5432
JWT_KEY=<insert custom token here>
```

### Independent modules

You can deploy hermes as independent modules (database, backend, frontend) separately.

To do so, ensure that you have PostgreSQL, Node, and Go installed on your machine.

Create a database named `hermes` and set the `.env` file to have the appropriate credentials for accessing this database.

Then, you can initialize the database.

```bash
psql -U postgres -d hermes -a -f sql/create.sql
```

Once the database has been initialized, you can start each component.

```bash
cd hermes-backend/
go run main.go
cd ../hermes-frontend/
yarn start
```

You can access the frontend at http://localhost:3000 and the backend at http://localhost:8080.

### Docker

hermes supports Docker for deployments. To use Docker for a local deployment, ensure Docker is installed and running on your machine.

```bash
docker compose build
docker compose up
```

Docker Compose handles the initialization and setup of the application.

## Core features

### Completed

- Thread tags
- Pinned threads
- Markdown rendering
- User authentication (admins/users)

### To do

- Upvotes
- Search thread
