FROM ubuntu:18.04

# RUN apt-get -y update && apt-get -y install sqlite3 
# RUN apt-get -y update && apt-get -y install iputils-ping && apt-get -y install curl && apt-get -y install mc

WORKDIR /app
COPY . .
COPY ./testapps/node_modules ./testapps/node_modules
COPY ./testapps/node_modules ./testapps/nodemodules


EXPOSE 4400

CMD bash -c "cd /app && ./auth-proxy -serve 4400 -env=prod"




