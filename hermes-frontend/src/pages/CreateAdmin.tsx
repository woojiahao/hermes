import React, {useRef, useState} from "react"
import {createUser} from "../services/user"
import Layout from "../components/Layout"
import {fieldRegex} from "../utility/general"

export default function CreateAdmin() {
  const adminUsernameRef = useRef<HTMLInputElement>()
  const adminPasswordRef = useRef<HTMLInputElement>()

  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [clickable, setClickable] = useState(true)

  function validateInput(): [boolean, string?, string?, string?] {
    const username = adminUsernameRef.current.value.trim()
    const password = adminPasswordRef.current.value.trim()

    if (!username.match(fieldRegex)) {
      return [
        false,
        'Admin username must start with a letter, contain at least three characters, and only include letters, digits, and _'
      ]
    }

    if (!password.match(fieldRegex)) {
      return [
        false,
        'Admin password must start with a letter, contain at least three characters, and only include letters, digits, and _'
      ]
    }

    return [true, null, username, password]
  }

  async function createAdmin() {
    setError('')
    setSuccess('')
    setClickable(false)

    const [status, error, username, password] = validateInput()
    if (!status) {
      setError(error)
      return
    }

    await createUser(
      username,
      password,
      'ADMIN',
      user => setSuccess(`New admin ${user.username} has been created successfully!`),
      message => setError(message),
      e => setError(e.message)
    )

    adminUsernameRef.current.value = ""
    adminPasswordRef.current.value = ""
    setClickable(true)
  }

  return (
    <Layout isProtected="admin">
      <div className="single">
        <div className="title">
          <h1 className="heading">New Admin</h1>
        </div>

        <div className="form card thick">
          {error && <p className="error">{error}</p>}

          {success && <p className="success">{success}</p>}

          <div className="field">
            <p>Admin's Username</p>
            <input type="text" name="username" id="username" ref={adminUsernameRef}/>
          </div>
          <div className="field">
            <p>Admin's Password</p>
            <input type="password" name="password" id="password" ref={adminPasswordRef}/>
          </div>

          <div className="flex justify-end">
            <button
              type="button"
              onClick={async () => await createAdmin()}
              className="button blue"
              disabled={!clickable}>Create
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
}
