import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
  Post,
  Comment,
  CreatePostRequest,
  UpdatePostRequest,
  CreateCommentRequest,
  UpdateUserRequest,
  FeedResponse,
  SearchResponse,
  Notification,
  Course,
  Module,
  GenerateTextRequest,
  GenerateTextResponse,
  GeneratePostRequest,
  GenerateCommentRequest,
  GenerateStudyNotesRequest,
  GenerateQuizRequest,
  ExplainConceptRequest,
} from '@/types/api'

class ApiClient {
  private baseURL: string
  private accessToken: string | null = null

  constructor() {
    // In production, use relative URLs to avoid CORS issues
    if (typeof window !== 'undefined' && window.location.hostname === 'bailanysta.nd-lab.space') {
      this.baseURL = '' // Use relative URLs
    } else {
      this.baseURL = import.meta.env?.VITE_API_URL || 'http://localhost:8080'
    }
    // Initialize token from localStorage if available
    this.accessToken = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null
    
    if (this.accessToken) {
      console.log('🔑 Token loaded from localStorage:', this.accessToken.substring(0, 20) + '...')
    } else {
      console.log('❌ No token found in localStorage')
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    if (this.accessToken) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${this.accessToken}`,
      }
      console.log('✅ Request with token to:', endpoint)
    } else {
      console.log('❌ Request without token to:', endpoint)
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        
        // Handle unauthorized - clear token
        if (response.status === 401) {
          this.setAccessToken(null)
        }
        
        throw new ApiError(response.status, errorData.error?.message || response.statusText)
      }

      const contentType = response.headers.get('content-type')
      if (contentType && contentType.includes('application/json')) {
        return await response.json()
      }
      
      return {} as T
    } catch (error) {
      if (error instanceof ApiError) {
        throw error
      }
      throw new ApiError(0, 'Network error')
    }
  }

  setAccessToken(token: string | null) {
    this.accessToken = token
    if (token) {
      localStorage.setItem('access_token', token)
      console.log('✅ Token set:', token.substring(0, 20) + '...')
    } else {
      localStorage.removeItem('access_token')
      console.log('❌ Token cleared')
    }
  }

  // Auth endpoints
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    
    if (response.tokens?.access_token) {
      this.setAccessToken(response.tokens.access_token)
    }
    
    return response
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    
    if (response.tokens?.access_token) {
      this.setAccessToken(response.tokens.access_token)
    }
    
    return response
  }

  async logout(): Promise<void> {
    try {
      await this.request('/api/v1/auth/logout', {
        method: 'POST',
      })
    } finally {
      this.setAccessToken(null)
    }
  }

  async refreshToken(): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/api/v1/auth/refresh', {
      method: 'POST',
    })
    
    if (response.tokens?.access_token) {
      this.setAccessToken(response.tokens.access_token)
    }
    
    return response
  }

  // User endpoints
  async getCurrentUser(): Promise<User> {
    return this.request<User>('/api/v1/me')
  }

  async updateCurrentUser(data: UpdateUserRequest): Promise<User> {
    return this.request<User>('/api/v1/me', {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  async getAllUsers(limit = 20, offset = 0): Promise<{ users: User[]; total: number }> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    })
    return this.request<{ users: User[]; total: number }>(`/api/v1/users?${params}`)
  }

  async getUserById(id: string): Promise<User> {
    return this.request<User>(`/api/v1/users/${id}`)
  }

  async followUser(id: string): Promise<void> {
    await this.request(`/api/v1/users/${id}/follow`, {
      method: 'POST',
    })
  }

  async unfollowUser(id: string): Promise<void> {
    await this.request(`/api/v1/users/${id}/follow`, {
      method: 'DELETE',
    })
  }

  // Posts endpoints
  async createPost(data: CreatePostRequest): Promise<Post> {
    return this.request<Post>('/api/v1/posts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getPostById(id: string): Promise<Post> {
    return this.request<Post>(`/api/v1/posts/${id}`)
  }

  async updatePost(id: string, data: UpdatePostRequest): Promise<Post> {
    return this.request<Post>(`/api/v1/posts/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  async deletePost(id: string): Promise<void> {
    await this.request(`/api/v1/posts/${id}`, {
      method: 'DELETE',
    })
  }

  async likePost(id: string): Promise<void> {
    await this.request(`/api/v1/posts/${id}/like`, {
      method: 'POST',
    })
  }

  async unlikePost(id: string): Promise<void> {
    await this.request(`/api/v1/posts/${id}/like`, {
      method: 'DELETE',
    })
  }

  // Comments endpoints
  async getComments(postId: string, limit = 20, offset = 0): Promise<{ comments: Comment[] }> {
    const response = await this.request<{ comments: Comment[] | null }>(
      `/api/v1/posts/${postId}/comments?limit=${limit}&offset=${offset}`
    )
    
    // Ensure comments is always an array, even if backend returns null
    if (!response.comments) {
      response.comments = []
    }
    
    return response as { comments: Comment[] }
  }

  async createComment(postId: string, data: CreateCommentRequest): Promise<Comment> {
    return this.request<Comment>(`/api/v1/posts/${postId}/comments`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Feed endpoints
  async getFeed(limit = 20, offset = 0): Promise<FeedResponse> {
    const response = await this.request<any>(`/api/v1/feed?limit=${limit}&offset=${offset}`)
    
    // Ensure posts is always an array
    return {
      posts: Array.isArray(response.posts) ? response.posts : [],
      limit: response.limit || limit,
      offset: response.offset || offset,
    }
  }

  async getTrendingPosts(limit = 20, offset = 0): Promise<FeedResponse> {
    const response = await this.request<any>(`/api/v1/feed?limit=${limit}&offset=${offset}`)
    
    // Ensure posts is always an array and sort by likes descending
    const posts = Array.isArray(response.posts) ? response.posts : []
    return {
      posts: posts.sort((a: any, b: any) => b.like_count - a.like_count),
      limit: response.limit || limit,
      offset: response.offset || offset,
    }
  }

  // Search endpoints
  async search(query: string, limit = 20, offset = 0): Promise<SearchResponse> {
    const params = new URLSearchParams({
      query: query,
      limit: limit.toString(),
      offset: offset.toString(),
    })
    return this.request<SearchResponse>(`/api/v1/search?${params}`)
  }

  // Courses endpoints
  async getCourses(): Promise<{ courses: Course[] }> {
    return this.request<{ courses: Course[] }>('/api/v1/courses')
  }

  async getModulesByCourse(courseId: string): Promise<{ modules: Module[] }> {
    return this.request<{ modules: Module[] }>(`/api/v1/courses/${courseId}/modules`)
  }

  // Notifications endpoints
  async getNotifications(unreadOnly = false): Promise<Notification[]> {
    const params = unreadOnly ? '?unread_only=true' : ''
    const response = await this.request<{ notifications: Notification[] }>(`/api/v1/notifications${params}`)
    return response.notifications || []
  }

  async getUnreadCount(): Promise<{ count: number }> {
    const response = await this.request<{ unread_count: number }>('/api/v1/notifications/unread-count')
    return { count: response.unread_count }
  }

  async markAllNotificationsAsRead(): Promise<void> {
    await this.request('/api/v1/notifications/mark-read', {
      method: 'POST',
    })
  }

  async markNotificationAsRead(id: string): Promise<void> {
    await this.request(`/api/v1/notifications/${id}/mark-read`, {
      method: 'POST',
    })
  }

  async deleteNotification(id: string): Promise<void> {
    await this.request(`/api/v1/notifications/${id}`, {
      method: 'DELETE',
    })
  }

  // AI endpoints
  async generateText(data: GenerateTextRequest): Promise<GenerateTextResponse> {
    return this.request<GenerateTextResponse>('/api/v1/ai/generate', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async generatePost(data: GeneratePostRequest): Promise<GenerateTextResponse> {
    return this.request<GenerateTextResponse>('/api/v1/ai/generate-post', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async generateComment(data: GenerateCommentRequest): Promise<GenerateTextResponse> {
    return this.request<GenerateTextResponse>('/api/v1/ai/generate-comment', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async generateStudyNotes(data: GenerateStudyNotesRequest): Promise<GenerateTextResponse> {
    return this.request<GenerateTextResponse>('/api/v1/ai/generate-study-notes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async generateQuiz(data: GenerateQuizRequest): Promise<GenerateTextResponse> {
    return this.request<GenerateTextResponse>('/api/v1/ai/generate-quiz', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async explainConcept(data: ExplainConceptRequest): Promise<GenerateTextResponse> {
    return this.request<GenerateTextResponse>('/api/v1/ai/explain-concept', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }
}

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export const apiClient = new ApiClient()