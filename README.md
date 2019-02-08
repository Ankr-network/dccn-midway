# dccn-midway



## How to build it in local:

To run this application, start a [redis server] on your local machine:

```sh
redis-server
```

Next, start the Go application:

```sh
go build
./dccn-midway
```

Now, using any HTTP client with support for cookies (like [Postman](https://www.getpostman.com/apps)) make a sign-up request with the appropriate credentials:

```
POST http://localhost:8000/signup

{"username":"user","password":"password", "name": "AnkrNetwork", "nickname": "Ankr""}
```


make a sign-in request with the appropriate credentials:

```
POST http://localhost:8000/signin

{"username":"user","password":"password"}
```

You can now try hitting the welcome route from the same client to get the welcome message:

```
GET http://localhost:8000/welcome
```

You can now try hitting the create route from the same client to add tasks:

```
GET http://localhost:8000/create

{
    "UserId": "123",
    "Name": "Ankr-network",
    "Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
    "Replica": 1,
    "DataCenter": "datacenter01",
    "DataCenterId": "10"}
```

You can now try hitting the list route from the same client to get the tasklist:

```
GET http://localhost:8000/list
```

Hit the refresh route,`session_token`'s length will be extended by 120s:

```
POST http://localhost:8000/refresh
```

## How to build it in k8s:

First, one need to deploy a redis-master service:
```
kubectl create -f https://k8s.io/examples/application/guestbook/redis-master-deployment.yaml
```

Then create a Redis service:
```
kubectl create -f https://k8s.io/examples/application/guestbook/redis-master-service.yaml
```
Finally, create the dccn-midway service:
```
cd dccn-midway/k8s
kubectl create -f dccn-midway.yaml
```

## How to test it in local:
To run the test, one must start a [redis server] on his/her local machine before the test:

```sh
redis-server
```

Next, start the Go application:

```sh
go build
./dccn-midway
```
Then, open a new terminal
go test -v
it will start to do the test.