CREATE TABLE users (

    -- Unique ID auto generated
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- USERNAME that is text
    username TEXT NOT NULL UNIQUE CHECK (length(username) > 0),

    -- PASSWORD HASH
    password_hash TEXT NOT NULL CHECK (length(password_hash) > 0)

);


-- WORKOUTS TABLE
-- Stores each account
CREATE TABLE exercises (
    -- UNIQUE ID automatically generated
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- Exercise Name
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),

    -- Add user ID
    user_id INTEGER NOT NULL REFERENCES users(id),

    UNIQUE(user_id,  name)

);

CREATE TABLE programs (
    -- UNIQUE ID automatically generated
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- Program name
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),

    -- Add user ID
    user_id INTEGER NOT NULL REFERENCES users(id),

    UNIQUE(user_id,  name)

);


CREATE TABLE workouts (
    -- UNIQUE ID automatically generated
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- Workout name
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),

    -- Workout must have a date
    date DATE NOT NULL,

    -- Workout must be inputted at a certain time
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Timestamp when workout was last updated
    updated_at TIMESTAMPTZ,

    -- Timestamp used for soft deletion
    deleted_at TIMESTAMPTZ,

    -- Program that is in the workout
    program_id INTEGER REFERENCES programs(id) ON DELETE SET NULL,

    -- Add user ID
    user_id INTEGER NOT NULL REFERENCES users(id)
    
);


CREATE TABLE workout_exercises (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    workout_id INTEGER NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,

    exercise_id INTEGER NOT NULL REFERENCES exercises(id)
);

CREATE TABLE program_exercises (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,

    exercise_id INTEGER NOT NULL REFERENCES exercises(id),

    UNIQUE (program_id, exercise_id)
);

CREATE TABLE exercise_sets (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    workout_exercise_id INTEGER NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    
    set_number INTEGER NOT NULL CHECK (set_number > 0),

    reps INTEGER CHECK (reps >= 0),

    weight NUMERIC(7,2) CHECK (weight >= 0),

    UNIQUE (workout_exercise_id, set_number)
);



CREATE TABLE sessions (
    user_id INTEGER NOT NULL REFERENCES users(id),
    token_hash TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    )

