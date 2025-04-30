import { render, screen } from '@testing-library/react'
import ProtectedRoute from '@/components/ProtectedRoute'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'

// mock useAppSelector
vi.mock('@/store/hook', () => ({
  useAppSelector: vi.fn(),
}))

// import หลังจาก mock
import { useAppSelector } from '@/store/hook'

describe('ProtectedRoute', () => {
  it('should render children if user exists', () => {
    (useAppSelector as any).mockReturnValue({ id: 1, name: 'User' })

    render(
      <MemoryRouter>
        <ProtectedRoute>
          <div>Protected Content</div>
        </ProtectedRoute>
      </MemoryRouter>
    )

    expect(screen.getByText('Protected Content')).toBeInTheDocument()
  })

  it('should not render children if user is not logged in', () => {
    (useAppSelector as any).mockReturnValue(null)

    render(
      <MemoryRouter initialEntries={['/protected']}>
        <ProtectedRoute>
          <div>Protected Content</div>
        </ProtectedRoute>
      </MemoryRouter>
    )

    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument()
  })
})
