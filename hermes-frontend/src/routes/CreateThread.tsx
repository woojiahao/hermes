import React, {createRef, useEffect, useState} from "react";
import {Editor} from 'react-draft-wysiwyg';
import "react-draft-wysiwyg/dist/react-draft-wysiwyg.css";
import {useNavigate} from "react-router-dom";
import {useAppSelector} from "../redux/hooks";
import {HermesRequest} from "../utility/request";
import TagSelection from "../components/TagSelection"
import {ContentState, convertToRaw} from "draft-js"
import {draftToMarkdown} from "markdown-draft-js"
import Tag from "../models/Tag"

export default function CreateThread() {
  const titleRef = createRef<HTMLInputElement>()
  const [contentState, setContentState] = useState(convertToRaw(ContentState.createFromText("")))
  const [selectedTags, setSelectedTags] = useState<Map<number, Tag>>(new Map())

  const isLoggedIn = useAppSelector((state) => state.auth.value)
  const user = useAppSelector((state) => state.user.user)
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLoggedIn)
      navigate('/login', {state: {message: 'You must be logged in before you can create a new thread'}})
  }, [])

  async function createThread() {
    await new HermesRequest()
      .POST()
      .endpoint("threads")
      .hasAuthorization()
      .body({
        "user_id": user.id,
        "title": titleRef.current.value,
        "content": draftToMarkdown(contentState),
        "tags": Array.from(selectedTags.values()),
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
          <input type="text" name="title" id="title" ref={titleRef}/>
        </div>

        <div className="field">
          <p>Thread content</p>
          <span><em>hermes</em> supports rich text editing!</span>
          <div className="editor">
            <Editor
              editorClassName="editor-field"
              toolbarClassName="editor-toolbar"
              defaultContentState={contentState}
              onContentStateChange={setContentState}>
            </Editor>
          </div>
        </div>

        <div className="field">
          <p>Tags</p>
          <span>Select tags for this thread to be easily identified.</span>
          <TagSelection selectedTags={selectedTags} setSelectedTags={setSelectedTags}/>
        </div>

        <div className="buttons">
          <button type="button">Cancel</button>
          <button type="button" className="primary-button" onClick={async () => await createThread()}>Submit</button>
        </div>
      </div>
    </div>
  )
}
