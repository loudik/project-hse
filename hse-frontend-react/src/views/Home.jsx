import { Link } from 'react-router'
import Card from '@/components/ui/Card'
import Avatar from '@/components/ui/Avatar'
import {
    PiPlusCircleDuotone,
    PiListChecksDuotone,
    PiMagnifyingGlassDuotone,
    PiCalendarPlusDuotone,
    PiChartBarDuotone,
    PiBuildingsDuotone,
} from 'react-icons/pi'

const BRAND_NAVY = '#114B6B'

const quickLinks = [
    {
        title: 'Submit a new entry',
        description: 'Start a vessel or helicopter entry application.',
        icon: <PiPlusCircleDuotone />,
        to: '/vessel',
    },
    {
        title: 'My applications',
        description: 'Track the status of applications you have submitted.',
        icon: <PiListChecksDuotone />,
        to: '/vessel/list',
    },
    {
        title: 'Review queue',
        description: 'Applications waiting for HSE review and decision.',
        icon: <PiMagnifyingGlassDuotone />,
        to: '/vessel/review',
    },
    {
        title: 'Extension requests',
        description: 'Requests to extend an approved entry authorisation.',
        icon: <PiCalendarPlusDuotone />,
        to: '/vessel/extension-requests',
    },
    {
        title: 'Dashboard',
        description: 'Overview of application volume and review activity.',
        icon: <PiChartBarDuotone />,
        to: '/dashboard',
    },
    {
        title: 'Organization profile',
        description: 'Manage your organization details and documents.',
        icon: <PiBuildingsDuotone />,
        to: '/organization/profile',
    },
]

const Home = () => {
    return (
        <div className="flex flex-col gap-6">
            <div>
                <h3>Welcome to ANP HSE</h3>
                <p className="text-gray-500">
                    Vessel and helicopter entry authorisation for Timor-Leste&apos;s
                    petroleum contract area.
                </p>
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {quickLinks.map((item) => (
                    <Link key={item.to} to={item.to}>
                        <Card className="h-full hover:shadow-lg transition-shadow cursor-pointer">
                            <div className="flex items-start gap-4">
                                <Avatar
                                    size={48}
                                    icon={item.icon}
                                    style={{
                                        backgroundColor: `${BRAND_NAVY}1A`,
                                        color: BRAND_NAVY,
                                    }}
                                />
                                <div>
                                    <div className="font-semibold text-gray-800 dark:text-gray-100">
                                        {item.title}
                                    </div>
                                    <div className="text-sm text-gray-500 mt-1">
                                        {item.description}
                                    </div>
                                </div>
                            </div>
                        </Card>
                    </Link>
                ))}
            </div>
        </div>
    )
}

export default Home