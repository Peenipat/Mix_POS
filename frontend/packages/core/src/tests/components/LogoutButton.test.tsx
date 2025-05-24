import { render, screen, fireEvent } from '@testing-library/react'
import LogoutButton from '../../components/LogoutButton'
import { beforeEach, describe, expect, it, vi } from 'vitest'

//  mock global confirm
vi.spyOn(window, 'confirm')

// mock useNavigate
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

//  mock useAppDispatch + logout
const mockDispatch = vi.fn()
vi.mock('@/store/hook', () => ({
  useAppDispatch: () => mockDispatch,
}))

vi.mock('@/store/authSlice', () => ({
  logout: () => ({ type: 'LOGOUT' }),
}))

describe('LogoutButton', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should confirm before logout, then dispatch and navigate', () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true) // mock ตรงนี้ให้ return true

    render(<LogoutButton />)

    fireEvent.click(screen.getByText(/logout/i))

    expect(window.confirm).toHaveBeenCalled()
    expect(mockDispatch).toHaveBeenCalledWith({ type: 'LOGOUT' })
    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('should not logout if user cancels confirm', () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)

    render(<LogoutButton />)

    fireEvent.click(screen.getByText(/logout/i))

    expect(window.confirm).toHaveBeenCalled()
    expect(mockDispatch).not.toHaveBeenCalled()
    expect(mockNavigate).not.toHaveBeenCalled()
  })
})
