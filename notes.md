docker build -t swezin55/dashboard:latest app/dashboard
docker build -t swezin55/counting:latest app/counting

docker login

docker push swezin55/counting:latest
docker push swezin55/dashboard:latest
docker pull swezin55/dashboard:latest

docker rm -f dashboard 

docker network create micro-net

docker run -d --name counting --network micro-net -p 9003:9003 -e PORT=9003 swezin55/counting:latest
docker run -d --name dashboard --network micro-net -p 9002:9002 -e PORT=9002 -e COUNTING_SERVICE_URL="http://counting:9003" swezin55/dashboard:latest

docker images -aq | xargs -r docker rmi -f
docker ps -aq | xargs -r docker rm -f