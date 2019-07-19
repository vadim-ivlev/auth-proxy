FROM ubuntu:18.04

RUN apt-get -y update && apt-get -y install sqlite3 

COPY configs  /app/configs
COPY migrations  /app/migrations
COPY templates  /app/templates
COPY test_apps  /app/test_apps


COPY auth-proxy  /app/auth-proxy

EXPOSE 4400

CMD bash -c "cd /app && ./auth-proxy -serve 4400 -env=prod"




