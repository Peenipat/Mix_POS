import { renderHook, act } from '@testing-library/react'
import { useTypewriter } from '@/hooks/useTypewriter'
import { vi } from 'vitest'

describe('useTypewriter', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.clearAllTimers()
    vi.useRealTimers()
  })

  it('should type and delete message', () => {
    const { result } = renderHook(() =>
      useTypewriter(['Hello'], 100, 50, 300)
    )

    expect(result.current).toBe('')

    act(() => vi.advanceTimersByTime(100))
    expect(result.current).toBe('H')

    act(() => vi.advanceTimersByTime(100))
    expect(result.current).toBe('He')

    act(() => vi.advanceTimersByTime(3000)) // simulate typing + deleting ทั้ง loop
    expect(result.current).toMatch(/[Helo]/) // ระหว่าง cycle
  })
})
