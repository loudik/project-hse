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
        <div className="flex items-center gap-2">
            <Button
                className="flex-1"
                type="button"
                onClick={handleGoogleSignIn}
            >
                <div className="flex items-center justify-center gap-2">
                    <img
                        className="h-[25px] w-[25px]"
                        src="/img/others/google.png"
                        alt="Google sign in"
                    />
                    <span>Google</span>
                </div>
            </Button>
            <Button
                className="flex-1"
                type="button"
                onClick={handleMicrosoftSignIn}
            >
                <div className="flex items-center justify-center gap-2">
                    <img
                        className="h-[25px] w-[25px]"
                        src="/img/others/microsoft.png"
                        alt="Microsoft 365 sign in"
                    />
                    <span>Microsoft 365</span>
                </div>
            </Button>
        </div>
    )
}

export default OauthSignIn
