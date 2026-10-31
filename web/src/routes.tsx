import { createRootRoute, createRoute } from "@tanstack/react-router"
import { Confirmation } from "./pages/Confirmation"
import { Root } from "./components/root"
import { Home } from "./pages/Home"

const rootRoute = createRootRoute({
    component: Root
})

const indexRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/',
    component: Home
})

const confirmRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/confirm/$token',
    component: Confirmation
})

export const routeTree = rootRoute.addChildren([indexRoute, confirmRoute])