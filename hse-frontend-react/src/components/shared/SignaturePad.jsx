import { useRef, useState, useEffect } from 'react'
import Button from '@/components/ui/Button'

// Renders a small canvas the user can draw on with mouse/touch/stylus, plus
// an "Upload instead" fallback for devices without a drawing surface.
// Exposes the final signature as a PNG Blob via onChange.
export default function SignaturePad({ onChange }) {
    const canvasRef = useRef(null)
    const drawingRef = useRef(false)
    const [hasDrawing, setHasDrawing] = useState(false)
    const [mode, setMode] = useState('draw') // 'draw' | 'upload'
    const [uploadPreview, setUploadPreview] = useState(null)

    useEffect(() => {
        const canvas = canvasRef.current
        if (!canvas) return
        const ctx = canvas.getContext('2d')
        ctx.fillStyle = '#ffffff'
        ctx.fillRect(0, 0, canvas.width, canvas.height)
        ctx.strokeStyle = '#111827'
        ctx.lineWidth = 2
        ctx.lineCap = 'round'
    }, [mode])

    const getPos = (e) => {
        const rect = canvasRef.current.getBoundingClientRect()
        const point = e.touches ? e.touches[0] : e
        return {
            x: point.clientX - rect.left,
            y: point.clientY - rect.top,
        }
    }

    const startDraw = (e) => {
        e.preventDefault()
        drawingRef.current = true
        const { x, y } = getPos(e)
        const ctx = canvasRef.current.getContext('2d')
        ctx.beginPath()
        ctx.moveTo(x, y)
    }

    const draw = (e) => {
        if (!drawingRef.current) return
        e.preventDefault()
        const { x, y } = getPos(e)
        const ctx = canvasRef.current.getContext('2d')
        ctx.lineTo(x, y)
        ctx.stroke()
        setHasDrawing(true)
    }

    const endDraw = () => {
        if (!drawingRef.current) return
        drawingRef.current = false
        canvasRef.current.toBlob((blob) => {
            if (blob) onChange?.(blob)
        }, 'image/png')
    }

    const clearCanvas = () => {
        const canvas = canvasRef.current
        const ctx = canvas.getContext('2d')
        ctx.fillStyle = '#ffffff'
        ctx.fillRect(0, 0, canvas.width, canvas.height)
        setHasDrawing(false)
        onChange?.(null)
    }

    const handleUpload = (e) => {
        const file = e.target.files?.[0]
        if (!file) return
        setUploadPreview(URL.createObjectURL(file))
        onChange?.(file)
    }

    return (
        <div>
            <div className="mb-2 flex gap-2">
                <Button
                    size="xs"
                    variant={mode === 'draw' ? 'solid' : 'default'}
                    onClick={() => setMode('draw')}
                >
                    Draw
                </Button>
                <Button
                    size="xs"
                    variant={mode === 'upload' ? 'solid' : 'default'}
                    onClick={() => setMode('upload')}
                >
                    Upload instead
                </Button>
            </div>

            {mode === 'draw' ? (
                <div>
                    <canvas
                        ref={canvasRef}
                        width={400}
                        height={150}
                        className="cursor-crosshair rounded-lg border border-gray-300 bg-white"
                        onMouseDown={startDraw}
                        onMouseMove={draw}
                        onMouseUp={endDraw}
                        onMouseLeave={endDraw}
                        onTouchStart={startDraw}
                        onTouchMove={draw}
                        onTouchEnd={endDraw}
                    />
                    <div className="mt-2">
                        <Button size="xs" onClick={clearCanvas} disabled={!hasDrawing}>
                            Clear
                        </Button>
                    </div>
                </div>
            ) : (
                <div>
                    <input type="file" accept="image/*" onChange={handleUpload} />
                    {uploadPreview && (
                        <img
                            src={uploadPreview}
                            alt="Signature preview"
                            className="mt-2 max-h-[100px] rounded border border-gray-200"
                        />
                    )}
                </div>
            )}
        </div>
    )
}