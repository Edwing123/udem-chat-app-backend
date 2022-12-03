# Documentation about the HTTP API

## Available routes

Routes under `/api/images`:

| Path               | Method(s) | Auth Required | Content-Type(Request) | Content-Type(Response) |
| :----------------- | :-------- | :------------ | :-------------------- | ---------------------- |
| /profile/:id<guid> | GET       | No            | None                  | image/{jpeg,webp,png}  |

Routes under `/api/user`:

| Path    | Method(s) | Auth Required | Content-Type(Request) | Content-Type(Response) |
| :------ | :-------- | :------------ | :-------------------- | ---------------------- |
| /signup | POST      | No            | application/json      | application/json       |
| /login  | POST      | No            | application/json      | application/json       |
| /logout | POST      | Yes           | None                  | application/json       |
| /status | GET       | No            | None                  | application/json       |
| /data   | GET       | Yes           | None                  | application/json       |
| /update | PATCH     | Yes           | multipart/form-data   | application/json       |
