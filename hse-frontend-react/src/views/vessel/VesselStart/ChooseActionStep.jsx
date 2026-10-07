import { useState } from 'react'
import { useNavigate } from 'react-router'
import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'
import Alert from '@/components/ui/Alert'
import Radio from '@/components/ui/Radio'
import {
    apiCreateVesselApplication,
    apiLookupVesselApplication,
    apiWithdrawVesselApplication,
    apiReopenVesselApplication,
    apiSubstituteVessel,
} from '@/services/VesselService'

const ACTIONS = {
    NEW: 'new',
    UPDATE: 'update',
    WITHDRAW_ENTRY: 'withdraw_entry',
    WITHDRAW_AUTH: 'withdraw_auth',
    SUBSTITUTE: 'substitute',
    CONTINUE: 'continue',
}

export default function ChooseActionStep() {
    const navigate = useNavigate()
    const [action, setAction] = useState(ACTIONS.NEW)
    const [appNumber, setAppNumber] = useState('')
    const [confirmWithdraw, setConfirmWithdraw] = useState(false)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    const needsNumber = action !== ACTIONS.NEW
    const needsWithdrawConfirm = action === ACTIONS.WITHDRAW_ENTRY || action === ACTIONS.WITHDRAW_AUTH

    const handleGo = async () => {
        setError('')

        if (action === ACTIONS.NEW) {
            setLoading(true)
            try {
                const app = await apiCreateVesselApplication({})
                navigate(`/vessel/apply/${app.id}`)
            } catch (err) {
                setError(err?.response?.data?.message || 'Failed to start a new application.')
            } finally {
                setLoading(false)
            }
            return
        }

        if (!appNumber.trim()) {
            setError('Please enter your application unique number.')
            return
        }
        if (needsWithdrawConfirm && !confirmWithdraw) {
            setError('Please confirm that you want to withdraw.')
            return
        }

        setLoading(true)
        try {
            if (action === ACTIONS.SUBSTITUTE) {
                const result = await apiSubstituteVessel(appNumber.trim())
                navigate(`/vessel/apply/${result.id}`)
                return
            }

            const app = await apiLookupVesselApplication(appNumber.trim())

            if (action === ACTIONS.CONTINUE) {
                navigate(`/vessel/apply/${app.id}`)
                return
            }

            if (action === ACTIONS.UPDATE) {
                if (app.status !== 'Approved') {
                    setError('Only an approved application can be updated.')
                    return
                }
                await apiReopenVesselApplication(app.id)
                navigate(`/vessel/apply/${app.id}`)
                return
            }

            if (action === ACTIONS.WITHDRAW_ENTRY || action === ACTIONS.WITHDRAW_AUTH) {
                await apiWithdrawVesselApplication(app.id)
                navigate('/vessel')
                return
            }
        } catch (err) {
            setError(err?.response?.data?.message || 'Something went wrong. Please check the application number.')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div>
            <h5 className="mb-4">Please choose the following options</h5>

            <div className="mb-6 flex flex-col gap-5">
                <label className="flex cursor-pointer items-start gap-3">
                    <Radio checked={action === ACTIONS.NEW} onChange={() => setAction(ACTIONS.NEW)} />
                    <div>
                        <div className="font-semibold">I intend to submit a new entry application</div>
                        <div className="text-sm text-gray-500">
                            Choose this if you wish to submit a new entry application or extension of entry
                            authorisation.
                        </div>
                    </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3">
                    <Radio checked={action === ACTIONS.UPDATE} onChange={() => setAction(ACTIONS.UPDATE)} />
                    <div>
                        <div className="font-semibold">
                            I intend to update entry application previously submitted to ANP
                        </div>
                        <div className="text-sm text-gray-500">
                            Choose this if you wish to update your entry application (e.g. adding/removing scope
                            of work, updating a certificate). Only applicable for a vessel with a valid entry
                            authorisation.
                        </div>
                    </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3">
                    <Radio
                        checked={action === ACTIONS.WITHDRAW_ENTRY}
                        onChange={() => setAction(ACTIONS.WITHDRAW_ENTRY)}
                    />
                    <div>
                        <div className="font-semibold">
                            I intend to withdraw entry application previously submitted to ANP
                        </div>
                        <div className="text-sm text-gray-500">
                            Choose this if you wish to withdraw your entry application while it's still being
                            processed (e.g. due to a change in vessel or helicopter).
                        </div>
                    </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3">
                    <Radio
                        checked={action === ACTIONS.WITHDRAW_AUTH}
                        onChange={() => setAction(ACTIONS.WITHDRAW_AUTH)}
                    />
                    <div>
                        <div className="font-semibold">I intend to withdraw from entry authorisation</div>
                        <div className="text-sm text-gray-500">
                            Choose this if you wish to withdraw the entry authorisation for a specific vessel
                            (e.g. it's no longer available).
                        </div>
                    </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3">
                    <Radio
                        checked={action === ACTIONS.SUBSTITUTE}
                        onChange={() => setAction(ACTIONS.SUBSTITUTE)}
                    />
                    <div>
                        <div className="font-semibold">Substitute Vessel</div>
                        <div className="text-sm text-gray-500">
                            Choose this if a vessel with a valid entry authorisation needs to
                            be substituted with another vessel due to operational reasons
                            (e.g. maintenance, drydock). Can only be used under the same
                            vessel contract - the new vessel will go through the normal
                            review process.
                        </div>
                    </div>
                </label>

                <label className="flex cursor-pointer items-start gap-3">
                    <Radio checked={action === ACTIONS.CONTINUE} onChange={() => setAction(ACTIONS.CONTINUE)} />
                    <div>
                        <div className="font-semibold">I intend to continue working with Entry application</div>
                        <div className="text-sm text-gray-500">
                            Choose this to continue working on a draft entry application.
                        </div>
                    </div>
                </label>
            </div>

            {needsNumber && (
                <div className="mb-4 max-w-md rounded-lg bg-amber-50 p-4">
                    <label className="mb-1.5 block text-sm font-semibold">
                        {action === ACTIONS.SUBSTITUTE
                            ? "Please insert the original vessel's application unique number"
                            : 'Please insert your application unique number'}
                    </label>
                    <Input
                        value={appNumber}
                        onChange={(e) => setAppNumber(e.target.value)}
                        placeholder="e.g. VEA-20260810-1234"
                    />

                    {needsWithdrawConfirm && (
                        <div className="mt-4">
                            <div className="mb-1.5 text-sm font-semibold">Are you sure you want to withdraw?</div>
                            <div className="flex gap-2">
                                <Button
                                    size="sm"
                                    variant={confirmWithdraw ? 'solid' : 'default'}
                                    onClick={() => setConfirmWithdraw(true)}
                                >
                                    Yes
                                </Button>
                                <Button
                                    size="sm"
                                    variant={!confirmWithdraw ? 'solid' : 'default'}
                                    onClick={() => setConfirmWithdraw(false)}
                                >
                                    No
                                </Button>
                            </div>
                        </div>
                    )}
                </div>
            )}

            {error && (
                <Alert showIcon className="mb-4" type="danger">
                    {error}
                </Alert>
            )}

            <div className="flex justify-end">
                <Button variant="solid" loading={loading} onClick={handleGo}>
                    Continue
                </Button>
            </div>
        </div>
    )
}