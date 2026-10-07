import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import { countryList } from '@/constants/countries.constant'
import Checkbox from '@/components/ui/Checkbox'
import Radio from '@/components/ui/Radio'
import { FormItem } from '@/components/ui/Form'

const VESSEL_TYPE_OPTIONS = [
    'Drilling rigs (Semi-submersible, Jack-up & Drillship)',
    'Platform Supply Vessel (PSVs)',
    'Anchor Handling Tug Supply (AHTS) Vessel',
    'Pipelay vessels',
    'Dive Support Vessels (DSVs)',
    'Offshore accommodation Vessels (ASV)',
    'Floating Production & Offloading (FSO)',
    'Production Storage & Offloading (FPSO)',
    'Infield Support Vessel (ISV)',
    'Diving Support Vessel (DSV)',
    'Remote Operated Vehicle Support Vessel',
    'Offtake tanker (LPG, Condensate, Crude Oil)',
    'Barge',
    'Walk-to-Walk Vessel (W2W)',
    'Towing vessels',
    'Seismic survey vessel',
    'Crane vessel',
    'Tugboats',
    'Multipurpose support vessel (MPSV)',
    'Chase vessels',
]

const OPERATION_MODES = [
    { value: 'DP1', label: 'Dynamic positioning vessel - DP1' },
    { value: 'DP2', label: 'Dynamic positioning vessel - DP2' },
    { value: 'DP3', label: 'Dynamic positioning vessel - DP3' },
    { value: 'Non-DP', label: 'Non-DP Vessel (Dynamic positioning)' },
]

export default function StepVesselDescription({ values, onChange }) {
    const scopeOfWork = values.scopeOfWork || []

    const toggleScope = (option) => {
        const next = scopeOfWork.includes(option)
            ? scopeOfWork.filter((o) => o !== option)
            : [...scopeOfWork, option]
        onChange({ scopeOfWork: next })
    }

    return (
        <div>
            <div className="mb-6 rounded-lg border border-gray-200 p-4">
                <p className="mb-1 text-sm text-gray-500">
                    Note: this online application is only applicable for a single vessel. If you wish to
                    apply for multiple vessels, please repeat the process for each vessel.
                </p>
                <p className="mb-3 font-semibold">
                    Please choose the following scope of work. You may tick more than one box if the
                    vessel intends to carry out more than one scope of work.
                </p>

                <div className="mb-4 grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2">
                    {VESSEL_TYPE_OPTIONS.map((option) => (
                        <Checkbox
                            key={option}
                            checked={scopeOfWork.includes(option)}
                            onChange={() => toggleScope(option)}
                        >
                            {option}
                        </Checkbox>
                    ))}
                </div>
                <div className="flex items-center gap-3">
                    <Checkbox checked={!!values.scopeOfWorkOther} readOnly>
                        Others, please specify
                    </Checkbox>
                    <Input
                        size="sm"
                        className="max-w-xs"
                        value={values.scopeOfWorkOther || ''}
                        onChange={(e) => onChange({ scopeOfWorkOther: e.target.value })}
                    />
                </div>

                <p className="mt-5 mb-2 font-semibold">
                    Please select the following vessel operation mode/classification, if applicable
                </p>
                <div className="flex flex-col gap-2">
                    {OPERATION_MODES.map((mode) => (
                        <label key={mode.value} className="flex cursor-pointer items-center gap-2">
                            <Radio
                                checked={values.operationMode === mode.value}
                                onChange={() => onChange({ operationMode: mode.value })}
                            />
                            {mode.label}
                        </label>
                    ))}
                </div>
            </div>

            <div className="rounded-lg border border-gray-200 p-4">
                <h5 className="mb-1">Vessel Description</h5>
                <p className="mb-4 text-sm text-gray-500">
                    The Authorised Person is mandatory to provide vessel description in the following.
                </p>

                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <FormItem label="Name of Vessel">
                        <Input value={values.vesselName || ''} onChange={(e) => onChange({ vesselName: e.target.value })} />
                    </FormItem>
                    <FormItem label="IMO Number">
                        <Input value={values.vesselImoNumber || ''} onChange={(e) => onChange({ vesselImoNumber: e.target.value })} />
                    </FormItem>
                    <FormItem label="Vessel Owner/Contractor">
                        <Input value={values.vesselOwner || ''} onChange={(e) => onChange({ vesselOwner: e.target.value })} />
                    </FormItem>
                    <FormItem label="Vessel Type">
                        <Input value={values.vesselType || ''} onChange={(e) => onChange({ vesselType: e.target.value })} />
                    </FormItem>
                    <FormItem label="Flag State">
                        <Select
                            options={countryList}
                            placeholder="Select flag state"
                            value={countryList.find((c) => c.value === values.flagState) || null}
                            onChange={(option) => onChange({ flagState: option?.value || '' })}
                        />
                    </FormItem>
                    {/* <FormItem label="Port of Registry">
                        <Input value={values.portOfRegistry || ''} onChange={(e) => onChange({ portOfRegistry: e.target.value })} />
                    </FormItem> */}
                    <FormItem label="Port of Registry">
                        <Select
                            options={countryList}
                            placeholder="Select Port of Registry"
                            value={countryList.find((c) => c.value === values.portOfRegistry) || null}
                            onChange={(option) => onChange({ portOfRegistry: option?.value || '' })}
                        />
                    </FormItem>
                    <FormItem label="Classification Society">
                        <Input value={values.classificationSociety || ''} onChange={(e) => onChange({ classificationSociety: e.target.value })} />
                    </FormItem>
                    <FormItem label="Class ID Number">
                        <Input value={values.classIdNumber || ''} onChange={(e) => onChange({ classIdNumber: e.target.value })} />
                    </FormItem>
                    <FormItem label="Length Overall (LOA)">
                        <Input value={values.lengthOverall || ''} onChange={(e) => onChange({ lengthOverall: e.target.value })} />
                    </FormItem>
                    <FormItem label="Draft">
                        <Input value={values.draftValue || ''} onChange={(e) => onChange({ draftValue: e.target.value })} />
                    </FormItem>
                    <FormItem label="Gross Tonnage">
                        <Input value={values.grossTonnage || ''} onChange={(e) => onChange({ grossTonnage: e.target.value })} />
                    </FormItem>
                    <FormItem label="Call Sign">
                        <Input value={values.callSign || ''} onChange={(e) => onChange({ callSign: e.target.value })} />
                    </FormItem>
                </div>
            </div>
        </div>
    )
}