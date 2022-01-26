import { reactive } from 'vue'
import { authService } from '../api' 

export const user = reactive({
    isLoggedIn: false,
    email: '',
    fecthed: false,
})

const updateUser = async () => {
    user.fecthed = false
    try {
        const status = await authService.Status()
        user.isLoggedIn = true
        user.email = status.Email
    } catch (e) {
        // TODO: check if auth error
        console.log(e)
        user.isLoggedIn = false
        user.email = ''
    }
    user.fecthed = true
}

updateUser()
