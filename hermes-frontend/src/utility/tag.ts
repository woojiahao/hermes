import Tag from "../models/Tag"
import React from "react"

export function isLight(hexCode: string): boolean {
  const code = hexCode.replace('#', '')
  const r = parseInt(code.substr(0, 2), 16)
  const g = parseInt(code.substr(2, 2), 16)
  const b = parseInt(code.substr(4, 2), 16)
  const brightness = ((r * 299) + (g * 587) + (b * 114)) / 1000
  return brightness > 155
}

export function tagStyle(tag: Tag): React.CSSProperties {
  return {
    backgroundColor: tag.hexCode,
    color: isLight(tag.hexCode) ? `var(--dark-primary)` : `var(--background-primary)`
  }
}
