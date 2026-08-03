FROM apache/kafka:4.0.0

COPY server.properties /opt/kafka/config/server.properties
COPY start.sh /start.sh

RUN chmod +x /start.sh

EXPOSE 9092

ENTRYPOINT ["/start.sh"]