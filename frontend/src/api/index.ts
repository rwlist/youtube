import { HTTPTransport } from '../rpc/http'
import { buildImpl } from '../rpc/proto_gen'

const rpcTransport = new HTTPTransport('/rpc')
const api = buildImpl(rpcTransport)

export default api
