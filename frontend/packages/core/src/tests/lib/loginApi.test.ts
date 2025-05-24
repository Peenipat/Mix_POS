/// <reference types="vitest/globals" />
import { describe, it, expect, vi } from 'vitest'
import { loginApi } from "../../lib/api/loginApi"
import api from '../../lib/axios'

vi.mock('../../lib/axios')

const mockApi = api as unknown as { post: ReturnType<typeof vi.fn> }

describe('loginApi', () => {
  it('should return parsed data when response is valid', async () => {
    const mockData = {
      user: {
        id: 1,
        username: 'test',    // ต้องมี username ตาม schema
        email: 'test@example.com',
        role: "USER",
        createdAt: "2024-04-01T12:00:00.000Z",
        updatedAt: "2024-04-01T12:00:00.000Z",
        deletedAt: null
      }
    }

    mockApi.post.mockResolvedValue({ data: mockData })

    const result = await loginApi({ email: 'test@example.com', password: '123456' })
    expect(result).toEqual(mockData)
  })

  it('should throw error if response is invalid format', async () => {
    const invalidData = {
      token: 123, // ผิด: token ควรเป็น string
      user: { name: 'test' } // ผิด schema: ควรมี id, username, email, role
    }

    mockApi.post.mockResolvedValue({ data: invalidData })

    await expect(
      loginApi({ email: 'test@example.com', password: '123456' })
    ).rejects.toThrow('Invalid response format')
  })

  it('should rethrow error if api.post fails', async () => {
    mockApi.post.mockRejectedValue(new Error('401 Unauthorized'))

    await expect(
      loginApi({ email: 'test@example.com', password: 'wrong' })
    ).rejects.toThrow('401 Unauthorized')
  })
})
