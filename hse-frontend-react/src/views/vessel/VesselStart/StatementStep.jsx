import { useState } from 'react'
import Checkbox from '@/components/ui/Checkbox'
import Button from '@/components/ui/Button'
import Alert from '@/components/ui/Alert'

const DUTIES = [
    'Entry into the Contract Area or area authorised for petroleum operations shall only be upon authorisation of ANP',
    'Contract Operator, Authorised Person, or Applicant is obliged to secure entry authorisation from ANP for its vessel, personnel, helicopter, goods, equipment and materials',
    'Contract Operator, Authorised Person, or Applicant is obliged to ensure that all information submitted to ANP is current',
    'Contract Operator, Authorised Person, or Applicant is obliged to ensure that movements in and out of the Contract Area for petroleum operations is controlled and changes in circumstance shall be reported to ANP and updated as necessary in this online system',
    'Contract Operator, Authorised Person, or Applicant can apply for blanket authorisation valid up to 1 year to avoid repeated application',
]

export default function StatementStep({ onNext }) {
    const [understood, setUnderstood] = useState(false)
    const [error, setError] = useState('')

    const handleNext = () => {
        if (!understood) {
            setError('You must confirm that you understand the statement above to continue.')
            return
        }
        onNext()
    }

    return (
        <div>
            <div className="mb-4 rounded-lg border-2 border-gray-800 p-4">
                <h5 className="mb-3">Duties of the Contract Operator, Authorised Person or Applicant</h5>
                <ul className="mb-4 list-disc space-y-2 pl-5 text-sm text-gray-600">
                    {DUTIES.map((duty) => (
                        <li key={duty}>{duty}</li>
                    ))}
                </ul>

                <div className="flex flex-col gap-2">
                    <Checkbox checked={understood} onChange={(checked) => setUnderstood(checked)}>
                        I understand the above statement and shall ensure accordingly
                    </Checkbox>
                    <Checkbox checked={!understood} onChange={(checked) => setUnderstood(!checked)}>
                        I do not understand the above statements
                    </Checkbox>
                </div>
            </div>

            {error && (
                <Alert showIcon className="mb-4" type="danger">
                    {error}
                </Alert>
            )}

            <div className="flex justify-end">
                <Button variant="solid" onClick={handleNext}>
                    Save and go to next page
                </Button>
            </div>
        </div>
    )
}