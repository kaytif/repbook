// Centralize API paths: adjust these to match the routes in your Go router.
export const routes = {
  login: '/api/login', register: '/api/register', logout: '/api/logout',
  workouts: '/api/workouts', exercises: '/api/exercises', programs: '/api/programs',
  program: (id: number) => `/api/programs?id=${id}`,
  programExercises: '/api/program-exercises',
  programExercise: (id: number) => `/api/program-exercises?id=${id}`,
  workoutProgram: '/api/workout-program',
  workout: (id: number) => `/api/workouts?id=${id}`,
  workoutExercises: '/api/workout-exercises',
  workoutExercise: (id: number) => `/api/workout-exercises?id=${id}`,
  exerciseSets: '/api/exercise-sets',
  exerciseSet: (id: number) => `/api/exercise-sets?id=${id}`,
};

export type Workout = { id: number; name: string; date: string; program_id?: number | null };
export type Exercise = { id: number; name: string };
export type Program = { id: number; name: string };
export type ProgramDetail = Program & { exercises: { id: number; exercise_id: number; name: string }[] };
export type ExerciseSet = { id: number; workout_exercise_id: number; set_number: number; reps: number | null; weight: number | null };
export type WorkoutExercise = { id: number; exercise_id: number; name: string; sets: ExerciseSet[] };
export type WorkoutDetail = Workout & { program_name?: string | null; exercises: WorkoutExercise[] };

export class ApiError extends Error {
  constructor(public status: number, message: string) { super(message); this.name = 'ApiError'; }
}

// JSON helper: sends the session cookie and converts non-2xx responses into errors.
export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(path, {
    ...init,
    credentials: 'same-origin',
    headers: { ...(init.body ? { 'Content-Type': 'application/json' } : {}), ...init.headers },
  });
  const body = await response.text();
  if (!response.ok) {
    let message = body || `Request failed (${response.status})`;
    try {
      const parsed = JSON.parse(body);
      if (typeof parsed.error === 'string') message = parsed.error;
    } catch { /* Plain text errors can be displayed directly. */ }
    throw new ApiError(response.status, message);
  }
  if (!body) return undefined as T;
  try { return JSON.parse(body) as T; }
  catch { return body as T; }
}

export function send<T>(path: string, method: 'POST' | 'PUT' | 'PATCH' | 'DELETE', body?: unknown) {
  return api<T>(path, { method, ...(body === undefined ? {} : { body: JSON.stringify(body) }) });
}

// Date-only values should use local calendar fields, not toISOString() (UTC).
export function localDate(date: Date) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
}

// Supports either an array or {workouts:[...]} / {exercises:[...]} responses.
export function asList<T>(value: unknown, key: string): T[] {
  if (Array.isArray(value)) return value as T[];
  if (value && typeof value === 'object' && key in value) {
    const nested = (value as Record<string, unknown>)[key];
    if (Array.isArray(nested)) return nested as T[];
  }
  return [];
}
