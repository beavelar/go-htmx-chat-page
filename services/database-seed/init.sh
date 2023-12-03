#!/bin/bash
echo "initilize database value: $INIT_DATABASE"

if [ $INIT_DATABASE = "true" ]
then
    echo "sleeping for 5 seconds before seeding database"
    sleep 5
    psql -h $DATABASE_HOST -d $DATABASE_NAME -U $DATABASE_USER -a -f /seed/sql/init.sql
fi
