import Input from '@/components/ui/Input'
import Radio from '@/components/ui/Radio'
import Checkbox from '@/components/ui/Checkbox'
import { FormItem } from '@/components/ui/Form'
import DocumentUploadRow from './DocumentUploadRow'

const CONDITIONS = [
    {
        value: 'Initial',
        title: 'Initial Entry Authorisation to the Contract Area or Area Authorised to carry out petroleum operations',
        description:
            'Applicable for: first time entry of a vessel to the Contract Area, or entry of a vessel previously off-hired that is being brought back for petroleum operations.',
    },
    {
        value: 'Extension',
        title: 'Extension of Entry Authorisation to the Contract Area or Area Authorised to carry out petroleum operations',
        description:
            'Applicable for vessels, aircraft, personnel and equipment already in the JPDA that will remain in the field - e.g. annual extension, or extension following contract extension.',
    },
    {
        value: 'Substitute',
        title: 'Substitute Vessel',
        description:
            'Only applicable when a vessel with valid entry authorisation needs to be substituted with another vessel due to operational reasons (maintenance, drydock, etc). Can only be exercised under the same vessel contract.',
    },
]

export default function StepEntryCondition({ applicationId, values, onChange }) {
    return (
        <div>
            <p className="mb-4 font-semibold">Please choose the most applicable condition:</p>

            <div className="mb-6 flex flex-col gap-4">
                {CONDITIONS.map((cond) => (
                    <label key={cond.value} className="flex cursor-pointer items-start gap-3">
                        <Radio
                            checked={values.entryCondition === cond.value}
                            onChange={() => onChange({ entryCondition: cond.value })}
                        />
                        <div>
                            <div className="font-semibold text-primary">{cond.title}</div>
                            <div className="text-sm text-gray-500">{cond.description}</div>
                        </div>
                    </label>
                ))}
            </div>

            <div className="rounded-lg border border-gray-200 p-4">
                <h5 className="mb-4">Vessel Entry Application</h5>

                <FormItem label="Purpose of vessel entry request">
                    <Input
                        textArea
                        rows={2}
                        value={values.purpose || ''}
                        onChange={(e) => onChange({ purpose: e.target.value })}
                    />
                </FormItem>

                <FormItem label="Information on Contract or Agreement">
                    <Input
                        value={values.contractInfo || ''}
                        onChange={(e) => onChange({ contractInfo: e.target.value })}
                        placeholder="e.g. PSC TL-SO-T 19-12 and PSC TL-SO-T 19-13"
                    />
                </FormItem>

                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <FormItem label="Estimated date of arrival">
                        <Input
                            type="date"
                            value={values.entryDateFrom || ''}
                            onChange={(e) => onChange({ entryDateFrom: e.target.value })}
                        />
                    </FormItem>
                    <FormItem label="Estimated date of departure">
                        <Input
                            type="date"
                            value={values.entryDateTo || ''}
                            onChange={(e) => onChange({ entryDateTo: e.target.value })}
                        />
                    </FormItem>
                </div>
                <p className="mb-4 text-xs text-gray-400">
                    Maximum duration of operating in the Contract Area is 1 year from the date of entry
                    authorisation.
                </p>

                <FormItem label="Single or multiple entry to Contract Area or Area Authorised for petroleum operations">
                    <div className="flex gap-6">
                        <label className="flex cursor-pointer items-center gap-2">
                            <Checkbox
                                checked={values.entryType === 'single'}
                                onChange={() => onChange({ entryType: 'single' })}
                            >
                                Single
                            </Checkbox>
                        </label>
                        <label className="flex cursor-pointer items-center gap-2">
                            <Checkbox
                                checked={values.entryType === 'multiple'}
                                onChange={() => onChange({ entryType: 'multiple' })}
                            >
                                Multiple
                            </Checkbox>
                        </label>
                    </div>
                </FormItem>

                <div className="mt-2">
                    <div className="mb-1.5 text-sm font-semibold">Please upload your cover letter</div>
                    <DocumentUploadRow
                        applicationId={applicationId}
                        category="cover_letter"
                        documentKey="cover_letter"
                        label="Cover Letter"
                    />
                </div>
            </div>
        </div>
    )
}