# golang exercise

## What is this?
This repo consists of a simple web application with a couple endpoints. Also my first go application.

## Endpoints

+ /hash
    + method: POST
    + parameter: password (string)
    + responds with an [index] (integer) which can be used to request a sha512 hash of the password in base64 after 5 seconds
+ /hash/[index]
    + method: GET
    + responds with a hash of the password if it is ready, otherwise an empty response
+ /stats
    + method: GET
    + responds with json of the total number of hashes made as well as the average time spent processing the /hash requests
+ /shutdown
    + method: GET
    + initiates a shutdown which will exit the process after all pending hashes have been processed

## Setup

+ Install Go Tools
+ Clone the repo in your $GOPATH

## Running the server

+ For an interactive shell to easily manually test, run `./bin/start.sh`
+ To simply run the server `go build && ./golang-exercise`
