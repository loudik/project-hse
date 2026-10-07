import { useEffect, useRef, useState } from 'react'
import { useParams, useNavigate, useLocation } from 'react-router'
import Card from '@/components/ui/Card'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Spinner from '@/components/ui/Spinner'
import Tag from '@/components/ui/Tag'
import Notification from '@/components/ui/Notification'
import toast from '@/components/ui/toast'
import HistoryTimeline from '@/components/shared/HistoryTimeline'
import Select from '@/components/ui/Select'
import Dialog from '@/components/ui/Dialog'
import SignaturePad from '@/components/shared/SignaturePad'
import {
    apiGetVesselApplication,
    apiListVesselDocuments,
    apiGetVesselApplicationForReview,
    apiListVesselDocumentsForReview,
    apiDecideVesselApplication,
    apiGetVesselApplicationHistory,
    apiDownloadVesselDocument,
    apiUploadVesselDocument,
    apiListHSEOfficers,
    apiAssignHSEOfficer,
    apiGetAuthorisationLetter,
    apiIssueAuthorisationLetter,
    apiDownloadAuthorisationLetter,
    apiDownloadVesselOutcomeDocument,
    apiViewVesselOutcomeDocument,
    apiRequestExtension,
    apiListExtensionRequestsForApplication,
    apiListStaffUsers,
    apiRequestStaffReview,
    apiListReviewRequests,
    apiAddReviewComment,
    apiListReviewComments,
    apiNotifyANPReady,
} from '@/services/VesselService'

const STATUS_TAG = {
    Draft: 'bg-gray-100 text-gray-700',
    Submitted: 'bg-amber-100 text-amber-700',
    Approved: 'bg-emerald-100 text-emerald-700',
    Rejected: 'bg-red-100 text-red-700',
    Withdrawn: 'bg-gray-200 text-gray-500',
}

function Field({ label, value }) {
    return (
        <div className="mb-2.5 grid grid-cols-1 gap-1 border-b border-gray-100 py-1.5 sm:grid-cols-3 sm:gap-4">
            <div className="text-sm font-semibold text-gray-500 sm:col-span-1">{label}</div>
            <div className="text-sm text-gray-800 sm:col-span-2">
                {value === true ? 'Yes' : value === false ? 'No' : value || '-'}
            </div>
        </div>
    )
}

function DocField({ doc, label }) {
    const { id: appId } = useParams()
    const location = useLocation()
    const isReviewMode = location.pathname.startsWith('/vessel/review')

    const [viewing, setViewing] = useState(false)
    const [showReplace, setShowReplace] = useState(false)
    const [replacing, setReplacing] = useState(false)
    const [newDateIssued, setNewDateIssued] = useState('')
    const [newDateExpired, setNewDateExpired] = useState('')
    const fileInputRef = useRef(null)

    if (!doc) return <Field label={label} value="Not provided" />
    if (doc.notApplicable) return <Field label={label} value={`N/A - ${doc.naReason || 'no reason given'}`} />
    const dateInfo = doc.dateIssued ? ` (${doc.dateIssued} to ${doc.dateExpired || 'N/A'})` : ''

        const [showPreview, setShowPreview] = useState(false)
    const [previewUrl, setPreviewUrl] = useState(null)

    const handleView = async () => {
        if (showPreview) {
            setShowPreview(false)
            return
        }
        if (previewUrl) {
            setShowPreview(true)
            return
        }
        setViewing(true)
        try {
            const resp = await apiDownloadVesselDocument(appId, doc.id)
            setPreviewUrl(resp.url)
            setShowPreview(true)
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to open document.'}
                </Notification>,
            )
        } finally {
            setViewing(false)
        }
    }

    const handleReplaceFile = async (e) => {
        const file = e.target.files?.[0]
        if (!file) return
        setReplacing(true)
        try {
            const formData = new FormData()
            formData.append('category', doc.category)
            formData.append('documentKey', doc.documentKey)
            formData.append('label', doc.label || '')
            formData.append('notApplicable', 'false')
            formData.append('naReason', '')
            formData.append('dateIssued', newDateIssued || doc.dateIssued || '')
            formData.append('dateExpired', newDateExpired || doc.dateExpired || '')
            formData.append('file', file)

            await apiUploadVesselDocument(appId, formData)
            toast.push(
                <Notification type="success" title="Uploaded">
                    Document replaced successfully.
                </Notification>,
            )
            // Simplest reliable way to refresh the document list, dates, and
            // history timeline after a replace - this action is rare enough
            // (certificate renewal) that a full reload is acceptable.
            window.location.reload()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to upload document.'}
                </Notification>,
            )
        } finally {
            setReplacing(false)
        }
    }

    return (
        <div className="mb-2.5 grid grid-cols-1 gap-1 border-b border-gray-100 py-1.5 sm:grid-cols-3 sm:gap-4">
            <div className="text-sm font-semibold text-gray-500 sm:col-span-1">{label}</div>
            <div className="sm:col-span-2">
                <div className="flex flex-wrap items-center gap-2 text-sm text-gray-800">
                    <span>
                        {doc.fileName ? `Submitted: ${doc.fileName}${dateInfo}` : 'Not submitted'}
                    </span>
                    {doc.fileName && doc.id && (
                        <Button
                            size="xs"
                            variant="plain"
                            loading={viewing}
                            onClick={handleView}
                            className="print:hidden"
                        >
                            {showPreview ? 'Hide preview' : 'View document'}
                        </Button>
                    )}
                    {/* Only the applicant can replace, and only when there's
                        already a file to replace - not shown on the ANP HSE
                        review screen. */}
                    {!isReviewMode && doc.fileName && (
                        <Button
                            size="xs"
                            variant="plain"
                            className="print:hidden"
                            onClick={() => setShowReplace((v) => !v)}
                        >
                            {showReplace ? 'Cancel' : 'Replace'}
                        </Button>
                    )}
                </div>

                {showReplace && (
                    <div className="mt-2 flex flex-wrap items-center gap-2 rounded-lg bg-gray-50 p-2 print:hidden">
                        <Input
                            size="sm"
                            type="date"
                            value={newDateIssued}
                            onChange={(e) => setNewDateIssued(e.target.value)}
                            placeholder="New issue date"
                        />
                        <Input
                            size="sm"
                            type="date"
                            value={newDateExpired}
                            onChange={(e) => setNewDateExpired(e.target.value)}
                            placeholder="New expiry date"
                        />
                        <input
                            ref={fileInputRef}
                            type="file"
                            className="hidden"
                            onChange={handleReplaceFile}
                        />
                        <Button
                            size="xs"
                            variant="solid"
                            loading={replacing}
                            onClick={() => fileInputRef.current?.click()}
                        >
                            Choose renewed file
                        </Button>
                    </div>
                )}
                {showPreview && previewUrl && (
                    <div className="mt-2 print:hidden">
                        <iframe
                            src={previewUrl}
                            title={doc.fileName}
                            className="h-[600px] w-full rounded-lg border border-gray-200"
                        />
                    </div>
                )}
            </div>
        </div>
    )
}

const STATUTORY_CERTS = [
    ['registration_certificate', 'Registration certificate'],
    ['insurance_certificate', 'Insurance certificate'],
    ['class_certificate', 'Class certificate'],
    ['ism_certificate', 'ISM certificate'],
    ['isps_declaration_certificate', 'ISPS declaration/certificate'],
    ['cargo_ship_safety_construction_certificate', 'Cargo ship safety construction certificate'],
    ['cargo_ship_safety_radio_certificate', 'Cargo ship safety radio certificate'],
    ['document_of_compliance_with_the_special_requirements_for_ship_carrying_dangerous_goods', 'Document of compliance (Dangerous Goods)'],
    ['minimum_safe_manning_documents', 'Minimum safe manning documents'],
    ['international_certificate_of_fitness_for_the_carriage_of_dangerous_chemicals_in_bulk', 'International certificate of fitness (dangerous chemicals)'],
    ['document_of_compliance_ism', 'Document of compliance (ISM)'],
    ['international_ship_security_certificate', 'International ship security certificate'],
    ['fire_control_and_safety_plan', 'Fire control and safety plan'],
    ['safety_data_sheet_sds', 'Safety data sheet (SDS)'],
    ['iopp_certificate', 'IOPP certificate'],
    ['iapp_certificate', 'IAPP certificate'],
    ['shipboard_oil_pollution_emergency_plan_sopep_or_smpep', 'Shipboard oil pollution emergency plan (SOPEP/SMPEP)'],
    ['certificate_of_fitness_for_offshore_support_vessel', 'Certificate of fitness for offshore support vessel'],
    ['international_sewage_pollution_prevention_certificate', 'International sewage pollution prevention certificate'],
    ['international_load_line_certificate', 'International load line certificate'],
    ['international_anti_fouling_certificate', 'International anti-fouling certificate'],
    ['dynamic_positioning_certificate', 'Dynamic positioning certificate'],
    ['ship_sanitation_control_exemption_certificate', 'Ship sanitation control exemption certificate'],
    ['radio_station_licence', 'Radio station licence'],
    ['medical_inventory_pharmacy_certificate', 'Medical inventory/pharmacy certificate'],
    ['crane_and_cargo_gear_certificate', 'Crane and cargo gear certificate'],
    ['seaworthiness_certificate', 'Seaworthiness certificate'],
    ['helideck_certificate', 'Helideck certificate'],
    ['asbestos_declaration_asbestos_condition_conformity_certificate', 'Asbestos declaration & conformity certificate'],
    ['ballast_water_exchange_certificate', 'Ballast Water Exchange Certificate'],
]

const SUPPORTING_DOCS = [
    ['class_status_survey_report', 'Class status survey report'],
    ['marine_vessel_vetting_summary', 'Marine vessel vetting summary'],
    ['latest_third_party_inspection_audit_report_ovid_or_imca_report', 'Latest third-party inspection/audit report'],
    ['fmea_report', 'FMEA report'],
    ['latest_dp_annual_trial_report', 'Latest DP Annual Trial Report'],
    ['vessel_crew_competency_matrix', 'Vessel crew competency matrix'],
    ['vessel_specification_sheet', 'Vessel Specification Sheet'],
    ['vessel_deck_and_profile_plan', 'Vessel deck and profile plan'],
]

const REGULATORY_DOCS = [
    ['vessel_safety_case', 'Vessel Safety Case'],
    ['safety_case_addendum', 'Safety Case Addendum or Bridging Document'],
    ['construction_installation_safety_case', 'Construction, Installation & Commissioning Safety Case'],
    ['production_operations_safety_case', 'Production/Operations Safety Case'],
    ['decommissioning_safety_case', 'Decommissioning Safety Case'],
    ['emergency_response_plan', 'Emergency Response Plan (ERP)'],
    ['erp_bridging_document', 'ERP Bridging Document'],
    ['diving_project_plan', 'Diving Project Plan & Diving Safety Management System'],
    ['health_safety_plan', 'Health & Safety Plan'],
    ['health_safety_plan_bridging', 'Health and Safety Plan Bridging Document'],
]

export default function VesselSummary() {
    const { id } = useParams()
    const navigate = useNavigate()
    const location = useLocation()
    const isReviewMode = location.pathname.startsWith('/vessel/review')

    const [app, setApp] = useState(null)
    const [docsByKey, setDocsByKey] = useState({})
    const [loading, setLoading] = useState(true)
    const [deciding, setDeciding] = useState(false)
    const [history, setHistory] = useState([])
    const [historyLoading, setHistoryLoading] = useState(true)
    const [officers, setOfficers] = useState([])
    const [staffUsers, setStaffUsers] = useState([])
    const [reviewRequests, setReviewRequests] = useState([])
    const [comments, setComments] = useState([])
    const [selectedStaffIds, setSelectedStaffIds] = useState([])
    const [requestingReview, setRequestingReview] = useState(false)
    const [newComment, setNewComment] = useState('')
    const [postingComment, setPostingComment] = useState(false)
    const [notifyingANP, setNotifyingANP] = useState(false)
    const [assigning, setAssigning] = useState(false)

    const [authLetter, setAuthLetter] = useState(null)
    const [extensionRequests, setExtensionRequests] = useState([])
    const [showExtensionDialog, setShowExtensionDialog] = useState(false)
    const [requestingExtension, setRequestingExtension] = useState(false)
    const [newExtensionDate, setNewExtensionDate] = useState('')
    const [extensionReason, setExtensionReason] = useState('')
    const [showIssueDialog, setShowIssueDialog] = useState(false)
    const [issuing, setIssuing] = useState(false)
    const [facilityName, setFacilityName] = useState('')
    const [decreeLawNumber, setDecreeLawNumber] = useState('')
    const [decreeLawClause, setDecreeLawClause] = useState('')
    const [inspectionDate, setInspectionDate] = useState('')
    const [signatureBlob, setSignatureBlob] = useState(null)
    const [downloadingAuth, setDownloadingAuth] = useState(false)
    const [downloadingOutcome, setDownloadingOutcome] = useState(false)


    const loadCollab = () => {
        apiListReviewRequests(id)
            .then((list) => setReviewRequests(list || []))
            .catch(() => setReviewRequests([]))
        apiListReviewComments(id)
            .then((list) => setComments(list || []))
            .catch(() => setComments([]))
    }

    const handleRequestReview = async () => {
        if (!selectedStaffIds.length === 0) {
            toast.push(
                <Notification type="danger" title="No staff selected">
                    Please select at least one staff member to request a review.
                </Notification>,
            )
            return
        }
        setRequestingReview(true)
        try {
            await apiRequestStaffReview(id, selectedStaffIds)
            toast.push(
                <Notification type="success" title="Submitted">
                    Review request submitted.
                </Notification>,
            )
            loadCollab()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to submit review request.'}
                </Notification>,
            )
        } finally {
            setRequestingReview(false)
        }
    }

    const handlePostComment = async () => {
        if (!newComment.trim()) return
        setPostingComment(true)
        try {
            await apiAddReviewComment(id, newComment.trim())
            setNewComment('')
            loadCollab()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to post comment.'}
                </Notification>,
            )
        } finally {
            setPostingComment(false)
        }
    }

    const handleNotifyANP = async () => {
        setNotifyingANP(true)
        try {
            await apiNotifyANPReady(id)
            toast.push(
                <Notification type="success" title="Sent">
                    ANP HSE has been notified.
                </Notification>,
            )
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to notify ANP HSE.'}
                </Notification>,
            )
        } finally {
            setNotifyingANP(false)
        }
    }


     const loadExtensionRequests = () => {
        apiListExtensionRequestsForApplication(id)
            .then((list) => setExtensionRequests(list || []))
            .catch(() => setExtensionRequests([]))
    }

    const handleRequestExtension = async () => {
        if (!newExtensionDate || !extensionReason) {
            toast.push(
                <Notification type="danger" title="Missing info">
                    Please fill in the new date and reason.
                </Notification>,
            )
            return
        }
        setRequestingExtension(true)
        try {
            await apiRequestExtension(id, newExtensionDate, extensionReason)
            toast.push(
                <Notification type="success" title="Submitted">
                    Extension request submitted.
                </Notification>,
            )
            setShowExtensionDialog(false)
            setNewExtensionDate('')
            setExtensionReason('')
            loadExtensionRequests()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to submit extension request.'}
                </Notification>,
            )
        } finally {
            setRequestingExtension(false)
        }
    }

    const handleViewOutcome = async () => {
        setDownloadingOutcome(true)
        try {
            const resp = await apiViewVesselOutcomeDocument(id)
            window.open(resp.url, '_blank', 'noopener,noreferrer')
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to open outcome document.'}
                </Notification>,
            )
        } finally {
            setDownloadingOutcome(false)
        }
    }

    const handleDownloadOutcome = async () => {
        setDownloadingOutcome(true)
        try {
            const resp = await apiDownloadVesselOutcomeDocument(id)
            window.open(resp.url, '_blank', 'noopener,noreferrer')
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to download outcome document.'}
                </Notification>,
            )
        } finally {
            setDownloadingOutcome(false)
        }
    }

    const loadAuthLetter = () => {
        apiGetAuthorisationLetter(id)
            .then((letter) => setAuthLetter(letter))
            .catch(() => setAuthLetter(null))
    }

    const handleIssueLetter = async () => {
        if (!facilityName || !decreeLawNumber || !decreeLawClause || !signatureBlob) {
            toast.push(
                <Notification type="danger" title="Missing info">
                    Please fill in all fields and provide a signature.
                </Notification>,
            )
            return
        }
        setIssuing(true)
        try {
            const formData = new FormData()
            formData.append('facilityName', facilityName)
            formData.append('decreeLawNumber', decreeLawNumber)
            formData.append('decreeLawClause', decreeLawClause)
            formData.append('inspectionDate', inspectionDate)
            formData.append('signature', signatureBlob, 'signature.png')

            await apiIssueAuthorisationLetter(id, formData)
            toast.push(
                <Notification type="success" title="Issued">
                    Authorisation letter issued successfully.
                </Notification>,
            )
            setShowIssueDialog(false)
            loadAuthLetter()
            loadExtensionRequests()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to issue authorisation letter.'}
                </Notification>,
            )
        } finally {
            setIssuing(false)
        }
    }

    const handleDownloadAuth = async (format) => {
        setDownloadingAuth(true)
        try {
            const resp = await apiDownloadAuthorisationLetter(id, format)
            window.open(resp.url, '_blank', 'noopener,noreferrer')
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to open authorisation letter.'}
                </Notification>,
            )
        } finally {
            setDownloadingAuth(false)
        }
    }

        const load = () => {
        setLoading(true)
        const getApp = isReviewMode ? apiGetVesselApplicationForReview(id) : apiGetVesselApplication(id)
        const getDocs = isReviewMode ? apiListVesselDocumentsForReview(id) : apiListVesselDocuments(id)

        Promise.all([getApp, getDocs])
            .then(([appData, docs]) => {
                setApp(appData)
                const byKey = {}
                ;(docs || []).forEach((d) => {
                    byKey[`${d.category}:${d.documentKey}`] = d
                })
                setDocsByKey(byKey)
            })
            .finally(() => setLoading(false))

        setHistoryLoading(true)
        apiGetVesselApplicationHistory(id)
            .then((entries) => setHistory(entries || []))
            .catch(() => setHistory([]))
            .finally(() => setHistoryLoading(false))

        if (isReviewMode) {
            apiListHSEOfficers()
                .then((list) => setOfficers(list || []))
                .catch(() => setOfficers([]))
            apiListStaffUsers()
                .then((list) => setStaff(list || []))
                .catch(() => setStaff([]))
        }
        loadAuthLetter()
        loadExtensionRequests()
    }

    const handleAssign = async (option) => {
        setAssigning(true)
        try {
            await apiAssignHSEOfficer(id, option ? option.value : null)
            toast.push(
                <Notification type="success" title="Saved">
                    HSE Officer assignment updated.
                </Notification>,
            )
            load()
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to update assignment.'}
                </Notification>,
            )
        } finally {
            setAssigning(false)
        }
    }

    useEffect(() => {
        load()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id])

    const handleDecision = async (status) => {
        setDeciding(true)
        try {
            await apiDecideVesselApplication(id, status)
            toast.push(
                <Notification type="success" title={status}>
                    Application has been {status.toLowerCase()}.
                </Notification>,
            )
            navigate('/vessel/review')
        } catch (err) {
            toast.push(
                <Notification type="danger" title="Error">
                    {err?.response?.data?.message || 'Failed to save decision.'}
                </Notification>,
            )
        } finally {
            setDeciding(false)
        }
    }

    if (loading) {
        return (
            <div className="flex justify-center py-10">
                <Spinner size={32} />
            </div>
        )
    }
    if (!app) return <p className="text-gray-500">Application not found.</p>

    return (
        <>

                      {(app.status === 'Approved' || app.status === 'Submitted') && (
                <Card className="mb-4 print:hidden">
                    <div className="flex flex-wrap items-center justify-between gap-3">
                        <div>
                            <h5>Outcome Document</h5>
                            <p className="text-sm text-gray-500">
                                Full summary of the application, generated
                                automatically when this decision was made.
                            </p>
                        </div>
                        <div className="flex gap-2">
                            <Button
                                size="sm"
                                loading={downloadingOutcome}
                                onClick={handleViewOutcome}
                            >
                                View PDF
                            </Button>
                            <Button
                                size="sm"
                                loading={downloadingOutcome}
                                onClick={handleDownloadOutcome}
                            >
                                Download Word
                            </Button>
                        </div>
                    </div>
                </Card>
            )}

            {app.status === 'Approved' && (
                <Card className="mb-4 print:hidden">
                    <div className="flex flex-wrap items-center justify-between gap-3">
                        <div>
                            <h5>Vessel Entry Authorisation Letter</h5>
                            {authLetter ? (
                                <p className="text-sm text-gray-500">
                                    {authLetter.referenceNo} - issued by{' '}
                                    {authLetter.issuedByName} on{' '}
                                    {new Date(authLetter.issuedAt).toLocaleDateString()}
                                </p>
                            ) : (
                                <p className="text-sm text-gray-500">
                                    No official authorisation letter has been
                                    issued yet.
                                </p>
                            )}
                        </div>
                        <div className="flex gap-2">
                            {authLetter && (
                                <>
                                    <Button
                                        size="sm"
                                        loading={downloadingAuth}
                                        onClick={() => handleDownloadAuth('pdf')}
                                    >
                                        View PDF
                                    </Button>
                                    <Button
                                        size="sm"
                                        loading={downloadingAuth}
                                        onClick={() => handleDownloadAuth('docx')}
                                    >
                                        Download Word
                                    </Button>
                                </>
                            )}
                            {isReviewMode && (
                                <Button
                                    size="sm"
                                    variant="solid"
                                    onClick={() => setShowIssueDialog(true)}
                                >
                                    {authLetter ? 'Reissue' : 'Issue Letter'}
                                </Button>
                            )}
                        </div>
                    </div>
                </Card>
            )}

                        {app.status === 'Approved' && !isReviewMode && (
                <Card className="mb-4 print:hidden">
                    <div className="flex flex-wrap items-center justify-between gap-3">
                        <div>
                            <h5>Extension of Entry Authorisation</h5>
                            <p className="text-sm text-gray-500">
                                Current end date: {app.entryDateTo || '-'}
                            </p>
                        </div>
                        <Button
                            size="sm"
                            variant="solid"
                            onClick={() => setShowExtensionDialog(true)}
                        >
                            Request Extension
                        </Button>
                    </div>

                    {extensionRequests.length > 0 && (
                        <div className="mt-4 border-t border-gray-100 pt-3">
                            {extensionRequests.map((r) => (
                                <div key={r.id} className="mb-2 text-sm">
                                    <Tag
                                        className={
                                            r.status === 'Approved'
                                                ? 'bg-emerald-100 text-emerald-700'
                                                : r.status === 'Rejected'
                                                  ? 'bg-red-100 text-red-700'
                                                  : 'bg-amber-100 text-amber-700'
                                        }
                                    >
                                        {r.status}
                                    </Tag>{' '}
                                    Requested {r.requestedEntryDateTo} -{' '}
                                    {r.reason}
                                    {r.status === 'Rejected' && r.rejectionReason && (
                                        <div className="text-gray-500">
                                            Reason: {r.rejectionReason}
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </Card>
            )}

            <Dialog
                isOpen={showExtensionDialog}
                onClose={() => setShowExtensionDialog(false)}
                onRequestClose={() => setShowExtensionDialog(false)}
            >
                <h4 className="mb-4">Request Extension</h4>
                <div className="flex flex-col gap-3">
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            New end date
                        </label>
                        <Input
                            type="date"
                            value={newExtensionDate}
                            onChange={(e) => setNewExtensionDate(e.target.value)}
                        />
                    </div>
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            Reason
                        </label>
                        <Input
                            textArea
                            value={extensionReason}
                            onChange={(e) => setExtensionReason(e.target.value)}
                        />
                    </div>
                </div>
                <div className="mt-6 flex justify-end gap-2">
                    <Button onClick={() => setShowExtensionDialog(false)}>
                        Cancel
                    </Button>
                    <Button
                        variant="solid"
                        loading={requestingExtension}
                        onClick={handleRequestExtension}
                    >
                        Submit Request
                    </Button>
                </div>
            </Dialog>

            <Dialog
                isOpen={showIssueDialog}
                onClose={() => setShowIssueDialog(false)}
                onRequestClose={() => setShowIssueDialog(false)}
            >
                <h4 className="mb-4">Issue Authorisation Letter</h4>
                <div className="flex max-h-[65vh] flex-col gap-3 overflow-y-auto pr-1">
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            Facility name
                        </label>
                        <Input
                            value={facilityName}
                            onChange={(e) => setFacilityName(e.target.value)}
                        />
                    </div>
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            Decree-Law number
                        </label>
                        <Input
                            value={decreeLawNumber}
                            onChange={(e) => setDecreeLawNumber(e.target.value)}
                        />
                    </div>
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            Decree-Law clause/article number
                        </label>
                        <Input
                            value={decreeLawClause}
                            onChange={(e) => setDecreeLawClause(e.target.value)}
                        />
                    </div>
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            ANP HSE inspection date (optional)
                        </label>
                        <Input
                            type="date"
                            value={inspectionDate}
                            onChange={(e) => setInspectionDate(e.target.value)}
                        />
                    </div>
                    <div>
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            Signature
                        </label>
                        <SignaturePad onChange={setSignatureBlob} />
                    </div>
                </div>
                <div className="mt-6 flex justify-end gap-2">
                    <Button onClick={() => setShowIssueDialog(false)}>
                        Cancel
                    </Button>
                    <Button
                        variant="solid"
                        loading={issuing}
                        onClick={handleIssueLetter}
                    >
                        Issue Letter
                    </Button>
                </div>
            </Dialog>

                        {isReviewMode && app.assignedHseOfficerId && (
                <Card className="mb-4 print:hidden">
                    <h5 className="mb-3">Internal Review</h5>

                    <div className="mb-4">
                        <label className="mb-1 block text-sm font-semibold text-gray-600">
                            Request review from Staff
                        </label>
                        <div className="flex flex-wrap items-center gap-2">
                            <div className="min-w-[260px] flex-1">
                                <Select
                                    isMulti
                                    placeholder="Select staff members..."
                                    options={staffUsers.map((s) => ({ value: s.id, label: s.name }))}
                                    value={staffUsers
                                        .map((s) => ({ value: s.id, label: s.name }))
                                        .filter((o) => selectedStaffIds.includes(o.value))}
                                    onChange={(options) =>
                                        setSelectedStaffIds((options || []).map((o) => o.value))
                                    }
                                />
                            </div>
                            <Button
                                size="sm"
                                variant="solid"
                                loading={requestingReview}
                                onClick={handleRequestReview}
                            >
                                Send Request
                            </Button>
                        </div>
                        {reviewRequests.length > 0 && (
                            <p className="mt-2 text-xs text-gray-500">
                                Already requested:{' '}
                                {reviewRequests.map((r) => r.staffName).join(', ')}
                            </p>
                        )}
                    </div>

                    <div className="mb-4 border-t border-gray-100 pt-4">
                        <label className="mb-2 block text-sm font-semibold text-gray-600">
                            Comments
                        </label>
                        <div className="mb-3 flex max-h-[300px] flex-col gap-3 overflow-y-auto">
                            {comments.length === 0 && (
                                <p className="text-sm text-gray-400">No comments yet.</p>
                            )}
                            {comments.map((c) => (
                                <div key={c.id} className="rounded-lg bg-gray-50 p-3">
                                    <div className="mb-1 flex items-center gap-2 text-xs text-gray-500">
                                        <span className="font-semibold text-gray-700">
                                            {c.authorName}
                                        </span>
                                        <Tag className="bg-gray-200 text-gray-600">
                                            {c.authorRole}
                                        </Tag>
                                        <span>
                                            {new Date(c.createdAt).toLocaleString()}
                                        </span>
                                    </div>
                                    <div className="text-sm text-gray-800">{c.comment}</div>
                                </div>
                            ))}
                        </div>
                        <div className="flex gap-2">
                            <Input
                                textArea
                                placeholder="Add a comment..."
                                value={newComment}
                                onChange={(e) => setNewComment(e.target.value)}
                            />
                            <Button
                                size="sm"
                                variant="solid"
                                loading={postingComment}
                                onClick={handlePostComment}
                            >
                                Post
                            </Button>
                        </div>
                    </div>

                    <div className="border-t border-gray-100 pt-4">
                        <Button
                            size="sm"
                            variant="solid"
                            loading={notifyingANP}
                            onClick={handleNotifyANP}
                        >
                            Notify ANP HSE — Ready for Decision
                        </Button>
                    </div>
                </Card>
            )}

                    {isReviewMode && (
                <Card className="mb-4 print:hidden">
                    <div className="flex flex-wrap items-center gap-3">
                        <span className="text-sm font-semibold text-gray-600">
                            Assigned HSE Officer:
                        </span>
                        <div className="min-w-[220px]">
                            <Select
                                isClearable
                                isDisabled={assigning}
                                placeholder="Unassigned"
                                options={officers.map((o) => ({ value: o.id, label: o.name }))}
                                value={
                                    app.assignedHseOfficerId
                                        ? officers
                                              .map((o) => ({ value: o.id, label: o.name }))
                                              .find((o) => o.value === app.assignedHseOfficerId) || null
                                        : null
                                }
                                onChange={handleAssign}
                            />
                        </div>
                        <span className="text-xs text-gray-400">
                            Notified in parallel at each stage - coordinates
                            the case but doesn&apos;t approve/reject.
                        </span>
                    </div>
                </Card>
            )}
            <div className="mb-4 flex items-center justify-between print:hidden">
                <Button onClick={() => navigate(isReviewMode ? '/vessel/review' : '/vessel/list')}>
                    ← Back to list
                </Button>
                <div className="flex gap-2">
                    <Button variant="solid" onClick={() => window.print()}>
                        Print / Export PDF
                    </Button>
                    {isReviewMode && app.status === 'Submitted' && app.assignedHseOfficerId && (
                        <>
                            <Button
                                variant="solid"
                                color="emerald-600"
                                loading={deciding}
                                onClick={() => handleDecision('Approved')}
                            >
                                Approve
                            </Button>
                            <Button
                                variant="solid"
                                color="red-600"
                                loading={deciding}
                                onClick={() => handleDecision('Rejected')}
                            >
                                Reject
                            </Button>
                        </>
                    )}
                </div>
            </div>

            <Card>
                <div id="vessel-summary-print">
                    <div className="mb-6 flex items-center justify-between border-b border-gray-200 pb-4">
                        <div>
                            <div className="font-mono text-xs text-gray-400">{app.applicationNumber}</div>
                            <h3>Vessel Entry Application - Outcome</h3>
                        </div>
                        <Tag className={STATUS_TAG[app.status]}>{app.status}</Tag>
                    </div>

                    <h5 className="mb-2 mt-4">Entry to Contract Area</h5>
                    <Field label="Condition" value={app.entryCondition} />
                    <Field label="Purpose of vessel entry request" value={app.purpose} />
                    <Field label="Estimated arrival" value={app.entryDateFrom} />
                    <Field label="Estimated departure" value={app.entryDateTo} />
                    <Field label="Entry type" value={app.entryType} />
                    <Field label="Scope of work" value={(app.scopeOfWork || []).join(', ')} />
                    <Field label="Vessel operation mode" value={app.operationMode} />
                    <DocField label="Cover letter" doc={docsByKey['cover_letter:cover_letter']} />

                    <h5 className="mb-2 mt-6">Vessel Description</h5>
                    <Field label="Name of Vessel" value={app.vesselName} />
                    <Field label="IMO Number" value={app.vesselImoNumber} />
                    <Field label="Vessel Owner/Contractor" value={app.vesselOwner} />
                    <Field label="Vessel Type" value={app.vesselType} />
                    <Field label="Flag State" value={app.flagState} />
                    <Field label="Classification Society" value={app.classificationSociety} />
                    <Field label="Port of Registry" value={app.portOfRegistry} />
                    <Field label="Class ID Number" value={app.classIdNumber} />
                    <Field label="Length Overall (LOA)" value={app.lengthOverall} />
                    <Field label="Draft" value={app.draftValue} />
                    <Field label="Gross Tonnage" value={app.grossTonnage} />
                    <Field label="Call Sign" value={app.callSign} />

                    <h5 className="mb-2 mt-6">Main Regulatory Documents</h5>
                    {REGULATORY_DOCS.map(([key, label]) => (
                        <DocField key={key} label={label} doc={docsByKey[`regulatory_document:${key}`]} />
                    ))}

                    <h5 className="mb-2 mt-6">Vessel Statutory Certificates</h5>
                    {STATUTORY_CERTS.map(([key, label]) => (
                        <DocField key={key} label={label} doc={docsByKey[`statutory_certificate:${key}`]} />
                    ))}

                    <h5 className="mb-2 mt-6">Other Necessary Vessel Documentation</h5>
                    {SUPPORTING_DOCS.map(([key, label]) => (
                        <DocField key={key} label={label} doc={docsByKey[`supporting_document:${key}`]} />
                    ))}

                    <h5 className="mb-2 mt-6">The Authorised Person Declares</h5>
                    <Field label="No major deficiencies" value={app.noMajorDeficiencies} />
                    <Field label="No detention within last 12 months" value={app.noDetention12Months} />
                    <Field label="Safety equipment operational" value={app.safetyEquipmentOperational} />
                    <Field label="Firefighting systems operational" value={app.firefightingOperational} />
                    <Field label="Lifesaving appliances operational" value={app.lifesavingOperational} />

                    <h5 className="mb-2 mt-6">Number of Crew</h5>
                    <Field label="Vessel crew" value={app.crewCount ? `${app.crewCount} persons` : null} />
                    <Field label="Survey/Project crew" value={app.surveyCrewCount ? `${app.surveyCrewCount} persons` : null} />
                    <DocField label="Crew list document" doc={docsByKey['crew_document:crew_list']} />

                    <h5 className="mb-2 mt-6">Cargo Information</h5>
                    <Field label="Cargo onboard" value={app.cargoOnboard} />
                    <Field label="Cargo Description" value={app.cargoDescription} />
                    <Field label="Hazardous Cargo" value={app.hazardousCargo} />
                    <DocField label="Dangerous Good Declaration" doc={docsByKey['cargo_document:dangerous_good_declaration']} />
                    <DocField label="Cargo Manifest" doc={docsByKey['cargo_document:cargo_manifest']} />

                    <h5 className="mb-2 mt-6">Declaration of Environmental Compliance</h5>
                    <Field label="Waste Discharge" value={app.wasteDischarge} />
                    <Field label="Oily Waste Onboard" value={app.oilyWasteOnboard} />
                    <Field label="Sewage Disposal Required" value={app.sewageDisposalRequired} />

                    <h5 className="mb-2 mt-6">Clearances from Timor-Leste Relevant Government Entities</h5>
                    <DocField label="Immigration" doc={docsByKey['clearance_document:immigration']} />
                    <DocField label="Customs" doc={docsByKey['clearance_document:customs']} />
                    <DocField label="Quarantine" doc={docsByKey['clearance_document:quarantine']} />
                </div>
            </Card>

            {app.status !== 'Draft' && (
                <Card className="mt-4 print:hidden">
                    <h4 className="mb-4">History</h4>
                    <HistoryTimeline entries={history} loading={historyLoading} />
                </Card>
            )}
        </>
    )
}