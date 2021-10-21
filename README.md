# Example Fibonacci server

###Example `config.yaml`
``` yaml
port: "8080"
redis:
  adr: "rdb:6379"
  password: ""
  db: 0
grpc_port: "50051"
```

##Run
use flag -v to set config file

use flag --network=host to use local Redis store
~~~
docker build --no-cache -t fib-server .
docker run -it --network=host -v $(pwd)/config.yaml:/src/fib/config.yaml fib-server
~~~

Or

~~~
docker-compose build
docker-compose up
~~~