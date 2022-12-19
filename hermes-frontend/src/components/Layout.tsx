import React, {PropsWithChildren, useEffect, useState} from "react"
import {useAppSelector} from "../redux/hooks"
import {useDispatch} from "react-redux"
import {clearJWT} from "../utility/jwt"
import {Link, useNavigate} from "react-router-dom"
import {logout} from "../redux/authSlice"
import {GiHamburgerMenu} from "react-icons/gi"

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

  const navStyles = "flex gap-x-8 items-center tablet:gap-x-4 phone:gap-y-4 phone:flex-wrap phone:justify-center"

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
      <div className="w-[80%] m-auto py-12 tablet:w-[90%] phablet:py-4">
        <header className="flex flex-row justify-between items-center pb-4 phablet:flex-col">
          <h1 className="text-primary phablet:items-center">hermes</h1>

          {isLoggedIn && user &&
            <nav className={navStyles}>
              <Link to="/">Home</Link>
              <Link to="/about">About</Link>
              {user.role === 'ADMIN' && <Link to="/admin/new">New admin</Link>}
              {<Link to="/threads/you">Your Threads</Link>}
              {<p>Welcome back {user.username}!</p>}
              <button type="button" onClick={logoutAction} className="button effect blue">Logout</button>
            </nav>
          }

          {(!isLoggedIn || !user) &&
            <nav className={navStyles}>
              <Link to="/">Home</Link>
              <Link to="/about">About</Link>
              <Link to="/login" className="button effect blue">Login</Link>
            </nav>
          }

        </header>

        {children}

        <footer className="flex justify-between pt-8 phablet:flex-col phablet:items-center phone:text-center">
          <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
          <p>hermes is a web forum designed with ❤️ using React and Go</p>
        </footer>
      </div>
    </div>
  )
}
