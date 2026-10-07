import Checkbox from '@/components/ui/Checkbox'

const TYPES = ['Vessel', 'Helicopter', 'Personnel', 'Goods, equipments and materials']

export default function StepEntryApplicationType({ values, onChange }) {
    const selected = values.entryApplicationTypes || []

    const toggle = (type) => {
        const next = selected.includes(type)
            ? selected.filter((t) => t !== type)
            : [...selected, type]
        onChange({ entryApplicationTypes: next })
    }

    return (
        <div>
            <h5 className="mb-4">
                Entry to Contract Area or Area Authorised for Petroleum Activities
            </h5>
            <p className="mb-4 font-semibold">
                Please choose the following application for entry to Contract Area or Area
                Authorised to carry out Petroleum Activities
            </p>

            <div className="flex flex-col gap-3 rounded-lg border border-gray-200 p-4">
                {TYPES.map((type) => (
                    <Checkbox key={type} checked={selected.includes(type)} onChange={() => toggle(type)}>
                        {type}
                    </Checkbox>
                ))}
            </div>

            {selected.length > 0 && !selected.includes('Vessel') && (
                <p className="mt-4 text-sm text-amber-600">
                    Note: this application form currently only supports the "Vessel" entry type.
                    Support for Helicopter/Personnel/Goods entry will be added separately.
                </p>
            )}
        </div>
    )
}