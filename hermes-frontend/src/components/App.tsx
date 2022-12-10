import React, { useEffect } from "react";
import { CookiesProvider } from "react-cookie";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import User from "../models/User";
import { toggle } from "../redux/authSlice";
import { useAppDispatch, useAppSelector } from "../redux/hooks";
import { load } from "../redux/userSlice";
import About from "../routes/about";
import CreateThread from "../routes/createThread";
import Home from "../routes/home";
import Login from "../routes/login";
import UserThreads from "../routes/userThreads";
import { clearJWT, hasValidJWT } from "../utility/jwt";
import { HermesRequest } from "../utility/request";

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
  },
  {
    path: "/your-threads",
    element: <UserThreads></UserThreads>
  },
  {
    path: "/create-thread",
    element: <CreateThread></CreateThread>
  }
])

export default function App() {
  const isLoggedIn = useAppSelector((state) => state.auth.value)
  const user = useAppSelector((state) => state.user.user)
  const dispatch = useAppDispatch()

  function logout() {
    dispatch(toggle())
    clearJWT()
  }

  useEffect(() => {
    if (hasValidJWT()) {
      (async () => {
        await new HermesRequest()
          .GET()
          .endpoint("/users/current")
          .hasAuthorization()
          .onSuccess((u: User) => {
            dispatch(load(u))
          })
          .call()
      })()
    }
  }, [])

  return (
    <CookiesProvider>
      <div className="container">
        <header>
          <h1>hermes</h1>
          <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
            {isLoggedIn && <a href="/your-threads">Your Threads</a>}
            {isLoggedIn && <span>Welcome back {user.username}!</span>}
            {!isLoggedIn ?
              <a href="/login" className="button">Login</a> :
              <a href="/" onClick={logout} className="button">Logout</a>
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