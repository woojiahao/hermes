import React from "react";
import Tag from "../models/Tag";

interface DisplayTagProps {
  tag: Tag
  onClick?: () => void
  style?: React.CSSProperties
}

export default function DisplayTag({tag, onClick, style}: DisplayTagProps) {
  return (
    <div className="display-tag" onClick={onClick}>
      <div style={{backgroundColor: tag.hexCode}}></div>
      <p style={style}>{tag.content}</p>
    </div>
  )
}
