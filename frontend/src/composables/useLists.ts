import api from '../api'
import { ref, watch } from 'vue'
import createAsyncProcess from '../utils/create-async-process'
import { ListInfo } from '../rpc/proto_gen'
import { user } from '../stores/user'

export function useLists() {
    const lists = ref<ListInfo[]>([])

    // TODO: prevent running multiple requests at the same time!!!
    async function fetchLists(): Promise<void> {
        const prepareLists = (fetched: ListInfo[]) => {
            return [
                {
                    ID: '$meta',
                    Name: 'Meta',
                    ListType: 'virtual',
                },
            ].concat(fetched)
        }
        lists.value = prepareLists([])

        if (user.isLoggedIn) {
            try {
                const fetched = await api.ListService.All()
                lists.value = prepareLists(fetched.Lists)
            } catch (e) {
                console.error('Failed to get all lists', e)
            }
        }
    }

    const { active: listsFetching, run: runWrappedFetchLists } =
        createAsyncProcess(fetchLists)

    watch(() => user.isLoggedIn, runWrappedFetchLists, { immediate: true })

    return {
        lists,
        listsFetching,
        runWrappedFetchLists,
    }
}
