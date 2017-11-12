# Gett


Task: implement web server with 2 rest api handlers.

In project structure you can find acceptance tests, that run server on
seperate test port and on test database.

## import handler

We have some problems with status in bulk change operations.
We do not know, how to treat, if some driver is not valid.
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

## checking duplicates before save

When we call import drivers with already exist id or license-number,
we must check response from database after insert. It's much easier to
make **only 1** request to database, instead of 2 queries: 1. selecting driver by id or
license, and 2. inserting new entry to database.
