import React, {useEffect} from "react";
import {useLocation, useNavigate} from "react-router-dom";
import {toggle} from "../redux/authSlice";
import {useAppDispatch} from "../redux/hooks";
import {loadCurrentUser} from "../redux/userSlice";
import {setJWT} from "../utility/jwt";
import {HermesRequest} from "../utility/request";

export default function Login() {
  const usernameRef = React.createRef<HTMLInputElement>();
  const passwordRef = React.createRef<HTMLInputElement>();

  const [clickable, setClickable] = React.useState(true);
  const [error, setError] = React.useState<string>();
  const [success, setSuccess] = React.useState<string>();
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const {state} = useLocation()

  useEffect(() => {
    if (state) setError(state.message)
  }, [])

  async function register() {
    setError("")
    setSuccess("")
    setClickable(false)
    await new HermesRequest()
      .POST()
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
      .onFailure((e) => {
        if (typeof e.message === 'string') {
          setError(e.message)
        } else {
          const formatted = e.message.map((e) => `${e.field} ${e.message}`).join(". ")
          setError(formatted)
        }
      })
      .call()

    setClickable(true)
  }

  async function login() {
    setError('')
    setSuccess('')
    setClickable(false)

    await new HermesRequest()
      .POST()
      .endpoint('login')
      .body({
        'username': usernameRef.current.value,
        'password': passwordRef.current.value
      })
      .onSuccess((r: { token: string }) => {
        setSuccess('Welcome back!')
        setJWT(r.token)
        dispatch(toggle())
        dispatch(loadCurrentUser())
        navigate('/')
      })
      .onFailure((e: { message: string }) => {
        setError(e.message)
      })
      .onError((f) => setError(f.message))
      .call()

    setClickable(true)
  }

  return (
    <div className="single">
      <div className="title">
        <h1 className="heading">Login to hermes</h1>
      </div>

      <div className="form thick-card">
        {error && <p className="error">{error}</p>}

        {success && <p className="success">{success}</p>}

        <div className="field">
          <p>Username</p>
          <input type="text" name="username" id="username" ref={usernameRef}/>
        </div>
        <div className="field">
          <p>Password</p>
          <input type="password" name="password" id="password" ref={passwordRef}/>
        </div>

        <div className="buttons">
          <button
            type="button"
            onClick={async () => await register()}
            className="static-button-plain"
            disabled={!clickable}>Register
          </button>
          <button
            type="button"
            onClick={async () => await login()}
            className="static-button-blue"
            disabled={!clickable}>Login
          </button>
        </div>
      </div>
    </div>
  )
}
