
# Power-Project
The "power-project" sample is a basic implementation of a load balancer with cache in Go. It consists of multiple components that work together to distribute incoming HTTP requests across multiple backend servers.


##  Special Features
1.**main.go**: This is the entry point of the application. It sets up the backend servers, creates a load balancer, defines the request handling logic, and starts the HTTP server. It acts as the coordinator of the load balancing process.

2.**server.go**: This file defines the Server interface and provides an implementation called SimpleServer. The Server interface represents a backend server that can handle HTTP requests. The SimpleServer struct holds information about a server's address and health check URL. It also implements methods to check if the server is alive and to serve HTTP requests.

3.**loadbalancer.go**: This file defines the LoadBalancer struct and the load balancing algorithms. The LoadBalancer struct represents the load balancer itself. It contains the port number, the load balancing algorithm to use, the list of backend servers, a mutex for synchronization, and a cache instance. The file also defines three load balancing algorithms: round-robin, weighted round-robin, and least connections. The load balancer's ServeProxy() method handles incoming requests by selecting an available backend server based on the configured load balancing algorithm and forwards the request to that server.

4.**cache.go**: This file provides a simple cache implementation. It defines the Cache struct, which internally uses a map and a doubly linked list to store key-value pairs. The cache supports setting and retrieving values by key, checking for key existence, and removing expired entries. It also tracks cache hit and miss counts for performance metrics.

5.**responsewriter.go**: This file contains a custom implementation of the http.ResponseWriter interface called ResponseWriter. It wraps the original http.ResponseWriter and adds additional methods and fields to track the response status, body content, and other information. This allows the load balancer to store successful responses in the cache.

*Overall, the "power-project" sample demonstrates the basic concepts of load balancing, including server health checks, load balancing algorithms, and caching. It provides a foundation for building more sophisticated load balancer implementations in Go.*





## Run Locally

Clone the project

```bash
  git clone https://github.com/RickDeb2004/power-project
```

Go to the project directory

```bash
  cd power-project
```

Install dependencies

```bash
  go mod init
```

Start the server

```bash
  go run main.go
```


## Tech Stack
**Go Lang **, **Computer Science Fundamentals**

##  Advantages

The "power-project" sample project offers several advantages:

1.**Load Distribution**: The load balancer evenly distributes incoming requests across multiple backend servers. This helps in distributing the workload and preventing any single server from becoming overwhelmed with traffic. By balancing the load, the project ensures better performance, improved response times, and higher availability.

2.**High Scalability**: The load balancer allows you to easily scale your application by adding or removing backend servers. As the traffic increases, you can add more servers to handle the load. Similarly, if the demand decreases, you can remove servers to optimize resource utilization. This flexibility enables your application to handle varying levels of traffic efficiently.

3.**Fault Tolerance**: The load balancer is capable of monitoring the health of backend servers. If a server becomes unresponsive or fails, the load balancer can automatically redirect traffic to other healthy servers. This enhances the fault tolerance of your application, ensuring continuous availability even in the presence of server failures.

4.**Improved Performance**: By distributing requests among multiple servers, the project can handle a higher number of concurrent connections and requests. This helps in improving the overall performance and responsiveness of your application. It reduces the chances of bottlenecks and optimizes resource utilization across the server infrastructure.

5.**Caching Mechanism**: The project includes a simple caching mechanism that stores responses from backend servers. This helps in reducing the response time for frequently accessed resources. The cache improves overall system performance by serving cached responses instead of forwarding requests to backend servers for the same resource. It can be especially beneficial for static or infrequently changing content.

6.**Flexibility**: The project provides flexibility in choosing load balancing algorithms. It supports round-robin, weighted round-robin, and least connections algorithms. This allows you to select the most suitable algorithm based on your specific requirements and the characterestics of the application
