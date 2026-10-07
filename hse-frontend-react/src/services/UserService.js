import ApiService from './ApiService'

export async function apiListUsers() {
    return ApiService.fetchDataWithAxios({ url: '/users', method: 'get' })
}

export async function apiListRoles() {
    return ApiService.fetchDataWithAxios({ url: '/roles', method: 'get' })
}

export async function apiUpdateUserRole(id, roleId) {
    return ApiService.fetchDataWithAxios({
        url: `/users/${id}/role`,
        method: 'patch',
        data: { roleId },
    })
}

export async function apiUpdateUserStatus(id, status) {
    return ApiService.fetchDataWithAxios({
        url: `/users/${id}/status`,
        method: 'patch',
        data: { status },
    })
}

export async function apiGetAccessMatrix() {
    return ApiService.fetchDataWithAxios({
        url: '/admin/access-matrix',
        method: 'get',
    })
}

export async function apiUpdateAccessMatrix(access) {
    return ApiService.fetchDataWithAxios({
        url: '/admin/access-matrix',
        method: 'patch',
        data: { access },
    })
}