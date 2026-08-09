import { useState } from 'react'
import Dialog from '@/components/ui/Dialog'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Alert from '@/components/ui/Alert'
import { apiDecideOrganization } from '@/services/OrganizationService'

const DecisionDialog = ({ organization, action, onClose, onDecided }) => {
    const [rejectionReason, setRejectionReason] = useState('')
    const [submitting, setSubmitting] = useState(false)
    const [error, setError] = useState('')

    if (!organization || !action) return null

    const isReject = action === 'Rejected'

    const handleConfirm = async () => {
        if (isReject && !rejectionReason.trim()) {
            setError('A rejection reason is required.')
            return
        }
        setSubmitting(true)
        setError('')
        try {
            const updated = await apiDecideOrganization(organization.id, {
                status: action,
                rejectionReason: isReject ? rejectionReason : undefined,
            })
            onDecided?.(updated)
            onClose?.()
        } catch (err) {
            setError(err?.response?.data?.message || 'Failed to save decision.')
        } finally {
            setSubmitting(false)
        }
    }

    return (
        <Dialog isOpen onClose={onClose} onRequestClose={onClose}>
            <h5 className="mb-4">
                {isReject ? 'Reject organization' : 'Approve organization'}
            </h5>
            <p className="mb-4 text-gray-500">
                {organization.name} ({organization.registrationNumber})
            </p>

            {isReject && (
                <div className="mb-4">
                    <label className="mb-1.5 block text-sm font-semibold">
                        Rejection reason
                    </label>
                    <Input
                        textArea
                        rows={3}
                        placeholder="Explain why this organization is being rejected..."
                        value={rejectionReason}
                        onChange={(e) => setRejectionReason(e.target.value)}
                    />
                </div>
            )}

            {error && (
                <Alert showIcon className="mb-4" type="danger">
                    {error}
                </Alert>
            )}

            <div className="flex justify-end gap-2">
                <Button onClick={onClose} disabled={submitting}>
                    Cancel
                </Button>
                <Button
                    variant="solid"
                    color={isReject ? 'red-600' : undefined}
                    loading={submitting}
                    onClick={handleConfirm}
                >
                    {isReject ? 'Reject' : 'Approve'}
                </Button>
            </div>
        </Dialog>
    )
}

export default DecisionDialog
