import React, {useRef, useState} from "react"
import {createUser} from "../services/user"
import Layout from "../components/Layout"

export default function CreateAdmin() {
  const adminUsernameRef = useRef<HTMLInputElement>()
  const adminPasswordRef = useRef<HTMLInputElement>()

  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [clickable, setClickable] = useState(true)

  async function createAdmin() {
    setError('')
    setSuccess('')
    setClickable(false)
    await createUser(
      adminUsernameRef.current.value,
      adminPasswordRef.current.value,
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

        <div className="form thick-card">
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

          <div className="buttons on-end">
            <button
              type="button"
              onClick={async () => await createAdmin()}
              className="static-button-blue on-end"
              disabled={!clickable}>Create
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
}
