import * as models from './models'
import { API } from './api'
import { Auth } from '../rpc/proto_gen'

export class RPCAPI extends API {
  constructor(private def: API, private auth: Auth) {
    super()
  }

  async login(request: models.LoginRequest): Promise<models.LoginResponse> {
    return await this.auth.Oauth()
  }

  getStatus(
    request: models.StatusRequest,
  ): Promise<models.StatusResponse> {
    return this.def.getStatus(request)
  }

  getLists(request: models.ListsRequest): Promise<models.ListsResponse> {
    return this.def.getLists(request)
  }

  getListInfo(listID: string): Promise<models.ListInfo> {
    return this.def.getListInfo(listID)
  }

  getListItems(listID: string): Promise<models.ListItems> {
    return this.def.getListItems(listID)
  }
}
