FROM apache/kafka:4.0.0

COPY server.properties /opt/kafka/config/server.properties
COPY --chmod=755 start.sh /start.sh

EXPOSE 9092

ENTRYPOINT ["/start.sh"]