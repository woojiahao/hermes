import React from "react";
import Layout from "../components/Layout"
import hermes from "../resources/hermes.jpeg"

export default function About() {
  return (
    <Layout>
      <div className="single">
        <div className="title">
          <h1 className="heading">About</h1>
        </div>
        <div className="flex gap-x-8">
          <img className="w-1/4 aspect-auto rounded-br shadow-bs drop-shadow-ds"
               src={hermes}
               alt="Hermes, messenger of the Greek Gods"/>
          <div className="flex flex-col gap-y-4 card">
            <h2>About <em className="text-2xl">hermes</em></h2>
            <p><em>hermes</em> is a web forum built with React, Golang, and PostgreSQL as part of the CVWO Assignment
              2022.
            </p>
            <p>The goal of <em>hermes</em> is to provide a simple and easy to use forum for users to start threads and
              create
              discussions.</p>
            <p>The name <em>hermes</em> is inspired by the messenger of the Greek Gods.</p>
            <p>Check out the Github repository for <em>hermes</em> <a href="https://github.com/woojiahao/hermes">here!</a></p>
          </div>
        </div>
      </div>
    </Layout>
  )
}
