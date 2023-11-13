FROM postgres:latest

WORKDIR /seed

COPY ./sql/init.sql ./sql/init.sql
COPY ./init.sh ./init.sh

RUN chmod +x ./init.sh
ENTRYPOINT "/seed/init.sh"
