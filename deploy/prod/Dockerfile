FROM alpine

RUN apk update
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN apk add libc6-compat

# Go App
WORKDIR /app
COPY . .
RUN ls

# Start app
CMD bash -c "./auth-proxy -env=prod -serve=4400"

EXPOSE 4400