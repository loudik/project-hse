import ApiService from './ApiService'

export async function apiGetMyOrganization() {
    return ApiService.fetchDataWithAxios({
        url: '/organizations/mine',
        method: 'get',
    })
}

export async function apiCreateOrganization(data) {
    return ApiService.fetchDataWithAxios({
        url: '/organizations',
        method: 'post',
        data,
    })
}

export async function apiListOrganizations(status) {
    return ApiService.fetchDataWithAxios({
        url: status ? `/organizations?status=${encodeURIComponent(status)}` : '/organizations',
        method: 'get',
    })
}

export async function apiDecideOrganization(id, data) {
    return ApiService.fetchDataWithAxios({
        url: `/organizations/${id}/decision`,
        method: 'patch',
        data,
    })
}
