import { useRef, useState } from 'react'
import Checkbox from '@/components/ui/Checkbox'
import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'
import { apiUploadVesselDocument } from '@/services/VesselService'

/**
 * One row for a single document/certificate requirement.
 * Handles: Not Applicable + reason, Date Issued/Expired, file upload.
 * Saves itself immediately on blur/file-select (no separate "save" step).
 */
export default function DocumentUploadRow({
    applicationId,
    category,
    documentKey,
    label,
    showDates = false,
    existingDoc,
}) {
    const fileInputRef = useRef(null)
    const [notApplicable, setNotApplicable] = useState(existingDoc?.notApplicable || false)
    const [naReason, setNaReason] = useState(existingDoc?.naReason || '')
    const [dateIssued, setDateIssued] = useState(existingDoc?.dateIssued || '')
    const [dateExpired, setDateExpired] = useState(existingDoc?.dateExpired || '')
    const [fileName, setFileName] = useState(existingDoc?.fileName || '')
    const [saving, setSaving] = useState(false)
    const [saved, setSaved] = useState(false)

    const save = async (file) => {
        setSaving(true)
        setSaved(false)
        try {
            const formData = new FormData()
            formData.append('category', category)
            formData.append('documentKey', documentKey)
            formData.append('label', label)
            formData.append('notApplicable', String(notApplicable))
            formData.append('naReason', naReason)
            formData.append('dateIssued', dateIssued)
            formData.append('dateExpired', dateExpired)
            if (file) formData.append('file', file)

            await apiUploadVesselDocument(applicationId, formData)
            setSaved(true)
        } finally {
            setSaving(false)
        }
    }

    const handleFileChange = (e) => {
        const file = e.target.files?.[0]
        if (!file) return
        setFileName(file.name)
        save(file)
    }

    const handleNotApplicableChange = (checked) => {
        setNotApplicable(checked)
        if (checked) save(null)
    }

    return (
        <div className="grid grid-cols-1 items-center gap-2 border-b border-gray-100 py-2.5 sm:grid-cols-12">
            <div className="sm:col-span-4 text-sm">{label}</div>

            <div className="sm:col-span-2">
                <Checkbox checked={notApplicable} onChange={handleNotApplicableChange}>
                    N/A
                </Checkbox>
            </div>

            {notApplicable ? (
                <div className="sm:col-span-4">
                    <Input
                        size="sm"
                        placeholder="Reason for N/A"
                        value={naReason}
                        onChange={(e) => setNaReason(e.target.value)}
                        onBlur={() => save(null)}
                    />
                </div>
            ) : (
                <>
                    {showDates && (
                        <>
                            <div className="sm:col-span-2">
                                <Input
                                    size="sm"
                                    type="date"
                                    value={dateIssued}
                                    onChange={(e) => setDateIssued(e.target.value)}
                                    onBlur={() => fileName && save(null)}
                                    placeholder="Issued"
                                />
                            </div>
                            <div className="sm:col-span-2">
                                <Input
                                    size="sm"
                                    type="date"
                                    value={dateExpired}
                                    onChange={(e) => setDateExpired(e.target.value)}
                                    onBlur={() => fileName && save(null)}
                                    placeholder="Expired"
                                />
                            </div>
                        </>
                    )}
                    <div className={showDates ? 'sm:col-span-2' : 'sm:col-span-6'}>
                        <input
                            ref={fileInputRef}
                            type="file"
                            className="hidden"
                            onChange={handleFileChange}
                        />
                        <Button
                            size="xs"
                            loading={saving}
                            onClick={() => fileInputRef.current?.click()}
                        >
                            {fileName ? '✓ ' + fileName.slice(0, 18) : 'Upload'}
                        </Button>
                    </div>
                </>
            )}

            {saved && <span className="text-xs text-emerald-600 sm:col-span-12">Saved</span>}
        </div>
    )
}