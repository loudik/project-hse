import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'react-router'
import { useAuth } from '@/auth'
import { apiMicrosoftOauthSignIn } from '@/services/OAuthServices'

export default function MicrosoftCallback() {
    const [searchParams] = useSearchParams()
    const { oAuthSignIn } = useAuth()
    const [error, setError] = useState('')

  
    const hasRun = useRef(false)

    useEffect(() => {
        if (hasRun.current) return
        hasRun.current = true

        const code = searchParams.get('code')
        const errorParam = searchParams.get('error_description')

        if (errorParam) {
            setError(errorParam)
            return
        }
        if (!code) {
            setError('Authorization code from Microsoft not found.')
            return
        }

        oAuthSignIn(async ({ redirect, onSignIn }) => {
            try {
                const resp = await apiMicrosoftOauthSignIn(code)
                if (resp) {
                    const { token, user } = resp
                    onSignIn({ accessToken: token }, user)
                    redirect()
                }
            } catch (err) {
                setError(
                    err?.response?.data?.message ||
                        err?.toString() ||
                        'Failed to sign in with Microsoft',
                )
            }
        })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <div className="flex min-h-screen flex-col items-center justify-center gap-3 px-4 text-center">
            {error ? (
                <>
                    <p className="font-semibold text-red-500">{error}</p>
                    <a href="/sign-in" className="underline">
                        Back to sign in
                    </a>
                </>
            ) : (
                <p>Signing in with Microsoft...</p>
            )}
        </div>
    )
}
