import React, {useEffect, useState} from "react";
import Thread from "../models/Thread"
import {useAppSelector} from "../redux/hooks"
import {HermesRequest, jsonConvert} from "../utility/request"
import ThreadList from "../components/ThreadList"

export default function UserThreads() {
  const [threads, setThreads] = useState<Thread[]>([])

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
    <div className="single">
      <div className="title">
        <h1 className="heading">Your Threads</h1>
        <a href="/create-thread" className="effect-button">New Thread</a>
      </div>
      <ThreadList threads={threads} />
    </div>
  )
}
