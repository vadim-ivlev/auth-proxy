FROM alpine

RUN apk update
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN apk add libc6-compat

# Go App
WORKDIR /app
COPY . .
RUN ls

# Start app
CMD bash -c "cd /app && ./auth-proxy -env=staging -serve=4402"

EXPOSE 4402