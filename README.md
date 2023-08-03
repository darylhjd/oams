<div align="center">

# OAMS: Online Attendance Management System

<i>Digitalised attendance taking, made easy!</i>

</div>

## Setting up your environment
### Backend Module
The backend is built on Go and PostgreSQL.

The first step is to set up the environment variables. The recommendation is to name your `.env` files appropriate names
to avoid confusion (e.g. `.env.test`, `.env.staging`, etc...).
- You may refer to the example `.env.example` for more information on what values to enter.
- The test suite also requires an `.env` file. This is to help accurately simulate the real conditions in the staging/production
environments.

The test suite requires set up. More precisely, you need a local PostgreSQL docker container. Remember to fill in the 
appropriate database values for connection in your `.env` file for tests.

When running tests, for example with the `go test ./...` command in the `backend` module, you must remember to set the appropriate
environment variables together with the programme. The recommendation is to use a proper IDE that can set up the appropriate
configurations for the testing (e.g. helping you to load a local `.env.test` file before running the `go test ./...` command).

If you have followed the above steps, you will be able to successfully run the test suite with no issues.

A note on connecting to the staging database, remember to download the SSL certificate, and appropriately point it to the file
location in the `.env` file. Otherwise, you will not be able to successfully connect.

### Frontend Module
The frontend is built on Flutter.

The `frontend` module requires its own set of environment variables that is defined separately from the `backend` module.
Refer to the `env.example.json` file for details on the formatting of the file formatting and values required.

As always, the recommendation is to use a proper IDE that can help you set up the environment variables before running any
build or run commands. For your information, you can use the `--dart-define-from-file=env.json` argument with `flutter`
to load the correct environment variables at compile time.
