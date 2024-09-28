#!/bin/bash
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
TIMEOUT_SETUP=30

# Load RabbitMQ  environment varibales
source "${SCRIPT_DIR}/testacc.env"

setup() {
    docker compose -f "${SCRIPT_DIR}/docker-compose.yml" up -d

    echo "Waiting for RabbitMQ to be up"
    i=0
    until curl -s "${RABBITMQ_ENDPOINT}/api" > /dev/null; do
        i=$((i + 1))
        if [ $i -eq ${TIMEOUT_SETUP} ]; then
            echo
            echo "Timeout while waiting for RabbitMQ to be up"
            exit 1
        fi
        printf "."
        sleep 2
    done
}

run() {
    if [ "$GITHUB_ACTIONS" = "true" ]; then
        echo "Running under GitHub Actions"
        echo "Workflow: $GITHUB_WORKFLOW"
        echo "Run ID: $GITHUB_RUN_ID"
    else
        echo "Running locally"
    fi

    go test -cover -count=1 ./internal/provider -v -timeout 120m -coverprofile coverage.out

    # keep the return value for the scripts to fail and clean properly
    return $?
}

cleanup() {
    docker compose -f "${SCRIPT_DIR}/docker-compose.yml" down
}

main() {
    setup
    run || (cleanup && exit 1)
    cleanup
}


case "$1" in
    "setup")
        setup
        ;;
    "cleanup")
        cleanup
        ;;
    *)
        main
        ;;
esac
