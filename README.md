# ethereum_parser

This repository provides the solution to the below problem as stated on the home exercise

## Problem
Users not able to receive push notifications for incoming/outgoing transactions.
By Implementing Parser interface we would be able to hook this up to notifications service to notify about any incoming/outgoing transactions

For this solution I choose to keep it simple and not adapt any particular architecture that would put additional complexity
as this is simply a home exercise.

I would like to make a note about one of the limitations "Avoid usage of external libraries" the only 2 libraries used in this project are
* env  (so that we can have configurations stored in the environment )
* testify (due to the fact that go does not provide a good testing mechanism)

## Storage 
To address the requirement that the storage should be easily extendable and changed in future the repository pattern has been used

## Running 

There is a Makefile provided to help you run the application bellow you will find instruction on how to run it, please look in to Make file as it provides more options
    
    make run

### Running with Docker

```go
make docker-build
```

     make docker-run

### Exposed endpoinds 
localhost is used below but if you run it from docker please replace with the appropriate address

Retrieving the current block that has been parsed or if there is no block at all parsed (first time run) it fetches it
```go
    localhost:8080/getCurrentBlock 
```

Adds the address provided in the subscribers list
```go
    localhost:8080/subscribe?address={address_goes_here} 
```
Retrieves all the parsed transactions for the given address
```go
   localhost:8080/getTransactions?address{address_goes_here}
```