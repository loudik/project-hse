const STEPS = [
    'Emails',
    'Contract Info',
    'Entry Type',
    'Regulatory Docs',
    'Entry Condition',
    'Vessel Description',
    'Statutory Certificates',
    'Declarations',
    'Review & Submit',
]

export default function StepIndicator({ current }) {
    return (
        <div className="mb-8 flex items-center overflow-x-auto pb-2">
            {STEPS.map((label, idx) => {
                const stepNum = idx + 1
                const isActive = stepNum === current
                const isDone = stepNum < current
                return (
                    <div key={label} className="flex min-w-[95px] flex-1 items-center">
                        <div className="flex flex-col items-center">
                            <div
                                className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-semibold ${
                                    isDone
                                        ? 'bg-emerald-500 text-white'
                                        : isActive
                                          ? 'bg-primary text-white'
                                          : 'bg-gray-200 text-gray-500'
                                }`}
                            >
                                {isDone ? '✓' : stepNum}
                            </div>
                            <span className="mt-1 text-center text-[10px] text-gray-500">{label}</span>
                        </div>
                        {idx < STEPS.length - 1 && (
                            <div className={`mx-1.5 h-0.5 flex-1 ${isDone ? 'bg-emerald-500' : 'bg-gray-200'}`} />
                        )}
                    </div>
                )
            })}
        </div>
    )
}