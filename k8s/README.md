# Kubernetes Deployment for Launchpad MCP

This directory contains Kubernetes manifests for deploying the launchpad-mcp application to a Kubernetes cluster.

## Files

- `deployment.yaml` - Main application deployment with health checks and resource limits
- `service.yaml` - ClusterIP service to expose the application internally
- `ingress.yaml` - Ingress configuration for external access
- `secrets.yaml` - Secret template for sensitive configuration
- `kustomization.yaml` - Kustomize configuration for managing resources

## Prerequisites

1. A running Kubernetes cluster
2. `kubectl` configured to access your cluster
3. NGINX Ingress Controller installed (for ingress)
4. cert-manager installed (for TLS certificates, optional)

## Setup

### 1. Create Secrets

Before deploying, you need to create the required secrets. The application requires a PostgreSQL database URL:

```bash
# Create the secret with your actual values
kubectl create secret generic launchpad-mcp-secrets \
  --from-literal=postgres-url="postgres://user:password@host:5432/dbname?sslmode=require"
```

Optional OAuth configuration (only if authentication is needed):
```bash
kubectl create secret generic launchpad-mcp-secrets \
  --from-literal=postgres-url="postgres://user:password@host:5432/dbname?sslmode=require" \
  --from-literal=scalekit-env-url="https://your-auth-provider.com/.well-known/jwks.json" \
  --from-literal=oauth-authentication-server="https://your-auth-provider.com" \
  --from-literal=oauth-resource-url="https://your-api.com" \
  --from-literal=oauth-resource-documentation-url="https://docs.your-api.com" \
  --from-literal=scalekit-resource-metadata-url="https://your-auth-provider.com/metadata"
```

### 2. Update Ingress Configuration

Edit `ingress.yaml` and replace `launchpad-mcp.example.com` with your actual domain:

```yaml
# Update both places in the file
- host: your-actual-domain.com
```

### 3. Deploy the Application

Using kubectl:
```bash
kubectl apply -f .
```

Using kustomize:
```bash
kubectl apply -k .
```

## Verification

Check if the deployment is successful:

```bash
# Check pods
kubectl get pods -l app=launchpad-mcp

# Check service
kubectl get svc launchpad-mcp-service

# Check ingress
kubectl get ingress launchpad-mcp-ingress

# Check deployment status
kubectl rollout status deployment/launchpad-mcp
```

## Scaling

To scale the deployment (currently set to 1 replica):

```bash
kubectl scale deployment launchpad-mcp --replicas=3
```

## Health Checks

The deployment includes:
- **Liveness Probe**: Checks `/health` endpoint every 10 seconds
- **Readiness Probe**: Checks `/health` endpoint every 5 seconds

## Troubleshooting

View application logs:
```bash
kubectl logs -l app=launchpad-mcp -f
```

Describe deployment for events:
```bash
kubectl describe deployment launchpad-mcp
```

Check pod events:
```bash
kubectl describe pod -l app=launchpad-mcp
```

## CI/CD Integration

The GitHub Actions workflow automatically deploys new releases when a release is created. The deployment process:

1. Builds and pushes Docker image to GitHub Container Registry
2. Updates the Kubernetes deployment with the new image tag
3. Waits for rollout to complete
4. Verifies deployment success

Required GitHub Secret:
- `K8S_CONFIG_FILE_B64`: Base64-encoded kubeconfig file for cluster access

## Security

The deployment runs with:
- Non-root user (UID 1001)
- Read-only root filesystem
- No privileged escalation
- Dropped capabilities
- Resource limits enforced