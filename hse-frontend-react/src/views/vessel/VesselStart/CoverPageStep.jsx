import { useState } from 'react'
import Checkbox from '@/components/ui/Checkbox'
import Button from '@/components/ui/Button'
import Alert from '@/components/ui/Alert'

const REQUIREMENTS = [
    'You are required to complete the Online Entry Application Form as accurately as possible.',
    'All submission is to be made within 30 days before the expected date of entry to Contract Area. This timeline is only applicable for complete submission.',
    'Urgent submission can only be made with justification.',
    'All submissions made to the ANP for the assessment must be accompanied by a cover letter from the Contract Operator and should have a clear subject, reflecting the content of the submission i.e., vessel "AAAA" entry application.',
    'Drawings in all submissions must be in full size to ensure scalability; for instance, provide A1 drawings at A1 size, not reduced down to A4.',
    'All submissions including supporting documents shall be provided in either paper or electronic form, the later preferably in Adobe Portable Document Format (PDF) i.e. text searchable and are not protected.',
]

export default function CoverPageStep({ onNext }) {
    const [understood, setUnderstood] = useState(false)
    const [error, setError] = useState('')

    const handleNext = () => {
        if (!understood) {
            setError('You must confirm that you understand the notice above to continue.')
            return
        }
        onNext()
    }

    return (
        <div>
            <h5 className="mb-3">
                Entry Application to Contract Area or Area Authorised for
                Petroleum Activities
            </h5>
            <div className="mb-4 rounded-lg border-2 border-gray-800 p-4">
                <p className="mb-4 text-sm text-gray-700">
                    Pursuant to clause ----- of the DL---/2016 on
                    ------------, the authorised person shall seek approval
                    from the Ministry for all entry of personnel, vessels
                    and aircrafts for the purposes of entering into contract
                    area for petroleum operations.
                </p>
                <p className="mb-4 text-sm font-semibold text-gray-700">
                    Thank you for using online entry application
                </p>
                <ul className="mb-4 list-disc space-y-2 pl-5 text-sm text-gray-600">
                    {REQUIREMENTS.map((req) => (
                        <li key={req}>{req}</li>
                    ))}
                </ul>
                <Checkbox checked={understood} onChange={(checked) => setUnderstood(checked)}>
                    I understand the above notice and shall ensure
                    accordingly
                </Checkbox>
            </div>

            {error && (
                <Alert showIcon className="mb-4" type="danger">
                    {error}
                </Alert>
            )}

            <div className="flex justify-end">
                <Button variant="solid" onClick={handleNext}>
                    I Understand, Continue
                </Button>
            </div>
        </div>
    )
}