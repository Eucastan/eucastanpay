#!/bin/bash
set -e

CONFIG=/opt/kafka/config/server.properties
DATA=/var/lib/kafka/data
META=$DATA/meta.properties

mkdir -p "$DATA"

if [ ! -f "$META" ]; then
    echo "Formatting Kafka storage..."

    CLUSTER_ID=$(/opt/kafka/bin/kafka-storage.sh random-uuid)

    /opt/kafka/bin/kafka-storage.sh format \
        --ignore-formatted \
        -t "$CLUSTER_ID" \
        -c "$CONFIG"
fi

echo "========== Kafka Configuration =========="
cat "$CONFIG"
echo "========================================="

exec /opt/kafka/bin/kafka-server-start.sh "$CONFIG"