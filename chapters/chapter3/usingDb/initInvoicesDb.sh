#!/bin/bash

rm -f ./invoices.db
sqlite3 invoices.db <<EOF
  CREATE TABLE invoices( 
    id INT PRIMARY KEY     NOT NULL,
    title           TEXT    NOT NULL,
    details        CHAR(50),
    billedamt         REAL,
    modified TIMESTAMP DEFAULT (STRFTIME('%s', 'now')) NOT NULL
  );
EOF
