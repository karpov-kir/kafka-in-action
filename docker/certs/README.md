# How to generate

The instructions are based on https://medium.com/jinternals/kafka-ssl-setup-with-self-signed-certificate-part-1-c2679a57e16c.

## Broker

Create own private Certificate Authority (CA):

- `openssl req -new -newkey rsa:4096 -days 100000 -x509 -subj "/CN=Kafka-in-action" -keyout ca-key -out ca-cert -nodes`

Create Kafka Server Certificate and store in KeyStore:

- `keytool -genkey -keystore kafka.server.keystore.jks -validity 100000 -storepass kafka123 -keypass kafka123 -dname "CN=broker1" -storetype pkcs12 -ext "SAN=DNS:broker2,DNS:broker3" -keyalg RSA`

Create Certificate signed request (CSR):

- `keytool -keystore kafka.server.keystore.jks -certreq -file broker-cert-file -storepass kafka123 -keypass kafka123`

Get CSR Signed with the CA:

- `openssl x509 -req -CA ca-cert -CAkey ca-key -in broker-cert-file -out broker-cert-file-signed -days 100000 -CAcreateserial -passin pass:kafka123`

Import CA certificate in KeyStore:

- `keytool -keystore kafka.server.keystore.jks -alias CARoot -import -file ca-cert -storepass kafka123 -keypass kafka123 -noprompt`

Import Signed CSR In KeyStore:

- `keytool -keystore kafka.server.keystore.jks -import -file broker-cert-file-signed -storepass kafka123 -keypass kafka123 -noprompt`

Import CA certificate In TrustStore:

- `keytool -keystore kafka.server.truststore.jks -alias CARoot -import -file ca-cert -storepass kafka123 -keypass kafka123 -noprompt`

## Client

Import CA certificate In TrustStore:

- `keytool -keystore kafka.client.truststore.jks -alias CARoot -import -file ca-cert -storepass kafka123 -keypass kafka123 -noprompt`
