import Timeline from '@/components/ui/Timeline'
import Spinner from '@/components/ui/Spinner'

const ACTION_COLOR = {
    submitted: 'bg-amber-500',
    approved: 'bg-emerald-500',
    rejected: 'bg-red-500',
    withdrawn: 'bg-gray-400',
    reopened: 'bg-sky-500',
}

const ACTION_LABEL = {
    submitted: 'Submitted',
    approved: 'Approved',
    rejected: 'Rejected',
    withdrawn: 'Withdrawn',
    reopened: 'Reopened',
}

function formatDateTime(isoString) {
    const date = new Date(isoString)
    if (Number.isNaN(date.getTime())) return isoString
    return date.toLocaleString(undefined, {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}

// Renders the audit-log history for a vessel application or organization
// as a vertical timeline. `entries` is the array returned by
// GET .../history (each: { action, oldStatus, newStatus, performedByName,
// notes, createdAt }).
export default function HistoryTimeline({ entries, loading }) {
    if (loading) {
        return (
            <div className="flex justify-center py-6">
                <Spinner size={28} />
            </div>
        )
    }

    if (!entries || entries.length === 0) {
        return <p className="text-sm text-gray-500">No history recorded yet.</p>
    }

    return (
        <Timeline>
            {entries.map((entry) => (
                <Timeline.Item
                    key={entry.id}
                    media={
                        <div
                            className={`h-3 w-3 rounded-full ${
                                ACTION_COLOR[entry.action] || 'bg-gray-400'
                            }`}
                        />
                    }
                >
                    <div className="text-sm font-semibold text-gray-800">
                        {ACTION_LABEL[entry.action] || entry.action}
                        {entry.performedByName && (
                            <span className="font-normal text-gray-500">
                                {' '}
                                by {entry.performedByName}
                            </span>
                        )}
                    </div>
                    {entry.notes && (
                        <div className="mt-0.5 text-sm text-gray-600">
                            {entry.notes}
                        </div>
                    )}
                    <div className="mt-0.5 text-xs text-gray-400">
                        {formatDateTime(entry.createdAt)}
                    </div>
                </Timeline.Item>
            ))}
        </Timeline>
    )
}