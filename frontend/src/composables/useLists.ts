import api from '../api'
import { ref, watch } from 'vue'
import createAsyncProcess from '../utils/create-async-process'
import { ListInfo } from '../rpc/proto_gen'

export function useLists() {
    const lists = ref<ListInfo[]>([])

    // TODO: prevent running multiple requests at the same time???
    async function fetchLists(): Promise<void> {
        lists.value = [
            {
                ID: '$meta',
                Name: 'Meta',
                ListType: 'virtual',
            },
        ]

        // TODO: check if not logged in before request?
        try {
            const fetched = await api.ListService.All()
            lists.value = lists.value.concat(fetched.Lists)
        } catch (e) {
            console.error('Failed to get all lists', e)
        }
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
