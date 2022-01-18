import * as models from './models'
import { API } from './api'

const defaultLists: models.ListInfo[] = [
  {
    id: 'liked',
    name: 'Liked videos',
    type: 'external',
  },
  {
    id: 'history',
    name: 'History',
    type: 'external',
  },
  {
    id: 'view-later',
    name: 'View later',
    type: 'external',
  },
]

const listContent: Map<string, models.ListItem[]> = new Map()
listContent.set('history', [
  {
    youtubeID: '5ESJH1NLMLs',
    title: 'Children Of The Magenta Line',
    author: 'Mossie Fly',
    channelID: 'UCGIkFNbztHRaX0GB78SWaZA',
    xord: 'aba',
    itemID: '0',
  },
  {
    youtubeID: '68T9EFlCsUc',
    title: 'Making Music Is Easy',
    author: 'Eliminate',
    channelID: 'UCI7kKmUuSQOHUvSWIYFDf1Q',
    xord: 'acc',
    itemID: '1',
  },
])

async function someSleep(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 400))
}

export class LocalAPI extends API {
  constructor(private lists: models.ListInfo[] = defaultLists) {
    super()
  }

  async login(request: models.LoginRequest): Promise<models.LoginResponse> {
    await someSleep()
    return {
      ok: true,
    }
  }

  async getStatus(
    request: models.StatusRequest,
  ): Promise<models.StatusResponse> {
    await someSleep()
    return {
      ok: true,
    }
  }

  async getLists(request: models.ListsRequest): Promise<models.ListsResponse> {
    await someSleep()
    return {
      lists: this.lists,
    }
  }

  async getListInfo(listID: string): Promise<models.ListInfo> {
    await someSleep()
    const res = this.lists.find((list) => list.id === listID)
    if (!res) {
      throw new Error(`List ${listID} not found`)
    }
    return res
  }

  async getListItems(listID: string): Promise<models.ListItems> {
    await someSleep()
    const items = listContent.get(listID)
    if (!items) {
      throw new Error(`List ${listID} not found`)
    }
    return {
      items,
    }
  }
}
