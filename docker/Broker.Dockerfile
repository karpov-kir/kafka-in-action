FROM confluentinc/cp-kafka:7.7.0

USER root

# Install Confluent Hub Client
RUN mkdir -p /usr/local/share/confluent-hub-client
RUN wget -c http://client.hub.confluent.io/confluent-hub-client-latest.tar.gz -O - | tar -xz -C /usr/local/share/confluent-hub-client
ENV PATH="/usr/local/share/confluent-hub-client/bin:$PATH"

# Install JDBC connector (to connect DBs e.g. Sqlite)
RUN confluent-hub install confluentinc/kafka-connect-jdbc:10.7.11 --no-prompt

# Install kafkacat (kcat)
COPY ./confluent.repo /etc/yum.repos.d/confluent.repo
RUN dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm \
    && yum install kcat -y

USER appuser
