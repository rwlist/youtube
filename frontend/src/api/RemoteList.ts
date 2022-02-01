import {
    readonly,
    reactive,
    DeepReadonly,
    UnwrapNestedRefs,
    watch,
    ref,
    Ref,
} from 'vue'
import createAsyncProcess from '../utils/create-async-process'
import { API, ListInfo, ListItem } from '../rpc/proto_gen'
import { RpcError } from '../rpc/http'
import { ListCtl, ListStatus } from './lists'

export class RemoteList implements ListCtl {
    constructor(private api: API, private listInfo: ListInfo) {}

    private fetchProcess: ReturnType<typeof createAsyncProcess> | null = null
    private items: UnwrapNestedRefs<ListItem[]> = reactive([])
    private status: UnwrapNestedRefs<ListStatus> = reactive({
        header: 'Status: READY',
        response: 'waiting for query...',
    })
    private listUpdater: Ref<number> = ref(0)

    allItemsRef(): DeepReadonly<UnwrapNestedRefs<ListItem[]>> {
        if (this.fetchProcess == null) {
            this.fetchProcess = createAsyncProcess(async () =>
                this.fetchItems(),
            )
            watch(
                () => this.listUpdater.value,
                () => this.fetchProcess?.run(),
                { immediate: true },
            )
        }

        return readonly(this.items)
    }

    statusRef(): DeepReadonly<UnwrapNestedRefs<ListStatus>> {
        return readonly(this.status)
    }

    async fetchItems(): Promise<void> {
        this.status.header = 'Status: FETCHING'
        this.status.response = 'fetching...'

        try {
            const listItems = await this.api.ListService.ListItems(
                this.listInfo.ID,
            )
            this.items.splice(0, this.items.length, ...listItems.Items)
            this.status.header = 'Status: READY'
            this.status.response = `fetched ${listItems.Items.length} items`
        } catch (e) {
            if (e instanceof RpcError) {
                console.log(e)
            } else {
                throw e
            }
            this.status.header = 'Status: ERROR'
            this.status.response = `error: ${e.message}`
        }
    }

    executeQuery(query: string): void {
        query = query.trim()

        // TODO: implement queries
        console.log(`query: ${query}`)
    }
}
