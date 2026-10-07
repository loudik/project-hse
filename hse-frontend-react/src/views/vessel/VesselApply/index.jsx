import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router'
import Card from '@/components/ui/Card'
import Button from '@/components/ui/Button'
import Alert from '@/components/ui/Alert'
import Spinner from '@/components/ui/Spinner'
import Notification from '@/components/ui/Notification'
import toast from '@/components/ui/toast'
import {
    apiGetVesselApplication,
    apiUpdateVesselApplication,
    apiSubmitVesselApplication,
} from '@/services/VesselService'
import StepIndicator from './components/StepIndicator'
import StepEmailAddresses from './components/StepEmailAddresses'
import StepContractInfo from './components/StepContractInfo'
import StepEntryApplicationType from './components/StepEntryApplicationType'
import StepRegulatoryDocuments from './components/StepRegulatoryDocuments'
import StepEntryCondition from './components/StepEntryCondition'
import StepVesselDescription from './components/StepVesselDescription'
import StepStatutoryCertificates from './components/StepStatutoryCertificates'
import StepDeclarations from './components/StepDeclarations'
import StepReview from './components/StepReview'

const EMPTY_VALUES = {
    notificationEmails: [''],
    contractStatus: '',
    contractType: '',
    contractTypeOther: '',
    contractNumber: '',
    contractContactEmail: '',
    entryApplicationTypes: [],
    entryCondition: '',
    purpose: '',
    entryDateFrom: '',
    entryDateTo: '',
    entryType: '',
    scopeOfWork: [],
    scopeOfWorkOther: '',
    operationMode: '',
    vesselName: '',
    vesselImoNumber: '',
    vesselOwner: '',
    vesselType: '',
    flagState: '',
    portOfRegistry: '',
    classificationSociety: '',
    classIdNumber: '',
    lengthOverall: '',
    draftValue: '',
    grossTonnage: '',
    callSign: '',
    noMajorDeficiencies: null,
    noDetention12Months: null,
    safetyEquipmentOperational: null,
    firefightingOperational: null,
    lifesavingOperational: null,
    crewCount: null,
    surveyCrewCount: null,
    cargoOnboard: null,
    cargoDescription: '',
    hazardousCargo: null,
    wasteDischarge: null,
    oilyWasteOnboard: null,
    sewageDisposalRequired: null,
}

export default function VesselApply() {
    const { id } = useParams()
    const navigate = useNavigate()

    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)
    const [submitting, setSubmitting] = useState(false)
    const [step, setStep] = useState(1)
    const [status, setStatus] = useState('Draft')
    const [values, setValues] = useState(EMPTY_VALUES)
    const [error, setError] = useState('')

    useEffect(() => {
        apiGetVesselApplication(id)
            .then((app) => {
                setStatus(app.status)
                setValues({
                    notificationEmails: app.notificationEmails?.length ? app.notificationEmails : [''],
                    contractStatus: app.contractStatus || '',
                    contractType: app.contractType || '',
                    contractTypeOther: app.contractTypeOther || '',
                    contractNumber: app.contractNumber || '',
                    contractContactEmail: app.contractContactEmail || '',
                    entryApplicationTypes: app.entryApplicationTypes || [],
                    entryCondition: app.entryCondition || '',
                    purpose: app.purpose || '',
                    entryDateFrom: app.entryDateFrom || '',
                    entryDateTo: app.entryDateTo || '',
                    entryType: app.entryType || '',
                    scopeOfWork: app.scopeOfWork || [],
                    scopeOfWorkOther: app.scopeOfWorkOther || '',
                    operationMode: app.operationMode || '',
                    vesselName: app.vesselName || '',
                    vesselImoNumber: app.vesselImoNumber || '',
                    vesselOwner: app.vesselOwner || '',
                    vesselType: app.vesselType || '',
                    flagState: app.flagState || '',
                    portOfRegistry: app.portOfRegistry || '',
                    classificationSociety: app.classificationSociety || '',
                    classIdNumber: app.classIdNumber || '',
                    lengthOverall: app.lengthOverall || '',
                    draftValue: app.draftValue || '',
                    grossTonnage: app.grossTonnage || '',
                    callSign: app.callSign || '',
                    noMajorDeficiencies: app.noMajorDeficiencies ?? null,
                    noDetention12Months: app.noDetention12Months ?? null,
                    safetyEquipmentOperational: app.safetyEquipmentOperational ?? null,
                    firefightingOperational: app.firefightingOperational ?? null,
                    lifesavingOperational: app.lifesavingOperational ?? null,
                    crewCount: app.crewCount ?? null,
                    surveyCrewCount: app.surveyCrewCount ?? null,
                    cargoOnboard: app.cargoOnboard ?? null,
                    cargoDescription: app.cargoDescription || '',
                    hazardousCargo: app.hazardousCargo ?? null,
                    wasteDischarge: app.wasteDischarge ?? null,
                    oilyWasteOnboard: app.oilyWasteOnboard ?? null,
                    sewageDisposalRequired: app.sewageDisposalRequired ?? null,
                })
            })
            .catch(() => setError('Failed to load this application.'))
            .finally(() => setLoading(false))
    }, [id])

    const handleChange = (patch) => setValues((prev) => ({ ...prev, ...patch }))

    const saveDraft = async () => {
        setSaving(true)
        setError('')
        try {
            await apiUpdateVesselApplication(id, values)
            return true
        } catch (err) {
            setError(err?.response?.data?.message || 'Failed to save.')
            return false
        } finally {
            setSaving(false)
        }
    }

    const handleNext = async () => {
        const ok = await saveDraft()
        if (ok) setStep((s) => Math.min(s + 1, 9))
    }

    const handleBack = () => setStep((s) => Math.max(s - 1, 1))

    const handleSubmit = async () => {
        const ok = await saveDraft()
        if (!ok) return

        setSubmitting(true)
        setError('')
        try {
            await apiSubmitVesselApplication(id)
            toast.push(
                <Notification type="success" title="Submitted">
                    Your vessel entry application has been submitted for review.
                </Notification>,
            )
            navigate('/vessel/list')
        } catch (err) {
            const missing = err?.response?.data?.missing
            setError(
                missing
                    ? `Missing required fields: ${missing.join(', ')}`
                    : err?.response?.data?.message || 'Failed to submit application.',
            )
        } finally {
            setSubmitting(false)
        }
    }

    if (loading) {
        return (
            <div className="flex justify-center py-10">
                <Spinner size={32} />
            </div>
        )
    }

    const isReadOnly = status !== 'Draft'

    return (
        <Card>
            <div className="mb-1 flex items-center justify-between">
                <h4>Vessel Entry Application</h4>
                {isReadOnly && (
                    <span className="text-sm text-gray-500">
                        This application is {status} and can no longer be edited.
                    </span>
                )}
            </div>

            <StepIndicator current={step} />

            {error && (
                <Alert showIcon className="mb-4" type="danger">
                    {error}
                </Alert>
            )}

            <fieldset disabled={isReadOnly}>
                {step === 1 && <StepEmailAddresses values={values} onChange={handleChange} />}
                {step === 2 && <StepContractInfo values={values} onChange={handleChange} />}
                {step === 3 && <StepEntryApplicationType values={values} onChange={handleChange} />}
                {step === 4 && <StepRegulatoryDocuments applicationId={id} />}
                {step === 5 && (
                    <StepEntryCondition applicationId={id} values={values} onChange={handleChange} />
                )}
                {step === 6 && <StepVesselDescription values={values} onChange={handleChange} />}
                {step === 7 && <StepStatutoryCertificates applicationId={id} />}
                {step === 8 && (
                    <StepDeclarations applicationId={id} values={values} onChange={handleChange} />
                )}
                {step === 9 && <StepReview values={values} />}
            </fieldset>

            <div className="mt-8 flex justify-between">
                <Button onClick={() => navigate('/vessel/list')}>Cancel</Button>

                <div className="flex gap-2">
                    {step > 1 && (
                        <Button onClick={handleBack} disabled={saving}>
                            Back
                        </Button>
                    )}
                    {step < 9 && !isReadOnly && (
                        <Button variant="solid" loading={saving} onClick={handleNext}>
                            Next
                        </Button>
                    )}
                    {step === 9 && !isReadOnly && (
                        <Button variant="solid" loading={submitting} onClick={handleSubmit}>
                            Submit Application
                        </Button>
                    )}
                </div>
            </div>
        </Card>
    )
}