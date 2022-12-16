import React from "react"
import Layout from "../components/Layout"

interface ErrorProps {
  message: string
}

export default function Error({message}: ErrorProps) {
  return (
    <Layout>
      <div className="single">
        <div className="title">
          <h1 className="heading">Error</h1>
          <p>The page could not be found</p>
        </div>
      </div>
    </Layout>
  )
}
