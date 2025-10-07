# 🧠 Technical Stack & Justifications

## ⚙️ Stack

- **Frontend:** Vue.js + DaisyUI
- **Backend:** Go + Gin
- **KPI Generator Service:** Python + FastAPI
- **Ops / Infrastructure:** Nginx + Docker + GitHub Workflow
- **Database:** PostgreSQL
- **Tests:** Playwright + Vitest + (Go?)

## 💬 Justifications

We took a **massive risk** by choosing the Go framework **Gin**, since none of us in the team had any prior experience with Go.  
We could have gone for easier or more familiar options, such as **Node.js**, **Symfony**, **Django**, or **Spring Boot**. But we decided to challenge ourselves.

Because of this bold choice, we opted for **more familiar technologies** for the other parts of the stack, like **Vue.js** and **FastAPI**, ensuring that at least one of us already had some experience with them.  
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
