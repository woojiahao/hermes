import React, {useEffect, useState} from "react";
import "react-draft-wysiwyg/dist/react-draft-wysiwyg.css";
import {useNavigate} from "react-router-dom";
import {useAppSelector} from "../redux/hooks";
import {errorFields, errorMessage, HermesRequest} from "../utility/request";
import ThreadForm from "../components/ThreadForm"
import {ThreadDto} from "../models/Thread"

export default function CreateThread() {
  const [error, setError] = useState("")
  const user = useAppSelector((state) => state.user.user)
  const navigate = useNavigate()

  useEffect(() => {
    if (!user)
      navigate('/login', {state: {message: 'You must be logged in before you can create a new thread'}})
  }, [])

  async function createThread(dto: ThreadDto) {
    await new HermesRequest()
      .POST()
      .endpoint("threads")
      .hasAuthorization()
      .body(dto)
      .onSuccess(_ => navigate('/'))
      .onFailure((e: errorMessage | errorFields) => {
        if (typeof e.message === 'string') {
          setError(e.message)
        } else {
          const formatted = e.message.map((e) => `${e.field} ${e.message}`).join(". ")
          setError(formatted)
        }
      })
      .onError(e => {
        setError(e.message)
      })
      .call()
  }

  return (
    <div className="single">
      <div className="title">
        <h1 className="heading">New Thread</h1>
      </div>

      <ThreadForm threadId="" action={createThread} error={error} setError={setError}/>
    </div>
  )
}
