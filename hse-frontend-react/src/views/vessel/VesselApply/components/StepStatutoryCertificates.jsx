import { useEffect, useState } from 'react'
import DocumentUploadRow from './DocumentUploadRow'
import { apiListVesselDocuments } from '@/services/VesselService'

const STATUTORY_CERTS = [
    'Registration certificate',
    'Insurance certificate',
    'Class certificate',
    'ISM certificate',
    'ISPS declaration/certificate',
    'Cargo ship safety construction certificate',
    'Cargo ship safety radio certificate',
    'Document of compliance with the special requirements for ship carrying Dangerous Goods',
    'Minimum safe manning documents',
    'International certificate of fitness for the carriage of dangerous chemicals in bulk',
    'Document of compliance (ISM)',
    'International ship security certificate',
    'Fire control and safety plan',
    'Safety data sheet (SDS)',
    'IOPP certificate',
    'IAPP certificate',
    'Shipboard oil pollution emergency plan (SOPEP or SMPEP)',
    'Certificate of fitness for offshore support vessel',
    'International sewage pollution prevention certificate',
    'International load line certificate',
    'International anti-fouling certificate',
    'Dynamic positioning certificate',
    'Ship sanitation control exemption certificate',
    'Radio station licence',
    'Medical inventory/pharmacy certificate',
    'Crane and cargo gear certificate',
    'Seaworthiness certificate',
    'Helideck certificate',
    'Asbestos declaration & asbestos condition conformity certificate',
    'Ballast Water Exchange Certificate',
]

const SUPPORTING_DOCS = [
    'Class status survey report',
    'Marine vessel vetting summary',
    'Latest third-party inspection/audit report (OVID or IMCA report)',
    'FMEA report',
    'Latest DP Annual Trial Report',
    'Vessel crew competency matrix',
    'Vessel Specification Sheet',
    'Vessel deck and profile plan',
]

function slugify(label) {
    return label
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '_')
        .replace(/^_|_$/g, '')
}

export default function StepStatutoryCertificates({ applicationId }) {
    const [documents, setDocuments] = useState({})
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        apiListVesselDocuments(applicationId)
            .then((list) => {
                const byKey = {}
                ;(list || []).forEach((d) => {
                    byKey[`${d.category}:${d.documentKey}`] = d
                })
                setDocuments(byKey)
            })
            .finally(() => setLoading(false))
    }, [applicationId])

    if (loading) return <p className="text-gray-500">Loading...</p>

    return (
        <div>
            <div className="mb-6">
                <p className="mb-3 font-semibold">
                    The Authorised Person shall submit the following list of minimum statutory
                    certificates which may be relevant to each vessel depending on the purpose, size and
                    type of the vessel required for vessel entry authorisation.
                </p>
                <div className="rounded-lg border border-gray-200">
                    {STATUTORY_CERTS.map((label) => {
                        const key = slugify(label)
                        return (
                            <DocumentUploadRow
                                key={key}
                                applicationId={applicationId}
                                category="statutory_certificate"
                                documentKey={key}
                                label={label}
                                showDates
                                existingDoc={documents[`statutory_certificate:${key}`]}
                            />
                        )
                    })}
                </div>
            </div>

            <div>
                <p className="mb-3 font-semibold">
                    Does the vessel have all the necessary supporting documentation? If Yes, please attach
                    the document.
                </p>
                <div className="rounded-lg border border-gray-200">
                    {SUPPORTING_DOCS.map((label) => {
                        const key = slugify(label)
                        return (
                            <DocumentUploadRow
                                key={key}
                                applicationId={applicationId}
                                category="supporting_document"
                                documentKey={key}
                                label={label}
                                existingDoc={documents[`supporting_document:${key}`]}
                            />
                        )
                    })}
                </div>
            </div>
        </div>
    )
}