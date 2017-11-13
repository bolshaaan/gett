# Gett

Task: implement web server with 2 rest api handlers.

In project structure you can find acceptance test, that runs server on
seperate test port with connection to test database.

If I have to do highload service, I certainly choose using batch insert
with no ORM.

## import handler

We have some problems with status in bulk change operations.
We do not know, how to treat situation when some driver is not valid.
- cancel all
- or create good drivers

This application creates good drivers, and if there are some invalid
drivers in JSON, system will return JSON report with bad drivers in
response body.

We have to think about 207 (mutlistatus) status code, when not
all drivers are created
from the list.

## get handler
For all invalid requests get handler will return 400 status.
If there is no driver with valid id, then response will contain 404 status.

## acceptance tests
Test in api_test.go starts web-sever and runs base scenario.
Of course, there should be more tests, checking failure logic.

## checking duplicates before save

When we call import drivers with already existing id or license-number,
we must check response from database after insert. It's much easier to
make **only 1** request to database, instead of 2 queries: 1. selecting driver by id or
license, and 2. inserting new entry to database.

## install

Just to go get github.com/bolshaaan/gett

## deployed instance

Deployed instance you can try here:
http://gett.fun-bunny.ru/driver/1
:+1:

```bash

root@magickserver:~# curl http://gett.fun-bunny.ru/driver/1 | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    61  100    61    0     0   6676      0 --:--:-- --:--:-- --:--:--  6777
{
  "license_number": "12-234-45",
  "name": "Johny b. Goode",
  "id": 1
}

```


## test coverage

```bash
C:\Users\Александр\go\src\github.com\bolshaaan\gett>go test -cover ./...
?       github.com/bolshaaan/gett       [no test files]
ok      github.com/bolshaaan/gett/acceptanceTests       1.130s  coverage: 0.0% of statements
?       github.com/bolshaaan/gett/db    [no test files]
ok      github.com/bolshaaan/gett/handlers      0.107s  coverage: 73.1% of statements
?       github.com/bolshaaan/gett/main  [no test files]
ok      github.com/bolshaaan/gett/models        0.111s  coverage: 92.0% of statements

```


## Mock...

Unit tests must not connect to real database. So we have to write code, which allows to
mock driver model. Interfaces are ought to be helpers for writing this code (see MockDriver in test_helpers.go)

