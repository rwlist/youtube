import api from '../api'
import { DeepReadonly, UnwrapNestedRefs } from 'vue'
import { ListInfo, ListItem } from '../rpc/proto_gen'
import { MetaList } from './MetaList'
import { RemoteList } from './RemoteList'

export interface ListStatus {
    header: string
    response: string
}

export interface SupportMatrix {
    // returns the instance if pagination is supported
    pagedList?: PagedList
}

export interface ListCtl {
    // Return full list content, updating it reactively.
    // Usually initialized on the first call.
    allItemsRef(): DeepReadonly<UnwrapNestedRefs<ListItem[]>>

    // Return object with supported features.
    supportMatrix(): SupportMatrix

    // Return reactive object with current fetch status and
    // query response.
    statusRef(): DeepReadonly<UnwrapNestedRefs<ListStatus>>

    // TODO: write description
    executeQuery(query: string): Promise<void>
}

export interface PagedList {
    getInfo(): ListInfo

    // TODO: better than offset/limit?
    fetchPage(offset: number, limit: number): Promise<ListItem[]>
}

export async function fetchListByID(listID: string): Promise<ListCtl> {
    if (listID == '$meta') {
        return new MetaList()
    }

    const listInfo = await api.ListService.Info(listID)
    return new RemoteList(api, listInfo)
}
