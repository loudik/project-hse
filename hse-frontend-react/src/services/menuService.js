import ApiService from '@/services/ApiService'

export const menuService = {
    getMyMenu: () =>
        ApiService.fetchDataWithAxios({
            url: '/menu',
            method: 'get',
        }),
}
