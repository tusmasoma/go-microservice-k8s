# go-microservice-k8s

![go-microservice-k8s drawio](https://github.com/user-attachments/assets/3e4e143b-3a15-4379-9485-cb39d8855679)

# How to Run the Application

This document explains how to run the application using Minikube. Follow the steps below to start Minikube, deploy the necessary services, and manage ingress settings.

## Prerequisites

- Minikube installed and set up on your local machine.
- Kubernetes command-line tool (`kubectl`) installed and configured.

## Step 1: Start Minikube

Start Minikube by running the following command:

```bash
minikube start
```

This initializes a local Kubernetes cluster using Minikube.

## Step 2: Apply Kubernetes Resources

Once Minikube is running, apply the necessary Kubernetes configurations in the following order:

1. **Apply shared configurations and MySQL services**:

   ```bash
   kubectl apply -f shared-configmap.yaml
   kubectl apply -f shared-mysql-statefulset.yaml
   kubectl apply -f shared-mysql-service.yaml
   ```

2. **Deploy the customer service**:

   ```bash
   kubectl apply -f customer-deployment.yaml
   kubectl apply -f customer-service.yaml
   ```

3. **Deploy the catalog service**:

   ```bash
   kubectl apply -f catalog-deployment.yaml
   kubectl apply -f catalog-service.yaml
   ```

4. **Deploy the order service**:

   ```bash
   kubectl apply -f order-deployment.yaml
   kubectl apply -f order-service.yaml
   ```

5. **Deploy the commerce-gateway service**:

   ```bash
   kubectl apply -f commerce-gateway-deployment.yaml
   kubectl apply -f commerce-gateway-service.yaml
   ```

6. **Apply ingress configuration**:

   ```bash
   kubectl apply -f ingress.yaml
   ```

This sets up the necessary deployments, services, and ingress rules for the application.

## Step 3: Access the Application

### Option 1: Use Ingress

If your ingress is working correctly, you should be able to access the application via the ingress settings. Check the ingress rules by running:

```bash
kubectl get ingress
```

### Option 2: Access Directly via Minikube (If Ingress Fails)

If the ingress does not work as expected, you can access the services directly using Minikube's `service` command. For example, to access the commerce-gateway service, run:

```bash
minikube service commerce-gateway-service
```

This will open the service in your default web browser.

### Troubleshooting Ingress

If you're encountering issues with ingress, ensure that the ingress controller is installed and running. You can check the status by running:

```bash
kubectl get pods -n ingress-nginx
```

If not installed, you can install the ingress controller using the following command:

```bash
minikube addons enable ingress
```
