/**
 * Utilities for making API requests
 */

import { getJWT } from "./app";

// TODO: Change this if necessary
const apiURL = 'http://localhost:8081'

function createURL(
  endpoint: string,
  queryParams: queryParams = {},
  pathParams: pathParams = []
) {
  let url = new URL(endpoint, apiURL);

  if (queryParams)
    for (const key in queryParams) url.searchParams.append(key, queryParams[key])

  if (pathParams)
    for (const path of pathParams) url = new URL(path, url.toString())

  return url.toString();
}

export type apiCallback<T> = (json: T) => void
export type queryParams = { [key: string]: string }
export type pathParams = string[]

const defaultOnError: apiCallback<Error> = (err) => console.log(err)

async function handleRequest(
  url: string,
  config: RequestInit,
  onSuccess: apiCallback<any>,
  onFailure: apiCallback<any>,
  onError: apiCallback<Error>
) {
  try {
    const result = await fetch(url, config)

    const json = await result.json()

    if (!result.ok) onFailure(json)
    else onSuccess(json)
  } catch (e) {
    onError(e as Error)
  }
}

export enum RequestType {
  GET, POST, PUT, DELETE
}

export class Request {
  private _requestType: RequestType = RequestType.GET
  private _endpoint: string
  private _queryParams: queryParams = {}
  private _pathParams: pathParams = []
  private _body: any = {}
  private _hasAuthorization: boolean = false
  private _onSuccess: apiCallback<any>
  private _onFailure: apiCallback<any>
  private _onError: apiCallback<Error> = defaultOnError

  constructor() {

  }

  requestType(rt: RequestType): Request {
    this._requestType = rt
    return this
  }

  endpoint(e: string): Request {
    this._endpoint = e
    return this
  }

  queryParams(qp: queryParams): Request {
    this._queryParams = qp
    return this
  }

  pathParams(pp: pathParams): Request {
    this._pathParams = pp
    return this
  }

  body(b: any): Request {
    this._body = b
    return this
  }

  hasAuthorization(ha: boolean): Request {
    this._hasAuthorization = ha
    return this
  }

  onSuccess(os: apiCallback<any>): Request {
    this._onSuccess = os
    return this
  }

  onFailure(of: apiCallback<any>): Request {
    this._onFailure = of
    return this
  }

  onError(oe: apiCallback<Error>): Request {
    this._onError = oe
    return this
  }

  async call() {
    const headers: HeadersInit = {}
    const config: RequestInit = {
      method: RequestType[this._requestType],
    }

    switch (this._requestType) {
      case RequestType.GET:
        break
      case RequestType.POST:
        config.body = JSON.stringify(this._body)
        headers['Content-Type'] = 'application/json'
        break
      case RequestType.PUT:
        throw new Error("not supported yet: PUT")
        break
      case RequestType.DELETE:
        throw new Error("not supported yet: DELETE")
        break
    }

    if (this._hasAuthorization) headers['Authorization:Bearer'] = getJWT()

    const url = createURL(this._endpoint, this._queryParams, this._pathParams)
    await handleRequest(url, config, this._onSuccess, this._onFailure, this._onError)
  }
}