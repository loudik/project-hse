import Button from '@/components/ui/Button'
import { apiGoogleOauthSignIn } from '@/services/OAuthServices'
import { useAuth } from '@/auth'

const OauthSignIn = ({ setMessage, disableSubmit }) => {
    const { oAuthSignIn } = useAuth()

    const handleGoogleSignIn = async () => {
        if (!disableSubmit) {
            oAuthSignIn(async ({ redirect, onSignIn }) => {
                try {
                    const resp = await apiGoogleOauthSignIn()
                    if (resp) {
                        const { token, user } = resp
                        onSignIn({ accessToken: token }, user)
                        redirect()
                    }
                } catch (error) {
                    setMessage?.(error?.toString() || '')
                }
            })
        }
    }

    // Microsoft pakai alur "redirect penuh": browser diarahkan ke halaman
    // login Microsoft, lalu Microsoft redirect balik ke /auth/callback
    // membawa "code". Proses tukar code -> token dilanjutkan di sana
    // (lihat src/views/auth/MicrosoftCallback/index.jsx).
    const handleMicrosoftSignIn = () => {
        if (disableSubmit) return

        const tenantId = import.meta.env.VITE_MS_TENANT_ID
        const clientId = import.meta.env.VITE_MS_CLIENT_ID
        const redirectUri = import.meta.env.VITE_MS_REDIRECT_URI

        const params = new URLSearchParams({
            client_id: clientId,
            response_type: 'code',
            redirect_uri: redirectUri,
            response_mode: 'query',
            scope: 'openid profile email User.Read',
            prompt: 'select_account',
        })

        window.location.href = `https://login.microsoftonline.com/${tenantId}/oauth2/v2.0/authorize?${params.toString()}`
    }

    return (
        <div className="flex flex-col gap-3">
            <button
                type="button"
                onClick={handleGoogleSignIn}
                className="flex items-center justify-center gap-3 w-full h-11 rounded-lg border border-gray-300 dark:border-gray-600 bg-white text-gray-700 font-medium hover:bg-gray-50 active:bg-gray-100 transition-colors"
            >
                <img
                    className="h-5 w-5"
                    src="/img/others/google.png"
                    alt=""
                    aria-hidden="true"
                />
                <span>Continue with Google</span>
            </button>
            <button
                type="button"
                onClick={handleMicrosoftSignIn}
                className="flex items-center justify-center gap-3 w-full h-11 rounded-lg bg-[#2F2F2F] text-white font-medium hover:bg-[#1f1f1f] active:bg-[#141414] transition-colors"
            >
                <img
                    className="h-5 w-5"
                    src="/img/others/microsoft.png"
                    alt=""
                    aria-hidden="true"
                />
                <span>Continue with Microsoft 365</span>
            </button>
        </div>
    )
}

export default OauthSignIn