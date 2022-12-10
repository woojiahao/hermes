import React, { useEffect, useState } from 'react';
import ThreadList from '../components/ThreadList';
import Thread from '../models/Thread';
import { HermesRequest } from '../utility/request';

export default function Home() {
  const [threads, setThreads] = useState<Thread[]>([])
  const [message, setMessage] = useState("Loading...")

  useEffect(() => {
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads")
        .onSuccess((threads: Thread[]) => {
          console.log(threads)
          setThreads(threads)
          setMessage("")
        })
        .onFailure((e: { message: string }) => {
          setMessage(e.message)
        })
        .onError((e: { message: string }) => {
          setMessage(e.message)
        })
        .call()
    })()
  }, [])

  return (
    <div className="content">
      <div className="menu">
        <h1 className="heading">Threads</h1>
        <a href="/create-thread" className='button'>New Thread</a>
      </div>
      {message && <p>{message}</p>}
      {threads.length > 0 ?
        <ThreadList threads={threads}></ThreadList> :
        <p>No threads created yet.</p>
      }
    </div>
  );
}
