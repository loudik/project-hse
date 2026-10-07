import { useEffect, useState } from 'react'
import Card from '@/components/ui/Card'
import Avatar from '@/components/ui/Avatar'
import Spinner from '@/components/ui/Spinner'
import Chart from '@/components/shared/Chart'
import {
    PiClipboardTextDuotone,
    PiHourglassMediumDuotone,
    PiCheckCircleDuotone,
    PiXCircleDuotone,
    PiBuildingsDuotone,
    PiTimerDuotone,
} from 'react-icons/pi'
import {
    apiGetDashboardSummary,
    apiGetDashboardTrend,
    apiGetDashboardAvgApprovalTime,
} from '@/services/CommonService'

const STATUS_ORDER = ['Draft', 'Submitted', 'Approved', 'Rejected', 'Withdrawn']

// Fixed status roles - reserved for state only, never reused as a generic
// series color elsewhere on the page.
const STATUS_DONUT_COLOR = {
    Draft: '#6B7280',
    Submitted: '#fab219',
    Approved: '#0ca30c',
    Rejected: '#d03b3b',
    Withdrawn: '#94A3B8',
}

// Brand navy, matching the sign-in/sign-up side panel - used for plain
// volume/info metrics that aren't a status.
const BRAND_NAVY = '#114B6B'

function HeroStat({ label, value, icon, iconColor, bg }) {
    return (
        <Card>
            <div className="flex items-center gap-4">
                <Avatar
                    size={50}
                    icon={icon}
                    style={{ backgroundColor: bg, color: iconColor }}
                />
                <div>
                    <div className="text-2xl font-bold text-gray-800 dark:text-gray-100">
                        {value}
                    </div>
                    <div className="text-sm text-gray-500">{label}</div>
                </div>
            </div>
        </Card>
    )
}

export default function Dashboard() {
    const [summary, setSummary] = useState(null)
    const [trend, setTrend] = useState([])
    const [avgApproval, setAvgApproval] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        Promise.all([
            apiGetDashboardSummary(),
            apiGetDashboardTrend(6),
            apiGetDashboardAvgApprovalTime(),
        ])
            .then(([summaryData, trendData, avgData]) => {
                setSummary(summaryData)
                setTrend(trendData || [])
                setAvgApproval(avgData)
            })
            .finally(() => setLoading(false))
    }, [])

    if (loading) {
        return (
            <div className="flex justify-center py-10">
                <Spinner size={32} />
            </div>
        )
    }

    const statusCounts = summary?.vesselApplications || {}
    const totalApplications = STATUS_ORDER.reduce(
        (sum, status) => sum + (statusCounts[status] || 0),
        0,
    )

    const donutLabels = STATUS_ORDER.filter((s) => (statusCounts[s] || 0) > 0)
    const donutSeries = donutLabels.map((s) => statusCounts[s] || 0)
    const donutColors = donutLabels.map((s) => STATUS_DONUT_COLOR[s])

    return (
        <div className="flex flex-col gap-4">
            <div>
                <h3>Dashboard</h3>
                <p className="text-gray-500">
                    Overview of vessel entry applications and review activity
                </p>
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <HeroStat
                    label="Total Applications"
                    value={totalApplications}
                    icon={<PiClipboardTextDuotone />}
                    bg={`${BRAND_NAVY}1A`}
                    iconColor={BRAND_NAVY}
                />
                <HeroStat
                    label="Awaiting Review"
                    value={statusCounts.Submitted || 0}
                    icon={<PiHourglassMediumDuotone />}
                    bg={`${STATUS_DONUT_COLOR.Submitted}22`}
                    iconColor="#92600d"
                />
                <HeroStat
                    label="Approved"
                    value={statusCounts.Approved || 0}
                    icon={<PiCheckCircleDuotone />}
                    bg={`${STATUS_DONUT_COLOR.Approved}1A`}
                    iconColor={STATUS_DONUT_COLOR.Approved}
                />
                <HeroStat
                    label="Rejected"
                    value={statusCounts.Rejected || 0}
                    icon={<PiXCircleDuotone />}
                    bg={`${STATUS_DONUT_COLOR.Rejected}1A`}
                    iconColor={STATUS_DONUT_COLOR.Rejected}
                />
            </div>

            <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
                <Card className="lg:col-span-2">
                    <h5 className="mb-4">Applications Submitted (Last 6 Months)</h5>
                    {trend.length > 0 ? (
                        <Chart
                            type="area"
                            series={[
                                {
                                    name: 'Applications',
                                    data: trend.map((p) => p.count),
                                },
                            ]}
                            xAxis={trend.map((p) => p.month)}
                            height={300}
                            customOptions={{
                                colors: [BRAND_NAVY],
                            }}
                        />
                    ) : (
                        <div className="flex h-[300px] items-center justify-center">
                            <p className="text-sm text-gray-500">
                                No submissions in this period yet.
                            </p>
                        </div>
                    )}
                </Card>

                <Card>
                    <h5 className="mb-4">Status Breakdown</h5>
                    {donutSeries.length > 0 ? (
                        <Chart
                            type="donut"
                            series={donutSeries}
                            height={230}
                            customOptions={{
                                labels: donutLabels,
                                colors: donutColors,
                                legend: { position: 'bottom' },
                            }}
                            donutTitle="Total"
                            donutText={String(totalApplications)}
                        />
                    ) : (
                        <div className="flex h-[230px] items-center justify-center">
                            <p className="text-sm text-gray-500">No data yet.</p>
                        </div>
                    )}
                </Card>
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <Card>
                    <div className="flex items-center gap-4">
                        <Avatar
                            size={50}
                            icon={<PiBuildingsDuotone />}
                            style={{ backgroundColor: `${BRAND_NAVY}1A`, color: BRAND_NAVY }}
                        />
                        <div>
                            <div className="text-2xl font-bold text-gray-800 dark:text-gray-100">
                                {summary?.pendingOrganizations || 0}
                            </div>
                            <div className="text-sm text-gray-500">
                                Organization profiles awaiting decision
                            </div>
                        </div>
                    </div>
                </Card>

                <Card>
                    <div className="flex items-center gap-4">
                        <Avatar
                            size={50}
                            icon={<PiTimerDuotone />}
                            style={{ backgroundColor: '#475569 1A'.replace(' ', ''), color: '#475569' }}
                        />
                        <div>
                            <div className="text-2xl font-bold text-gray-800 dark:text-gray-100">
                                {avgApproval?.avgDays != null
                                    ? `${avgApproval.avgDays.toFixed(1)} days`
                                    : '-'}
                            </div>
                            <div className="text-sm text-gray-500">
                                Average approval time
                            </div>
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    )
}