# K8Sman

This repo contains logic to perform operations on k8s cluster

```sh
kubectl create clusterrolebinding k8sman-app-cluster-admin \
  --clusterrole=cluster-admin \
  --serviceaccount=k8sman:k8sman-app

kubectl create clusterrole k8sman --verb=get,list,watch,create,delete,patch,update --resource=deployments.apps

kubectl create clusterrolebinding k8sman-binding --clusterrole=k8sman --serviceaccount=k8sman:k8sman-app
```
