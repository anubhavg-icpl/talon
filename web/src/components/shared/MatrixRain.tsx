'use client'

// React Imports
import { useEffect, useRef } from 'react'

// Util Imports
import { cn } from '@/lib/utils'

const GLYPHS = 'アカサタナハマヤラワ0123456789ABCDEF<>/\\$#*+=-'.split('')

/**
 * Canvas-based falling-character rain. Electric cyan operator accent,
 * resize-aware, animation cancelled on unmount, disabled when the
 * user prefers reduced motion.
 */
const MatrixRain = ({ className, opacity = 0.08 }: { className?: string; opacity?: number }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current

    if (!canvas) return
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return

    const ctx = canvas.getContext('2d')

    if (!ctx) return

    const fontSize = 14
    let columns = 0
    let drops: number[] = []
    let raf = 0
    let last = 0

    const resize = () => {
      const parent = canvas.parentElement

      if (!parent) return

      canvas.width = parent.clientWidth
      canvas.height = parent.clientHeight
      columns = Math.floor(canvas.width / fontSize)
      drops = Array.from({ length: columns }, () => Math.floor(Math.random() * (canvas.height / fontSize)))
    }

    resize()

    const observer = new ResizeObserver(resize)

    if (canvas.parentElement) observer.observe(canvas.parentElement)

    const draw = (time: number) => {
      raf = requestAnimationFrame(draw)

      // throttle to ~20fps
      if (time - last < 50) return
      last = time

      ctx.fillStyle = 'rgba(6, 10, 14, 0.2)'
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      ctx.font = `${fontSize}px monospace`
      ctx.fillStyle = '#22d3ee'

      for (let i = 0; i < columns; i++) {
        const char = GLYPHS[Math.floor(Math.random() * GLYPHS.length)]

        ctx.fillText(char, i * fontSize, drops[i] * fontSize)

        if (drops[i] * fontSize > canvas.height && Math.random() > 0.975) drops[i] = 0
        drops[i]++
      }
    }

    raf = requestAnimationFrame(draw)

    return () => {
      cancelAnimationFrame(raf)
      observer.disconnect()
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      aria-hidden
      className={cn('pointer-events-none absolute inset-0 size-full', className)}
      style={{ opacity }}
    />
  )
}

export default MatrixRain
