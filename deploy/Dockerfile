FROM ubuntu:18.04

RUN apt-get -y update 
RUN apt-get -y install curl
# RUN apt-get -y install iputils-ping
# RUN apt-get -y install mc
# RUN apt-get -y install nodejs
# RUN apt-get -y install sqlite3 

WORKDIR /app
COPY . .

EXPOSE 4400

CMD bash -c "cd /app && ./auth-proxy -serve 4400 -env=prod"




