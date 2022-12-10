import React from "react";
import { Cookies } from "react-cookie";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { toggle } from "../redux/authSlice";
import { partialPOST } from "../utility/request";

export default function Login() {
  const usernameRef = React.createRef<HTMLInputElement>();
  const passwordRef = React.createRef<HTMLInputElement>();

  const [clickable, setClickable] = React.useState(true);
  const [error, setError] = React.useState<string>();
  const [success, setSuccess] = React.useState<string>();
  const navigate = useNavigate()
  const dispatch = useDispatch()

  async function register() {
    setError("")
    setSuccess("")
    setClickable(false)
    await partialPOST(
      'register',
      {
        'username': usernameRef.current.value,
        'password': passwordRef.current.value,
        'role': 'USER',
      },
      (r: { username: string }) => {
        setSuccess(`Welcome ${r.username}. Please login to proceed!`)
      },
      (e: { message: string }) => {
        setError(e.message)
      },
      (f) => setError(f.message)
    )

    setClickable(true)
  }

  async function login() {
    setError('')
    setSuccess('')
    setClickable(false)

    await partialPOST(
      'login',
      {
        'username': usernameRef.current.value,
        'password': passwordRef.current.value
      },
      (r: { token: string }) => {
        setSuccess('Welcome back!')
        const cookies = new Cookies()
        cookies.set("token", r.token)
        dispatch(toggle())
        navigate('/')
      },
      (e: { message: string }) => {
        setError(e.message)
      },
      (f) => setError(f.message)
    )
  }

  return (
    <div className="content">
      <h1 className="heading">Login to hermes</h1>
      <div className="login-form">
        {error &&
          <div className="error">
            <p>{error}</p>
          </div>
        }

        {success &&
          <div className="success">
            <p>{success}</p>
          </div>
        }

        <div>
          <p>Username</p>
          <input type="text" name="username" id="username" ref={usernameRef} />
        </div>
        <div>
          <p>Password</p>
          <input type="password" name="password" id="password" ref={passwordRef} />
        </div>
        <div className="buttons">
          <button
            type="button"
            onClick={async () => await register()}
            disabled={!clickable}>Register</button>
          <button
            type="button"
            className="login-button"
            onClick={async () => await login()}
            disabled={!clickable}>Login</button>
        </div>
      </div>
    </div>
  )
}