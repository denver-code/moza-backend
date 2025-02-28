# Create test database  
```bash
docker-compose -f postgres.docker-compose.yml exec -T postgres psql -U postgres -c "CREATE DATABASE moza_test;"
```