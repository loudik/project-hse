import { useEffect, useState } from 'react'
import Input from '@/components/ui/Input'
import Checkbox from '@/components/ui/Checkbox'
import Alert from '@/components/ui/Alert'
import DocumentUploadRow from './DocumentUploadRow'
import { apiListVesselDocuments } from '@/services/VesselService'

const REGULATORY_DOCS = [
    { key: 'vessel_safety_case', label: 'Vessel Safety Case' },
    { key: 'safety_case_addendum', label: 'Safety Case Addendum or Bridging Document' },
    { key: 'construction_installation_safety_case', label: 'Construction, Installation & Commissioning Safety Case' },
    { key: 'production_operations_safety_case', label: 'Production/Operations Safety Case' },
    { key: 'decommissioning_safety_case', label: 'Decommissioning Safety Case' },
    { key: 'emergency_response_plan', label: 'Emergency Response Plan (ERP)' },
    { key: 'erp_bridging_document', label: 'Emergency Response Plan Bridging Document' },
    { key: 'diving_project_plan', label: 'Diving Project Plan & Diving Safety Management System' },
    { key: 'health_safety_plan', label: 'Health & Safety Plan' },
    { key: 'health_safety_plan_bridging', label: 'Health and Safety Plan Bridging Document' },
]

const CATEGORY = 'regulatory_document'

export default function StepRegulatoryDocuments({ applicationId }) {
    const [documents, setDocuments] = useState({})
    const [loading, setLoading] = useState(true)
    const [otherSpecify, setOtherSpecify] = useState('')

    useEffect(() => {
        apiListVesselDocuments(applicationId)
            .then((list) => {
                const byKey = {}
                ;(list || [])
                    .filter((d) => d.category === CATEGORY)
                    .forEach((d) => {
                        byKey[d.documentKey] = d
                        if (d.documentKey === 'other' && d.label) setOtherSpecify(d.label)
                    })
                setDocuments(byKey)
            })
            .finally(() => setLoading(false))
    }, [applicationId])

    if (loading) return <p className="text-gray-500">Loading...</p>

    return (
        <div>
            <Alert showIcon type="info" className="mb-5">
                <strong>Note:</strong> Operator is permitted to submit a vessel entry application online
                only after having submitted the main regulatory document to the ANP as part of an ongoing
                assessment, or after the ANP has approved the document.
            </Alert>

            <p className="mb-3 font-semibold">
                Have the main regulatory documents required by the relevant legal framework for the
                proposed operations/activities been submitted to the ANP for assessment and
                decision-making?
            </p>

            <div className="rounded-lg border border-gray-200">
                {REGULATORY_DOCS.map((doc) => (
                    <DocumentUploadRow
                        key={doc.key}
                        applicationId={applicationId}
                        category={CATEGORY}
                        documentKey={doc.key}
                        label={doc.label}
                        existingDoc={documents[doc.key]}
                    />
                ))}
            </div>

            <div className="mt-4 flex items-center gap-3">
                <Checkbox checked={!!otherSpecify} readOnly>
                    Other, please specify
                </Checkbox>
                <Input
                    size="sm"
                    className="max-w-xs"
                    value={otherSpecify}
                    onChange={(e) => setOtherSpecify(e.target.value)}
                    placeholder="Specify other document"
                />
            </div>
        </div>
    )
}