import { reactive } from 'vue'
import api from '../api'
import { toast } from '../main'

export const user = reactive({
    isLoggedIn: false,
    email: '',
    fecthed: false,
})

export const logout = async () => {
    if (!user.isLoggedIn) {
        return
    }
    user.isLoggedIn = false
    user.email = ''
    user.fecthed = true
    toast.info("You're been logged out, please login again.")
}

export const updateUser = async () => {
    user.fecthed = false
    try {
        const status = await api.Auth.Status()
        user.isLoggedIn = true
        user.email = status.Email
    } catch (e) {
        // TODO: check if auth error
        console.log('auth status error', e)
        user.isLoggedIn = false
        user.email = ''
    }
    user.fecthed = true
}
