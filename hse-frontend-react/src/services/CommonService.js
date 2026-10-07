import ApiService from './ApiService'

export async function apiGetNotificationCount() {
    return ApiService.fetchDataWithAxios({
        url: '/notifications/unread-count',
        method: 'get',
    })
}

export async function apiGetNotificationList() {
    return ApiService.fetchDataWithAxios({
        url: '/notifications',
        method: 'get',
    })
}

export async function apiMarkNotificationRead(id) {
    return ApiService.fetchDataWithAxios({
        url: `/notifications/${id}/read`,
        method: 'patch',
    })
}

export async function apiMarkAllNotificationsRead() {
    return ApiService.fetchDataWithAxios({
        url: '/notifications/read-all',
        method: 'patch',
    })
}

export async function apiGetDashboardSummary() {
    return ApiService.fetchDataWithAxios({
        url: '/dashboard/summary',
        method: 'get',
    })
}

export async function apiGetDashboardTrend(months = 6) {
    return ApiService.fetchDataWithAxios({
        url: '/dashboard/trend',
        method: 'get',
        params: { months },
    })
}

export async function apiGetDashboardAvgApprovalTime() {
    return ApiService.fetchDataWithAxios({
        url: '/dashboard/avg-approval-time',
        method: 'get',
    })
}

export async function apiGetSearchResult(params) {
    return ApiService.fetchDataWithAxios({
        url: '/search/query',
        method: 'get',
        params,
    })
}