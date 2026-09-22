CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- USERS
CREATE TABLE users (
    id            uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    email         varchar(254) NOT NULL UNIQUE,
    first_name    varchar(50)  NOT NULL,
    last_name     varchar(50)  NOT NULL DEFAULT '',
    password_hash varchar(255) NOT NULL,
    is_admin      boolean      NOT NULL DEFAULT false,
    is_active     boolean      NOT NULL DEFAULT true,
    created_at    timestamptz  NOT NULL DEFAULT now(),
    updated_at    timestamptz  NOT NULL DEFAULT now(),
    deleted_at    timestamptz
);

-- THEMES
CREATE TABLE themes (
    id            uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    name          varchar(500) NOT NULL,
    description   text         NOT NULL DEFAULT '',
    is_active     boolean      NOT NULL DEFAULT false,
    created_by    uuid         NOT NULL REFERENCES users(id),
    max_points    int          NOT NULL DEFAULT 0 CHECK (max_points   >= 0),
    check_points  int          NOT NULL DEFAULT 0 CHECK (check_points >= 0),
    image_path    varchar(255) NOT NULL DEFAULT '',
    attempt_count int          NOT NULL DEFAULT 1 CHECK (attempt_count > 0),
    created_at    timestamptz  NOT NULL DEFAULT now(),

    CONSTRAINT themes_check_points_le_max CHECK (check_points <= max_points)
);

CREATE INDEX idx_themes_created_by ON themes(created_by);
CREATE INDEX idx_themes_active     ON themes(id) WHERE is_active;

-- QUESTIONS
CREATE TABLE questions (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    theme_id      uuid        NOT NULL REFERENCES themes(id) ON DELETE CASCADE,
    name          text        NOT NULL,
    points        int         NOT NULL DEFAULT 0 CHECK (points >= 0),
    sort_order    int         NOT NULL DEFAULT 0,
    question_type varchar(20) NOT NULL
        CHECK (question_type IN ('single', 'multi')),
    created_at    timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT questions_theme_sort_unique UNIQUE (theme_id, sort_order)
);

CREATE INDEX idx_questions_theme_id ON questions(theme_id);

-- ANSWERS
CREATE TABLE answers (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id uuid        NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    name        text        NOT NULL,
    is_correct  boolean     NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_answers_question_id ON answers(question_id);

-- TEST ATTEMPTS
CREATE TABLE test_attempts (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      uuid        NOT NULL REFERENCES users(id),
    theme_id     uuid        NOT NULL REFERENCES themes(id),
    attempt_num  int         NOT NULL CHECK (attempt_num > 0),
    status       varchar(20) NOT NULL
        CHECK (status IN ('in_progress', 'completed', 'expired')),
    max_points   int         NOT NULL DEFAULT 0 CHECK (max_points   >= 0),
    check_points int         NOT NULL DEFAULT 0 CHECK (check_points >= 0),
    points       int         NOT NULL DEFAULT 0 CHECK (points       >= 0),
    is_passed    boolean     NOT NULL DEFAULT false,
    started_at   timestamptz NOT NULL DEFAULT now(),
    completed_at timestamptz,

    CONSTRAINT test_attempts_user_theme_num_unique
        UNIQUE (user_id, theme_id, attempt_num)
);

CREATE INDEX idx_test_attempts_user_id  ON test_attempts(user_id);
CREATE INDEX idx_test_attempts_theme_id ON test_attempts(theme_id);

-- USER ANSWERS

CREATE TABLE user_answers (
    id              uuid    PRIMARY KEY DEFAULT gen_random_uuid(),
    test_attempt_id uuid    NOT NULL REFERENCES test_attempts(id) ON DELETE CASCADE,
    question_id     uuid    NOT NULL REFERENCES questions(id),
    is_correct      boolean NOT NULL DEFAULT false,
    points          int     NOT NULL DEFAULT 0 CHECK (points >= 0),

    CONSTRAINT user_answers_attempt_question_unique
        UNIQUE (test_attempt_id, question_id)
);

CREATE INDEX idx_user_answers_attempt_id ON user_answers(test_attempt_id);

CREATE TABLE user_answer_choices (
    user_answer_id uuid NOT NULL REFERENCES user_answers(id) ON DELETE CASCADE,
    answer_id      uuid NOT NULL REFERENCES answers(id),
    PRIMARY KEY (user_answer_id, answer_id)
);
