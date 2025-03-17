#!/bin/bash

echo "Executing MongoDB Setup..."
mongosh -u root -p pass.123 < /etc/mongo/setup.js
echo "Setup Completed!"

echo "Executing MongoDB Init..."
for file in $(ls /etc/mongo/init-data/*.js); do mongosh -u root -p pass.123 < $file; done
echo INIT_DATA_DONE >> /tmp/dbstatus
echo "Init Completed!"