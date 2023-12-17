# Fruits API
This is a simple Fruits API service that offers three endpoints, as described below.

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
   `POST /api/fruits` Add a fruit

The complete API documentation & client examples is available in http://localhost:3000/swagger

## Requirements
To build the project

- Makefile
- Golang (1.21.5)

## How to run
### Run tests
`make test`

Example:

```shell
fruits-api git:(main) make test
?       github.com/dlouvier/fruits-api/src/docs [no test files]
=== RUN   TestFruitsApi
=== RUN   TestFruitsApi/should_return_all_fruits_in_JSON_format
=== RUN   TestFruitsApi/should_return_a_specific_fruit_in_JSON_format
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_1
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_2
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/should_return_contain_both_fruits
--- PASS: TestFruitsApi (0.00s)
    --- PASS: TestFruitsApi/should_return_all_fruits_in_JSON_format (0.00s)
    --- PASS: TestFruitsApi/should_return_a_specific_fruit_in_JSON_format (0.00s)
    --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests (0.00s)
        --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_1 (0.00s)
        --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_2 (0.00s)
        --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/should_return_contain_both_fruits (0.00s)
PASS
coverage: 58.3% of statements
ok      github.com/dlouvier/fruits-api/src      0.006s  coverage: 58.3% of statements

  fruits-api git:(main) make test
?       github.com/dlouvier/fruits-api/src/docs [no test files]
=== RUN   TestFruitsApi
=== RUN   TestFruitsApi/should_return_all_fruits_in_JSON_format
=== RUN   TestFruitsApi/should_return_a_specific_fruit_in_JSON_format
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_1
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_2
=== RUN   TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/should_return_contain_both_fruits
--- PASS: TestFruitsApi (0.00s)
    --- PASS: TestFruitsApi/should_return_all_fruits_in_JSON_format (0.00s)
    --- PASS: TestFruitsApi/should_return_a_specific_fruit_in_JSON_format (0.00s)
    --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests (0.00s)
        --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_1 (0.00s)
        --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/current_size_should_be_2 (0.00s)
        --- PASS: TestFruitsApi/should_add_a_fruit_(multiple_times)_while_keeping_persistency_across_requests/should_return_contain_both_fruits (0.00s)
PASS
coverage: 58.3% of statements
ok      github.com/dlouvier/fruits-api/src      0.006s  coverage: 58.3% of statements

```

### Run the API (requires golang installed)
`make run`

```shell
➜  fruits-api git:(main) ✗ make run

 ┌───────────────────────────────────────────────────┐
 │                   Fiber v2.51.0                   │
 │               http://127.0.0.1:3000               │
 │       (bound on host 0.0.0.0 and port 3000)       │
 │                                                   │
 │ Handlers ............. 7  Processes ........... 1 │
 │ Prefork ....... Disabled  PID ............ 300639 │
 └───────────────────────────────────────────────────┘
```
