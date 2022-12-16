import React, {useEffect, useState} from "react"
import {ThreadDto} from "../models/Thread"
import ThreadForm from "../components/ThreadForm"
import {errorFields, errorMessage, HermesRequest} from "../utility/request"
import {useNavigate, useParams} from "react-router-dom"
import {useAppSelector} from "../redux/hooks"

export default function EditThread() {
  const {id} = useParams()
  const [error, setError] = useState("")
  const user = useAppSelector((state) => state.user.user)
  const navigate = useNavigate()

  useEffect(() => {
    if (!user)
      navigate('/login', {state: {message: 'You must be logged in before you can edit a thread'}})
  }, [])

  async function editThread(dto: ThreadDto) {
    await new HermesRequest()
      .PUT()
      .endpoint(`/threads/${id}`)
      .hasAuthorization()
      .body(dto)
      .onSuccess(_ => navigate(-1))
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
    <div>
      <div className="single">
        <div className="title">
          <h1 className="heading">Edit Thread</h1>
        </div>

        <ThreadForm threadId={id}
                    action={editThread}
                    error={error}
                    setError={setError}/>
      </div>
    </div>
  )
}
