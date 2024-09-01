#!/bin/bash

id="$1"

echo "Inserting invoice with id $id"

sqlite3 invoices.db <<EOF
  INSERT INTO invoices (id,title,details,billedamt)  VALUES ($id, 'book-$id', 'Franz Kafka $id', $id.00 );
EOF
