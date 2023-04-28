<div align="center">
  <p></p>
  <img src=".github/assets/favicon.svg" alt="favicon.svg" height="100">

# Oats: Online Attendance Taking System

<i>Digitalised attendance taking, made easy!</i>

</div>

## Getting Started

### Environment Variables
This project has a strict requirement on environment variables. For example, it is impossible to pass this project's test suite on your local computer without first setting up your environment!

While this may be frustrating to some, this allows us to know that our programme is truly working as how we set it up to be once it does run.

For more information on the environment variables, study the `backend/env` package.

### Tests
It is not possible to run and pass the project's test suite without some set up.

#### Test Database
This project uses **Postgresql** as its database interface. To help accurately simulate behaviour in production, our tests are created to run in actual databases as well. We can accomplish this by using **Docker**.

Set up a persistent Postgresql container. Make sure to set the following environment variables for the container:
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`

The values you use for these 2 will also be used for the project's `DATABASE_USER` and `DATABASE_PASSWORD` environment variables respectively. Also take note of the database port!

The steps to set up the container will be left as an exercise for the reader.