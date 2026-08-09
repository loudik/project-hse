import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router'
import Logo from '@/components/template/Logo'
import ActionLink from '@/components/shared/ActionLink'
import { useThemeStore } from '@/store/themeStore'
import { apiVerifyEmail } from '@/services/AuthService'

export default function VerifyEmail() {
    const [searchParams] = useSearchParams()
    const mode = useThemeStore((state) => state.mode)
    const [status, setStatus] = useState('loading') // loading | success | error
    const [message, setMessage] = useState('')

    useEffect(() => {
        const token = searchParams.get('token')
        if (!token) {
            setStatus('error')
            setMessage('Verification token is missing from the link.')
            return
        }

        apiVerifyEmail(token)
            .then((resp) => {
                setStatus('success')
                setMessage(resp?.message || 'Your email has been verified.')
            })
            .catch((err) => {
                setStatus('error')
                setMessage(err?.response?.data?.message || 'Failed to verify email.')
            })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <div className="flex min-h-screen flex-col items-center justify-center gap-4 px-4 text-center">
            <Logo type="streamline" mode={mode} imgClass="mx-auto" logoWidth={60} />
            {status === 'loading' && <p>Verifying your email...</p>}
            {status === 'success' && (
                <>
                    <h3 className="mb-1">Email verified</h3>
                    <p className="font-semibold heading-text">{message}</p>
                </>
            )}
            {status === 'error' && (
                <>
                    <h3 className="mb-1 text-red-500">Verification failed</h3>
                    <p className="font-semibold heading-text">{message}</p>
                </>
            )}
            <ActionLink to="/sign-in" className="heading-text mt-2 font-bold underline" themeColor={false}>
                Back to sign in
            </ActionLink>
        </div>
    )
}
