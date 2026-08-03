FROM docker.redpanda.com/redpandadata/redpanda:v25.1.4

EXPOSE 9092
EXPOSE 9644

ENTRYPOINT ["redpanda", "start", "--overprovisioned", "--smp=1", "--memory=512M", "--reserve-memory=0M", "--node-id=0", "--check=false", "--kafka-addr=PLAINTEXT://0.0.0.0:9092", "--advertise-kafka-addr=PLAINTEXT://redpanda:9092"]