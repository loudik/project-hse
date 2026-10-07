function Row({ label, value }) {
    return (
        <div className="mb-3">
            <div className="text-xs font-semibold tracking-wide text-gray-400 uppercase">{label}</div>
            <div className="text-gray-700">
                {value === true ? 'Yes' : value === false ? 'No' : value || '-'}
            </div>
        </div>
    )
}

export default function StepReview({ values }) {
    return (
        <div>
            <p className="mb-4 text-gray-500">
                Please review the information below before submitting. Once submitted, this
                application becomes read-only unless the reviewer requests changes.
            </p>

            <h5 className="mb-2">Entry Condition</h5>
            <div className="mb-4 grid grid-cols-1 gap-3 rounded-lg border border-gray-200 p-4 sm:grid-cols-2">
                <Row label="Condition" value={values.entryCondition} />
                <Row label="Entry type" value={values.entryType} />
                <Row label="Arrival date" value={values.entryDateFrom} />
                <Row label="Departure date" value={values.entryDateTo} />
                <Row label="Purpose" value={values.purpose} />
            </div>

            <h5 className="mb-2">Vessel Description</h5>
            <div className="mb-4 grid grid-cols-1 gap-3 rounded-lg border border-gray-200 p-4 sm:grid-cols-2">
                <Row label="Vessel name" value={values.vesselName} />
                <Row label="IMO number" value={values.vesselImoNumber} />
                <Row label="Vessel type" value={values.vesselType} />
                <Row label="Owner/Contractor" value={values.vesselOwner} />
                <Row label="Flag state" value={values.flagState} />
                <Row label="Operation mode" value={values.operationMode} />
                <Row label="Scope of work" value={(values.scopeOfWork || []).join(', ')} />
            </div>

            <h5 className="mb-2">Declarations</h5>
            <div className="grid grid-cols-1 gap-3 rounded-lg border border-gray-200 p-4 sm:grid-cols-2">
                <Row label="No major deficiencies" value={values.noMajorDeficiencies} />
                <Row label="No detention (12 months)" value={values.noDetention12Months} />
                <Row label="Crew count" value={values.crewCount} />
                <Row label="Cargo onboard" value={values.cargoOnboard} />
                <Row label="Hazardous cargo" value={values.hazardousCargo} />
            </div>
        </div>
    )
}