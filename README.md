# weather-app

Professional starter structure for a weather forecasting web application built with Go, `html/template`, and vanilla frontend assets.

## Tech stack

- Go HTTP server
- HTML5 + CSS3 + Vanilla JavaScript
- Clean layered architecture

## Project structure

```text
weather-app/
├── cmd/server/main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── services/
│   └── weather/
├── templates/
│   ├── layouts/
│   ├── partials/
│   ├── error.html
│   ├── index.html
│   └── weather.html
├── static/
│   ├── css/
│   ├── js/
│   ├── icons/
│   └── images/
├── api/
└── pkg/
```

## Run

```bash
go run ./cmd/server
```

Then open: `http://localhost:8080`

## Routes

- `GET /`
- `GET /weather`
- `GET /forecast`
- `GET /about`

Note: external weather API integration is intentionally not implemented yet.
