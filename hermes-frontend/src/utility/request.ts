/**
 * Utilities for making API requests
 */

import { JsonConvert } from "json2typescript";
import urlJoin from 'url-join';
import { getJWT, refreshJWT } from "./jwt";

export const jsonConvert = new JsonConvert()

function createURL(
  endpoint: string,
  queryParams: queryParams = {},
  pathParams: pathParams = []
) {
  let url: URL
  if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
    url = new URL(endpoint, 'http://localhost:8080')
  } else {
    url = new URL(urlJoin('https://hermes.woojiahao.com', 'api', endpoint))
  }

  if (queryParams)
    for (const key in queryParams) url.searchParams.append(key, queryParams[key])

  if (pathParams)
    for (const path of pathParams) url = new URL(path, url.toString())

  return url.toString();
}

export type apiCallback<T> = (json: T) => void
export type queryParams = { [key: string]: string }
export type pathParams = string[]
export type errorMessage = { message: string }
export type errorFields = { message: { field: string, message: string }[] }

const defaultOnError: apiCallback<Error> = (err) => console.log(err)

const defaultOnFailure: apiCallback<errorMessage | errorFields> = err => console.log(err)

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
  private _onFailure: apiCallback<errorMessage | errorFields> = defaultOnFailure
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

  onFailure(of: apiCallback<errorMessage | errorFields>): HermesRequest {
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

      if (result.status === 204) {
        this._onSuccess("")
        return
      }

      const json = await result.json()

      if (result.ok) {
        this._onSuccess(json)
        return
      }

      if (!this._hasAuthorization) {
        this._onFailure(json)
        return
      }

      if (typeof json.message !== 'string') {
        this._onFailure(json)
        return
      }

      const err = json as { message: string }
      const tryRefresh = result.status === 401 && err.message === "Token is expired"

      if (!tryRefresh) {
        this._onFailure(err)
        return
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
      await this.handleRequest(url, copyConfig)
    } catch (e) {
      this._onError(e as Error)
    }
  }

  async call() {
    if (!this._endpoint)
      throw new Error("Specify endpoint for Hermes request")

    const headers: HeadersInit = {}
    const config: RequestInit = {
      method: RequestType[this._requestType],
    }

    if ([RequestType.POST, RequestType.PUT].includes(this._requestType)) {
      config.body = JSON.stringify(this._body)
      headers['Content-Type'] = 'application/json'
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
