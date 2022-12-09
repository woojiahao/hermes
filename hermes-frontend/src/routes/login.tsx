import React from "react";

export default function Login() {
  const usernameRef = React.createRef<HTMLInputElement>();
  const passwordRef = React.createRef<HTMLInputElement>();

  const [clickable, setClickable] = React.useState(true);
  const [error, setError] = React.useState<string>();
  const [success, setSuccess] = React.useState<string>();

  async function register() {
    setError("")
    setSuccess("")
    setClickable(false)

    try {
      const result = await fetch('http://localhost:8081/register', {
        method: 'post',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'username': usernameRef.current.value,
          'password': passwordRef.current.value,
          'role': 'USER',
        })
      })

      const json = await result.json()

      if (!result.ok) {
        const err = json as { message: string }
        setError(err.message)
      } else {
        const s = json as { username: string }
        setSuccess(`Welcome ${s.username}. Please login to proceed!`)
      }
    } catch (e) {
      console.log(e)
      setError(e)
    }

    setClickable(true)
  }

  async function login() {
    setError('')
    setSuccess('')
    setClickable(false)

    try {
      const result = await fetch('http://localhost:8081/login', {
        method: "POST",
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'username': usernameRef.current.value,
          'password': passwordRef.current.value
        })
      })

      const json = await result.json()

      if (!result.ok) {
        console.log(json)
        const err = json as { message: string }
        setError(err.message)
      } else {
        const s = json as { token: string }
        setSuccess(s.token)
      }
    } catch (e) {
      setError(e)
    }
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