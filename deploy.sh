
#!/usr/bin/env bash
set -euo pipefail

PROJECT="task-manager"
APP="task-manager"

# Verify OpenShift access
oc whoami >/dev/null

# Create the project if necessary
if ! oc get namespace "$PROJECT" >/dev/null 2>&1; then
    oc new-project "$PROJECT"
else
    oc project "$PROJECT"
fi

# Create the database Secret on first deployment
if ! oc get secret task-manager-secret >/dev/null 2>&1; then
    read -rsp "PostgreSQL password: " DB_PASSWORD
    echo

    oc create secret generic task-manager-secret \
        --from-literal=postgres-password="$DB_PASSWORD" \
        --from-literal=database-url="postgresql://taskmanager:${DB_PASSWORD}@db:5432/taskmanager?sslmode=disable"

    unset DB_PASSWORD
fi

# Create the application build configuration if missing
if ! oc get buildconfig "$APP" >/dev/null 2>&1; then
    oc new-build --name="$APP" --binary --strategy=docker
fi

# Build the application image
oc start-build "$APP" --from-dir=. --follow --wait

# Deploy PostgreSQL and wait until it is ready
oc apply -f deploy/database-service.yaml
oc apply -f deploy/database-deployment.yaml
oc rollout status deployment/db --timeout=180s

# Deploy the application, Service, and Route
oc apply -f deploy/deployment.yaml
oc apply -f deploy/service.yaml
oc apply -f deploy/route.yaml

# Restart existing Pods to pick up the newly built image
oc rollout restart deployment/"$APP"
oc rollout status deployment/"$APP" --timeout=180s

echo
echo "Deployment complete!"
echo
oc get deployments
echo

ROUTE_HOST=$(oc get route "$APP" -o jsonpath='{.spec.host}')
echo "API URL: http://${ROUTE_HOST}/tasks"
