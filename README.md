# Zero-Downtime deployments in Kubernetes

Repository contains code and presentation of: "Zero-Downtime deployments in Kubernetes"

## Presentation

[zero-downtime.pptx](presentation/zero-downtime.pptx)

## Requirements

* minikube
* kubectl
* helm
* sql
* JDK8 or higher for Java examples
* Golang 1.8 or higher for Go examples

### Demo Users Service

Start minikube:
```
minikube start
```

Install MySQL(with Helm):
```
helm install stable/mysql
```

Create db:
```
mysql/mysql.sh create
```

Run Users Service in Kubernetes:
```
kubectl apply -f kube/java/01_users_srv.yaml
```

Run Users Deployment in Kubernetes:
```
kubectl apply -f kube/java/02_deployment.yaml
```

Verification:

Open:
[http://192.168.99.100:30001/api/users/1](http://192.168.99.100:30001/api/users/1)
[http://192.168.99.100:30001/actuator/info](http://192.168.99.100:30001/actuator/info)