import React from "react";
import {CookiesProvider} from "react-cookie";
import {useDispatch} from "react-redux";
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import {toggle} from "./redux/authSlice";
import {useAppSelector} from "./redux/hooks";
import About from "./routes/About";
import CreateThread from "./routes/CreateThread";
import Home from "./routes/Home";
import Login from "./routes/Login";
import UserThreads from "./routes/UserThreads";
import {clearJWT} from "./utility/jwt";
import ExpandedThread from "./routes/ExpandedThread"
import EditThread from "./routes/EditThread"

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
  },
  {
    path: "/threads/:id",
    element: <ExpandedThread></ExpandedThread>
  },
  {
    path: "/threads/:id/edit",
    element: <EditThread></EditThread>
  }
])

export default function App() {
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
            <a href="about">About</a>
            {user && <a href="/your-threads">Your Threads</a>}
            {user && user && <p>Welcome back {user.username}!</p>}
            {!user ?
              <a href="login" className="effect-button">Login</a> :
              <a href="/" onClick={logout} className="effect-button">Logout</a>
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
