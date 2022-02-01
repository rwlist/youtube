import api from '.'
import {
    readonly,
    reactive,
    DeepReadonly,
    UnwrapNestedRefs,
    watch,
    ref,
} from 'vue'
import { ListItem } from '../rpc/proto_gen'
import { RpcError } from '../rpc/http'
import { ListCtl, ListStatus } from './lists'
import { logout } from '../stores/user'

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
            (header) => (this.display.header = header),
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
            case '$logout':
                this.queryLogout()
                break
            default:
                this.queryHelp()
        }
    }

    private queryHelp() {
        this.status.value = 'HELP'
        this.display.response = `$help: show this help
$login: login
$logout: logout`
    }

    private async queryLogin() {
        this.status.value = 'FETCHING'

        try {
            const res = await api.Auth.Oauth()
            this.status.value = 'LOGIN'
            document.location.href = res.RedirectURL
            // unreachable
        } catch (e) {
            console.error(e)
            if (e instanceof RpcError) {
                console.log(e)
            } else {
                throw e
            }
            this.status.value = 'ERROR'
            this.display.response = `Error: ${e.message}`
        }
    }

    private async queryLogout() {
        this.status.value = 'THINKING'

        try {
            await logout()
            this.status.value = 'LOGOUT'
            this.display.response = 'Logged out.'
        } catch (e) {
            console.error(e)
            this.status.value = 'ERROR'
            this.display.response = `Cannot log out`
        }
    }
}
