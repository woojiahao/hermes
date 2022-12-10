import React from "react";
import { CookiesProvider } from "react-cookie";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { useAppSelector } from "../redux/hooks";
import About from "../routes/about";
import Home from "../routes/home";
import Login from "../routes/login";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home></Home>,
  },
  {
    path: "/about",
    element: <About></About>,
  },
  {
    path: "/login",
    element: <Login></Login>,
  }
])

export default function App() {
  const isLoggedIn = useAppSelector((state) => state.auth.value)

  return (
    <CookiesProvider>
      <div className="container">
        <header>
          <h1>hermes</h1>
          {/* TODO: Dynamically generate navigation depending on whether they are logged in or not */}
          <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
            {!isLoggedIn ?
              <a href="/login" className="button">Login</a> :
              <a href="/" className="button">Logout</a>
            }
          </nav>
        </header>
        <RouterProvider router={router} />
        <footer>
          <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
          <p>hermes is a web forum designed with ❤️ using React and Go</p>
        </footer>
      </div>
    </CookiesProvider>
  )
}