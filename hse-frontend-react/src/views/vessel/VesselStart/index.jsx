import { useState } from 'react'
import Card from '@/components/ui/Card'
import CoverPageStep from './CoverPageStep'
import StatementStep from './StatementStep'
import ChooseActionStep from './ChooseActionStep'

export default function VesselStart() {
    const [step, setStep] = useState('cover') // 'cover' | 'statement' | 'action'
    return (
        <Card>
            <h4 className="mb-6">Vessel Entry Application</h4>
            {step === 'cover' && <CoverPageStep onNext={() => setStep('statement')} />}
            {step === 'statement' && <StatementStep onNext={() => setStep('action')} />}
            {step === 'action' && <ChooseActionStep />}
        </Card>
    )
}