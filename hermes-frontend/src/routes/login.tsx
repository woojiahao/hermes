import React from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { toggle } from "../redux/authSlice";
import { setJWT } from "../utility/app";
import { Request } from "../utility/request";

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
    await new Request()
      .endpoint('register')
      .body({
        'username': usernameRef.current.value,
        'password': passwordRef.current.value,
        'role': 'USER',
      })
      .onSuccess((r: { username: string }) => {
        setSuccess(`Welcome ${r.username}. Please login to proceed!`)
      })
      .onFailure((e: { message: string }) => {
        setError(e.message)
      })
      .onFailure((f) => setError(f.message))
      .call()

    setClickable(true)
  }

  async function login() {
    setError('')
    setSuccess('')
    setClickable(false)

    await new Request()
      .endpoint('login')
      .body({
        'username': usernameRef.current.value,
        'password': passwordRef.current.value
      })
      .onSuccess((r: { token: string }) => {
        setSuccess('Welcome back!')
        setJWT(r.token)
        dispatch(toggle())
        navigate('/')
      })
      .onFailure((e: { message: string }) => {
        setError(e.message)
      })
      .onError((f) => setError(f.message))
      .call()
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