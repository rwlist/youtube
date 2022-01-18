import * as models from './models'

export abstract class API {
    // API without authentication
    abstract login(request: models.LoginRequest): Promise<models.LoginResponse>;

    // API with authentication
    abstract getStatus(request: models.StatusRequest): Promise<models.StatusResponse>;
    abstract getLists(request: models.ListsRequest): Promise<models.ListsResponse>;
    abstract getListInfo(listID: string): Promise<models.ListInfo>;
    abstract getListItems(listID: string): Promise<models.ListItems>;
}
