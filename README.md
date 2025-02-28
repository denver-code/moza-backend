# Moza Backend  
Simple backend for Moza bank application, inspired by Monzo, build in `Go`.

## Download  
```bash
git clone https://github.com/denver-code/moza-backend.git
cd moza-backend
```

## Environment Variables  
To start, configure the environment variables in the `.env` file and docker-compose.yml file.
```bash
cp .env.example .env
```  

## SSL Certificates (NGINX Optional)  
You don't need to do this if you are running the application in development mode.  
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./nginx/ssl/moza.key -out ./nginx/ssl/moza.crt
```

## Installation (Docker)  
We highly recommend using Docker for both development and production.    
```bash
docker-compose up
```  
or to run in detached mode:  
```bash
docker-compose up -d
```  


## Installation (Manual)
You will need to have a PostgreSQL database running on your machine or a remote server.  
Again, easiest way to get started is to use Docker, we have a `postgres.docker-compose.yml` file that you can use to start a PostgreSQL database.  
```bash
docker-compose -f postgres.docker-compose.yml up -d
```  

After setting up the database, it will be available on `postgres:5432` within docker `moza-network` or `localhost:5432` on your host machine, feel free to adjust ports in compose file.  

To start the application, you will need to have Go installed on your machine.  
```bash
go mod download
go run main.go
```

## API Documentation  
Coming soon...