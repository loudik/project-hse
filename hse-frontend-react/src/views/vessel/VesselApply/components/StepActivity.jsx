import Input from '@/components/ui/Input'
import { FormItem } from '@/components/ui/Form'

export default function StepActivity({ values, onChange }) {
    return (
        <div>
            <FormItem label="Proposed activity">
                <Input
                    textArea
                    rows={4}
                    value={values.proposedActivity}
                    onChange={(e) => onChange({ proposedActivity: e.target.value })}
                    placeholder="Describe the activity this vessel will perform..."
                />
            </FormItem>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <FormItem label="Planned arrival date">
                    <Input
                        type="date"
                        value={values.plannedArrivalDate}
                        onChange={(e) => onChange({ plannedArrivalDate: e.target.value })}
                    />
                </FormItem>
                <FormItem label="Planned departure date">
                    <Input
                        type="date"
                        value={values.plannedDepartureDate}
                        onChange={(e) => onChange({ plannedDepartureDate: e.target.value })}
                    />
                </FormItem>
            </div>
        </div>
    )
}