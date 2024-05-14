# GoBlncr ⚖️
A Round Robin style, simple Golang load balancer. Prior to sending the user's request to the test server, it first verifies if the server is accessible. A basic web server built on Golang is used to mimic live web servers. The simple Golang server's three docker containers are spun up by the docker compose command. 

## Flow Diagram 🏗️
![Untitled-2024-02-22-1020](https://github.com/bishalr0y/GoBlncr/assets/56751927/8b210b18-d0c8-4bcd-938f-9d8e9102d582)


## How to run the project?🪄
- Open the terminal and fire command below
```
git clone <the URL of the GitHub repo>
```
- ``cd`` into the project directory and use docker-compose in the `test-servers` directory to use ports `3001`,`3002` and `3003` to run the test servers 
```
cd test-servers
docker compose up -d
```
- After that, in the root directory run the `main.go` file to start the load balancer server
```
go run main.go
```
- Finally visit `http://localhost:8000` in your browser to see the load balancer in action
