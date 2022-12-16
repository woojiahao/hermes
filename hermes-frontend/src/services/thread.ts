import Thread from "../models/Thread"
import {errorMessage, HermesRequest, jsonConvert} from "../utility/request"

export async function asyncGetThread(
  id: string,
  onSuccess: (thread: Thread) => void,
  onFailure: (e: errorMessage) => void,
  onError: (e: Error) => void
) {
  await new HermesRequest()
    .GET()
    .endpoint(`/threads/${id}`)
    .onSuccess(json => {
      const t = jsonConvert.deserializeObject(json, Thread)
      onSuccess(t)
    })
    .onFailure(onFailure)
    .onError(onError)
    .call()
}
