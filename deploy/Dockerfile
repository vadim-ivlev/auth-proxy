FROM ubuntu:18.04

RUN apt-get -y update 
RUN apt-get -y install iputils-ping && apt-get -y install curl && apt-get -y install mc
# RUN apt-get -y install nodejs
# RUN apt-get -y install sqlite3 

WORKDIR /app
COPY . .
# COPY testapps/nodemodules ./testapps/nodemodules
# COPY ./testapps/nodemodules ./testapps/node_modules
# RUN /bin/bash -c 'cd /app/testapps && npm install'

EXPOSE 4400

CMD bash -c "cd /app && ./auth-proxy -serve 4400 -env=prod"




