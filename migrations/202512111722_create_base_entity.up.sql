CREATE TABLE IF NOT EXISTS public.teams (
    team_name   VARCHAR(100) PRIMARY KEY,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS public.users (
    user_id    VARCHAR(100) PRIMARY KEY,
    username   VARCHAR(100) NOT NULL,
    team_name  VARCHAR(100) REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE TYPE public."pr_status" AS ENUM ('OPEN', 'MERGED');


CREATE TABLE IF NOT EXISTS public.prs (
    pr_id               VARCHAR(100) PRIMARY KEY,
    title               VARCHAR(100) NOT NULL,
    author_id           VARCHAR(100) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status              public."pr_status" NOT NULL DEFAULT 'OPEN',
    need_more_reviewers BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    merged_at           TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.pr_reviewers (
    pr_id    VARCHAR(100) NOT NULL REFERENCES prs(pr_id) ON DELETE CASCADE,
    user_id  VARCHAR(100) NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (pr_id, user_id)
);



CREATE INDEX users_team_idx ON public.users(team_name);
CREATE INDEX users_team_active_idx ON public.users(team_name, is_active);

CREATE INDEX prs_author_idx ON public.prs(author_id);