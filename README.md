# GRPC Hello Client and Server
This repository is a simple client and headless server that can be
used in testing and demonstration. This is loosely based on the GRPC
greeting example.

The main difference between this example and the greeting example is the
injection of hostname/ip and timing information into the messages. This
allows logging of this information for introspection purposes.

## Deploy
`kubectl apply -f ./deployments/server.yaml -f ./deployments/client.yaml`

## Delete
`kubectl delete --ignore-not-found -f ./deployments/server.yaml -f ./deployments/client.yaml`

## The Server
Receives a "hello" message responds

## The Client
Periodically sends an "hello" message
