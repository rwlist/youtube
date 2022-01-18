import type { ListItem } from '../api/models'
import { Ref, ref, watch } from 'vue'
import { fetchListByID, ListCtl, ListStatus } from '../api/lists'

export function useListCtl(listID: Ref<string>) {
  const status = ref<ListStatus>({
    header: '',
    response: '',
  })
  const allItems = ref<ReadonlyArray<ListItem>>([])

  let listCtl: ListCtl | undefined = undefined

  watch(
    () => listID.value,
    async (newListID) => {
      status.value = {
        header: '',
        response: '',
      }
      allItems.value = []

      if (newListID === '') {
        listCtl = undefined
        return
      }

      listCtl = await fetchListByID(newListID)
      status.value = listCtl.statusRef()
      allItems.value = listCtl.allItemsRef()
    },
  )

  const query = ref('')
  const executeQuery = () => listCtl?.executeQuery(query.value)

  return {
    status,
    allItems,
    query,
    executeQuery,
  }
}
