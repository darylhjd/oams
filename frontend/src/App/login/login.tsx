import { RootRoute, Route } from "@tanstack/router";
import { Link } from "@chakra-ui/react";

const path = "/login"

function component() {
  return (
    <>
      <Link onClick={() => {
        fetch('http://localhost:8080/api/v1/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
        })
      }}>Login</Link>
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
