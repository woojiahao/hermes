import React, { createRef, useEffect, useState } from "react";
import { Editor } from 'react-draft-wysiwyg';
import "react-draft-wysiwyg/dist/react-draft-wysiwyg.css";
import { useNavigate } from "react-router-dom";
import Tag from "../models/Tag";
import { useAppSelector } from "../redux/hooks";
import { HermesRequest, jsonConvert } from "../utility/request";

export default function CreateThread() {
  const titleRef = createRef<HTMLInputElement>()
  const contentRef = createRef<HTMLInputElement>()

  const [tags, setTags] = useState<Tag[]>([])
  const isLoggedIn = useAppSelector((state) => state.auth.value)
  const user = useAppSelector((state) => state.user.user)
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLoggedIn)
      navigate('/login', { state: { message: 'You must be logged in before you can create a new thread' } });

    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads/tags")
        .onSuccess((json) => {
          const tags = jsonConvert.deserializeArray(json, Tag)
          setTags(tags)
        })
        .call()
    })()
  }, [])

  async function createThread() {
    await new HermesRequest()
      .POST()
      .endpoint("threads")
      .hasAuthorization()
      .body({
        "user_id": user.id,
        "title": titleRef.current.value,
        "content": "loream ipsumloream ipsumloream ipsumloream ipsumloream ipsum",
        "tags": [],
      })
      .onSuccess(_ => navigate('/'))
      .onFailure(e => console.log(e.message))
      .onError(e => console.log(e))
      .call()
  }

  return (
    <div className="content">
      <h1 className="heading">New Thread</h1>
      <div className="form">
        <div className="field">
          <p>Thread title</p>
          <input type="text" name="title" id="title" ref={titleRef} />
        </div>

        <div className="field">
          <p>Thread content</p>
          <span><em>hermes</em> supports <a href="https://www.markdownguide.org/getting-started/">markdown</a> rendering!</span>
          <div className="editor">
            <Editor editorClassName="editor-field" toolbarClassName="editor-toolbar">
            </Editor>
          </div>
        </div>

        <div className="field">
          <p>Tags</p>
          <span>Select tags for this thread to be easily identified.</span>
          {tags && tags.map(tag => <p style={{ backgroundColor: tag.hexCode }}>{tag.content}</p>)}
        </div>

        <div className="buttons">
          <button type="button">Cancel</button>
          <button type="button" className="primary-button" onClick={async () => await createThread()}>Submit</button>
        </div>
      </div>
    </div>
  )
}