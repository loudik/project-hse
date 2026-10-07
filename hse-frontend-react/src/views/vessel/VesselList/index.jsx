import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import Card from '@/components/ui/Card'
import Table from '@/components/ui/Table'
import Tag from '@/components/ui/Tag'
import Button from '@/components/ui/Button'
import Spinner from '@/components/ui/Spinner'
import { apiListVesselApplications, apiCreateVesselApplication } from '@/services/VesselService'

const { Tr, Th, Td, THead, TBody } = Table

const STATUS_TAG = {
    Draft: 'bg-gray-100 text-gray-700',
    Submitted: 'bg-amber-100 text-amber-700',
    Approved: 'bg-emerald-100 text-emerald-700',
    Rejected: 'bg-red-100 text-red-700',
    Withdrawn: 'bg-gray-200 text-gray-500',
}

export default function VesselList() {
    const navigate = useNavigate()
    const [loading, setLoading] = useState(true)
    const [creating, setCreating] = useState(false)
    const [applications, setApplications] = useState([])

    useEffect(() => {
        apiListVesselApplications()
            .then((data) => setApplications(data || []))
            .catch(() => setApplications([]))
            .finally(() => setLoading(false))
    }, [])

    const handleNew = async () => {
        setCreating(true)
        try {
            const app = await apiCreateVesselApplication({})
            navigate(`/vessel/apply/${app.id}`)
        } catch (err) {
            console.error(err)
        } finally {
            setCreating(false)
        }
    }

    return (
        <Card>
            <div className="mb-4 flex items-center justify-between">
                <h4>Vessel Entry Applications</h4>
                <Button variant="solid" loading={creating} onClick={handleNew}>
                    + New Application
                </Button>
            </div>

            {loading ? (
                <div className="flex justify-center py-10">
                    <Spinner size={32} />
                </div>
            ) : applications.length === 0 ? (
                <p className="py-8 text-center text-gray-500">
                    No applications yet. Click "New Application" to start.
                </p>
            ) : (
                <Table>
                    <THead>
                        <Tr>
                            <Th>Application #</Th>
                            <Th>Vessel</Th>
                            <Th>Status</Th>
                            <Th>Created</Th>
                            <Th></Th>
                        </Tr>
                    </THead>
                    <TBody>
                        {applications.map((app) => (
                            <Tr key={app.id}>
                                <Td>{app.applicationNumber}</Td>
                                <Td>{app.vesselName || '-'}</Td>
                                <Td>
                                    <Tag className={STATUS_TAG[app.status]}>{app.status}</Tag>
                                </Td>
                                <Td>{new Date(app.createdAt).toLocaleDateString()}</Td>
                                <Td>
                                    <Button
                                        size="xs"
                                        onClick={() =>
                                            navigate(
                                                app.status === 'Draft' ? `/vessel/apply/${app.id}` : `/vessel/summary/${app.id}`
                                            )
                                        }
                                    >
                                        {app.status === 'Draft' ? 'Continue' : 'View'}
                                    </Button>
                                </Td>
                            </Tr>
                        ))}
                    </TBody>
                </Table>
            )}
        </Card>
    )
}