# gongfu

Kungfu Wechat public account Frontend and Backend end integration, including the following components:

- Front-End: Angular(v16.1)、Material
- Back-End: Golang(v1.20)、Gin、Redis、JWT、MySQL

Using the Gin static files proxy Angular html, the Golang section takes a domain-driven model development approach.

## Quick start

Backend:

```shell
$ git clone
$ cd gongfu/server
$ cp config/dev.toml config/local.toml
$ go run cmd/main.go --config=config/local.toml
```

Frontend:

```shell
$ cd gongfu/web
$ yarn run start
```

Than access `http://localhost:4200` to see the web side.

## Config file

- Backend: `gongfu/server/config/default.toml`
- Frontend: `gongfu/web/src/environments/environment.ts`

## Want to help?

Create a issue or pull request.