import Input from '@/components/ui/Input'
import Radio from '@/components/ui/Radio'
import Checkbox from '@/components/ui/Checkbox'

const CONTRACT_TYPES = [
    { value: 'PSC', label: 'Production Sharing Contract (PSC)' },
    { value: 'Prospecting Authorisation', label: 'Prospecting Authorisation' },
    { value: 'Access Authorisation', label: 'Access Authorisation' },
    { value: 'Others', label: 'Others' },
]

export default function StepContractInfo({ values, onChange }) {
    const status = values.contractStatus || ''

    return (
        <div>
            <h5 className="mb-4">Information on Contract or Agreement</h5>
            <p className="mb-4 font-semibold">Do you have a valid contract with ANP?</p>

            <div className="mb-5 flex flex-col gap-4">
                {/* YES - has contract */}
                <div>
                    <label className="flex cursor-pointer items-center gap-2">
                        <Radio
                            checked={status === 'has_contract'}
                            onChange={() => onChange({ contractStatus: 'has_contract' })}
                        />
                        <span className="font-semibold">Yes</span>
                    </label>

                    {status === 'has_contract' && (
                        <div className="mt-3 ml-7 flex flex-col gap-2">
                            <p className="text-sm text-gray-500">
                                Please choose the applicable one and write the number (e.g. TL xx-xx)
                            </p>
                            {CONTRACT_TYPES.map((type) => (
                                <div key={type.value} className="flex items-center gap-3">
                                    <Checkbox
                                        checked={values.contractType === type.value}
                                        onChange={() => onChange({ contractType: type.value })}
                                    >
                                        {type.label}
                                    </Checkbox>
                                    {values.contractType === type.value && (
                                        <Input
                                            size="sm"
                                            className="max-w-xs"
                                            value={
                                                type.value === 'Others'
                                                    ? values.contractTypeOther || ''
                                                    : values.contractNumber || ''
                                            }
                                            onChange={(e) =>
                                                type.value === 'Others'
                                                    ? onChange({ contractTypeOther: e.target.value })
                                                    : onChange({ contractNumber: e.target.value })
                                            }
                                            placeholder={type.value === 'Others' ? 'Specify' : 'e.g. TL xx-xx'}
                                        />
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* INTEND TO OBTAIN */}
                <div>
                    <label className="flex cursor-pointer items-start gap-2">
                        <Radio
                            checked={status === 'intend_to_obtain'}
                            onChange={() => onChange({ contractStatus: 'intend_to_obtain' })}
                        />
                        <span className="font-semibold">
                            Intend to obtain Contract or Authorisation or Agreement from ANP and
                            currently under discussion in parallel with this entry application
                        </span>
                    </label>

                    {status === 'intend_to_obtain' && (
                        <div className="mt-3 ml-7 flex flex-col gap-2">
                            <p className="text-sm text-gray-500">
                                Please choose one type which is currently under discussion with ANP
                            </p>
                            {CONTRACT_TYPES.map((type) => (
                                <Checkbox
                                    key={type.value}
                                    checked={values.contractType === type.value}
                                    onChange={() => onChange({ contractType: type.value })}
                                >
                                    {type.label}
                                </Checkbox>
                            ))}
                        </div>
                    )}
                </div>

                {/* NO */}
                <div>
                    <label className="flex cursor-pointer items-center gap-2">
                        <Radio
                            checked={status === 'no_contract'}
                            onChange={() => onChange({ contractStatus: 'no_contract' })}
                        />
                        <span className="font-semibold">No</span>
                    </label>

                    {status === 'no_contract' && (
                        <div className="mt-3 ml-7">
                            <p className="mb-1.5 text-sm text-gray-500">
                                Please contact ANP on email address:
                            </p>
                            <Input
                                type="email"
                                className="max-w-xs"
                                value={values.contractContactEmail || ''}
                                onChange={(e) => onChange({ contractContactEmail: e.target.value })}
                                placeholder="Your contact email"
                            />
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}