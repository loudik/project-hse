import { useEffect, useState } from 'react'
import Card from '@/components/ui/Card'
import Table from '@/components/ui/Table'
import Tag from '@/components/ui/Tag'
import Select from '@/components/ui/Select'
import Button from '@/components/ui/Button'
import Spinner from '@/components/ui/Spinner'
import Notification from '@/components/ui/Notification'
import toast from '@/components/ui/toast'
import {
    apiListUsers,
    apiListRoles,
    apiUpdateUserRole,
    apiUpdateUserStatus,
} from '@/services/UserService'
import { useSessionUser } from '@/store/authStore'

const { Tr, Th, Td, THead, TBody } = Table

const STATUS_TAG = {
    Active: 'bg-emerald-100 text-emerald-700',
    Pending: 'bg-amber-100 text-amber-700',
    Suspended: 'bg-red-100 text-red-700',
}

export default function UserManagement() {
    const currentEmail = useSessionUser((state) => state.user.email)
    const [loading, setLoading] = useState(true)
    const [users, setUsers] = useState([])
    const [roles, setRoles] = useState([])
    const [savingId, setSavingId] = useState(null)

    const load = () => {
        setLoading(true)
        Promise.all([apiListUsers(), apiListRoles()])
            .then(([userList, roleList]) => {
                setUsers(userList || [])
                setRoles(roleList || [])
            })
            .catch(() => {
                setUsers([])
                setRoles([])
            })
            .finally(() => setLoading(false))
    }

    useEffect(() => {
        load()
    }, [])

    const notifyError = (err, fallback) => {
        toast.push(
            <Notification type="danger" title="Error">
                {err?.response?.data?.message || fallback}
            </Notification>,
        )
    }

    const handleRoleChange = async (user, option) => {
        setSavingId(user.id)
        try {
            await apiUpdateUserRole(user.id, option.value)
            setUsers((prev) =>
                prev.map((u) =>
                    u.id === user.id ? { ...u, roleId: option.value, roleName: option.label } : u,
                ),
            )
        } catch (err) {
            notifyError(err, 'Failed to update role.')
        } finally {
            setSavingId(null)
        }
    }

    const handleToggleStatus = async (user) => {
        const nextStatus = user.status === 'Suspended' ? 'Active' : 'Suspended'
        setSavingId(user.id)
        try {
            await apiUpdateUserStatus(user.id, nextStatus)
            setUsers((prev) =>
                prev.map((u) => (u.id === user.id ? { ...u, status: nextStatus } : u)),
            )
        } catch (err) {
            notifyError(err, 'Failed to update status.')
        } finally {
            setSavingId(null)
        }
    }

    const roleOptions = roles.map((r) => ({ value: r.id, label: r.name }))

    return (
        <Card>
            <h4 className="mb-1">User Management</h4>
            <p className="mb-6 text-gray-500">
                Assign roles and manage account status. New Microsoft sign-ins start as
                "Unassigned" with no menu access until given a role here.
            </p>

            {loading ? (
                <div className="flex justify-center py-10">
                    <Spinner size={32} />
                </div>
            ) : (
                <Table>
                    <THead>
                        <Tr>
                            <Th>Name</Th>
                            <Th>Email</Th>
                            <Th>Organization</Th>
                            <Th>Role</Th>
                            <Th>Status</Th>
                            <Th></Th>
                        </Tr>
                    </THead>
                    <TBody>
                        {users.map((user) => {
                            const isSelf = user.email === currentEmail
                            return (
                                <Tr key={user.id}>
                                    <Td>{user.name}</Td>
                                    <Td>{user.email}</Td>
                                    <Td>{user.organizationName || '-'}</Td>
                                    <Td className="min-w-[180px]">
                                        <Select
                                            isDisabled={isSelf || savingId === user.id}
                                            options={roleOptions}
                                            value={roleOptions.find((o) => o.value === user.roleId) || null}
                                            onChange={(option) => handleRoleChange(user, option)}
                                        />
                                    </Td>
                                    <Td>
                                        <Tag className={STATUS_TAG[user.status]}>{user.status}</Tag>
                                    </Td>
                                    <Td>
                                        {!isSelf && (
                                            <Button
                                                size="xs"
                                                loading={savingId === user.id}
                                                onClick={() => handleToggleStatus(user)}
                                            >
                                                {user.status === 'Suspended' ? 'Activate' : 'Suspend'}
                                            </Button>
                                        )}
                                    </Td>
                                </Tr>
                            )
                        })}
                    </TBody>
                </Table>
            )}
        </Card>
    )
}
