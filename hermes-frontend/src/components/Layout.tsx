import React, {PropsWithChildren, useEffect} from "react"
import {useAppSelector} from "../redux/hooks"
import {useDispatch} from "react-redux"
import {clearJWT} from "../utility/jwt"
import {Link, useNavigate} from "react-router-dom"
import {logout} from "../redux/authSlice"

interface LayoutProps extends PropsWithChildren {
  isProtected?: "none" | "user" | "admin"
}

Layout.defaultProps = {
  isProtected: "none"
}

export default function Layout({isProtected, children}: LayoutProps) {
  const {isLoggedIn, userState, user} = useAppSelector(state => state.auth)
  const dispatch = useDispatch()
  const navigate = useNavigate()

  useEffect(() => {
    if (isProtected === "none" || userState === "loading") return

    if (isProtected === "user" && userState === "unauthorized")
      navigate('/login', {state: {message: "Please log in first"}})

    if (isProtected === "admin" && userState === "authorized" && user.role !== "ADMIN")
      navigate('/', {state: {message: "Invalid access"}})
  }, [isLoggedIn, userState, user])

  function logoutAction() {
    clearJWT()
    dispatch(logout())
    navigate('/')
  }

  return (
    <div className="container">
      <header>
        <h1>hermes</h1>

        {isLoggedIn && user &&
          <nav>
            <Link to="/">Home</Link>
            <Link to="/about">About</Link>
            {user.role === 'ADMIN' && <Link to="/admin/new">New admin</Link>}
            {<Link to="/threads/you">Your Threads</Link>}
            {<p>Welcome back {user.username}!</p>}
            <button type="button" onClick={logoutAction} className="effect-button">Logout</button>
          </nav>
        }

        {(!isLoggedIn || !user) &&
          <nav>
            <Link to="/">Home</Link>
            <Link to="/about">About</Link>
            <Link to="/login" className="effect-button">Login</Link>
          </nav>
        }

      </header>

      {children}

      <footer>
        <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
        <p>hermes is a web forum designed with ❤️ using React and Go</p>
      </footer>
    </div>
  )
}
