FROM golang:latest AS build

# Создаем рабочую директорию и компилим
COPY . /task
WORKDIR /task
RUN go build ./cmd/main.go

FROM ubuntu:20.04
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER zenehu WITH SUPERUSER PASSWORD 'zenehu';" &&\
    createdb -O zenehu forum-task &&\
    psql -f ./db/db.sql -d forum-task &&\
    /etc/init.d/postgresql stop

COPY --from=build /task/main .

EXPOSE 5000

USER root

CMD service postgresql start && ./main