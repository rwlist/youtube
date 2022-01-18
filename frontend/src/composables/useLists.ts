import type { ListInfo } from '../api/models'
import api from '../api'
import { ref, watch } from 'vue'
import createAsyncProcess from '../utils/create-async-process'

export function useLists() {
  const lists = ref<ListInfo[]>([])

  // TODO: prevent running multiple requests at the same time???
  async function fetchLists(): Promise<void> {
    lists.value = [
      {
        id: '$meta',
        name: 'Meta',
        type: 'virtual',
      },
    ]

    // TODO: check if not logged in
    const fetched = await api.getLists({})
    lists.value = lists.value.concat(fetched.lists)
  }

  const { active: listsFetching, run: runWrappedFetchLists } =
    createAsyncProcess(fetchLists)

  // TODO: bind watch to something useful
  watch(() => null, runWrappedFetchLists, { immediate: true })

  return {
    lists,
    listsFetching,
    runWrappedFetchLists,
  }
}
