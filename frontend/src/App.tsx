import { useCallback, useEffect, useRef, useState } from 'react';
import type { FormEvent, ReactNode } from 'react';
import { api, ApiError, localDate, routes, send } from './api';
import type { Exercise, Program, ProgramDetail, Workout, WorkoutDetail, WorkoutExercise } from './api';
import './styles.css';

type Page = 'workouts' | 'programs' | 'exercises';
type Data = { workouts: WorkoutDetail[]; exercises: Exercise[]; programs: Program[] };
type Edit = { kind: Page; item?: Workout | Exercise | Program };
type Run = (action: () => Promise<void>, message: string) => Promise<boolean>;
type ConfirmDelete = (label: string, action: () => Promise<void>) => void;
const titles = { workouts: 'Workouts', programs: 'Programs', exercises: 'Exercise library' };

const singular = { workouts: 'workout', programs: 'program', exercises: 'exercise' };
const messageOf = (error: unknown) => error instanceof Error ? error.message : 'Something went wrong. Please try again.';
const dateLabel = (date: string) => new Date(`${date}T12:00:00`).toLocaleDateString(undefined, { day: 'numeric', month: 'short', year: 'numeric' });

function Modal({ title, children, close }: { title: string; children: ReactNode; close: () => void }) {
  const ref = useRef<HTMLDialogElement>(null);
  useEffect(() => { const dialog = ref.current!; dialog.showModal(); return () => dialog.close(); }, []);
  return <dialog ref={ref} aria-label={title} onCancel={event => { event.preventDefault(); close(); }} onClick={event => { if (event.target === event.currentTarget) close(); }}>
    <div className="dialog-heading"><h2>{title}</h2><button className="icon-button" type="button" aria-label="Close dialog" onClick={close}>×</button></div>{children}
  </dialog>;
}

function Brand() { return <div className="brand"><span className="brand-mark" aria-hidden="true">▥</span><span>WORKOUT<span className="brand-light"> / TRACKER</span></span></div>; }

function Auth({ onReady }: { onReady: () => void }) {
  const [register, setRegister] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const credentials = { username: String(form.get('username')).trim(), password: String(form.get('password')) };
    setBusy(true); setError(''); setNotice('');
    let created = false;
    try {
      if (register) { await send(routes.register, 'POST', credentials); created = true; setRegister(false); }
      await send(routes.login, 'POST', credentials);
      onReady();
    } catch (error) { setError(messageOf(error)); if (created) setNotice('Account created. Sign in to continue.'); }
    finally { setBusy(false); }
  }
  return <div className="auth-layout"><section className="auth-story"><Brand /><div><h1>Workout Tracker</h1></div></section>
    <main className="auth-main"><div className="auth-card"><h2>{register ? 'Start your training log' : 'Welcome back'}</h2>
      {error && <div className="notice error" role="alert">{error}</div>}{notice && <p role="status">{notice}</p>}
      <form onSubmit={submit}><fieldset disabled={busy}><label>Username<input name="username" autoComplete="username" required maxLength={100} autoFocus /></label><label>Password<input name="password" type="password" autoComplete={register ? 'new-password' : 'current-password'} required maxLength={72} /></label><button className="primary full">{busy ? 'Please wait…' : register ? 'Create account' : 'Sign in'} <span aria-hidden="true">↗</span></button></fieldset></form>
      <p className="auth-toggle">{register ? 'Already have an account?' : 'New here?'} <button className="text-button" disabled={busy} onClick={() => { setRegister(!register); setError(''); setNotice(''); }}>{register ? 'Sign in' : 'Create an account'}</button></p>
    </div></main></div>;
}

function App() {
  const [authenticated, setAuthenticated] = useState<boolean | null>(null);
  const [data, setData] = useState<Data>({ workouts: [], programs: [], exercises: [] });
  const [page, setPage] = useState<Page>('workouts');
  const [calendarMonth, setCalendarMonth] = useState(() => new Date());
  const [selected, setSelected] = useState<number | null>(null);
  const [revision, setRevision] = useState(0);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [edit, setEdit] = useState<Edit | null>(null);
  const [deletion, setDeletion] = useState<{ label: string; action: () => Promise<void> } | null>(null);
  const [query, setQuery] = useState('');
  const onAuthError = useCallback(() => {
    setAuthenticated(false); setData({ workouts: [], programs: [], exercises: [] });
    setSelected(null); setEdit(null); setDeletion(null);
  }, []);

  useEffect(() => {
    let active = true;
    Promise.all([api<WorkoutDetail[]>(routes.workouts), api<Exercise[]>(routes.exercises), api<Program[]>(routes.programs)])
      .then(([workouts, exercises, programs]) => { if (active) { setData({ workouts, exercises, programs }); setAuthenticated(true); } })
      .catch(error => { if (active) { if (error instanceof ApiError && error.status === 401) onAuthError(); else setError(messageOf(error)); } })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, [revision, onAuthError]);

  const run: Run = useCallback(async (action, message) => {
    setBusy(true); setError(''); setNotice('');
    try { await action(); setNotice(message); setRevision(value => value + 1); return true; }
    catch (error) { if (error instanceof ApiError && error.status === 401) onAuthError(); setError(messageOf(error)); setRevision(value => value + 1); return false; }
    finally { setBusy(false); }
  }, [onAuthError]);

  const confirmDelete: ConfirmDelete = (label, action) => { setError(''); setDeletion({ label, action }); };
  function navigate(next: Page) { setPage(next); setSelected(null); setQuery(''); setError(''); setNotice(''); }

  if (authenticated === false) return <Auth onReady={() => { setError(''); setLoading(true); setAuthenticated(null); setRevision(value => value + 1); }} />;
  if (authenticated === null) return <main className="initial-state"><Brand />{error ? <><p className="notice error" role="alert">{error}</p><button onClick={() => { setError(''); setRevision(value => value + 1); }}>Retry connection</button></> : <p role="status">Opening your training log…</p>}</main>;

  const items = data[page].filter(item => item.name.toLowerCase().includes(query.toLowerCase()));
  const monthWorkouts = data.workouts.filter(workout => workout.date.slice(0, 7) === `${calendarMonth.getFullYear()}-${String(calendarMonth.getMonth() + 1).padStart(2, '0')}`);
  return <div className="app-shell"><aside className="sidebar"><Brand /><nav aria-label="Main navigation">{(['workouts', 'programs', 'exercises'] as Page[]).map((item, index) => <button key={item} disabled={busy} aria-current={page === item ? 'page' : undefined} className={page === item ? 'nav-item active' : 'nav-item'} onClick={() => navigate(item)}><span className="nav-symbol" aria-hidden="true">{['▦', '▤', '◇'][index]}</span><span>{item[0].toUpperCase() + item.slice(1)}</span><span className="nav-count">{data[item].length}</span></button>)}</nav></aside>
    <main className="workspace"><header className="topbar"><time dateTime={localDate(new Date())}>{new Date().toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })}</time></header>
      <div className="page-heading"><div><h1>{titles[page]}</h1></div><button className="primary" disabled={busy} onClick={() => { setError(''); setEdit({ kind: page }); }}>＋ New {singular[page]}</button></div>
      {page === 'workouts' && <section className="panel calendar-panel"><div className="calendar-top"><div><h2>{calendarMonth.toLocaleDateString(undefined,{month:'long',year:'numeric'})}</h2></div><div><button aria-label="Previous month" onClick={()=>setCalendarMonth(new Date(calendarMonth.getFullYear(),calendarMonth.getMonth()-1,1))}>‹</button><button aria-label="Next month" onClick={()=>setCalendarMonth(new Date(calendarMonth.getFullYear(),calendarMonth.getMonth()+1,1))}>›</button></div></div><Calendar month={calendarMonth} workouts={data.workouts} onDate={date=>{setError('');setNotice('');const found=data.workouts.find(workout=>workout.date===date);if(found)setSelected(found.id);else setEdit({kind:'workouts',item:{id:0,name:'',date,program_id:null}})}}/><p className="calendar-hint">● Workout logged  Click any date to log a session</p>{monthWorkouts.length>0&&<div className="calendar-days">{monthWorkouts.map(workout=><button key={workout.id} onClick={()=>setSelected(workout.id)}><strong>{Number(workout.date.slice(8))}</strong>{workout.name || 'Workout'}</button>)}</div>}</section>}
      {error && <div className="notice error" role="alert">{error}<button className="text-button" disabled={busy} onClick={() => { setError(''); setRevision(value => value + 1); }}>Retry</button></div>}
      {notice && <div className="notice success" role="status">✓ {notice}<button className="icon-button" aria-label="Dismiss notification" onClick={() => setNotice('')}>×</button></div>}
      <fieldset className="content-fieldset" disabled={busy}><div className={page === 'exercises' ? 'content-grid library-layout' : 'content-grid'}>
        <section className="panel collection"><div className="section-heading"><h2>{page === 'workouts' ? 'Session history' : page === 'programs' ? 'All programs' : 'All exercises'}</h2><span className="badge">{data[page].length}</span></div><label className="search-label"><span className="sr-only">Search {page}</span><input type="search" placeholder={`Search ${page}…`} value={query} onChange={event => setQuery(event.target.value)} /></label>
          <div className="item-list">{loading ? <p className="empty">Loading…</p> : items.length === 0 ? <div className="empty"><span className="empty-symbol" aria-hidden="true">◇</span><h3>{query ? 'No matches' : `Your ${page} start here`}</h3><p>{query ? 'Try another name.' : `Create your first ${singular[page]} to get started.`}</p>{!query && <button className="text-button" onClick={() => setEdit({ kind: page })}>＋ Add {singular[page]}</button>}</div> : items.map(item => <div key={item.id} className={`list-row ${selected === item.id ? 'selected' : ''}`}>
            {page === 'exercises' ? <><div className="exercise-monogram" aria-hidden="true">{item.name.slice(0, 2).toUpperCase()}</div><strong className="row-title">{item.name}</strong><button className="text-button" onClick={() => { setError(''); setEdit({ kind: page, item }); }}>Edit</button><button className="text-button danger" onClick={() => confirmDelete(item.name, async () => { await send(`${routes.exercises}?id=${item.id}`, 'DELETE'); })}>Delete</button></> : <button className="select-row" onClick={() => { setSelected(item.id); setError(''); setNotice(''); }}><span className="row-icon" aria-hidden="true">{page === 'workouts' ? '▦' : '▤'}</span><span className="row-text"><strong>{item.name || 'Workout'}</strong><span className="muted">{page === 'workouts' ? dateLabel((item as Workout).date) : 'Exercise program'}</span></span><span className="row-arrow" aria-hidden="true">↗</span></button>}
          </div>)}</div>
        </section>
      {page !== 'exercises' && <section className="panel detail-panel" aria-label={`${singular[page]} details`}>{selected === null ? <div className="empty detail-empty"><span className="empty-symbol" aria-hidden="true">{page === 'workouts' ? '▦' : '▤'}</span><h2>Select a {singular[page]}</h2><p>{page === 'workouts' ? 'Click any calendar date to log a workout, or choose a session from the calendar.' : 'Open a program to manage its exercises.'}</p></div> : <Detail key={`${page}-${selected}-${revision}`} page={page} id={selected} data={data} run={run} onEdit={item => { setError(''); setEdit({ kind: page, item }); }} onDelete={confirmDelete} onDeleted={() => setSelected(null)} onAuthError={onAuthError} />}</section>}
      </div></fieldset>
      <footer className="page-footer">Make the next session count.<span>WORKOUT / TRACKER</span></footer>
    </main>
    {edit && <EditModal edit={edit} programs={data.programs} busy={busy} error={error} close={() => { if (!busy) { setEdit(null); setError(''); } }} save={async body => {
      const target = routes[edit.kind] + (edit.item?.id ? `?id=${edit.item.id}` : '');
      const ok = await run(async () => { const result = await send<Workout | Program | Exercise>(target, edit.item?.id ? 'PUT' : 'POST', body); if (edit.kind === page && edit.kind !== 'exercises') setSelected(result.id); }, `${singular[edit.kind][0].toUpperCase() + singular[edit.kind].slice(1)} saved.`);
      if (ok) setEdit(null);
    }} />}
    {deletion && <Modal title={`Delete ${deletion.label}?`} close={() => { if (!busy) setDeletion(null); }}><p className="muted">This removes the selected record from your log. Removing an exercise from a workout also deletes its logged sets.</p>{error && <p className="notice error" role="alert">{error}</p>}<div className="dialog-actions"><button disabled={busy} onClick={() => setDeletion(null)}>Cancel</button><button className="danger-button" disabled={busy} onClick={async () => { if (await run(deletion.action, 'Record deleted.')) setDeletion(null); }}>{busy ? 'Deleting…' : 'Delete'}</button></div></Modal>}
  </div>;
}

function EditModal({ edit, programs, busy, error, close, save }: { edit: Edit; programs: Program[]; busy: boolean; error: string; close: () => void; save: (body: unknown) => Promise<void> }) {
  const workout = edit.item as Workout | undefined;
  return <Modal title={`${edit.item?.id ? 'Edit' : 'New'} ${singular[edit.kind]}`} close={close}>
    <form onSubmit={event => { event.preventDefault(); const form = new FormData(event.currentTarget); const body = { name: edit.kind === 'workouts' ? '' : String(form.get('name')).trim(), ...(edit.kind === 'workouts' ? { date: form.get('date'), program_id: form.get('program') ? Number(form.get('program')) : null } : {}) }; void save(body); }}>
      <fieldset disabled={busy}>{edit.kind !== 'workouts' && <label>Name<input name="name" required autoFocus defaultValue={edit.item?.name ?? ''} placeholder={edit.kind === 'programs' ? 'e.g. Push day' : 'e.g. Bench press'} maxLength={200} pattern=".*\S.*" /></label>}
      {edit.kind === 'workouts' && <><label>Date<input name="date" type="date" required defaultValue={workout?.date ?? localDate(new Date())} /></label><label>Program <span className="muted">(optional)</span><select name="program" defaultValue={workout?.program_id ?? ''}><option value="">No program</option>{programs.map(program => <option key={program.id} value={program.id}>{program.name}</option>)}</select></label><p className="field-help">Add exercises and sets after saving.</p></>}
      {error && <p className="notice error" role="alert">{error}</p>}<div className="dialog-actions"><button type="button" onClick={close}>Cancel</button><button className="primary">{busy ? 'Saving…' : `Save ${singular[edit.kind]}`}</button></div></fieldset>
    </form>
  </Modal>;
}

function Detail({ page, id, data, run, onEdit, onDelete, onDeleted, onAuthError }: { page: 'workouts' | 'programs'; id: number; data: Data; run: Run; onEdit: (item: Workout | Program) => void; onDelete: ConfirmDelete; onDeleted: () => void; onAuthError: () => void }) {
  const [detail, setDetail] = useState<WorkoutDetail | ProgramDetail | null>(null);
  const [error, setError] = useState('');
  const [retry, setRetry] = useState(0);
  const [programID, setProgramID] = useState('');
  useEffect(() => {
    let active = true;
    api<WorkoutDetail | ProgramDetail>(page === 'workouts' ? routes.workout(id) : routes.program(id))
      .then(value => { if (active) setDetail(value); })
      .catch(error => { if (active) { if (error instanceof ApiError && error.status === 401) onAuthError(); else setError(messageOf(error)); } });
    return () => { active = false; };
  }, [page, id, retry, onAuthError]);
  if (error) return <div className="empty"><p role="alert">{error}</p><button onClick={() => { setError(''); setRetry(value => value + 1); }}>Retry</button></div>;
  if (!detail) return <p className="empty" role="status">Loading details…</p>;
  const workout = page === 'workouts' ? detail as WorkoutDetail : null;
  const available = data.exercises.filter(exercise => !detail.exercises.some(entry => entry.exercise_id === exercise.id));
  return <><div className="detail-heading"><div className="detail-title"><h2>{workout ? dateLabel(workout.date) : detail.name}</h2><div className="row-actions"><button className="text-button" onClick={() => onEdit(detail)}>Edit</button><button className="text-button danger" onClick={() => onDelete(workout ? `workout on ${dateLabel(workout.date)}` : detail.name, async () => { await send(page === 'workouts' ? routes.workout(id) : routes.program(id), 'DELETE'); onDeleted(); })}>Delete</button></div></div>{workout?.program_name && <span className="program-tag">▤ {workout.program_name}</span>}<p className="muted">{detail.exercises.length} exercises{workout ? ` · ${workout.exercises.reduce((count, exercise) => count + exercise.sets.length, 0)} sets logged` : ''}</p></div>
    <div className="detail-body">{detail.exercises.length === 0 && <div className="empty compact"><h3>No exercises yet</h3><p>Add an exercise below to {workout ? 'start logging your sets' : 'build this program'}.</p></div>}
      {detail.exercises.map(exercise => workout ? <ExerciseCard key={exercise.id} exercise={exercise as WorkoutExercise} run={run} onDelete={onDelete} /> : <div className="program-exercise" key={exercise.id}><span className="exercise-monogram" aria-hidden="true">◇</span><strong>{exercise.name}</strong><button className="text-button danger" onClick={() => onDelete(exercise.name, async () => { await send(routes.programExercise(exercise.id), 'DELETE'); })}>Remove</button></div>)}
      <div className="add-exercise"><h3>Add an exercise</h3>{data.exercises.length === 0 ? <p className="muted">Create an exercise in the Exercise library first, then add it here.</p> : available.length === 0 ? <p className="muted">All your library exercises are already included.</p> : <form className="inline-form" onSubmit={event => {
        event.preventDefault(); const exerciseID = Number(new FormData(event.currentTarget).get('exercise'));
        void run(async () => { await send(workout ? routes.workoutExercises : routes.programExercises, 'POST', { [workout ? 'workout_id' : 'program_id']: id, exercise_id: exerciseID }); }, 'Exercise added.');
      }}><label className="grow"><span className="sr-only">Exercise to add</span><select name="exercise" required defaultValue=""><option value="" disabled>Choose an exercise</option>{available.map(exercise => <option value={exercise.id} key={exercise.id}>{exercise.name}</option>)}</select></label><button className="primary">＋ Add</button></form>}
      {workout && data.programs.length > 0 && <><label>Program<select aria-label="Program to add" value={programID} onChange={event => setProgramID(event.target.value)}><option value="">Choose a program</option>{data.programs.map(program => <option key={program.id} value={program.id}>{program.name}</option>)}</select></label><button disabled={!programID} className="text-button import-program" onClick={() => void run(async () => {
        const program = await api<ProgramDetail>(routes.program(Number(programID)));
        const missing = program.exercises.filter(exercise => !workout.exercises.some(entry => entry.exercise_id === exercise.exercise_id));
        if (!missing.length) throw new Error('No new exercises to add from this program.');
        let added = 0;
        try { for (const exercise of missing) { await send(routes.workoutExercises, 'POST', { workout_id: id, exercise_id: exercise.exercise_id }); added++; } }
        catch (error) { throw new Error(`${added} exercises added. ${messageOf(error)} Retry to add the remaining exercises.`, { cause: error }); }
      }, 'Program exercises added to this workout.')}>＋ Add program exercises</button></>}
      </div>
    </div>
  </>;
}

function ExerciseCard({ exercise, run, onDelete }: { exercise: WorkoutExercise; run: Run; onDelete: ConfirmDelete }) {
  return <article className="exercise-card"><div className="exercise-card-heading"><h3>{exercise.name}</h3><button className="text-button danger" onClick={() => onDelete(exercise.name, async () => { await send(routes.workoutExercise(exercise.id), 'DELETE'); })}>Remove</button></div>
    {exercise.sets.length > 0 ? <table><caption className="sr-only">Sets for {exercise.name}</caption><thead><tr><th scope="col">Set</th><th scope="col">Reps</th><th scope="col">Weight</th><th scope="col"><span className="sr-only">Actions</span></th></tr></thead><tbody>{exercise.sets.map(set => <tr key={set.id}><td><span className="set-number">{set.set_number}</span></td><td>{set.reps ?? '—'}</td><td>{set.weight ?? '—'}</td><td><button className="text-button danger" aria-label={`Delete set ${set.set_number} of ${exercise.name}`} onClick={() => onDelete(`set ${set.set_number}`, async () => { await send(routes.exerciseSet(set.id), 'DELETE'); })}>Delete</button></td></tr>)}</tbody></table> : <p className="no-sets">No sets logged yet. Start with your first set.</p>}
    <form className="set-form" onSubmit={event => { event.preventDefault(); const form = new FormData(event.currentTarget); void run(async () => { await send(routes.exerciseSets, 'POST', { workout_exercise_id: exercise.id, reps: form.get('reps') === '' ? null : Number(form.get('reps')), weight: form.get('weight') === '' ? null : Number(form.get('weight')) }); }, 'Set logged.'); }}>
      <label>Reps<input aria-label={`Reps for ${exercise.name}`} name="reps" type="number" min="0" step="1" placeholder="e.g. 10" /></label><label>Weight<input aria-label={`Weight for ${exercise.name}`} name="weight" type="number" min="0" max="99999.99" step="0.01" placeholder="e.g. 40" /></label><button className="primary">＋ Log set</button>
    </form><p className="field-help">Use the same weight unit across your sessions. Blank values are saved as unrecorded.</p>
  </article>;
}

function Calendar({month,workouts,onDate}:{month:Date;workouts:WorkoutDetail[];onDate:(date:string)=>void}) {
  const first=new Date(month.getFullYear(),month.getMonth(),1);
  const offset=(first.getDay()+6)%7;
  const days=new Date(month.getFullYear(),month.getMonth()+1,0).getDate();
  const marked=new Set(workouts.map(workout=>workout.date));
  return <div className="calendar-grid">{['Mon','Tue','Wed','Thu','Fri','Sat','Sun'].map(day=><span className="calendar-weekday" key={day}>{day}</span>)}{Array.from({length:offset},(_,i)=><span key={`empty-${i}`}/>)}{Array.from({length:days},(_,i)=>{const day=i+1;const date=`${month.getFullYear()}-${String(month.getMonth()+1).padStart(2,'0')}-${String(day).padStart(2,'0')}`;return <button key={date} aria-label={`${day} ${month.toLocaleDateString(undefined,{month:'long'})}${marked.has(date)?', workout logged':''}`} className={`calendar-day${marked.has(date)?' has-workout':''}${date===localDate(new Date())?' today':''}`} onClick={()=>onDate(date)}>{day}</button>})}</div>;
}

export default App;
