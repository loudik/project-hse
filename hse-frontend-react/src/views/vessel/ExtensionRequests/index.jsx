import { useEffect, useState } from 'react'
import Card from '@/components/ui/Card'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Tag from '@/components/ui/Tag'
import Spinner from '@/components/ui/Spinner'
import Dialog from '@/components/ui/Dialog'
import Notification from '@/components/ui/Notification'
import toast from '@/components/ui/toast'
import {
    apiListExtensionRequests,
    apiDecideExtension,
} from '@/services/VesselService'

export default function ExtensionRequests() {
    const [requests, setRequests] = useState([])
    const [loading, setLoading] = useState(true)
    const [deciding, setDeciding] = useState(false)
    const [rejectTarget, setRejectTarget] = useState(null)
    const [rejectionReason, setRejectionReason] = useState('')

    const load = () => {
        setLoading(true)
        apiListExtensionRequests('Pending')
            .then((list) => setRequests(list || []))
            .finally(() => setLoading(false))
    }

    useEffect(() => {
        load()
    }, [])

    const handleApprove = async (id) => {
        setDeciding(true)
        try {
            await apiDecideExtension(id, 'Approved')
            toast.push(
                <Notification type="success" title="Approved">
                    Extension approved.
                </Notification>,
            )
            load()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to approve.'}
                </Notification>,
            )
        } finally {
            setDeciding(false)
        }
    }

    const handleReject = async () => {
        if (!rejectionReason) {
            toast.push(
                <Notification type="danger" title="Missing reason">
                    Please provide a rejection reason.
                </Notification>,
            )
            return
        }
        setDeciding(true)
        try {
            await apiDecideExtension(rejectTarget, 'Rejected', rejectionReason)
            toast.push(
                <Notification type="success" title="Rejected">
                    Extension rejected.
                </Notification>,
            )
            setRejectTarget(null)
            setRejectionReason('')
            load()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to reject.'}
                </Notification>,
            )
        } finally {
            setDeciding(false)
        }
    }

    if (loading) {
        return (
            <div className="flex justify-center py-10">
                <Spinner size={32} />
            </div>
        )
    }

    return (
        <Card>
            <h3 className="mb-4">Extension Requests</h3>

            {requests.length === 0 ? (
                <p className="text-gray-500">No pending extension requests.</p>
            ) : (
                <div className="flex flex-col gap-3">
                    {requests.map((r) => (
                        <Card key={r.id} className="border border-gray-100">
                            <div className="flex flex-wrap items-center justify-between gap-3">
                                <div>
                                    <div className="font-mono text-xs text-gray-400">
                                        {r.applicationNumber}
                                    </div>
                                    <div className="font-semibold">
                                        {r.vesselName || 'Unnamed vessel'}
                                    </div>
                                    <div className="text-sm text-gray-600">
                                        {r.currentEntryDateTo} →{' '}
                                        <span className="font-semibold">
                                            {r.requestedEntryDateTo}
                                        </span>
                                    </div>
                                    <div className="text-sm text-gray-500">
                                        {r.reason}
                                    </div>
                                    <div className="text-xs text-gray-400">
                                        Requested by {r.requestedByName}
                                    </div>
                                </div>
                                <div className="flex gap-2">
                                    <Button
                                        size="sm"
                                        variant="solid"
                                        color="emerald-600"
                                        loading={deciding}
                                        onClick={() => handleApprove(r.id)}
                                    >
                                        Approve
                                    </Button>
                                    <Button
                                        size="sm"
                                        variant="solid"
                                        color="red-600"
                                        loading={deciding}
                                        onClick={() => setRejectTarget(r.id)}
                                    >
                                        Reject
                                    </Button>
                                </div>
                            </div>
                        </Card>
                    ))}
                </div>
            )}

            <Dialog
                isOpen={!!rejectTarget}
                onClose={() => setRejectTarget(null)}
                onRequestClose={() => setRejectTarget(null)}
            >
                <h4 className="mb-4">Reject Extension Request</h4>
                <Input
                    textArea
                    placeholder="Reason for rejection"
                    value={rejectionReason}
                    onChange={(e) => setRejectionReason(e.target.value)}
                />
                <div className="mt-6 flex justify-end gap-2">
                    <Button onClick={() => setRejectTarget(null)}>
                        Cancel
                    </Button>
                    <Button
                        variant="solid"
                        color="red-600"
                        loading={deciding}
                        onClick={handleReject}
                    >
                        Reject
                    </Button>
                </div>
            </Dialog>
        </Card>
    )
}