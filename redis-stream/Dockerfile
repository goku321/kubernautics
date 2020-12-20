FROM golang

WORKDIR /app

COPY ./app .

RUN go build -o main .

CMD [ "./main" ]