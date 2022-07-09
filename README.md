# Purpose 
- Idea is to create an application which listens to kafka and is zero downtime.
- Since the app's functionality is dependent on the availability of kafka, we shall create liveliness probes for kafka
and configure kubernetes to check for the health check endpoint of service, which in turn will try to connect to kafka after certain intervals,
if connection couldn't be created, then the restartPolicy of the service shall be executed by k8

Tested with minikube `v1.25.2`


# Steps to launch the application:
1. start zookeeper and kafka server
2. build image inside the minikube docker daemon - `minikube image build . -t listener-app-image`
3. launch service - `kubectl apply -f deployment.yaml`


