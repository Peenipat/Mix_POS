import { render, screen } from '@testing-library/react'
import RequireRole from '../../components/RequireRole'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'

// mock useAppSelector
vi.mock('@/store/hook', () => ({
  useAppSelector: vi.fn(),
}))

import { useAppSelector } from '../../store/hook'

describe('RequireRole', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should redirect to "/" if user is not logged in', () => {
    ;(useAppSelector as any).mockReturnValue(null)

    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <RequireRole roles={['ADMIN']}>
          <div>Admin Only Content</div>
        </RequireRole>
      </MemoryRouter>
    )

    expect(screen.queryByText('Admin Only Content')).not.toBeInTheDocument()
  })

  it('should redirect to "/unauthorized" if user role is not allowed', () => {
    ;(useAppSelector as any).mockReturnValue({
      id: 1,
      name: 'User',
      role: 'USER',
    })

    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <RequireRole roles={['ADMIN']}>
          <div>Admin Only Content</div>
        </RequireRole>
      </MemoryRouter>
    )

    expect(screen.queryByText('Admin Only Content')).not.toBeInTheDocument()
  })

  it('should render children if user has allowed role', () => {
    ;(useAppSelector as any).mockReturnValue({
      id: 1,
      name: 'Admin',
      role: 'ADMIN',
    })

    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <RequireRole roles={['ADMIN']}>
          <div>Admin Only Content</div>
        </RequireRole>
      </MemoryRouter>
    )

    expect(screen.getByText('Admin Only Content')).toBeInTheDocument()
  })
})
