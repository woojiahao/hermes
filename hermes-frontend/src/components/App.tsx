import React from "react";
import {CookiesProvider} from "react-cookie";
import {useDispatch} from "react-redux";
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import {toggle} from "../redux/authSlice";
import {useAppSelector} from "../redux/hooks";
import About from "../routes/About";
import CreateThread from "../routes/CreateThread";
import Home from "../routes/Home";
import Login from "../routes/Login";
import UserThreads from "../routes/UserThreads";
import {clearJWT} from "../utility/jwt";

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
  // TODO: Remove isLoggedIn flag since the user's existence would indicate if the user is logged in
  const isLoggedIn = useAppSelector((state) => state.auth.value)
  const user = useAppSelector((state) => state.user.user)
  const dispatch = useDispatch()

  function logout() {
    dispatch(toggle())
    clearJWT()
  }

  return (
    <CookiesProvider>
      <div className="container">
        <header>
          <h1>hermes</h1>
          <nav>
            <a href="/">Home</a>
            <a href="/About.tsx">About</a>
            {isLoggedIn && <a href="/your-threads">Your Threads</a>}
            {isLoggedIn && user && <p>Welcome back {user.username}!</p>}
            {!isLoggedIn ?
              <a href="/Login.tsx" className="button">Login</a> :
              <a href="/" onClick={logout} className="button">Logout</a>
            }
          </nav>
        </header>
        <RouterProvider router={router}/>
        <footer>
          <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
          <p>hermes is a web forum designed with ❤️ using React and Go</p>
        </footer>
      </div>
    </CookiesProvider>
  )
}
