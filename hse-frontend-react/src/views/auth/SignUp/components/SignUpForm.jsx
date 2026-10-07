import { useState } from 'react'
import Input from '@/components/ui/Input'
import Button from '@/components/ui/Button'
import Checkbox from '@/components/ui/Checkbox'
import { FormItem, Form } from '@/components/ui/Form'
import { useAuth } from '@/auth'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const validationSchema = z
    .object({
        email: z.email({ message: 'Please enter a valid email' }),
        userName: z.string().min(1, { message: 'Please enter your name' }),
        password: z.string().min(6, { message: 'Password must be at least 6 characters' }),
        confirmPassword: z.string().min(1, { message: 'Confirm Password Required' }),
        phoneNumber: z.string().min(1, { message: 'Please enter your phone number' }),
        position: z.string().min(1, { message: 'Please enter your position' }),
        organizationName: z.string().min(1, { message: 'Please enter your organization' }),
        location: z.string().min(1, { message: 'Please enter your location' }),
        website: z.string().optional(),
        acceptTerms: z.boolean().refine((val) => val === true, {
            message: 'You must accept the Terms and Conditions to continue',
        }),
    })
    .refine((data) => data.password === data.confirmPassword, {
        message: 'Password not match',
        path: ['confirmPassword'],
    })

const SignUpForm = (props) => {
    const { disableSubmit = false, className, setMessage, onRegistered } = props
    const [isSubmitting, setSubmitting] = useState(false)
    const { signUp } = useAuth()
    const {
        handleSubmit,
        formState: { errors },
        control,
    } = useForm({
        defaultValues: {
            acceptTerms: false,
        },
        resolver: zodResolver(validationSchema),
    })

    const onSignUp = async (values) => {
        const {
            userName,
            password,
            email,
            phoneNumber,
            position,
            organizationName,
            location,
            website,
            acceptTerms,
        } = values

        if (!disableSubmit) {
            setSubmitting(true)
            const result = await signUp({
                name: userName,
                password,
                email,
                phoneNumber,
                position,
                organizationName,
                location,
                website,
                acceptTerms,
            })

            if (result?.status === 'failed') {
                setMessage?.(result.message)
            } else if (result?.status === 'success') {
                onRegistered?.(result.message)
            }

            setSubmitting(false)
        }
    }

    return (
        <div className={className}>
            <Form onSubmit={handleSubmit(onSignUp)}>
                <h6 className="mb-3 heading-text">Account</h6>
                <FormItem
                    label="Email"
                    invalid={Boolean(errors.email)}
                    errorMessage={errors.email?.message}
                >
                    <Controller
                        name="email"
                        control={control}
                        render={({ field }) => (
                            <Input type="email" placeholder="Email" autoComplete="off" {...field} />
                        )}
                    />
                </FormItem>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-4">
                    <FormItem
                        label="Password"
                        invalid={Boolean(errors.password)}
                        errorMessage={errors.password?.message}
                    >
                        <Controller
                            name="password"
                            control={control}
                            render={({ field }) => (
                                <Input type="password" autoComplete="off" placeholder="Password" {...field} />
                            )}
                        />
                    </FormItem>
                    <FormItem
                        label="Confirm password"
                        invalid={Boolean(errors.confirmPassword)}
                        errorMessage={errors.confirmPassword?.message}
                    >
                        <Controller
                            name="confirmPassword"
                            control={control}
                            render={({ field }) => (
                                <Input type="password" autoComplete="off" placeholder="Confirm password" {...field} />
                            )}
                        />
                    </FormItem>
                </div>

                <h6 className="mb-3 mt-2 heading-text">Your details</h6>
                <FormItem
                    label="Full name"
                    invalid={Boolean(errors.userName)}
                    errorMessage={errors.userName?.message}
                >
                    <Controller
                        name="userName"
                        control={control}
                        render={({ field }) => (
                            <Input type="text" placeholder="Full name" autoComplete="off" {...field} />
                        )}
                    />
                </FormItem>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-4">
                    <FormItem
                        label="Phone number"
                        invalid={Boolean(errors.phoneNumber)}
                        errorMessage={errors.phoneNumber?.message}
                    >
                        <Controller
                            name="phoneNumber"
                            control={control}
                            render={({ field }) => (
                                <Input type="text" placeholder="Phone number" autoComplete="off" {...field} />
                            )}
                        />
                    </FormItem>
                    <FormItem
                        label="Position"
                        invalid={Boolean(errors.position)}
                        errorMessage={errors.position?.message}
                    >
                        <Controller
                            name="position"
                            control={control}
                            render={({ field }) => (
                                <Input type="text" placeholder="e.g. HSE Officer" autoComplete="off" {...field} />
                            )}
                        />
                    </FormItem>
                </div>

                <h6 className="mb-3 mt-2 heading-text">Organization</h6>
                <FormItem
                    label="Organization"
                    invalid={Boolean(errors.organizationName)}
                    errorMessage={errors.organizationName?.message}
                >
                    <Controller
                        name="organizationName"
                        control={control}
                        render={({ field }) => (
                            <Input type="text" placeholder="Organization name" autoComplete="off" {...field} />
                        )}
                    />
                </FormItem>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-4">
                    <FormItem
                        label="Location"
                        invalid={Boolean(errors.location)}
                        errorMessage={errors.location?.message}
                    >
                        <Controller
                            name="location"
                            control={control}
                            render={({ field }) => (
                                <Input type="text" placeholder="City, country" autoComplete="off" {...field} />
                            )}
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
                            render={({ field }) => (
                                <Input type="text" placeholder="https://..." autoComplete="off" {...field} />
                            )}
                        />
                    </FormItem>
                </div>

                <FormItem
                    invalid={Boolean(errors.acceptTerms)}
                    errorMessage={errors.acceptTerms?.message}
                    className="mt-2"
                >
                    <Controller
                        name="acceptTerms"
                        control={control}
                        render={({ field }) => (
                            <Checkbox
                                checked={field.value}
                                onChange={(checked) => field.onChange(checked)}
                            >
                                I agree to the{' '}
                                <a href="/terms" target="_blank" rel="noreferrer" className="underline">
                                    Terms and Conditions
                                </a>
                            </Checkbox>
                        )}
                    />
                </FormItem>
                <Button block loading={isSubmitting} variant="solid" type="submit">
                    {isSubmitting ? 'Creating Account...' : 'Sign Up'}
                </Button>
            </Form>
        </div>
    )
}

export default SignUpForm