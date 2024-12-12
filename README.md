### HARMONY AUTH API Microservice

The purpose on the microservice is maintain the main users system, auth and roles. Follow the instructions bellow to install and use the Microservice

The idea with the AUTH microservice is provide an API to the local subsystems, it means that the only way to use the AUTH microservice is an API router in the same system.

Other systems must use the local env HEADER on the API requests to use the AUTH microservice.

It will be build on AWS serverless instance.

## Install

To install and execute the harmony auth api you must use docker to build and compose the micro service. Follow the next instructions to install the service

    docker build -t harmony_auth_api:<VERSION>

## Run docker compose

Run the command on the external package to rebuild all microservices or only the ones who changed.

    docker-compose up -d

## Use the service

Import the requests json to your api requests manager (Insomnia, Postman) and use the connection on the stablished port, 8825

    api_requests_collection.json

# REST API

Use the authentication API to validate Role, Create sessions with users and obtain basic information about roles and users.

## Status

### Request

`GET /status`

Get the state of the application, it will return true when it is working normally

### Response

    {
        "state": true
    }

# Users

## Get all users

### Request

`GET /users`

Get the list of all active users on the database

### Response

    {
        "data": [
            {
                "name": "Sebastian",
                "email": "sebastianto1999@gmail.com",
                "password": "",
                "role": "DEFAULT"
            }
        ],
        "error": false
    }

## Create user

### Request

`POST /create-user`

Add new user to the database, you only can add users with the DEFAULT role

    {
        "name": "Sebastian",
        "email": "sebastianto1999@gmail.com",
        "password": "123456789"
    }

### Response

    {
        "error": false
    }

## Has role

### Request

`GET /has-role/:role/:email`

Validate the user role, if the user has the correct role it will return true

### Response

    {
        "data": true,
        "error": false
    }

## Get Role

### Request

`GET /role/:email`

Get the user role on the database

### Response

    {
        "data": "DEFAULT",
        "error": false
    }

## Set role

### Request

`POST /set-role`

Change or update the user role, you only can change it if you are an ADMIN

    {
        "email": "sebastianto1999@gmail.com",
        "role": "ADMIN",
        "fromUser": "admin@gmail.com",
        "sessionCode": "123456789"
    }

### Response

    {
        "error": false
    }

## Login

### Request

`POST /login`

Create session by user email and password, it will return the session code that will maintain the relation between user, role and auth

    {
        "email": "sebastianto1999@gmail.com",
        "password": "123456789"
    }

### Response

    {
        "data": "d6238aab-0300-48b3-81ac-79078959d3c5",
        "error": false
    }

## Validate session

### Request

`POST /validate-session`

Validate if the session is still active on the system, every session has 10 minutes to be active

    {
        "email": "admin@gmail.com",
        "sessionCode": "f4b84d31-fd82-46f7-96b3-ea7c5496f7a0"
    }

### Response

    {
        "error": false
    }

## Add new role

### Request

`POST /add-role`

Add new role to the system, only the admin can do it

    {
        "name": "DEFAULT",
        "fromUser": "sebastianto1999@gmail.com",
        "sessionCode": "380475b5-b511-42c6-932c-b295d78788df"
    }

### Response

    {
        "error": false
    }

## Get all available roles

### Request

`POST /roles`

Get the list of avaialble role on the system

### Response

    {
        "data": [
            {
                "name": "ADMIN",
                "config": {}
            },
            {
                "name": "DEFAULT",
                "config": {}
            }
        ],
        "error": false
    }

## Times

- Saturday 1 Hours
- Thuesday 2 Hours
- Wednesday 2 Hours
