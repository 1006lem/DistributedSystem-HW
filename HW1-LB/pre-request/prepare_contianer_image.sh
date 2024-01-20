echo "Client image installing .."
sudo docker build -t distributed-system-client:1.0 -f Dockerfile_client .
echo "Client image installed .."

echo "Load-Balancer image installing .."
sudo docker build -t distributed-system-load-balancer:1.0 -f Dockerfile_LB .
echo "Load-Balancer image installed .."

echo "Load-Balancer image installing .."
sudo docker build -t distributed-system-service:1.0 -f Dockerfile_service .
echo "Load-Balancer image installed .."

