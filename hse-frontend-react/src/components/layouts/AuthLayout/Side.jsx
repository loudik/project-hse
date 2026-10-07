import { cloneElement, useEffect, useState } from 'react'

const slides = [
    {
        eyebrow: 'Autoridade Nacional do Petróleo',
        title: 'Entry authorisation for vessels',
        body: "HSE review and approval for every vessel movement into Timor-Leste's petroleum contract area.",
        icon: 'vessel',
    },
    {
        eyebrow: 'Autoridade Nacional do Petróleo',
        title: 'Entry authorisation for helicopters',
        body: "HSE review and approval for every helicopter movement into Timor-Leste's petroleum contract area.",
        icon: 'helicopter',
    },
    {
        eyebrow: 'Autoridade Nacional do Petróleo',
        title: 'HSE review, start to finish',
        body: 'Submission, officer assignment, staff review and decision — tracked in one place for every application.',
        icon: 'both',
    },
]

// Fixed positions so stars don't reshuffle on every render
const stars = [
    [60, 60, 1.4, 0.7], [120, 110, 1, 0.4], [220, 50, 1.6, 0.6], [310, 90, 1, 0.3],
    [40, 180, 1.2, 0.5], [160, 40, 1, 0.45], [260, 150, 1.4, 0.55], [380, 60, 1, 0.35],
    [440, 150, 1.6, 0.6], [70, 260, 1, 0.3], [480, 230, 1.2, 0.45], [200, 220, 1, 0.3],
    [350, 190, 1.3, 0.5], [500, 90, 1, 0.4], [20, 100, 1, 0.35],
]

const Side = ({ children, ...rest }) => {
    const [active, setActive] = useState(0)

    useEffect(() => {
        const timer = setInterval(() => {
            setActive((prev) => (prev + 1) % slides.length)
        }, 5000)
        return () => clearInterval(timer)
    }, [])

    const slide = slides[active]
    const showVessel = slide.icon === 'vessel' || slide.icon === 'both'
    const showHeli = slide.icon === 'helicopter' || slide.icon === 'both'

    return (
        <div className="flex h-full p-6 bg-white dark:bg-gray-800">
            <div className="flex flex-col flex-1 items-center overflow-y-auto py-10">
                <div className="m-auto w-full xl:max-w-[450px] px-8 max-w-[380px]">
                    {children
                        ? cloneElement(children, {
                              ...rest,
                          })
                        : null}
                </div>
            </div>
            <div className="py-6 px-10 lg:flex flex-col flex-1 justify-between hidden rounded-3xl items-end relative xl:max-w-[520px] 2xl:max-w-[720px] overflow-hidden bg-gradient-to-b from-[#071A2C] via-[#0E3252] to-[#134569]">
                <svg
                    className="absolute inset-0 h-full w-full"
                    viewBox="0 0 600 800"
                    preserveAspectRatio="xMidYMid slice"
                    fill="none"
                >
                    <style>
                        {`
                            @keyframes pulseRing {
                                0% { r: 10; opacity: 0.7; }
                                75%, 100% { r: 32; opacity: 0; }
                            }
                            @keyframes driftWave {
                                0%, 100% { transform: translateX(0); }
                                50% { transform: translateX(16px); }
                            }
                            @keyframes twinkle {
                                0%, 100% { opacity: var(--base-op); }
                                50% { opacity: calc(var(--base-op) * 0.35); }
                            }
                            .pulse-ring { animation: pulseRing 2.6s ease-out infinite; }
                            .drift-wave { animation: driftWave 11s ease-in-out infinite; }
                            .star { animation: twinkle 4s ease-in-out infinite; }
                            @media (prefers-reduced-motion: reduce) {
                                .pulse-ring, .drift-wave, .star { animation: none !important; }
                            }
                        `}
                    </style>

                    <defs>
                        <radialGradient id="sunrise" cx="50%" cy="42%" r="65%">
                            <stop offset="0%" stopColor="#F2B705" stopOpacity="0.4" />
                            <stop offset="35%" stopColor="#F2B705" stopOpacity="0.15" />
                            <stop offset="100%" stopColor="#F2B705" stopOpacity="0" />
                        </radialGradient>
                    </defs>

                    {/* stars, upper sky */}
                    {stars.map(([x, y, r, op], i) => (
                        <circle
                            key={i}
                            className="star"
                            style={{ '--base-op': op, animationDelay: `${(i % 5) * 0.6}s` }}
                            cx={x}
                            cy={y}
                            r={r}
                            fill="#FFFFFF"
                            fillOpacity={op}
                        />
                    ))}

                    {/* sunrise glow behind the entry point */}
                    <rect width="600" height="800" fill="url(#sunrise)" />

                    {/* faint nautical-chart grid */}
                    <line x1="0" y1="200" x2="600" y2="200" stroke="#FFFFFF" strokeOpacity="0.05" strokeWidth="1" />
                    <line x1="150" y1="0" x2="150" y2="800" stroke="#FFFFFF" strokeOpacity="0.04" strokeWidth="1" />
                    <line x1="450" y1="0" x2="450" y2="800" stroke="#FFFFFF" strokeOpacity="0.04" strokeWidth="1" />

                    <line x1="0" y1="340" x2="600" y2="340" stroke="#F2B705" strokeOpacity="0.22" strokeWidth="1" />
                    <path d="M0 420 C 150 400, 220 460, 380 430 S 560 400, 600 420" stroke="#FFFFFF" strokeOpacity="0.1" strokeWidth="1" />
                    <g className="drift-wave">
                        <path d="M0 480 C 140 460, 240 510, 400 485 S 560 460, 600 480" stroke="#FFFFFF" strokeOpacity="0.09" strokeWidth="1" />
                    </g>
                    <path d="M0 540 C 160 520, 260 565, 420 540 S 560 520, 600 540" stroke="#FFFFFF" strokeOpacity="0.08" strokeWidth="1" />
                    <g className="drift-wave" style={{ animationDelay: '-4s' }}>
                        <path d="M0 610 C 150 595, 250 630, 400 612 S 560 595, 600 610" stroke="#FFFFFF" strokeOpacity="0.07" strokeWidth="1" />
                    </g>
                    <path d="M0 670 C 160 655, 260 690, 410 672 S 560 655, 600 670" stroke="#FFFFFF" strokeOpacity="0.06" strokeWidth="1" />
                    <path d="M0 730 C 150 718, 260 748, 400 732 S 560 718, 600 730" stroke="#FFFFFF" strokeOpacity="0.05" strokeWidth="1" />

                    <g transform="translate(500,110)" stroke="#F2B705" strokeOpacity="0.6" strokeWidth="1">
                        <circle r="34" fill="none" />
                        <circle r="4" fill="#F2B705" fillOpacity="0.8" stroke="none" />
                        <line x1="0" y1="-34" x2="0" y2="34" />
                        <line x1="-34" y1="0" x2="34" y2="0" />
                        <text x="0" y="-42" textAnchor="middle" fill="#F2B705" fillOpacity="0.8" fontSize="11" stroke="none">N</text>
                    </g>

                    {showVessel && (
                        <>
                            <path d="M480 500 Q 410 410 330 330" stroke="#F2B705" strokeOpacity="0.55" strokeWidth="1.5" strokeDasharray="4 5" />
                            <g transform="translate(480,500)" stroke="#FFFFFF" strokeOpacity="0.85" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M-40 10 L40 10 L30 26 L-28 26 Z" />
                                <line x1="-10" y1="10" x2="-10" y2="-24" />
                                <line x1="-10" y1="-24" x2="14" y2="-24" />
                                <line x1="12" y1="10" x2="12" y2="-10" />
                                <path d="M-48 28 Q -20 36 10 28" strokeOpacity="0.4" />
                            </g>
                        </>
                    )}

                    {showHeli && (
                        <>
                            <path d="M160 150 Q 300 230 330 330" stroke="#F2B705" strokeOpacity="0.55" strokeWidth="1.5" strokeDasharray="4 5" />
                            <g transform="translate(160,150)" stroke="#FFFFFF" strokeOpacity="0.85" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                                <line x1="-32" y1="0" x2="32" y2="0" strokeOpacity="0.5" />
                                <line x1="-26" y1="0" x2="26" y2="0" />
                                <line x1="-3" y1="-2" x2="-3" y2="-16" />
                                <line x1="-13" y1="-16" x2="7" y2="-16" />
                                <path d="M-6 0 C -6 10, 20 10, 24 2" />
                                <path d="M24 2 L34 2" />
                                <path d="M30 -3 L30 7" />
                            </g>
                        </>
                    )}

                    {/* entry point marker with radar pulse */}
                    <circle cx="330" cy="330" r="22" fill="#F2B705" fillOpacity="0.15" />
                    <circle className="pulse-ring" cx="330" cy="330" stroke="#F2B705" strokeWidth="1" fill="none" />
                    <circle className="pulse-ring" cx="330" cy="330" stroke="#F2B705" strokeWidth="1" fill="none" style={{ animationDelay: '1.3s' }} />
                    <circle cx="330" cy="330" r="12" stroke="#F2B705" strokeOpacity="0.6" strokeWidth="1" fill="none" />
                    <circle cx="330" cy="330" r="5" fill="#F2B705" />

                    <path d="M330 330 L300 560" stroke="#F2B705" strokeOpacity="0.4" strokeWidth="1.5" strokeDasharray="4 5" />

                    <g transform="translate(300,560)" stroke="#FFFFFF" strokeOpacity="0.8" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <line x1="-34" y1="40" x2="-24" y2="0" />
                        <line x1="34" y1="40" x2="24" y2="0" />
                        <line x1="0" y1="40" x2="0" y2="0" />
                        <line x1="-30" y1="20" x2="30" y2="20" />
                        <rect x="-38" y="-14" width="76" height="14" fill="#F2B705" fillOpacity="0.12" />
                    </g>
                    <text x="300" y="625" textAnchor="middle" fill="#FFFFFF" fillOpacity="0.55" fontSize="12" stroke="none">
                        Contract area
                    </text>
                </svg>

                <div className="relative z-10 max-w-[380px]">
                    <p className="text-[#F2B705] font-semibold mb-2">
                        {slide.eyebrow}
                    </p>
                    <h1 className="text-neutral mb-4">{slide.title}</h1>
                    <p className="text-neutral opacity-80 font-semibold">
                        {slide.body}
                    </p>
                    <div className="flex gap-2 mt-5 flex-wrap">
                        {['Vessels', 'Helicopters', 'HSE review'].map((tag) => (
                            <span
                                key={tag}
                                className="text-xs font-semibold text-neutral border border-white/25 rounded-full px-3 py-1 opacity-80"
                            >
                                {tag}
                            </span>
                        ))}
                    </div>
                </div>

                <div className="relative z-10 flex gap-2 mt-8">
                    {slides.map((_, i) => (
                        <button
                            key={i}
                            type="button"
                            aria-label={`Show slide ${i + 1}`}
                            onClick={() => setActive(i)}
                            className={`h-1.5 rounded-full transition-all ${
                                i === active
                                    ? 'w-6 bg-[#F2B705]'
                                    : 'w-1.5 bg-white/30 hover:bg-white/50'
                            }`}
                        />
                    ))}
                </div>
            </div>
        </div>
    )
}
export default Side