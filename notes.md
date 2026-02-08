# Go Microservices on EKS with Vault & GitLab CI/CD

## Docker Commands

### Build Images
docker build -t swezin55/dashboard:latest app/dashboard
docker build -t swezin55/counting:latest app/counting

### Push to Registry
docker login

docker push swezin55/counting:latest
docker push swezin55/dashboard:latest
docker pull swezin55/dashboard:latest

### Manual Container Management
docker rm -f dashboard 

docker network create micro-net

docker run -d --name counting --network micro-net -p 9003:9003 -e PORT=9003 swezin55/counting:latest
docker run -d --name dashboard --network micro-net -p 9002:9002 -e PORT=9002 -e COUNTING_SERVICE_URL="http://counting:9003" swezin55/dashboard:latest

### Cleanup
docker images -aq | xargs -r docker rmi -f
docker ps -aq | xargs -r docker rm -f

## Docker Compose with Consul

Use `docker-compose up -d` for local development with Consul service registry.

Services:
- **Consul**: Service registry & discovery (port 8500)
- **3 Counting Services**: Registered with Consul
- **3 Dashboard Services**: Discover counting via Consul

# another way to run as a single line
PORT=9002 COUNTING_SERVICE_URL="http://localhost:9003" ./dashboard-service

consul services register counting1.json