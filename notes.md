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
docker rm -f $(docker ps -aq)
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


docker-compose down && docker-compose up -d --scale counting=3 --scale dashboard=3 && docker ps

docker ps --format "{{.Names}}" | grep -E "counting|dashboard" | xargs -I {} sh -c 'echo "{}: $(docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" {})"'

Register the docker Ip to consul in json

```
go-microservices-on-eks-vault-gitlab-cicd-dashboard-2: 172.20.0.5
go-microservices-on-eks-vault-gitlab-cicd-dashboard-3: 172.20.0.6
go-microservices-on-eks-vault-gitlab-cicd-dashboard-1: 172.20.0.7
go-microservices-on-eks-vault-gitlab-cicd-counting-2: 172.20.0.2
go-microservices-on-eks-vault-gitlab-cicd-counting-1: 172.20.0.3
go-microservices-on-eks-vault-gitlab-cicd-counting-3: 172.20.0.4
```


consul services register dashboard1.json
consul services register dashboard2.json
consul services register dashboard3.json
consul services register counting1.json
consul services register counting2.json
consul services register counting3.json