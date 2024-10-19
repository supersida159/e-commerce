FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates 

WORKDIR /app/
ADD ./app /app/


ENTRYPOINT [ "./app" ]

# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app
# docker build -t e-commerce-test .
 
# docker run -d --name e-commerce-test-container -e DATABASE_URI="root:Tungpro123@@tcp(e-commerce2:3306)/e-commerce?parseTime=true" -e REDIS_URI="my-redis:6379" --network=fd-network e-commerce-test