# deploy 

```cmd
gcloud config set project go-basic-web-db-app
gcloud app deploy
```

# development

## run web server

```cmd
go run cmd/web/main.go
```

## run db server

```cmd
docker-compose up -d
```