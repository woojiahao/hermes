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

  function validateAuthDetails(): [boolean, string?, string?, string?] {
    const authRegex = /^[a-zA-Z]\w{2,}$/
    const username = usernameRef.current.value.trim()
    const password = passwordRef.current.value.trim()

    if (!username.match(authRegex)) {
      return [
        false,
        'Username must start with a letter, contain at least three characters, and only include letters, digits, and _'
      ]
    }

    if (!password.match(authRegex)) {
      return [
        false,
        'Password must start with a letter, contain at least three characters, and only include letters, digits, and _'
      ]
    }

    return [true, null, username, password]
  }

  async function register() {
    setError("")
    setSuccess("")
    setClickable(false)

    const [status, error, username, password] = validateAuthDetails()
    if (!status) {
      setError(error)
      setClickable(true)
      return
    }

    await createUser(
      username,
      password,
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

    const [status, error, username, password] = validateAuthDetails()
    if (!status) {
      setError(error)
      setClickable(true)
      return
    }

    await new HermesRequest()
      .POST()
      .endpoint('login')
      .body({
        'username': username,
        'password': password
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

        <div className="form card thick">
          {error && <p className="error">{error}</p>}

          {success && <p className="success">{success}</p>}

          <div className="field">
            <p>Username</p>
            <span className="italic">No spaces, $, #, !</span>
            <input type="text" name="username" id="username" ref={usernameRef}/>
          </div>
          <div className="field">
            <p>Password</p>
            <span className="italic">No spaces, $, #, !</span>
            <input type="password" name="password" id="password" ref={passwordRef}/>
          </div>

          <div className="flex justify-between">
            <button
              type="button"
              onClick={async () => await register()}
              className="button"
              disabled={!clickable}>Register
            </button>
            <button
              type="button"
              onClick={async () => await loginAction()}
              className="button blue"
              disabled={!clickable}>Login
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
}
