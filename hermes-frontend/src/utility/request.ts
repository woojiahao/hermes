/**
 * Utilities for making API requests
 */

import { getJWT, refreshJWT } from "./jwt";

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


export enum RequestType {
  GET, POST, PUT, DELETE
}

export class HermesRequest {
  private _requestType: RequestType = RequestType.GET
  private _endpoint: string
  private _queryParams: queryParams = {}
  private _pathParams: pathParams = []
  private _body: any = {}
  private _hasAuthorization: boolean = false
  private _onSuccess: apiCallback<any>
  private _onFailure: apiCallback<{ message: string }>
  private _onError: apiCallback<Error> = defaultOnError

  GET(): HermesRequest {
    this._requestType = RequestType.GET
    return this
  }

  POST(): HermesRequest {
    this._requestType = RequestType.POST
    return this
  }

  PUT(): HermesRequest {
    this._requestType = RequestType.PUT
    return this
  }

  DELETE(): HermesRequest {
    this._requestType = RequestType.DELETE
    return this
  }

  requestType(rt: RequestType): HermesRequest {
    this._requestType = rt
    return this
  }

  endpoint(e: string): HermesRequest {
    this._endpoint = e
    return this
  }

  queryParams(qp: queryParams): HermesRequest {
    this._queryParams = qp
    return this
  }

  pathParams(pp: pathParams): HermesRequest {
    this._pathParams = pp
    return this
  }

  body(b: any): HermesRequest {
    this._body = b
    return this
  }

  hasAuthorization(): HermesRequest {
    this._hasAuthorization = true
    return this
  }

  onSuccess(os: apiCallback<any>): HermesRequest {
    this._onSuccess = os
    return this
  }

  onFailure(of: apiCallback<{ message: string }>): HermesRequest {
    this._onFailure = of
    return this
  }

  onError(oe: apiCallback<Error>): HermesRequest {
    this._onError = oe
    return this
  }

  private bearerToken(): string {
    return `Bearer ${getJWT()}`;
  }

  private async handleRequest(url: string, config: RequestInit) {
    try {
      const result = await fetch(url, config)

      const json = await result.json()

      if (result.ok) {
        this._onSuccess(json)
        return
      }

      if (!this._hasAuthorization) {
        this._onFailure(json)
        return
      }

      const err = json as { message: string }
      const tryRefresh = result.status === 401 && err.message === "Token is expired"

      if (!tryRefresh) {
        this._onFailure(err)
      }

      const refreshed = await refreshJWT()
      if (!refreshed) {
        // If refreshing failed, this means that the refresh token has been expired and
        // the user MUST login again
        this._onFailure(err)
        return
      }

      const headers: HeadersInit = {}
      Object.assign(headers, config.headers)
      headers['Authorization'] = this.bearerToken()

      const copyConfig: RequestInit = {}
      Object.assign(copyConfig, config)
      copyConfig.headers = headers
      this.handleRequest(url, copyConfig)
    } catch (e) {
      this._onError(e as Error)
    }
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

    if (this._hasAuthorization) {
      const jwtToken = getJWT()
      if (!jwtToken) {
        this._onError(new Error("JWT Token does not exist"))
        return
      }
      headers['Authorization'] = this.bearerToken();
    }

    config.headers = headers

    const url = createURL(this._endpoint, this._queryParams, this._pathParams)
    await this.handleRequest(url, config)
  }
}