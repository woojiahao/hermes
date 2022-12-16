import React, {useEffect} from "react";
import {useLocation, useNavigate} from "react-router-dom";
import {login} from "../redux/authSlice";
import {useAppDispatch} from "../redux/hooks";
import {setJWT} from "../utility/jwt";
import {HermesRequest} from "../utility/request";
import {createUser} from "../services/user"
import Layout from "../components/Layout"

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

    await createUser(
      usernameRef.current.value,
      passwordRef.current.value,
      'USER',
      user => setSuccess(`Welcome ${user.username}. Please login to proceed!`),
      message => setError(message),
      e => setError(e.message)
    )

    setClickable(true)
  }

  async function loginAction() {
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
        dispatch(login())
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
    <Layout>
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
              onClick={async () => await loginAction()}
              className="static-button-blue"
              disabled={!clickable}>Login
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
}
