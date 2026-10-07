import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import Card from '@/components/ui/Card'
import Table from '@/components/ui/Table'
import Tag from '@/components/ui/Tag'
import Button from '@/components/ui/Button'
import Spinner from '@/components/ui/Spinner'
import Tabs from '@/components/ui/Tabs'
import { apiListVesselApplicationsForReview } from '@/services/VesselService'

const { Tr, Th, Td, THead, TBody } = Table
const { TabNav, TabList } = Tabs

const STATUS_TAG = {
    Submitted: 'bg-amber-100 text-amber-700',
    Approved: 'bg-emerald-100 text-emerald-700',
    Rejected: 'bg-red-100 text-red-700',
    Withdrawn: 'bg-gray-200 text-gray-500',
}

const TABS = ['Submitted', 'Approved', 'Rejected', 'All']

export default function VesselReview() {
    const navigate = useNavigate()
    const [tab, setTab] = useState('Submitted')
    const [loading, setLoading] = useState(true)
    const [applications, setApplications] = useState([])

    useEffect(() => {
        setLoading(true)
        apiListVesselApplicationsForReview(tab === 'All' ? undefined : tab)
            .then((data) => setApplications(data || []))
            .catch(() => setApplications([]))
            .finally(() => setLoading(false))
    }, [tab])

    return (
        <Card>
            <h4 className="mb-4">Vessel Entry Application Review</h4>

            <Tabs value={tab} onChange={setTab}>
                <TabList>
                    {TABS.map((t) => (
                        <TabNav key={t} value={t}>
                            {t}
                        </TabNav>
                    ))}
                </TabList>
            </Tabs>

            <div className="mt-4">
                {loading ? (
                    <div className="flex justify-center py-10">
                        <Spinner size={32} />
                    </div>
                ) : applications.length === 0 ? (
                    <p className="py-8 text-center text-gray-500">No applications in this category.</p>
                ) : (
                    <Table>
                        <THead>
                            <Tr>
                                <Th>Application #</Th>
                                <Th>Vessel</Th>
                                <Th>Status</Th>
                                <Th>Submitted</Th>
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
                                    <Td>
                                        {app.submittedAt ? new Date(app.submittedAt).toLocaleDateString() : '-'}
                                    </Td>
                                    <Td>
                                        <Button
                                            size="xs"
                                            onClick={() => navigate(`/vessel/review/${app.id}`)}
                                        >
                                            {app.status === 'Submitted' ? 'Review' : 'View'}
                                        </Button>
                                    </Td>
                                </Tr>
                            ))}
                        </TBody>
                    </Table>
                )}
            </div>
        </Card>
    )
}