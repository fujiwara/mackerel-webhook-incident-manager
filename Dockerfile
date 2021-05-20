FROM golang:1.16

RUN mkdir /var/myapp
COPY . /var/myapp
WORKDIR /var/myapp
RUN go build -o app .

EXPOSE 8000
CMD /var/myapp/app
