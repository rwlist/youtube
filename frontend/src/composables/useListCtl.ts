import { Ref, ref, watch } from 'vue'
import { fetchListByID, ListCtl, ListStatus, PagedList } from '../api/lists'
import { ListItem } from '../rpc/proto_gen'

export function useListCtl(listID: Ref<string>) {
    const status = ref<ListStatus>({
        header: '',
        response: '',
    })
    const allItems = ref<ReadonlyArray<ListItem>>([])

    const supportsPages = ref(false)

    let listCtl: ListCtl | undefined = undefined

    watch(
        () => listID.value,
        async (newListID) => {
            status.value = {
                header: 'Loading list...',
                response: '',
            }
            allItems.value = []

            if (newListID === '') {
                listCtl = undefined
                supportsPages.value = false
                return
            }

            listCtl = await fetchListByID(newListID)
            status.value = listCtl.statusRef()

            const matrix = listCtl.supportMatrix()
            if (matrix.pagedList) {
                supportsPages.value = true
                allItems.value = []
            } else {
                supportsPages.value = false
                allItems.value = listCtl.allItemsRef()
            }
        },
        { immediate: true },
    )

    const query = ref('')
    const executeQuery = () => listCtl?.executeQuery(query.value)

    const getListCtl = (): ListCtl => listCtl!
    const getPagedList = (): PagedList => {
        const ctl = getListCtl()
        const matrix = ctl.supportMatrix()
        if (!matrix.pagedList) {
            throw new Error('List does not support pagination')
        }
        return matrix.pagedList
    }

    return {
        status,
        allItems,
        query,
        executeQuery,
        supportsPages,
        getListCtl,
        getPagedList,
    }
}
