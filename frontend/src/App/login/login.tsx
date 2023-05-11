import { RootRoute, Route } from "@tanstack/router";
import { Link } from "@chakra-ui/react";

const path = "/login"

interface loginResponse {
  redirect_url: string
}

function component() {
  let initAuthCodeFlow = async () => {
    let response = await fetch('http://localhost:8080/api/v1/login?' +
      new URLSearchParams({
        redirect_url: 'http://localhost:3000/',
      }), {
      method: 'GET',
    })

    const body = await response.text()
    let resp: loginResponse = JSON.parse(body)
    console.log(resp)
    window.location.href = resp.redirect_url
  }

  return (
    <>
      <Link onClick={initAuthCodeFlow}>Login</Link>
    </>
  )
}

function LoginPage(root: RootRoute) {
  return new Route({
    getParentRoute: () => root,
    path: path,
    component: component,
  })
}

export default LoginPage
