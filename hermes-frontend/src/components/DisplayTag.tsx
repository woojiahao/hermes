import React from "react";
import Tag from "../models/Tag";

interface DisplayTagProps {
  tag: Tag
  onClick?: () => void
  style?: React.CSSProperties
}

export default function DisplayTag({tag, onClick, style}: DisplayTagProps) {
  return (
    <div className="flex items-center gap-2 hover:cursor-pointer" onClick={onClick}>
      <div className="w-[15px] h-[15px] border border-dark" style={{backgroundColor: tag.hexCode}}></div>
      <p className="break-words" style={style}>{tag.content}</p>
    </div>
  )
}
