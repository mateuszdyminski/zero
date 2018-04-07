# zero
zero-downtime deployments in Kubernetes

# Local environment

Run minikube
```
$ minikube start 
```


```
wrk -c 10 -d 40s http://192.168.99.100:31831/api/users/1  & sleep 1 && kubectl apply -f 02_deployment.yaml


http://192.168.99.100:31831/actuator/health


```