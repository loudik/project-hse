import { useEffect, useState } from 'react'
import Input from '@/components/ui/Input'
import Radio from '@/components/ui/Radio'
import { FormItem } from '@/components/ui/Form'
import DocumentUploadRow from './DocumentUploadRow'
import { apiListVesselDocuments } from '@/services/VesselService'

function YesNoRadio({ label, value, onChange }) {
    return (
        <div className="mb-3 flex items-center justify-between border-b border-gray-100 py-2">
            <span className="text-sm">{label}</span>
            <div className="flex gap-4">
                <label className="flex cursor-pointer items-center gap-1.5">
                    <Radio checked={value === true} onChange={() => onChange(true)} />
                    <span className="text-sm">Yes</span>
                </label>
                <label className="flex cursor-pointer items-center gap-1.5">
                    <Radio checked={value === false} onChange={() => onChange(false)} />
                    <span className="text-sm">No</span>
                </label>
            </div>
        </div>
    )
}

export default function StepDeclarations({ applicationId, values, onChange }) {
    const [documents, setDocuments] = useState({})

    useEffect(() => {
        apiListVesselDocuments(applicationId).then((list) => {
            const byKey = {}
            ;(list || []).forEach((d) => {
                byKey[`${d.category}:${d.documentKey}`] = d
            })
            setDocuments(byKey)
        })
    }, [applicationId])

    return (
        <div>
            <div className="mb-6 rounded-lg border border-gray-200 p-4">
                <h5 className="mb-3">The Authorised Person must declare the following:</h5>
                <YesNoRadio
                    label="No major deficiencies"
                    value={values.noMajorDeficiencies}
                    onChange={(v) => onChange({ noMajorDeficiencies: v })}
                />
                <YesNoRadio
                    label="No detention within last 12 months"
                    value={values.noDetention12Months}
                    onChange={(v) => onChange({ noDetention12Months: v })}
                />
                <YesNoRadio
                    label="Safety equipment operational"
                    value={values.safetyEquipmentOperational}
                    onChange={(v) => onChange({ safetyEquipmentOperational: v })}
                />
                <YesNoRadio
                    label="Firefighting systems operational"
                    value={values.firefightingOperational}
                    onChange={(v) => onChange({ firefightingOperational: v })}
                />
                <YesNoRadio
                    label="Lifesaving appliances operational"
                    value={values.lifesavingOperational}
                    onChange={(v) => onChange({ lifesavingOperational: v })}
                />

                <p className="mt-4 mb-2 text-sm font-semibold">
                    Please provide the number of crew members available on board the vessel and attach the
                    document.
                </p>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <FormItem label="Number of Vessel Crew">
                        <Input
                            type="number"
                            value={values.crewCount ?? ''}
                            onChange={(e) => onChange({ crewCount: e.target.value ? Number(e.target.value) : null })}
                        />
                    </FormItem>
                    <FormItem label="Number of Survey/Project Crew">
                        <Input
                            type="number"
                            value={values.surveyCrewCount ?? ''}
                            onChange={(e) => onChange({ surveyCrewCount: e.target.value ? Number(e.target.value) : null })}
                        />
                    </FormItem>
                </div>
                <DocumentUploadRow
                    applicationId={applicationId}
                    category="crew_document"
                    documentKey="crew_list"
                    label="Crew list document"
                    existingDoc={documents['crew_document:crew_list']}
                />
            </div>

            <div className="mb-6 rounded-lg border border-gray-200 p-4">
                <h5 className="mb-3">Please provide cargo information and submit the document</h5>
                <YesNoRadio
                    label="Cargo onboard"
                    value={values.cargoOnboard}
                    onChange={(v) => onChange({ cargoOnboard: v })}
                />
                <FormItem label="Cargo Description">
                    <Input
                        value={values.cargoDescription || ''}
                        onChange={(e) => onChange({ cargoDescription: e.target.value })}
                    />
                </FormItem>
                <YesNoRadio
                    label="Hazardous Cargo"
                    value={values.hazardousCargo}
                    onChange={(v) => onChange({ hazardousCargo: v })}
                />
                <DocumentUploadRow
                    applicationId={applicationId}
                    category="cargo_document"
                    documentKey="dangerous_good_declaration"
                    label="Dangerous Good Declaration"
                    existingDoc={documents['cargo_document:dangerous_good_declaration']}
                />
                <DocumentUploadRow
                    applicationId={applicationId}
                    category="cargo_document"
                    documentKey="cargo_manifest"
                    label="Cargo Manifest"
                    existingDoc={documents['cargo_document:cargo_manifest']}
                />
            </div>

            <div className="mb-6 rounded-lg border border-gray-200 p-4">
                <h5 className="mb-3">Please declare the following environmental compliance</h5>
                <YesNoRadio
                    label="Waste Discharge"
                    value={values.wasteDischarge}
                    onChange={(v) => onChange({ wasteDischarge: v })}
                />
                <YesNoRadio
                    label="Oily Waste Onboard"
                    value={values.oilyWasteOnboard}
                    onChange={(v) => onChange({ oilyWasteOnboard: v })}
                />
                <YesNoRadio
                    label="Sewage Disposal Required"
                    value={values.sewageDisposalRequired}
                    onChange={(v) => onChange({ sewageDisposalRequired: v })}
                />
            </div>

            <div className="rounded-lg border border-gray-200 p-4">
                <h5 className="mb-3">
                    Does the authorised person have clearance from the relevant government entity for
                    personnel, vessels, equipment, to enter Timor-Leste?
                </h5>
                <DocumentUploadRow
                    applicationId={applicationId}
                    category="clearance_document"
                    documentKey="immigration"
                    label="Immigration"
                    existingDoc={documents['clearance_document:immigration']}
                />
                <DocumentUploadRow
                    applicationId={applicationId}
                    category="clearance_document"
                    documentKey="customs"
                    label="Custom"
                    existingDoc={documents['clearance_document:customs']}
                />
                <DocumentUploadRow
                    applicationId={applicationId}
                    category="clearance_document"
                    documentKey="quarantine"
                    label="Quarantine"
                    existingDoc={documents['clearance_document:quarantine']}
                />
            </div>
        </div>
    )
}