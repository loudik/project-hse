import { lazy } from 'react'
import authRoute from './authRoute'
import othersRoute from './othersRoute'

export const publicRoutes = [...authRoute]

export const protectedRoutes = [
    {
        key: 'home',
        path: '/home',
        component: lazy(() => import('@/views/Home')),
        authority: [],
    },
    {
        key: 'organizationProfile',
        path: '/organization/profile',
        component: lazy(() => import('@/views/organization/OrganizationProfile')),
        authority: [],
    },
    {
        key: 'organizationReview',
        path: '/organization/review',
        component: lazy(() => import('@/views/organization/OrganizationReview')),
        authority: [],
    },
    {
        key: 'userManagement',
        path: '/admin/users',
        component: lazy(() => import('@/views/admin/UserManagement')), 
        authority: [],
    },
    {
        key: 'roleMenuManagement',
        path: '/admin/access',
        component: lazy(() => import('@/views/admin/RoleMenuManagement')),
        authority: [],
    },
    /** Example purpose only, please remove */
    {
        key: 'singleMenuItem',
        path: '/single-menu-view',
        component: lazy(() => import('@/views/demo/SingleMenuView')),
        authority: [],
    },
    {
        key: 'collapseMenu.item1',
        path: '/collapse-menu-item-view-1',
        component: lazy(() => import('@/views/demo/CollapseMenuItemView1')),
        authority: [],
    },
    {
        key: 'collapseMenu.item2',
        path: '/collapse-menu-item-view-2',
        component: lazy(() => import('@/views/demo/CollapseMenuItemView2')),
        authority: [],
    },
    {
        key: 'groupMenu.single',
        path: '/group-single-menu-item-view',
        component: lazy(() => import('@/views/demo/GroupSingleMenuItemView')),
        authority: [],
    },
    {
        key: 'groupMenu.collapse.item1',
        path: '/group-collapse-menu-item-view-1',
        component: lazy(
            () => import('@/views/demo/GroupCollapseMenuItemView1'),
        ),
        authority: [],
    },
    {
        key: 'groupMenu.collapse.item2',
        path: '/group-collapse-menu-item-view-2',
        component: lazy(
            () => import('@/views/demo/GroupCollapseMenuItemView2'),
        ),
        authority: [],
    },
     {
        key: 'vesselStart',
        path: '/vessel',
        component: lazy(() => import('@/views/vessel/VesselStart')),
        authority: [],
    },
    {
        key: 'vesselList',
        path: '/vessel/list',
        component: lazy(() => import('@/views/vessel/VesselList')),
        authority: [],
    },
    {
        key: 'vesselApply',
        path: '/vessel/apply/:id',
        component: lazy(() => import('@/views/vessel/VesselApply')),
        authority: [],
    },
    {
        key: 'vesselSummary',
        path: '/vessel/summary/:id',
        component: lazy(() => import('@/views/vessel/VesselSummary')),
        authority: [],
    },
    {
        key: 'vesselReview',
        path: '/vessel/review',
        component: lazy(() => import('@/views/vessel/VesselReview')),
        authority: [],
    },
    {
        key: 'vesselReviewDetail',
        path: '/vessel/review/:id',
        component: lazy(() => import('@/views/vessel/VesselSummary')),
        authority: [],
    },
        {
        key: 'extensionRequests',
        path: '/vessel/extension-requests',
        component: lazy(() => import('@/views/vessel/ExtensionRequests')),
        authority: [],
    },
     {
        key: 'dashboard',
        path: '/dashboard',
        component: lazy(() => import('@/views/dashboard')),
        authority: ['Admin', 'ANP HSE'],
    },

    ...othersRoute,
]
