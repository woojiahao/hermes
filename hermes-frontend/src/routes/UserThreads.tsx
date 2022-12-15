import React, {useEffect, useState} from "react";
import Thread from "../models/Thread"
import {useAppSelector} from "../redux/hooks"
import {HermesRequest} from "../utility/request"

export default function UserThreads() {
  const [threads, setThreads] = useState<Thread[]>([])
  const user = useAppSelector((state) => state.user.user)

  useEffect(() => {
    // TODO: Finish implementation
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads/current")
        .hasAuthorization()
        .onSuccess(s => setThreads(s))
    })()
  }, [])

  return (
    <div className="single">
      <div className="menu">
        <h1 className="heading">Your Threads</h1>
        <a href="/create-thread" className="effect-button">New Thread</a>
      </div>


    </div>
  )
}
