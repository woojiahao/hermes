/**
 * Utilities for making API requests
 */

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

async function handleRequest<T, E>(
  url: string,
  config: RequestInit,
  onSuccess: apiCallback<T>,
  onFailure: apiCallback<E>,
  onError: apiCallback<Error>
) {
  try {
    const result = await fetch(url, config)

    const json = await result.json()

    if (!result.ok) onFailure(json as E)
    else onSuccess(json as T)
  } catch (e) {
    onError(e as Error)
  }
}

export async function partialGET<T, E>(
  endpoint: string,
  onSuccess: apiCallback<T>,
  onFailure: apiCallback<E>,
  onError: apiCallback<Error> = defaultOnError
) {
  await GET(endpoint, {}, [], onSuccess, onFailure, onError)
}

export async function GET<T, E>(
  endpoint: string,
  queryParams: queryParams,
  pathParams: pathParams,
  onSuccess: apiCallback<T>,
  onFailure: apiCallback<E>,
  onError: apiCallback<Error> = defaultOnError
) {
  const url = createURL(endpoint, queryParams, pathParams)
  await handleRequest(url, { method: "GET" }, onSuccess, onFailure, onError)
}

export async function partialPOST<B, T, E>(
  endpoint: string,
  body: B,
  onSuccess: apiCallback<T>,
  onFailure: apiCallback<E>,
  onError: apiCallback<Error> = defaultOnError
) {
  await POST(endpoint, {}, [], body, onSuccess, onFailure, onError)
}

export async function POST<B, T, E>(
  endpoint: string,
  queryParams: queryParams,
  pathParams: pathParams,
  body: B,
  onSuccess: apiCallback<T>,
  onFailure: apiCallback<E>,
  onError: apiCallback<Error> = defaultOnError
) {
  const url = createURL(endpoint, queryParams, pathParams)
  await handleRequest(
    url,
    {
      method: "POST",
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    },
    onSuccess,
    onFailure,
    onError
  )
}