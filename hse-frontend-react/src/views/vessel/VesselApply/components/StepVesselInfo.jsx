import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import { FormItem } from '@/components/ui/Form'

const VESSEL_TYPE_OPTIONS = [
    { value: 'Cargo', label: 'Cargo' },
    { value: 'Tanker', label: 'Tanker' },
    { value: 'Supply Vessel', label: 'Supply Vessel' },
    { value: 'Passenger', label: 'Passenger' },
    { value: 'Other', label: 'Other' },
]

export default function StepVesselInfo({ values, onChange }) {
    return (
        <div>
            <FormItem label="Vessel name">
                <Input
                    value={values.vesselName}
                    onChange={(e) => onChange({ vesselName: e.target.value })}
                    placeholder="Vessel name"
                />
            </FormItem>
            <FormItem label="Vessel type">
                <Select
                    options={VESSEL_TYPE_OPTIONS}
                    value={VESSEL_TYPE_OPTIONS.find((o) => o.value === values.vesselType) || null}
                    onChange={(option) => onChange({ vesselType: option?.value || '' })}
                    placeholder="Select vessel type"
                />
            </FormItem>
        </div>
    )
}