# tekton-webhook-admission-webhook
This is a webhook admission webhook that adds validation for Tekton pipelines and tasks.

## Installation
This project can fully run locally and includes automation to deploy a local Kubernetes cluster (using Kind).

### Requirements
* Docker
* kubectl
* [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* Go >=1.19

## Usage
### Create Cluster
First, we need to create a Kubernetes cluster:
```
❯ make cluster

🔧 Creating Kubernetes cluster...
kind create cluster --config dev/manifests/kind/kind.cluster.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.21.1) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

Make sure that the Kubernetes node is ready:
```
❯ kubectl get nodes
NAME                 STATUS   ROLES                  AGE     VERSION
kind-control-plane   Ready    control-plane,master   3m25s   v1.21.1
```

And that system pods are running happily:
```
❯ kubectl -n kube-system get pods
NAME                                         READY   STATUS    RESTARTS   AGE
coredns-558bd4d5db-thwvj                     1/1     Running   0          3m39s
coredns-558bd4d5db-w85ks                     1/1     Running   0          3m39s
etcd-kind-control-plane                      1/1     Running   0          3m56s
kindnet-84slq                                1/1     Running   0          3m40s
kube-apiserver-kind-control-plane            1/1     Running   0          3m54s
kube-controller-manager-kind-control-plane   1/1     Running   0          3m56s
kube-proxy-4h6sj                             1/1     Running   0          3m40s
kube-scheduler-kind-control-plane            1/1     Running   0          3m54s
```

### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```
❯ make deploy

📦 Building tekton-webhook Docker image...
docker build -t tekton-webhook:latest .
[+] Building 14.3s (13/13) FINISHED
...

📦 Pushing admission-webhook image into Kind's Docker daemon...
kind load docker-image tekton-webhook:latest
Image: "tekton-webhook:latest" with ID "sha256:46b8603bcc11a8fa1825190d3ed99c099096395b22a709e13ec6e7ae2f54014d" not yet present on node "kind-control-plane", loading...

⚙️  Applying cluster config...
kubectl apply -f dev/manifests/cluster-config/
namespace/apps created
mutatingwebhookconfiguration.admissionregistration.k8s.io/tekton.webhook.config created
validatingwebhookconfiguration.admissionregistration.k8s.io/tekton.webhook.config created

🚀 Deploying tekton-webhook...
kubectl apply -f dev/manifests/webhook/
deployment.apps/tekton-webhook created
service/tekton-webhook created
secret/tekton-webhook-tls created
```

Then, make sure the admission webhook pod is running (in the `default` namespace):
```
❯ kubectl get pods
NAME                                        READY   STATUS    RESTARTS   AGE
tekton-webhook-77444566b7-wzwmx   1/1     Running   0          2m21s
```

You can stream logs from it:
```
❯ make logs

🔍 Streaming tekton-webhook logs...
kubectl logs -l app=tekton-webhook -f
time="2021-09-03T04:59:10Z" level=info msg="Listening on port 443..."
time="2021-09-03T05:02:21Z" level=debug msg=healthy uri=/health
```

And hit it's health endpoint from your local machine:
```
❯ curl -k https://localhost:8443/health
OK
```

### Deploying tasks
Deploy a valid task that gets successfully created:
```
❯ make valid-task

🚀 Deploying valid pod...
kubectl apply -f dev/manifests/tasks/valid-task.yaml
tasks/valid-task created
```
You should see in the admission webhook logs that the task was validated and created.

Deploy an invalid task that gets rejected:
```
❯ make invalid-task

🚀 Deploying "invalid" task...
kubectl apply -f dev/manifests/tasks/invalid-task.yaml
Error from server: error when creating "dev/manifests/tasks/invalid-task.yaml": admission webhook "tekton.webhook.config" denied the request: pod name contains "offensive"
```
You should see in the admission webhook logs that the pod validation failed.


## Admission Logic
A set of validations and mutations are implemented in an extensible framework. Those happen on the fly when a pod is deployed and no further resources are tracked and updated (ie. no controller logic).

### Validating Webhooks
#### Implemented
- [pipeline name validation](pkg/validation/name_validator.go): validates that a pipeline name doesn't contain any offensive string
- [task name validation](pkg/validation/name_validator.go): validates that a task name doesn't contain any offensive string


