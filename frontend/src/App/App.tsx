import React from 'react';
import {Outlet, RootRoute, Router, RouterProvider} from '@tanstack/router';
import IndexPage from './index';
import './App.css'
import LoginPage from "./login/login";

// appRouter is the router for the Oats backend.
let appRouter: Router;

function app() {
  return (
    <div className='app'>
      <Outlet/>
    </div>
  )
}

function AppRouter() {
  const rootRoute = new RootRoute({
    component: app
  })

  const tree = rootRoute.addChildren([
    IndexPage(rootRoute),
    LoginPage(rootRoute),
  ])
  appRouter = new Router({routeTree: tree})

  return <RouterProvider router={appRouter}/>
}

declare module '@tanstack/router' {
  interface Register {
    router: typeof appRouter
  }
}

export default AppRouter;
