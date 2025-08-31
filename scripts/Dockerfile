ARG RABBITMQ_VERSION=latest
FROM rabbitmq:${RABBITMQ_VERSION}-management-alpine

RUN echo "Based on Version: $RABBITMQ_VERSION"

COPY --chmod=555 download_plugins.sh .
RUN cd plugins && /download_plugins.sh && rm /download_plugins.sh