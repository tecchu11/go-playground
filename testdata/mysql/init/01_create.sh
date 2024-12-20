#!/bin/sh

echo "### Create keycloack db and user ###"

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<EOF
CREATE DATABASE IF NOT EXISTS $KC_DB_URL_DATABASE;
EOF

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<EOF
CREATE USER IF NOT EXISTS '$KC_DB_USERNAME'@'%' IDENTIFIED BY '$KC_DB_PASSWORD';
GRANT ALL PRIVILEGES ON $KC_DB_URL_DATABASE.* TO '$KC_DB_USERNAME'@'%';
FLUSH PRIVILEGES;
EOF

echo "### END ###"
