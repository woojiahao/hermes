import React, {useState} from "react"
import {ThreadDto} from "../models/Thread"
import ThreadForm from "../components/ThreadForm"
import {errorFields, errorMessage, HermesRequest} from "../utility/request"
import {useNavigate, useParams} from "react-router-dom"
import Layout from "../components/Layout"
import {useAppSelector} from "../redux/hooks"

export default function EditThread() {
  const {id} = useParams()
  const [error, setError] = useState("")
  const user = useAppSelector(state => state.auth.user)
  const navigate = useNavigate()

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
    <Layout isProtected="user">
      <div className="single">
        <div className="title">
          <h1 className="heading">Edit Thread</h1>
        </div>

        <ThreadForm threadId={id}
                    action={editThread}
                    error={error}
                    setError={setError}/>
      </div>
    </Layout>
  )
}
