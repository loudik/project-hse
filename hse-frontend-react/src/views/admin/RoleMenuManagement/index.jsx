import { useEffect, useState } from 'react'
import Card from '@/components/ui/Card'
import Button from '@/components/ui/Button'
import Spinner from '@/components/ui/Spinner'
import Checkbox from '@/components/ui/Checkbox'
import toast from '@/components/ui/toast'
import Notification from '@/components/ui/Notification'
import { apiGetAccessMatrix, apiUpdateAccessMatrix } from '@/services/UserService'

export default function RoleMenuManagement() {
    const [menus, setMenus] = useState([])
    const [roles, setRoles] = useState([])
    // access is a map keyed "roleId:menuId" -> boolean, so toggling a
    // checkbox is a simple lookup instead of scanning an array each time.
    const [access, setAccess] = useState({})
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)

    const load = () => {
        setLoading(true)
        apiGetAccessMatrix()
            .then((data) => {
                setMenus(data.menus || [])
                setRoles(data.roles || [])
                const map = {}
                ;(data.access || []).forEach((a) => {
                    map[`${a.roleId}:${a.menuId}`] = a.canView
                })
                setAccess(map)
            })
            .finally(() => setLoading(false))
    }

    useEffect(() => {
        load()
    }, [])

    const toggle = (roleId, menuId) => {
        const key = `${roleId}:${menuId}`
        setAccess((prev) => ({ ...prev, [key]: !prev[key] }))
    }

    const handleSave = async () => {
        setSaving(true)
        try {
            const payload = []
            roles.forEach((role) => {
                menus.forEach((menu) => {
                    payload.push({
                        roleId: role.id,
                        menuId: menu.id,
                        canView: !!access[`${role.id}:${menu.id}`],
                    })
                })
            })
            await apiUpdateAccessMatrix(payload)
            toast.push(
                <Notification type="success" title="Saved">
                    Menu access updated.
                </Notification>,
            )
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to save changes.'}
                </Notification>,
            )
        } finally {
            setSaving(false)
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
            <div className="mb-4 flex items-center justify-between">
                <div>
                    <h3>Role & Menu Management</h3>
                    <p className="text-gray-500">
                        Choose which menus each role can see. Tick a box to
                        grant access, then click Save.
                    </p>
                </div>
                <Button variant="solid" loading={saving} onClick={handleSave}>
                    Save Changes
                </Button>
            </div>

            <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                    <thead>
                        <tr className="border-b border-gray-200">
                            <th className="py-2 pr-4 font-semibold text-gray-600">
                                Menu
                            </th>
                            {roles.map((role) => (
                                <th
                                    key={role.id}
                                    className="whitespace-nowrap px-3 py-2 text-center font-semibold text-gray-600"
                                >
                                    {role.name}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {menus.map((menu) => (
                            <tr
                                key={menu.id}
                                className="border-b border-gray-100"
                            >
                                <td className="py-2 pr-4 text-gray-800">
                                    {menu.name}
                                    {menu.path && (
                                        <span className="ml-2 text-xs text-gray-400">
                                            {menu.path}
                                        </span>
                                    )}
                                </td>
                                {roles.map((role) => (
                                    <td
                                        key={role.id}
                                        className="px-3 py-2 text-center"
                                    >
                                        <Checkbox
                                            checked={
                                                !!access[
                                                    `${role.id}:${menu.id}`
                                                ]
                                            }
                                            onChange={() =>
                                                toggle(role.id, menu.id)
                                            }
                                        />
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </Card>
    )
}