#!/bin/bash
@REM create-topics.bat
# Wait for Kafka to be ready
echo "Waiting for Kafka to be ready..."
sleep 20

# Create topics from the producer configuration
echo "Creating Saga topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic Saga --partitions 3 --replication-factor 1

echo "Creating Order topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic Order --partitions 3 --replication-factor 1

echo "Creating Inventory topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic Inventory --partitions 3 --replication-factor 1

echo "Creating Cart topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic Cart --partitions 3 --replication-factor 1

echo "Creating Central topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic Central --partitions 3 --replication-factor 1

# Create topics from the consumer configuration
echo "Creating CREATE_ORDER_SAGA topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic CREATE_ORDER_SAGA --partitions 3 --replication-factor 1

echo "Creating UPDATE_SAGA_TRACKER topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic UPDATE_SAGA_TRACKER --partitions 3 --replication-factor 1

echo "Creating UPDATE_ROLLBACK topic..."
docker compose exec kafka kafka-topics --create --if-not-exists --bootstrap-server kafka:29092 --topic UPDATE_ROLLBACK --partitions 3 --replication-factor 1

# List all topics
echo -e "\nListing all topics:"
docker compose exec kafka kafka-topics --list --bootstrap-server kafka:29092

# Keep terminal open
read -p "Press Enter to continue..."