### REST server written in go

[![Go Report Card](https://goreportcard.com/badge/github.com/LanPavletic/go-rest-server)](https://goreportcard.com/report/github.com/LanPavletic/go-rest-server)

## Table of Contents

- [Overview](#overview)
- [Built With](#built-with)
- [Features](#features)

## Overview
REST server API, built using Go and Gin, facilitates task management by connecting to MongoDB for efficient data storage. Has a robust authentication system based on JWT.
This project was made with the purpose of learning production level authetication and deployment of Rest api-s.


### Built With
[Gin Web Framework](https://gin-gonic.com): a web framework written in Golang.

[Mongodb](https://www.mongodb.com/): Robust and easy to use database for storing all our data.

[JWT](https://jwt.io/): A industry standard for securely transmitting information between parties

## Features
The api has end points for both registering and logging in users. The register handler salts and hashed passwords before storing them in the users collection.
The login route verifes the users credetials and returns a JWT token that lasts for an hour. 
The api also servers protected endpoints that require a valid JWT token.

The protected routes can:
 - Create new tasks
 - Get a specific task via Id
 - Update or delete tasks
 - get all tasks
