import ApiService from './ApiService'
import endpointConfig from '@/configs/endpoint.config'

async function placeholderFunction() {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve({
                token: 'placeholder_token',
                user: {
                    id: 'placeholder_id',
                    name: 'Placeholder User',
                    email: 'user@example.com',
                },
            })
        }, 500)
    })
}

export async function apiGoogleOauthSignIn() {
    return await placeholderFunction()
}

export async function apiMicrosoftOauthSignIn(code) {
    return ApiService.fetchDataWithAxios({
        url: endpointConfig.microsoftOauthSignIn,
        method: 'post',
        data: { code },
    })
}
