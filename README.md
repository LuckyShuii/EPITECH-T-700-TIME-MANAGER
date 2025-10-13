# 🧠 Technical Stack & Justifications

## ⚙️ Stack

- **Frontend:** Vue.js + DaisyUI
- **Backend:** Go + Gin
- **Ops / Infrastructure:** Nginx + Docker + GitHub Workflow
- **Database:** PostgreSQL
- **Tests:** Playwright + Vitest + (Go?)

## 💬 Justifications

We took a **massive risk** by choosing the Go framework **Gin**, since none of us in the team had any prior experience with Go.  
We could have gone for easier or more familiar options, such as **Node.js**, **Symfony**, **Django**, or **Spring Boot**. But we decided to challenge ourselves.

Because of this bold choice, we opted for **more familiar technologies** for the other parts of the stack, like **Vue.js**, ensuring that at least one of us already had some experience with them.  
This balance allows us to experiment with Go while maintaining productivity and stability in the rest of the project.

# To Launch the project

## To launch the project in development mode:

Don't forget to copy the `.env.sample` file into a `.env` file and change the values. Then you can use the following command to startup & build the project.

```bash
docker compose -f dev.docker-compose.yml up --build
```

## To launch the project without seing any logs

```bash
docker compose -f dev.docker-compose.yml up --build -d
```

## To stop the containers

```bash
docker compose -f dev.docker-compose.yml down
```

or simply go with `CTRL + C` (**not CMD** for mac, but **CTRL**)

### For more commands about usual usage of Docker CLI: https://docs.docker.com/get-started/docker_cheatsheet.pdf

# Project Architecture

Each folder works as a package

`backend/internal/app/` contains the core features of the application, categorized by major themes or entities.

Each folder inside it is composed of 4 subfolders:

- `/handler` acts as a **controller**. It’s responsible for validating incoming data before processing — for example, checking that all required fields are present, ensuring the password meets security standards, verifying that the email is valid, and so on.
- `/service` handles the business logic. This is the second layer, right after validation in the handler. Here, the actual logic of the feature is implemented before reaching the repository.
- `/repository` manages all direct interactions with the database. This is the final step — all data reaching this layer should already be verified and processed upstream.
- `/model` stores all the data types and structures used throughout the backend, ensuring consistency and coherence across all layers.

`backend/internal/config/` contains the package that loads environment variables for the backend service. Every variable must be declared here to be accessible later in the backend with `config.LoadConfig().variableName`, after importing "app/internal/config" at the top of your file.

`backend/internal/db/` contains the package — split into multiple files — responsible for centralizing database connections, handling migration scripts [**COMING SOON**], and managing the Redis connection for caching [**COMING SOON**].

`backend/internal/middleware/` contains the package dedicated to the API’s middleware. It’s divided into three files:

- `auth.go`is the middleware used in the router to check if users are authenticated and have the proper roles.
- `logger.go` [**COMING SOON**] will centralize backend log management.

`backend/internal/router/` contains the package responsible for declaring all API routes, both public and protected ones that require authentication.

`backend/cmd/server/` contains the main entry point of the backend and API — the main.go file — which starts the router and applies the CORS middleware.
