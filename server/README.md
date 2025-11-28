# How to run
## Install dependencies
`go mod tidy`

## Run the server
`go run cmd/main.go`

# Build the image
docker build -t k8s-visualizer-backend:latest ./server

# Local Test
- run `deploy.sh`
-  Get the service URL
minikube service k8s-visualizer-backend -n k8s-visualizer --url

- Or use port-forward
kubectl port-forward -n k8s-visualizer \
  svc/k8s-visualizer-backend 8080:8080

- Test the API  
```
curl http://localhost:8080/health  
curl http://localhost:8080/api/nodes  
curl http://localhost:8080/api/namespaces  
curl http://localhost:8080/api/pods?namespace=all
```


- Alternative: Test Locally First  
Make sure minikube is running
minikube status  

- Your local Go app will use ~/.kube/config automatically
```
cd backend
go run cmd/server/main.go
```

-  In another terminal, test the API
```
curl http://localhost:8080/health
curl http://localhost:8080/api/nodes
```

# useful commands
## View logs
kubectl logs -n k8s-visualizer -l app=k8s-visualizer -f

## Check pod status
kubectl get pods -n k8s-visualizer

## Describe deployment
kubectl describe deployment k8s-visualizer-backend -n k8s-visualizer

## Delete and redeploy
kubectl delete -f k8s/
./deploy.sh

## Rebuild just the image
eval $(minikube docker-env)
docker build -t k8s-visualizer-backend:latest ./backend
kubectl rollout restart deployment/k8s-visualizer-backend -n k8s-visualizer

## SSH into minikube
minikube ssh

## Stop minikube
minikube stop

## Delete minikube cluster
minikube delete