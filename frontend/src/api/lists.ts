import api from '../api'
import { ListInfo, ListItem } from './models'
import {
  readonly,
  reactive,
  DeepReadonly,
  UnwrapNestedRefs,
  watch,
  ref,
  Ref,
  computed,
  watchEffect,
} from 'vue'
import { API } from './api'
import createAsyncProcess from '../utils/create-async-process'

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

export class MetaList implements ListCtl {
  private items: UnwrapNestedRefs<ListItem[]> = reactive([])
  private status = ref('META')
  private display: UnwrapNestedRefs<ListStatus> = reactive({
    header: '',
    response: 'Enter a query and press enter.',
  })

  constructor() {
    watch(
      () => 'Status: ' + this.status.value,
      header => this.display.header = header,
      { immediate: true },
    )
  }

  allItemsRef(): DeepReadonly<UnwrapNestedRefs<ListItem[]>> {
    return readonly(this.items)
  }

  statusRef(): DeepReadonly<UnwrapNestedRefs<ListStatus>> {
    return readonly(this.display)
  }

  executeQuery(query: string): void {
    query = query.trim()

    switch (query) {
      case '$help':
        this.queryHelp()
        break
      case '$login':
        this.queryLogin()
        break
      default:
        this.queryHelp()
    }
  }

  private queryHelp() {
    this.status.value = 'HELP'
    this.display.response = 
`$help: show this help
$login: login`
  }

  private async queryLogin() {
    this.status.value = 'FETCHING'

    try {
      const res = await api.login({})
      this.status.value = 'LOGIN'
      document.location.href = res.RedirectURL
      // unreachable
    } catch (err) {
      this.status.value = 'ERROR'
      this.display.response = `Error: ${err.message}`
    }
  }
}

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
      this.fetchProcess = createAsyncProcess(async () => this.fetchItems())
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
      const listItems = await this.api.getListItems(this.listInfo.id)
      this.items.splice(0, this.items.length, ...listItems.items)
      this.status.header = 'Status: READY'
      this.status.response = `fetched ${listItems.items.length} items`
    } catch (e) {
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

export async function fetchListByID(listID: string): Promise<ListCtl> {
  if (listID == '$meta') {
    return new MetaList()
  }

  const listInfo = await api.getListInfo(listID)
  return new RemoteList(api, listInfo)
}
