import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'

export default function StepEmailAddresses({ values, onChange }) {
    const emails = values.notificationEmails?.length ? values.notificationEmails : ['']

    const updateEmail = (index, value) => {
        const next = [...emails]
        next[index] = value
        onChange({ notificationEmails: next })
    }

    const addEmail = () => {
        if (emails.length >= 3) return
        onChange({ notificationEmails: [...emails, ''] })
    }

    const removeEmail = (index) => {
        const next = emails.filter((_, i) => i !== index)
        onChange({ notificationEmails: next.length ? next : [''] })
    }

    return (
        <div>
            <p className="mb-1 text-gray-600">
                Please enter valid email addresses for the purpose of receiving notification and
                establishing future communication with ANP regarding this application.
            </p>
            <p className="mb-5 text-gray-600">
                Information and queries will be sent to the following email addresses.
            </p>

            {emails.map((email, index) => (
                <div key={index} className="mb-3 flex items-center gap-2">
                    <div className="flex-1">
                        <label className="mb-1.5 block text-sm font-semibold">
                            Email address {index + 1}
                        </label>
                        <Input
                            type="email"
                            value={email}
                            onChange={(e) => updateEmail(index, e.target.value)}
                            placeholder="name@example.com"
                        />
                    </div>
                    {emails.length > 1 && (
                        <Button size="sm" className="mt-6" onClick={() => removeEmail(index)}>
                            Remove
                        </Button>
                    )}
                </div>
            ))}

            {emails.length < 3 && (
                <Button onClick={addEmail}>+ Add email address</Button>
            )}
        </div>
    )
}