import React from "react";

export default function LoginPopup() {
  return (
    <div className="login">
      <h3 className="heading">Login to hermes</h3>
      <form action="POST">
        <div>
          <label htmlFor="username">Username</label>
          <input type="text" name="username" id="username" />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input type="password" name="password" id="password" />
        </div>
        <div className="buttons">
          <button>Register</button>
          <button className="login-button">Login</button>
        </div>
      </form>
    </div>
  )
}