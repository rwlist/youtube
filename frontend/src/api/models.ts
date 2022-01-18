export interface LoginRequest {
  someData?: string
}

export interface LoginResponse {
  ok: boolean
}

export interface StatusRequest {
  someData?: string
}

export interface StatusResponse {
  ok: boolean
}

export interface ListsRequest {
  someData?: string
}

export interface ListsResponse {
  lists: ListInfo[]
}

// external list: cannot be modified, only can be synced
// custom list: can be modified, cannot be synced
// virtual list: fully handled by the client
type ListType = 'external' | 'custom' | 'virtual'

export interface ListInfo {
  id: string
  name: string
  type: ListType
}

export interface ListItem {
  youtubeID: string
  title: string
  author: string
  channelID: string
  itemID: string
  xord: string
}

export interface ListItems {
  items: ListItem[]
}
