![fruits-api](https://github.com/dlouvier/fruits-api/assets/13359249/072a2baa-e4e0-4e64-9edf-7f419318a2dc)

# Fruits API
This is a simple Fruits API service that offers four endpoints, as described below.

## Endpoints
1. Return all fruits in JSON format:

   `GET /api/fruits/` Get the list of all fruits
   
   ```json
   [{
   "id": <ID>,
   "fruit": "apple",
   "color": "red"
   }]
   ```

1. Return a specific fruit in JSON format

   `GET /api/fruits/{id}` Get a fruit by its Id

1. Add a fruit to the basket by sending a JSON payload 

   `POST /api/fruits`

   with a Payload like: `{"fruit": "banana", "color": "yellow"}`

1. Search a fruit(s) providing a payload

   `POST /api/fruits/search`

   with a Payload like: `{"color": "yellow"}` should return all the fruits which matches the color yellow.

The complete API documentation & client examples are available in [localhost:3000/swagger](http://localhost:3000/swagger)

## Requirements
To build the project

- Makefile
- Golang (v1.21.5)

## How to start the API locally
- Before of anything, ensure the port `3000` is free in your host (`ss -snltp | grep ':3000'`)
- Please also notice: **I only tested in x86_64/Linux** so the docker image & binary built might not run. In that case install `golang 1.21.5` and `make run`
- If during the use of the API any fruit is added, it will be store in a local file. This file will be automatically read the next time it starts.

### Docker
To run the API using Docker, execute the following command:
`docker run -p 3000:3000 ghcr.io/dlouvier/fruits_api:v0.4.0`

Note:
- You must first authenticate with the GitHub Container Registry using docker login ghcr.io.
- Persistence of the objects is not guaranteed if the container is removed. To preserve the data, you can use docker commit.

### Run the binary
Download the binary [the artifacts list](https://github.com/dlouvier/fruits-api/actions/runs/7239380543), unzip & `chmod +x` & run.

### Run directly from Makefile
Simply type `make run` (golang v1.21.5 is required)

### Interacting with the API
Once the program is running, the API will be available at [localhost:3000/api/fruits](http://localhost:3000/api/fruits)
You can use `curl` or any other HTTP Client. Request needs to be set `Parameter content type: application/json`

For an extensive list of examples, curl commands, and a live demo, visit the [Swagger UI](http://localhost:3000/swagger).

![Screenshot from 2023-12-17 16-48-18](https://github.com/dlouvier/fruits-api/assets/13359249/94859b4d-a9e0-4281-a72e-93671452d047)



