import { create } from 'zustand'
import { menuService } from '@/services/menuService'
import {
    NAV_ITEM_TYPE_ITEM,
    NAV_ITEM_TYPE_COLLAPSE,
} from '@/constants/navigation.constant'

/**
 * HSE sidebar menu - purely from the database (menus + role_menu_access
 * tables), NOT dependent on src/configs/navigation.config.js at all.
 *
 * VerticalMenuContent (Ecme's built-in sidebar component) is still used
 * for rendering (already good: expand/collapse animation, active-menu
 * highlighting, accessible) - only the DATA SOURCE changes, from a
 * static file to this store.
 */

function toNavItem(node, parentKey = '') {
    const key = parentKey ? `${parentKey}.${node.id}` : `menu-${node.id}`
    const hasChildren = node.children && node.children.length > 0

    return {
        key,
        path: node.path || '',
        title: node.name,
        translateKey: '',
        icon: node.icon || '',
        type: hasChildren ? NAV_ITEM_TYPE_COLLAPSE : NAV_ITEM_TYPE_ITEM,
        // Empty because /api/menu ALREADY filters based on the logged-in
        // user's role (the backend decides this, not the frontend)
        authority: [],
        subMenu: hasChildren ? node.children.map((c) => toNavItem(c, key)) : [],
    }
}

export const useMenuStore = create((set) => ({
    items: [],
    loading: false,

    fetchMenu: async () => {
        set({ loading: true })
        try {
            const tree = await menuService.getMyMenu()
            const items = (tree || []).map((node) => toNavItem(node))
            set({ items, loading: false })
        } catch (err) {
            console.error('Gagal memuat menu:', err)
            set({ items: [], loading: false })
        }
    },

    clearMenu: () => set({ items: [] }),
}))
