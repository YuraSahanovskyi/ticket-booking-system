#!/bin/bash

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BENCHMARK_DIR="$PROJECT_DIR/benchmarks"
RESULTS_DIR="$BENCHMARK_DIR/results"

mkdir -p "$RESULTS_DIR"

run_configuration() {
    local config_name=$1
    local enable_cache=$2

    echo "================================================================="
    echo "Running: $config_name"
    echo "================================================================="

    cd "$PROJECT_DIR"

    echo "=> Stopping previous containers and cleaning volumes..."
    docker compose down -v --remove-orphans

    echo "=> Building and starting containers in background..."
    echo "   [ENABLE_CACHE=$enable_cache]"
    ENABLE_CACHE=$enable_cache docker compose up -d

    echo "=> Waiting for app to be ready (polling every 2s)..."
    local attempt=1
    until docker exec app wget --spider -q http://localhost:8080/api/v1/events > /dev/null 2>&1; do
        echo "   ...attempt #$attempt: service not ready yet, waiting..."
        attempt=$((attempt + 1))
        sleep 2
    done
    echo "App is up and responding on port 8080."

    if [ -f "$PROJECT_DIR/benchmarks/seed_test_data.sql" ]; then
        echo "=> Seeding database..."
        docker exec -i postgres psql -U postgres -d postgres < "$PROJECT_DIR/benchmarks/seed_test_data.sql" > /dev/null 2>&1
        echo "Database seeded."
    fi

    echo "=> Waiting 3 seconds for system to stabilize..."
    sleep 3

    echo "=> Running k6 load test..."
    cd "$BENCHMARK_DIR"
    k6 run --quiet load_test.js > "$RESULTS_DIR/${config_name}.txt" 2>&1

    echo "Test complete. Results: results/${config_name}.txt"
    echo ""
}

run_configuration "1_no_cache" "false"
run_configuration "2_with_cache" "true"

cd "$PROJECT_DIR"
echo "=> Cleaning up: stopping test environment..."
docker compose down -v

echo "================================================================="
echo "All tests finished. Results are in benchmarks/results/"
echo "================================================================="