import { useEffect, useState } from 'react'
import Card from '@/components/ui/Card'
import Spinner from '@/components/ui/Spinner'
import { apiGetMyOrganization } from '@/services/OrganizationService'
import OrganizationForm from './components/OrganizationForm'

const STATUS_STYLE = {
    Pending: { label: 'Pending review', className: 'bg-amber-100 text-amber-700' },
    Approved: { label: 'Approved', className: 'bg-emerald-100 text-emerald-700' },
    Rejected: { label: 'Rejected', className: 'bg-red-100 text-red-700' },
}

function StatusBadge({ status }) {
    const cfg = STATUS_STYLE[status] || STATUS_STYLE.Pending
    return (
        <span className={`inline-block rounded-full px-3 py-1 text-xs font-semibold ${cfg.className}`}>
            {cfg.label}
        </span>
    )
}

function OrganizationSummary({ org }) {
    const rows = [
        ['Organization name', org.name],
        ['Registration number', org.registrationNumber],
        ['Type', org.type],
        ['Address', org.address],
        ['Country', org.country],
        ['Phone number', org.phoneNumber],
        ['Email', org.email],
        ['Website', org.website || '-'],
    ]

    return (
        <Card>
            <div className="mb-4 flex items-center justify-between">
                <h4>Organization profile</h4>
                <StatusBadge status={org.status} />
            </div>

            {org.status === 'Pending' && (
                <p className="mb-4 text-gray-500">
                    Your organization profile is awaiting review by ANP HSE. You will be
                    notified once it has been approved.
                </p>
            )}
            {org.status === 'Rejected' && (
                <div className="mb-4 rounded-lg bg-red-50 p-3 text-red-700">
                    <p className="font-semibold">This profile was rejected.</p>
                    {org.rejectionReason && <p className="mt-1">{org.rejectionReason}</p>}
                </div>
            )}
            {org.status === 'Approved' && (
                <p className="mb-4 text-gray-500">
                    Your organization has been approved. You can now access Vessel Entry
                    Application from the sidebar.
                </p>
            )}

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                {rows.map(([label, value]) => (
                    <div key={label}>
                        <div className="text-xs font-semibold tracking-wide text-gray-400 uppercase">
                            {label}
                        </div>
                        <div className="text-gray-700">{value}</div>
                    </div>
                ))}
            </div>
        </Card>
    )
}

export default function OrganizationProfile() {
    const [loading, setLoading] = useState(true)
    const [org, setOrg] = useState(null)

    useEffect(() => {
        apiGetMyOrganization()
            .then((data) => setOrg(data))
            .catch(() => setOrg(null))
            .finally(() => setLoading(false))
    }, [])

    if (loading) {
        return (
            <div className="flex justify-center py-10">
                <Spinner size={32} />
            </div>
        )
    }

    if (org) {
        return <OrganizationSummary org={org} />
    }

    return (
        <Card>
            <h4 className="mb-1">Create your organization profile</h4>
            <p className="mb-6 text-gray-500">
                This information will be reviewed by ANP HSE. Once approved, you'll be
                able to submit Vessel Entry Applications.
            </p>
            <OrganizationForm onCreated={setOrg} />
        </Card>
    )
}
