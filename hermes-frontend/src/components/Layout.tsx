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
    <div className="bg-background text-dark min-h-full">
      <div className="container py-12">
        <header className="flex flex-row justify-between items-center pb-4">
          <h1 className="text-primary">hermes</h1>

          {isLoggedIn && user &&
            <nav className="flex gap-x-8 items-center">
              <Link to="/">Home</Link>
              <Link to="/about">About</Link>
              {user.role === 'ADMIN' && <Link to="/admin/new">New admin</Link>}
              {<Link to="/threads/you">Your Threads</Link>}
              {<p>Welcome back {user.username}!</p>}
              <button type="button" onClick={logoutAction} className="button effect blue">Logout</button>
            </nav>
          }

          {(!isLoggedIn || !user) &&
            <nav className="flex gap-x-8 items-center">
              <Link to="/">Home</Link>
              <Link to="/about">About</Link>
              <Link to="/login" className="button effect blue">Login</Link>
            </nav>
          }

        </header>

        {children}

        <footer className="flex justify-between pt-8">
          <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
          <p>hermes is a web forum designed with ❤️ using React and Go</p>
        </footer>
      </div>
    </div>
  )
}
