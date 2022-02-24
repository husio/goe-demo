# demo



    ┌────────┐                 ┌────────┐ Store data   ┌─────┐
    │Producer├────────────────►│Consumer├─────────────►│Redis│
    └────────┘     gRPC        └────────┘              └─────┘
                  request


This is a monorepo containing two applications.

`producer` application is a gRPC client that periodically submits random data to the `consumer` service.

`consumer` application is a gRPC server that stores data in Redis.


Because this is a demo, I took shortcuts. There are no tests, code comments are minimal, logging is greatly simplified, etc.
Code is written the way that makes testing easy (i.e. now and random are not hardcoded, storage is implemented separately).

`make gen-proto` is using a specific to my setup protobuf library path. This is because I cannot install it via package manager and I do not want to put it to `include/` directory. In a real project, `include/` directory must be used.


## Run locally

Most functionality is declared in the makefile.

To run locally, first build docker images. Run `make build-images`.

To start the project, run `docker-compose up`, which will run all components (`producer`, `consumer` and `redis` applications).


## Local development

Whenever a protobuf declaration is changed, run `make gen-proto`.


## Redis database state

To watch the latest entries inserted into the database, run `make latest-store-entries`.
