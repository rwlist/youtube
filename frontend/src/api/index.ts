import { HTTPTransport } from "../rpc/http"
import { Auth } from "../rpc/proto_gen"
import { LocalAPI } from "./api-mock"
import { RPCAPI } from "./api-rpc"

const rpcTransport = new HTTPTransport("/rpc")
export const authService = new Auth(rpcTransport) // TODO: don't export
const localAPI = new LocalAPI()
const api = new RPCAPI(localAPI, authService)

export default api
