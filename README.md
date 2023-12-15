# Kaze Project 

Kaze Project is a cool, modern app designed to manage power plants. It's got this neat GraphQL API that lets you do stuff like adding and updating power plant info. It even connects to Open-Meteo to grab weather data, which is super handy for planning. The app is built using Go and packed into Docker. Plus, it's got a bunch of tests written in Go to keep everything running smoothly. Perfect for anyone into energy management and tech!

## Run the project

* make compose/build
* make compose/up
* make test/integration # in a separate console

## Test localy

* Create a New Power Plant:

```bash
curl -X POST http://localhost:8080/graphql \
-H "Content-Type: application/json" \
-d '{"query":"mutation ($input: NewPowerPlantInput!) { createPowerPlant(input: $input) { id name latitude longitude } }","variables": {"input": {"name": "Berlin/Pankow Wind Farm","latitude": 52.636083,"longitude": 13.42977}}}'
```

* Get Power Plant by ID:

```bash
curl -X POST http://localhost:8080/graphql \
-H "Content-Type: application/json" \
-d '{"query":"query getPowerPlant($id: ID!) { powerPlant(id: $id) { id name latitude longitude hasPrecipitationToday elevation weatherForecasts { time temperature precipitation windSpeed windDirection } } }","variables": {"id": "1"}}'
```

* List Power Plants (with elevation and weatherForecasts):

```bash
curl -X POST http://localhost:8080/graphql \
-H "Content-Type: application/json" \
-d '{"query":"query ListPowerPlants($page: Int, $pageSize: Int) { listPowerPlants(page: $page, pageSize: $pageSize) { powerPlants { id name latitude longitude elevation weatherForecasts { time temperature precipitation windSpeed windDirection } } totalCount } }","variables": {"page": 1,"pageSize": 10}}'
```

* List Power Plants (basic info):

```bash
curl -X POST http://localhost:8080/graphql \
-H "Content-Type: application/json" \
-d '{"query":"query ListPowerPlants($page: Int, $pageSize: Int) { listPowerPlants(page: $page, pageSize: $pageSize) { powerPlants { id name latitude longitude } totalCount } }","variables": {"page": 1,"pageSize": 10}}'
```

* Update a Power Plant:

```bash
curl -X POST http://localhost:8080/graphql \
-H "Content-Type: application/json" \
-d '{"query":"mutation UpdatePlant($id: ID!, $input: UpdatePowerPlantInput!) { updatePowerPlant(id: $id, input: $input) { id name latitude longitude } }","variables": {"id": "30","input": {"name": "Berlin Pankow Wind Farm"}}}'
```
