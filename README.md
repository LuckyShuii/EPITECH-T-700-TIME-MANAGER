# Project Overview [Time Manager]

This project is a time management, HR application designed to help users track and analyze their work sessions. It provides features such as starting and stopping work sessions, viewing active sessions, and generating key performance indicators (KPIs) based on the recorded data.
Managers can also monitor their team's productivity through aggregated KPIs and Admins have the ability to manage users and view overall system statistics.

This application is built using a modern tech stack, including Vue.js for the frontend, Go with the Gin framework for the backend, and PostgreSQL for data storage. The application is containerized using Docker and orchestrated with Docker Compose for easy deployment and scalability.

## Found it online!

This school project is accessible on: https://wildtransfer.fr/
You can log in with the following admin credentials:

- username: `admin1` | `manager1` | `employee1`
- password: `user`

# Technical Stack & Justifications

## Stack

- **Frontend:** Vue.js + DaisyUI
- **Backend:** Go + Gin
- **Ops / Infrastructure:** Nginx + Docker + GitHub Workflow (CI)
- **Database:** PostgreSQL
- **Tests:** Playwright + Vitest + Testing (Go) + Sonarqube/Sonarcloud

## Swagger

Local URL: `http://localhost:8081/api/swagger/index.html`

## Stack Justifications

We took a **risk** by choosing the Go framework **Gin**, since none of us in the team had any prior experience with Go.  
We could have gone for easier or more familiar options, such as **Node.js**, **Symfony**, **Django**, or **Spring Boot**. But we decided to challenge ourselves.

Because of this bold choice, we opted for **more familiar technologies** for the other parts of the stack, like **Vue.js**, ensuring that at least one of us already had some experience with them.  
This balance allows us to experiment with Go while maintaining productivity and stability in the rest of the project.

# To Launch the project

Go in the scripts directory to create all the logs files

```bash
cd scripts/
```

Execute this script

```bash
./init_log_files.sh
```

⚠️ Don't forget to create `.env` files in the following services: `frontend`, `backend` and at the root of the project.
⚠️ Use the `.env.sample` in each of these directories as model and change the values accordingly to your needs.

## To launch the project in development mode (fully working):

Don't forget to copy the `.env.sample` file into a `.env` file and change the values. Then you can use the following command to startup & build the project.

```bash
docker compose -f dev.docker-compose.yml up --build
```

⚠️ 🛑 You will not be able to create a new user unless you get a proper BREVO API KEY into the ./backend/.env

## To launch the project in production mode (almost fully working):

```bash
docker compose -f prod.docker-compose.yml up --build
```

⚠️ Make sure you change the PROJECT_STATUS in the ./backend/.env to "PROD" and the VITE_PROJECT_STATUS in ./frontend/.env to "PROD"

---

✅ If you started the application in DEV mode, you can access the application at: `http://localhost:8081`
With the following default admin credentials:

- username: `admin1`
- password: `user` (if you changed the `FIXTURES_PASSWORD` variable in the `./backend/.env`, use that one instead)

✅ If you started the application in PROD mode, you can access the application at: `http://localhost:8081`
With the following default admin credentials:

- username: `root` (ROOT_USERNAME in the ./backend/.env)
- password: `changeme` (ROOT_PASSWORD in the ./backend/.env)

## To stop the containers

```bash
docker compose -f dev.docker-compose.yml down
```

> For more commands about usual usage of Docker CLI: https://docs.docker.com/get-started/docker_cheatsheet.pdf

# Project Architecture

## Tests end-to-end

The structure is work in progress through `e2e.docker-compose.yml` file and `./e2e` directory.
It does not work fully yet, but the idea is to have a separate environment to run the e2e tests with Playwright, using a separate database and backend instance.

## ./scripts

`archive_logs.sh`: All logs from containers are persistent within the `./logs/...` directory. This script detects for each file is they are bigger than 5MB and compress them with the timestamp if that's the case, and put them in `./logs/archives/...`

`archived_work_session_active.sh`: The table in the DB `work_session_active` is meant to hold only the data from the last 30 days to ensure low traffic and better performences. This is a script that should be run with cron jobs on a daily basis, to put the data from this table, that are older than 30d into the `work_session_archived` table.

`archived_work_session_archived.sh`: Same principle as before, but hold the data from 30d. old to 2years maximum, then these data can be converted with this script into the table `work_session_history`.

`init_log_files.sh`: This files creates files and directories that doesn't exists and are required to startup the application. This should be run once before the first `docker compose up`.

`init_migration.sh`: usage - `./ini_migration [migration name]`: Initialize a migration file with the given name in the CLI. This will init an empty SQL file in the proper directory, to fill yourself so the migration can be made.

`migration.sh`: usage - `./migrate.sh [up | down | down all | version | force $number]` this file handles the migrations based on the version stored in the db.

- `./migrate.sh up`: runs all the migration from the current version found in the db
- `./migrate.sh down`: revert the last migration based on the current version in the db
- `./migration.sh down all`: revert all the migrations ran before
- `./migration.sh version`: get the current migration version in the db for debug purpose
- `./migration.sh force 2`: force the migration version in the db to version 2 ⚠️ use with caution ⚠️

## ./logs

This directory holds all the logs from the different containers, to ensure persistence of data even if the containers are destroyed.

## ./backend

### ./backend/data/kpi

This directory holds all the .csv export files for the KPI data exports.

### ./backend/migrations

This directory holds all the migration files for the backend database. Each file is named with a version number and a description of the migration with two files: one for the "up" migration and one for the "down" migration.

### ./backend/fixtures.sql

This file holds all the initial data to populate the database with default values. This is useful for development and testing purposes. This file is executed after the migrations are done when the database is created for the first time in DEV mode.

### ./backend/internal/config/config.go

This file holds all the configuration settings for the backend application. It reads from environment variables and provides a structured way to access these settings throughout the application.

### ./database/init.sql

This file is executed when the PostgreSQL container is created for the first time. It sets up the initial database schema and any required extensions.

# Useful Links

- [Github Project - Backlog](https://github.com/orgs/EpitechMscProPromo2027/projects/158/views/1)
- [Swagger Documentation](http://localhost:8081/api/swagger/index.html) (when the project is running)
- [Database MLD](https://app.diagrams.net/#G1BNLfQfv4mWb3HOw3CycKxpYVYmbaP2zb#%7B%22pageId%22%3A%229Q_FTJikAci536onxmBX%22%7D)
- [Notion Documentaton](https://www.notion.so/284b7b530843801d85d1c78275c809aa?v=284b7b53084381ef994d000cbb8e2988)
- [Application Architecture](https://miro.com/app/board/uXjVN3RIUgE=/)
- [SonarCloud Report](https://sonarcloud.io/project/overview?id=LuckyShuii_EPITECH-T-700-TIME-MANAGER)
- [Github Personnal Repository - CI Working](https://github.com/LuckyShuii/EPITECH-T-700-TIME-MANAGER/tree/staging)
