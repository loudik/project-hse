import { useState } from 'react'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import Button from '@/components/ui/Button'
import Alert from '@/components/ui/Alert'
import { countryList } from '@/constants/countries.constant'
import { FormItem, Form } from '@/components/ui/Form'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { apiCreateOrganization } from '@/services/OrganizationService'

const ORG_TYPE_OPTIONS = [
    { value: 'Operator', label: 'Operator' },
    { value: 'Contractor', label: 'Contractor' },
    { value: 'Subcontractor', label: 'Subcontractor' },
    { value: 'Supplier', label: 'Supplier' },
    { value: 'Other', label: 'Other' },
]

const validationSchema = z.object({
    name: z.string().min(1, { message: 'Please enter the organization name' }),
    registrationNumber: z.string().min(1, { message: 'Please enter the registration number' }),
    type: z.string().min(1, { message: 'Please select an organization type' }),
    address: z.string().min(1, { message: 'Please enter the address' }),
    country: z.string().min(1, { message: 'Please enter the country' }),
    phoneNumber: z.string().min(1, { message: 'Please enter a phone number' }),
    email: z.email({ message: 'Please enter a valid email' }),
    website: z.string().optional(),
})

const OrganizationForm = ({ onCreated }) => {
    const [isSubmitting, setSubmitting] = useState(false)
    const [errorMessage, setErrorMessage] = useState('')

    const {
        handleSubmit,
        formState: { errors },
        control,
    } = useForm({
        resolver: zodResolver(validationSchema),
    })

    const onSubmit = async (values) => {
        setSubmitting(true)
        setErrorMessage('')
        try {
            const org = await apiCreateOrganization(values)
            onCreated?.(org)
        } catch (err) {
            setErrorMessage(err?.response?.data?.message || 'Failed to save organization profile.')
        } finally {
            setSubmitting(false)
        }
    }

    return (
        <Form onSubmit={handleSubmit(onSubmit)}>
            {errorMessage && (
                <Alert showIcon className="mb-4" type="danger">
                    <span className="break-all">{errorMessage}</span>
                </Alert>
            )}

            <FormItem
                label="Organization name"
                invalid={Boolean(errors.name)}
                errorMessage={errors.name?.message}
            >
                <Controller
                    name="name"
                    control={control}
                    render={({ field }) => (
                        <Input type="text" placeholder="Organization name" {...field} />
                    )}
                />
            </FormItem>

            <FormItem
                label="Registration number"
                invalid={Boolean(errors.registrationNumber)}
                errorMessage={errors.registrationNumber?.message}
            >
                <Controller
                    name="registrationNumber"
                    control={control}
                    render={({ field }) => (
                        <Input type="text" placeholder="Registration number" {...field} />
                    )}
                />
            </FormItem>

            <FormItem
                label="Type"
                invalid={Boolean(errors.type)}
                errorMessage={errors.type?.message}
            >
                <Controller
                    name="type"
                    control={control}
                    render={({ field }) => (
                        <Select
                            options={ORG_TYPE_OPTIONS}
                            placeholder="Select organization type"
                            value={ORG_TYPE_OPTIONS.find((o) => o.value === field.value) || null}
                            onChange={(option) => field.onChange(option?.value || '')}
                        />
                    )}
                />
            </FormItem>

            <FormItem
                label="Address"
                invalid={Boolean(errors.address)}
                errorMessage={errors.address?.message}
            >
                <Controller
                    name="address"
                    control={control}
                    render={({ field }) => (
                        <Input textArea rows={2} placeholder="Full address" {...field} />
                    )}
                />
            </FormItem>

            <FormItem
                label="Country"
                invalid={Boolean(errors.country)}
                errorMessage={errors.country?.message}
            >
                <Controller
                    name="country"
                    control={control}
                    render={({ field }) => (
                        <Select
                            options={countryList}
                            placeholder="Select country"
                            value={countryList.find((c) => c.value === field.value) || null}
                            onChange={(option) => field.onChange(option?.value || '')}
                        />
                    )}
                />
            </FormItem>

            <FormItem
                label="Phone number"
                invalid={Boolean(errors.phoneNumber)}
                errorMessage={errors.phoneNumber?.message}
            >
                <Controller
                    name="phoneNumber"
                    control={control}
                    render={({ field }) => <Input type="text" placeholder="Phone number" {...field} />}
                />
            </FormItem>

            <FormItem
                label="Email"
                invalid={Boolean(errors.email)}
                errorMessage={errors.email?.message}
            >
                <Controller
                    name="email"
                    control={control}
                    render={({ field }) => <Input type="email" placeholder="Organization email" {...field} />}
                />
            </FormItem>

            <FormItem
                label="Website (optional)"
                invalid={Boolean(errors.website)}
                errorMessage={errors.website?.message}
            >
                <Controller
                    name="website"
                    control={control}
                    render={({ field }) => <Input type="text" placeholder="https://..." {...field} />}
                />
            </FormItem>

            <Button block loading={isSubmitting} variant="solid" type="submit">
                {isSubmitting ? 'Submitting...' : 'Submit for review'}
            </Button>
        </Form>
    )
}

export default OrganizationForm
