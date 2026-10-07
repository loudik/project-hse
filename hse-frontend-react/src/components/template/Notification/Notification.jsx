import { useEffect, useState, useRef } from 'react'
import classNames from 'classnames'
import withHeaderItem from '@/utils/hoc/withHeaderItem'
import Dropdown from '@/components/ui/Dropdown'
import ScrollBar from '@/components/ui/ScrollBar'
import Spinner from '@/components/ui/Spinner'
import Badge from '@/components/ui/Badge'
import Button from '@/components/ui/Button'
import NotificationAvatar from './NotificationAvatar'
import NotificationToggle from './NotificationToggle'
import { HiOutlineMailOpen } from 'react-icons/hi'
import {
    apiGetNotificationList,
    apiGetNotificationCount,
    apiMarkNotificationRead,
    apiMarkAllNotificationsRead,
} from '@/services/CommonService'
import isLastChild from '@/utils/isLastChild'
import useResponsive from '@/utils/hooks/useResponsive'
import { useNavigate } from 'react-router'

const notificationHeight = 'h-[280px]'

// Backend sends { id, title, message, link, isRead, createdAt }. Map that
// to the avatar type this component's NotificationAvatar expects, so
// approve/reject notifications get a clear success/fail icon instead of a
// generic avatar.
const toAvatarProps = (title = '') => {
    const lower = title.toLowerCase()
    if (lower.includes('approved')) {
        return { type: 2, status: 'succeed' }
    }
    if (lower.includes('rejected')) {
        return { type: 2, status: 'failed' }
    }
    if (lower.includes('submitted')) {
        return { type: 1 }
    }
    return { type: 0, target: title }
}

const formatRelativeDate = (isoString) => {
    if (!isoString) return ''
    const date = new Date(isoString)
    if (Number.isNaN(date.getTime())) return ''
    return date.toLocaleString(undefined, {
        day: '2-digit',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit',
    })
}

const _Notification = ({ className }) => {
    const [notificationList, setNotificationList] = useState([])
    const [unreadNotification, setUnreadNotification] = useState(false)
    const [noResult, setNoResult] = useState(false)
    const [loading, setLoading] = useState(false)

    const { larger } = useResponsive()

    const navigate = useNavigate()

    const getNotificationCount = async () => {
        const resp = await apiGetNotificationCount()
        if (resp.count > 0) {
            setNoResult(false)
            setUnreadNotification(true)
        } else {
            setNoResult(true)
        }
    }

    useEffect(() => {
        getNotificationCount()
    }, [])

    const onNotificationOpen = async () => {
        setLoading(true)
        const resp = await apiGetNotificationList()
        setLoading(false)
        setNotificationList(resp)
        setNoResult(resp.length === 0)
    }

    const onMarkAllAsRead = async () => {
        const hadUnread = notificationList.some((item) => !item.isRead)
        if (!hadUnread) return

        const list = notificationList.map((item) => ({
            ...item,
            isRead: true,
        }))
        setNotificationList(list)
        setUnreadNotification(false)
        await apiMarkAllNotificationsRead()
    }

    const notificationDropdownRef = useRef(null)

    const onMarkAsRead = async (item) => {
        if (!item.isRead) {
            const list = notificationList.map((n) =>
                n.id === item.id ? { ...n, isRead: true } : n,
            )
            setNotificationList(list)
            const stillUnread = list.some((n) => !n.isRead)
            setUnreadNotification(stillUnread)
            apiMarkNotificationRead(item.id)
        }
        if (item.link) {
            navigate(item.link)
            notificationDropdownRef.current?.handleDropdownClose()
        }
    }

    return (
        <Dropdown
            ref={notificationDropdownRef}
            renderTitle={
                <NotificationToggle
                    dot={unreadNotification}
                    className={className}
                />
            }
            menuClass="min-w-[280px] md:min-w-[340px]"
            placement={larger.md ? 'bottom-end' : 'bottom'}
            onOpen={onNotificationOpen}
        >
            <Dropdown.Item variant="header">
                <div className="dark:border-gray-700 px-2 flex items-center justify-between mb-1">
                    <h6>Notifications</h6>
                    <Button
                        variant="plain"
                        shape="circle"
                        size="sm"
                        icon={<HiOutlineMailOpen className="text-xl" />}
                        title="Mark all as read"
                        onClick={onMarkAllAsRead}
                    />
                </div>
            </Dropdown.Item>
            <ScrollBar
                className={classNames('overflow-y-auto', notificationHeight)}
            >
                {notificationList.length > 0 &&
                    notificationList.map((item, index) => (
                        <div key={item.id}>
                            <div
                                className={`relative rounded-xl flex px-4 py-3 cursor-pointer hover:bg-gray-100 active:bg-gray-100 dark:hover:bg-gray-700`}
                                onClick={() => onMarkAsRead(item)}
                            >
                                <div>
                                    <NotificationAvatar
                                        {...toAvatarProps(item.title)}
                                    />
                                </div>
                                <div className="mx-3">
                                    <div>
                                        <span className="font-semibold heading-text">
                                            {item.title}{' '}
                                        </span>
                                        {item.message && (
                                            <span>{item.message}</span>
                                        )}
                                    </div>
                                    <span className="text-xs">
                                        {formatRelativeDate(item.createdAt)}
                                    </span>
                                </div>
                                <Badge
                                    className="absolute top-4 ltr:right-4 rtl:left-4 mt-1.5"
                                    innerClass={`${
                                        item.isRead
                                            ? 'bg-gray-300 dark:bg-gray-600'
                                            : 'bg-primary'
                                    } `}
                                />
                            </div>
                            {!isLastChild(notificationList, index) ? (
                                <div className="border-b border-gray-200 dark:border-gray-700 my-2" />
                            ) : (
                                ''
                            )}
                        </div>
                    ))}
                {loading && (
                    <div
                        className={classNames(
                            'flex items-center justify-center',
                            notificationHeight,
                        )}
                    >
                        <Spinner size={40} />
                    </div>
                )}
                {noResult && notificationList.length === 0 && !loading && (
                    <div
                        className={classNames(
                            'flex items-center justify-center',
                            notificationHeight,
                        )}
                    >
                        <div className="text-center">
                            <img
                                className="mx-auto mb-2 max-w-[150px]"
                                src="/img/others/no-notification.png"
                                alt="no-notification"
                            />
                            <h6 className="font-semibold">
                                No notifications!
                            </h6>
                            <p className="mt-1">
                                You&apos;re all caught up for now
                            </p>
                        </div>
                    </div>
                )}
            </ScrollBar>
        </Dropdown>
    )
}

const Notification = withHeaderItem(_Notification)

export default Notification