import { createRouter, RouterProvider } from '@tanstack/react-router'
import './App.css'
import { routeTree } from './routes'

export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/v1"

const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

function App() {
  return <RouterProvider router={router} />;
}

export default App
