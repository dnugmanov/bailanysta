-- Initial database setup for development
-- This file runs after migrations

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS posts_author_id_idx ON posts(author_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS posts_created_at_idx ON posts(created_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS comments_post_id_idx ON comments(post_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS likes_post_id_idx ON likes(post_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS follows_follower_id_idx ON follows(follower_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS follows_followee_id_idx ON follows(followee_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS notifications_user_id_idx ON notifications(user_id);

-- Create text search indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS posts_text_trgm_idx ON posts USING gin(text gin_trgm_ops);
CREATE INDEX CONCURRENTLY IF NOT EXISTS users_username_trgm_idx ON users USING gin(username gin_trgm_ops);

-- Insert test data (only in development)
DO $$
DECLARE
    user1_id UUID;
    user2_id UUID;
    user3_id UUID;
    course1_id UUID;
    module1_id UUID;
    post1_id UUID;
    post2_id UUID;
BEGIN
    -- Insert test users
    INSERT INTO users (username, email, password_hash, bio, avatar_url)
    VALUES ('alice', 'alice@example.com', '$2a$10$example.hash.here', 'Hello, I am Alice! I love learning and teaching.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=alice')
    RETURNING id INTO user1_id;

    INSERT INTO users (username, email, password_hash, bio, avatar_url)
    VALUES ('bob', 'bob@example.com', '$2a$10$example.hash.here', 'Bob here! Software developer and AI enthusiast.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=bob')
    RETURNING id INTO user2_id;

    INSERT INTO users (username, email, password_hash, bio, avatar_url)
    VALUES ('charlie', 'charlie@example.com', '$2a$10$example.hash.here', 'Charlie - Data Science and Machine Learning expert.', 'https://api.dicebear.com/7.x/avataaars/svg?seed=charlie')
    RETURNING id INTO user3_id;

    -- Insert test courses
    INSERT INTO courses (title, description)
    VALUES ('Introduction to Machine Learning', 'Learn the basics of machine learning algorithms and techniques.')
    RETURNING id INTO course1_id;

    -- Insert test modules
    INSERT INTO modules (course_id, title, "order")
    VALUES (course1_id, 'Linear Regression Fundamentals', 1)
    RETURNING id INTO module1_id;

    -- Insert test posts
    INSERT INTO posts (author_id, text, course_id, module_id)
    VALUES (user1_id, 'Just finished an amazing course on #machine_learning! Highly recommend to anyone interested in AI. The practical examples were incredibly helpful. #learning #AI', course1_id, module1_id)
    RETURNING id INTO post1_id;

    INSERT INTO posts (author_id, text)
    VALUES (user2_id, 'Excited to share my latest project on neural networks! Check it out and let me know what you think. #deep_learning #neural_networks #AI')
    RETURNING id INTO post2_id;

    INSERT INTO posts (author_id, text)
    VALUES (user3_id, 'Great discussion today about the future of data science. The field is evolving so rapidly! What are your thoughts on automation in DS? #datascience #automation #future')
    RETURNING id INTO post1_id;

    -- Insert test comments
    INSERT INTO comments (post_id, author_id, text)
    VALUES (post1_id, user2_id, 'Congrats on completing the course! Which parts did you find most challenging?');

    INSERT INTO comments (post_id, author_id, text)
    VALUES (post2_id, user1_id, 'This looks amazing! Would love to see the code implementation.');

    -- Insert test likes
    INSERT INTO likes (user_id, post_id) VALUES (user2_id, post1_id);
    INSERT INTO likes (user_id, post_id) VALUES (user3_id, post1_id);
    INSERT INTO likes (user_id, post_id) VALUES (user1_id, post2_id);

    -- Insert test follows
    INSERT INTO follows (follower_id, followee_id) VALUES (user2_id, user1_id);
    INSERT INTO follows (follower_id, followee_id) VALUES (user3_id, user1_id);
    INSERT INTO follows (follower_id, followee_id) VALUES (user1_id, user2_id);

    -- Insert test hashtags
    INSERT INTO hashtags (tag) VALUES ('machine_learning') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('learning') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('AI') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('deep_learning') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('neural_networks') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('datascience') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('automation') ON CONFLICT (tag) DO NOTHING;
    INSERT INTO hashtags (tag) VALUES ('future') ON CONFLICT (tag) DO NOTHING;

    RAISE NOTICE 'Test data inserted successfully!';
END $$;
