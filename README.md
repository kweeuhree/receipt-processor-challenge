# 🧮 Receipt Processor

This repository contains a simple webservice that processes receipts and provides an ability to calculate bonus points based on the receipt data.

## 🪧 Overview

The Receipt Processor is a simple RESTful API service that accepts receipt data, returns a newly generated id and calculates bonus points.
As in-memory storage solution the program is using a Go map object. The data does not persist restarts as per reqiurements.

#### Third-Party Packages

- The `httprouter` package is used for fast and efficient routing.
- The `alice` package is used for clear and readable middleware chaining.
- The `uuid` package is used to generate new ids.

## 🔍 Prerequisites

- Go 1.22.5
- Windows, macOS, or Linux operating system

## 🖥️ Run the Program

- Make sure Go is installed locally;
- Clone or fork this repository;
- In the project folder, run the following command:

```sh
 go run ./cmd/web
```

## ▶️ Usage

- Send a POST request to `/receipts/process` with a body of a receipt to be processed;
- Receive a JSON response containing a unique id of the receipt;
- Send a GET request to `/receipts/{id}/points`, replacing `{id}` with the id from the previous step;
- Receive a JSON response containing the calculated points for the receipt.

## 💡API Specification

### Endpoint: Process Receipts

- Path: `/receipts/process`
- Method: `POST`
- Payload: Receipt JSON
- Response: JSON containing an id for the receipt.

Description:

Takes in a JSON receipt and returns a JSON object with an ID generated by uuid package.

Example Response:

```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

### Endpoint: Get Points

- Path: `/receipts/{id}/points`
- Method: `GET`
- Response: A JSON object containing the number of points awarded.

A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.

Example Response:

```json
{ "points": 32 }
```

## 🧱 Application Architecture

This project is organized into packages, such as web, handlers, helpers and utils.
Main package contains main point of entry, routes and middleware.

- **Main package**:

  - **Routes** maps incoming HTTP requests to their corresponding handler functions.
  - **Middleware**:
    - **logRequest** logs each incoming HTTP request with details such as IP, method, and URL;
    - **recoverPanic** catches any panics during request processing, closes the connection, and returns an internal server error response.

- **Handlers package**:

  - **ProcessReceipt** decodes JSON payload, validates input, sends the input to ReceiptFactory, inserts new receipt into the Go map, and returns the newly created id;
  - **GetReceiptPoints** gets receipt id from the request, sends the id into the CalculatePoints function, and encodes the received points to send back to the user;
  - **ReceiptFactory** constructs a new receipt object and returns it.

- **Helpers package**:
  - **ServerError** handles internal server errors;
  - **ClientError** handles client-side errors;
  - **NotFound** sends a 404 Not Found response;
  - **DecodeJSON** parses JSON data from HTTP requests into Go structs;
  - **EncodeJSON** serializes Go structs into JSON format for HTTP responses;
  - **GetIdFromParams** extracts and returns an identifier from URL parameters.

**Utils package** includes all functions necessary to calculate bonus points.

## 🚀 Testing

Go's built-in `testing` package is used for testing the utility functions that are responsible for calculating receipt points, as well as routes, middleware and helper functions. Test cases are organized into tables for better readability and maintainability.

### Run the Tests

In the project folder, run the following command:

```sh
 go test ./...
```

### Continuous Integration

In order to automate the testing, a GitHub Actions workflow is set up. The workflow executes all existing tests on every push to the repository. This helps to ensure code quality and catch issues early on.

**Workflow overview**:

- get the latest code from the repository;
- configure the Go environment with the specified version;
- install project dependencies using `go mod tidy`;
- run all tests using `go test ./...`.
