import api from '../api'
import { DeepReadonly, UnwrapNestedRefs } from 'vue'
import { ListItem } from '../rpc/proto_gen'
import { MetaList } from './MetaList'
import { RemoteList } from './RemoteList'

export interface ListStatus {
    header: string
    response: string
}

export interface ListCtl {
    // Return full list content, updating it reactively.
    // Usually initialized on the first call.
    allItemsRef(): DeepReadonly<UnwrapNestedRefs<ListItem[]>>

    // Return reactive object with current fetch status and
    // query response.
    statusRef(): DeepReadonly<UnwrapNestedRefs<ListStatus>>

    // TODO: return promise?
    executeQuery(query: string): void
}

export async function fetchListByID(listID: string): Promise<ListCtl> {
    if (listID == '$meta') {
        return new MetaList()
    }

    const listInfo = await api.ListService.Info(listID)
    return new RemoteList(api, listInfo)
}
