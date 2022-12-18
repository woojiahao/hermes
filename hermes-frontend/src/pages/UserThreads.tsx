import React, {useEffect, useState} from "react";
import Thread from "../models/Thread"
import {useAppSelector} from "../redux/hooks"
import {HermesRequest, jsonConvert} from "../utility/request"
import ThreadList from "../components/ThreadList"
import {useNavigate} from "react-router-dom"
import Layout from "../components/Layout"

export default function UserThreads() {
  const [threads, setThreads] = useState<Thread[]>([])
  const user = useAppSelector((state) => state.auth.user)
  const navigate = useNavigate()

  useEffect(() => {
    if (!user)
      navigate('/')
  }, [])

  useEffect(() => {
    // TODO: Finish implementation
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads/current")
        .hasAuthorization()
        .onSuccess(json => {
          const t = jsonConvert.deserializeArray(json, Thread)
          setThreads(t)
        })
        .call()
    })()
  }, [])

  return (
    <Layout>
      <div className="single">
        <div className="title">
          <h1 className="heading">Your Threads</h1>
          <a href="/create-thread" className="button effect blue">New Thread</a>
        </div>
        {threads.length > 0 ?
          <ThreadList threads={threads}/> :
          <p>You have not created any threads yet!</p>
        }
      </div>
    </Layout>
  )
}
