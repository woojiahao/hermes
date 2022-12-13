import React, {useEffect, useState} from 'react';
import ThreadList from '../components/ThreadList';
import Thread from '../models/Thread';
import {errorMessage, HermesRequest, jsonConvert} from '../utility/request';
import Tag from "../models/Tag"
import DisplayTag from "../components/DisplayTag"

export default function Home() {
  const [threads, setThreads] = useState<Thread[]>([])
  const [message, setMessage] = useState("Loading...")
  const [tags, setTags] = useState<Tag[]>([])
  const [filteredTags, setFilteredTags] = useState<Tag[]>([])

  useEffect(() => {
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads")
        .onSuccess((json) => {
          const threads = jsonConvert.deserializeArray(json, Thread)
          setThreads(threads)
          setMessage("")
        })
        .onFailure((e: errorMessage) => {
          setMessage(e.message)
        })
        .onError((e: { message: string }) => {
          setMessage(e.message)
        })
        .call()

      await new HermesRequest()
        .GET()
        .endpoint("threads/tags")
        .onSuccess(json => {
          const t = jsonConvert.deserializeArray(json, Tag)
          setTags(t)
        })
        .onFailure((e: errorMessage) => setMessage(e.message))
        .onError(e => setMessage(e.message))
        .call()
    })()
  }, [])

  return (
    <div className="split">
      <div className="actions">
        <h1 className="heading">Threads</h1>
        <a href="/create-thread" className='button'>New Thread</a>
      </div>

      <div className="content">
        <main>
          {message && <p>{message}</p>}
          {threads.length > 0 ?
            <ThreadList threads={threads.filter(thread => {
              for (const tag of filteredTags) {
                if (!thread.tags.flatMap(t => t.content).includes(tag.content)) return false
              }
              return true
            })}/> :
            <p>No threads created yet.</p>
          }
        </main>
        <aside className="tags-filter">
          <p><strong>Tags</strong> (click to filter)</p>
          <div hidden={filteredTags.length === 0}>
            {filteredTags.length > 0 && <p><em>Filtered tags</em></p>}
            {filteredTags
              .map((tag, i) => <DisplayTag key={i}
                                           tag={tag}
                                           onClick={() => setFilteredTags(prevState => prevState.filter(t => t !== tag))}
                                           style={{fontWeight: "bold"}}/>)}
          </div>
          <div>
            {filteredTags.length > 0 && <p><em>Other tags</em></p>}
            {tags
              .filter(tag => !(filteredTags.includes(tag)))
              .map((tag, i) => <DisplayTag key={i}
                                           tag={tag}
                                           onClick={() => setFilteredTags(prevState => [...prevState, tag])}/>)}
          </div>
        </aside>
      </div>
    </div>
  );
}
