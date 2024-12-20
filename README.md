CXTracker
An opinioned tool to evaluate different techniques for Customer Experience tracking and measurement.

# Customer Experience Representation
The experience of a given user/customer is represented as a vector. Each vector is supposed to be a collection of different measurements, in their raw OR processed format. We then use a set of such a vectors to represent the end user experience of a system, at a given moment in time.

The system is assumed to be a microservice environment. 
* Each microservice has a type
* Every type has a predefined set of metrics
* A metric for a microservice type has a predefined range

Additionally:
Not every customer interacts with every microservice. This is represented by -1.0 as a value for all the metrics for a microservice for a give customer/user.

# 1. Generate Randomized Test Data
You can use the cxgenerator package to generate a dataset of arbitrary size.
Use the command below:
```
go run cmd/gendata/main.go
```
this will produce a csv file by the name `user_behavior_vectors.csv`, to be used by other packages. 

# 2. Check Cosine Similarity
Cosine similarity check, in package cxsimilarity.
```
go run cmd/simcheck/main.go
```