import React from "react"
import {Editor} from "react-draft-wysiwyg"
import {ContentState} from "draft-js"

interface RichTextEditorProps {
  // @ts-ignore
  contentState: ContentState
  // @ts-ignore
  setContentState: React.Dispatch<React.SetStateAction<ContentState>>
}

export default function RichTextEditor({contentState, setContentState}: RichTextEditorProps) {
  return (
    <div className="rich-text-editor">
      <Editor
        editorClassName="rich-text-editor-field"
        toolbarClassName="rich-text-editor-toolbar"
        defaultContentState={contentState}
        onContentStateChange={setContentState}
        toolbar={{
          options: ['inline', 'blockType', 'list', 'link', 'emoji', 'image'],
          blockType: {
            options: ['Normal', 'H1', 'H2', 'H3']
          }
        }}>
      </Editor>
    </div>
  )
}
