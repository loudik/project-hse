import { useEffect, useState } from 'react'
import Card from '@/components/ui/Card'
import Table from '@/components/ui/Table'
import Tag from '@/components/ui/Tag'
import Button from '@/components/ui/Button'
import Spinner from '@/components/ui/Spinner'
import Tabs from '@/components/ui/Tabs'
import { apiListOrganizations } from '@/services/OrganizationService'
import DecisionDialog from './components/DecisionDialog'

const { Tr, Th, Td, THead, TBody } = Table
const { TabNav, TabList } = Tabs

const STATUS_TAG = {
    Pending: 'bg-amber-100 text-amber-700',
    Approved: 'bg-emerald-100 text-emerald-700',
    Rejected: 'bg-red-100 text-red-700',
}

const TABS = ['Pending', 'Approved', 'Rejected', 'All']

export default function OrganizationReview() {
    const [tab, setTab] = useState('Pending')
    const [loading, setLoading] = useState(true)
    const [organizations, setOrganizations] = useState([])
    const [dialogState, setDialogState] = useState({ organization: null, action: null })

    const load = () => {
        setLoading(true)
        apiListOrganizations(tab === 'All' ? undefined : tab)
            .then((data) => setOrganizations(data || []))
            .catch(() => setOrganizations([]))
            .finally(() => setLoading(false))
    }

    useEffect(() => {
        load()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [tab])

    const handleDecided = () => {
        load()
    }

    return (
        <Card>
            <h4 className="mb-4">Organization review</h4>

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
                ) : organizations.length === 0 ? (
                    <p className="py-8 text-center text-gray-500">
                        No organizations in this category.
                    </p>
                ) : (
                    <Table>
                        <THead>
                            <Tr>
                                <Th>Name</Th>
                                <Th>Registration #</Th>
                                <Th>Type</Th>
                                <Th>Country</Th>
                                <Th>Status</Th>
                                <Th>Submitted</Th>
                                <Th></Th>
                            </Tr>
                        </THead>
                        <TBody>
                            {organizations.map((org) => (
                                <Tr key={org.id}>
                                    <Td>{org.name}</Td>
                                    <Td>{org.registrationNumber}</Td>
                                    <Td>{org.type}</Td>
                                    <Td>{org.country || '-'}</Td>
                                    <Td>
                                        <Tag className={STATUS_TAG[org.status]}>{org.status}</Tag>
                                    </Td>
                                    <Td>{new Date(org.createdAt).toLocaleDateString()}</Td>
                                    <Td>
                                        {org.status === 'Pending' && (
                                            <div className="flex gap-2">
                                                <Button
                                                    size="xs"
                                                    variant="solid"
                                                    onClick={() =>
                                                        setDialogState({ organization: org, action: 'Approved' })
                                                    }
                                                >
                                                    Approve
                                                </Button>
                                                <Button
                                                    size="xs"
                                                    onClick={() =>
                                                        setDialogState({ organization: org, action: 'Rejected' })
                                                    }
                                                >
                                                    Reject
                                                </Button>
                                            </div>
                                        )}
                                    </Td>
                                </Tr>
                            ))}
                        </TBody>
                    </Table>
                )}
            </div>

            <DecisionDialog
                organization={dialogState.organization}
                action={dialogState.action}
                onClose={() => setDialogState({ organization: null, action: null })}
                onDecided={handleDecided}
            />
        </Card>
    )
}
