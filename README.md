# k8s-visualizer

## What the App Does
- Local UI wrapper for Kubernetes cluster inspection (view pods, nodes, logs, metrics)
- Runs on developer's own machine with their kubeconfig
- Behaves like kubectl with a GUI instead of CLI
- Uses docker-compose for local dev environment

## Test commands 
### Create a test pod
kubectl run test-nginx --image=nginx -n default

### Watch it appear in the WebSocket messages

### Delete it
kubectl delete pod test-nginx -n default

### Watch the delete event in WebSocket messages

### Create multiple pods
kubectl create deployment test-app --image=nginx --replicas=3 -n default

### Scale it
kubectl scale deployment test-app --replicas=5 -n default

### Delete it
kubectl delete deployment test-app -n default
