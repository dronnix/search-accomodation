## Geolocation - coding challenge for Awesome Company

This is a coding challenge for Awesome Company (name changed to not spoil the solution).
The full task description can be found at [TASK.md](TASK.md).

## Quickstart
Requirements: *make* utility, *docker-engine* and *docker-compose* installed.
To run, use: `make compose-run`. It downloads/builds necessary images, and runs it all together. PostgreSQL, CSV import
utility and the API will be started. You can find import statistics and the API logs in docker-compose output:
```
postgres_1        | 2022-07-17 08:27:44.505 UTC [1] LOG:  database system is ready to accept connections
iploc-importer_1  | Time spent(sec): 20
iploc-importer_1  | Total records found: 1000000
iploc-importer_1  | Non-valid records: 100569
iploc-importer_1  | Duplicated records: 47516
iploc-importer_1  | Imported records: 851915
docker-compose_iploc-importer_1 exited with code 0
iploc-server_1    | 2022/07/17 08:31:29 "GET http://localhost:8080/v1/iplocation?ip=70.95.73.73 HTTP/1.1" from 172.19.0.1:60884 - 200 127B in 1.439118ms

```
API available at http://localhost:8080/v1/iplocation?ip=<ip>

## Design approach
For this task a lightweight version of domain-driven design is used:
- It matches with the task, helping to represent the separate model.
- Demonstrates my typical approach for the last projects :)

So, `/model` contains the main business entity`IPLocation` record and two main business scenarios: importing data (`ImportIPLocations()`)
and locate user by IP (`PredictIPLocation()`). It uses repository interfaces describing the data access layer.
Concrete implementation of the repository can be found in `/storage` (based on PostgreSQL).
API layer was designed using code-generation from OpenAPI specification, see `/api`. Some useful internal libraries can be found in `/internal`.
Finally, the code for the import utility and API, that used the above are placed in `/cmd`.

## Tradeoffs and edge-cases
1. Since input data is completely randomized, there are no way to use normalized forms to store the data, hence one-table
   approach is used.
2. We don't have a confidence level for records, so if one IP address points to different locations,
   we can't know which one is the correct one. So "not found" response is returned.

## Other tooling
* `make start-testing` runs infrastructure (mainly PostgreSQL) for integration tests.
* `make test` runs all tests (both unit and integrational) with coverage reporting and race detection. A coverage report can be found
  at [coverage.html](coverage.html)
* `make stop-testing` stops and removes the infrastructure.
  'make lint'
* `make lint` runs all linters.
* `make build` just builds binaries to `/bin` directory.
* `make install-tools` installs development tools.
* `make generate-api` generates API from OpenAPI specification.