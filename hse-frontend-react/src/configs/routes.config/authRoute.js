import { lazy } from 'react'

const authRoute = [
    {
        key: 'signIn',
        path: `/sign-in`,
        component: lazy(() => import('@/views/auth/SignIn')),
        authority: [],
    },
    {
        key: 'signUp',
        path: `/sign-up`,
        component: lazy(() => import('@/views/auth/SignUp')),
        authority: [],
    },
    {
        key: 'forgotPassword',
        path: `/forgot-password`,
        component: lazy(() => import('@/views/auth/ForgotPassword')),
        authority: [],
    },
    {
        key: 'resetPassword',
        path: `/reset-password`,
        component: lazy(() => import('@/views/auth/ResetPassword')),
        authority: [],
    },
    {
        key: 'microsoftCallback',
        path: `/auth/callback`,
        component: lazy(() => import('@/views/auth/MicrosoftCallback')),
        authority: [],
    },
    {
        key: 'verifyEmail',
        path: `/verify-email`,
        component: lazy(() => import('@/views/auth/VerifyEmail')),
        authority: [],
    },
]

export default authRoute
