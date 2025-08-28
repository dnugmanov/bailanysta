export interface User {
  id: string
  username: string
  email: string
  bio?: string
  avatar_url?: string | null
  followers_count?: number
  following_count?: number
  is_following?: boolean
  created_at: string
}

export interface Course {
  id: string
  title: string
  description: string
}

export interface Module {
  id: string
  course_id: string
  title: string
  order: number
}

export interface Post {
  id: string
  author_id: string
  author?: User
  text: string
  course_id?: string | null
  module_id?: string | null
  course?: Course | null
  module?: Module | null
  hashtags?: string[]
  like_count: number
  comment_count: number
  is_liked?: boolean
  created_at: string
  updated_at: string
}

export interface Comment {
  id: string
  post_id: string
  author_id: string
  author?: User
  text: string
  created_at: string
}

export interface Notification {
  id: string
  user_id: string
  type: 'like' | 'comment' | 'follow' | 'new_post'
  entity_id: string
  payload: Record<string, any>
  read_at?: string | null
  created_at: string
  actor?: User
  post?: Post
}

// Auth types
export interface AuthTokens {
  access_token: string
  refresh_token?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface AuthResponse {
  user: User
  tokens: {
    access_token: string
    refresh_token: string
  }
}

// Request types
export interface CreatePostRequest {
  text: string
  course_id?: string
  module_id?: string
}

export interface UpdatePostRequest {
  text: string
  course_id?: string
  module_id?: string
}

export interface CreateCommentRequest {
  text: string
}

export interface UpdateUserRequest {
  bio?: string
  avatar_url?: string
}

// AI types
export interface GenerateTextRequest {
  prompt: string
  max_tokens?: number
  temperature?: number
}

export interface GenerateTextResponse {
  text: string
}

export interface GeneratePostRequest {
  topic: string
  style?: 'casual' | 'academic' | 'professional'
  length?: 'short' | 'medium' | 'long'
  format?: 'markdown' | 'plain'
}

export interface GenerateCommentRequest {
  post_text: string
  tone?: 'supportive' | 'questioning' | 'informative'
}

export interface GenerateStudyNotesRequest {
  topic: string
  course_context?: string
}

export interface GenerateQuizRequest {
  topic: string
  difficulty?: 'easy' | 'medium' | 'hard'
  question_count?: number
}

export interface ExplainConceptRequest {
  concept: string
  context?: string
  complexity_level?: 'beginner' | 'intermediate' | 'advanced'
}

// Response wrappers
export interface PaginatedResponse<T> {
  data: T[]
  limit: number
  offset: number
  total?: number
  has_more?: boolean
}

export interface ApiResponse<T> {
  data?: T
  error?: {
    code: string
    message: string
  }
}

// Feed types
export interface FeedResponse {
  posts: Post[]
  limit: number
  offset: number
}

export interface SearchResponse {
  posts: Post[]
  users: User[]
  limit: number
  offset: number
}

// Error types
export interface ApiError {
  error: {
    code: string
    message: string
  }
}

export type ErrorCode = 
  | 'BAD_REQUEST'
  | 'UNAUTHORIZED'
  | 'FORBIDDEN'
  | 'NOT_FOUND'
  | 'CONFLICT'
  | 'INTERNAL_SERVER_ERROR'
  | 'UNKNOWN_ERROR'