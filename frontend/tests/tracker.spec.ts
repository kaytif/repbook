import { expect, test } from '@playwright/test';
import type { Page } from '@playwright/test';

async function createNamed(page: Page, kind: string, name: string) {
  await page.getByRole('button', { name: `＋ New ${kind}`, exact: true }).click();
  const dialog = page.getByRole('dialog');
  await dialog.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
  await dialog.getByRole('button', { name: `Save ${kind}` }).click();
  await expect(dialog).toBeHidden();
}

async function confirmDelete(page: Page) {
  const dialog = page.getByRole('dialog');
  await dialog.getByRole('button', { name: 'Delete', exact: true }).click();
  await expect(dialog).toBeHidden();
}

test('create account, manage programs, log sets, reload, edit and delete', async ({ page }, testInfo) => {
  const errors: string[] = [];
  page.on('pageerror', error => errors.push(error.message));
  await page.goto('/');
  await page.getByRole('button', { name: 'Create an account' }).click();
  const username = `e2e_${testInfo.project.name}_${Date.now()}`;
  await page.getByRole('textbox', { name: 'Username', exact: true }).fill(username);
  await page.getByLabel('Password', { exact: true }).fill('test-password-123');
  await page.getByRole('button', { name: 'Create account', exact: false }).click();
  await expect(page.getByRole('heading', { name: 'Workouts', exact: true })).toBeVisible();
  await page.getByRole('button', { name: /Exercises/ }).click();
  await createNamed(page, 'exercise', 'Bench press');
  await expect(page.getByText('Bench press', { exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'Edit', exact: true }).click();
  await page.getByRole('dialog').getByRole('textbox', { name: 'Name', exact: true }).fill('Barbell bench press');
  await page.getByRole('button', { name: 'Save exercise' }).click();
  await expect(page.getByRole('dialog')).toBeHidden();
  await page.getByRole('button', { name: /Programs/ }).click();
  await createNamed(page, 'program', 'Push day');
  const details = page.getByRole('region', { name: 'program details' });
  await details.getByRole('combobox', { name: 'Exercise to add' }).selectOption({ label: 'Barbell bench press' });
  await details.getByRole('button', { name: '＋ Add', exact: true }).click();
  await expect(details.getByText('Barbell bench press', { exact: true })).toBeVisible();
  await page.getByRole('button', { name: /Workouts/ }).click();
  await page.getByRole('button', { name: '9 November', exact: true }).click();
  let dialog = page.getByRole('dialog');
  await expect(dialog.getByRole('textbox', { name: 'Name', exact: true })).toHaveCount(0);
  await dialog.getByLabel('Date', { exact: true }).fill('2026-11-09');
  await dialog.getByRole('combobox').selectOption({ label: 'Push day' });
  await dialog.getByRole('button', { name: 'Save workout' }).click();
  await expect(dialog).toBeHidden();
  const workout = page.getByRole('region', { name: 'workout details' });
  await workout.getByRole('combobox', { name: 'Program to add' }).selectOption({ label: 'Push day' });
  await workout.getByRole('button', { name: '＋ Add program exercises' }).click();
  await expect(workout.getByRole('heading', { name: 'Barbell bench press' })).toBeVisible();
  await workout.getByRole('spinbutton', { name: 'Reps for Barbell bench press' }).fill('10');
  await workout.getByRole('spinbutton', { name: 'Weight for Barbell bench press' }).fill('42.5');
  await workout.getByRole('button', { name: '＋ Log set' }).click();
  await expect(workout.getByRole('cell', { name: '42.5', exact: true })).toBeVisible();
  await page.reload();
  await page.getByRole('button', { name: /9 Nov 2026|Nov 9, 2026/ }).click();
  await expect(workout.getByRole('cell', { name: '42.5', exact: true })).toBeVisible();
  await expect(workout.getByRole('cell', { name: '10', exact: true })).toBeVisible();
  await expect(workout.getByRole('combobox', { name: 'Exercise to add' })).toHaveCount(0);
  await workout.getByRole('button', { name: 'Edit', exact: true }).click();
  dialog = page.getByRole('dialog');
  await expect(dialog.getByRole('combobox')).toHaveValue(/\d+/);
  await expect(dialog.getByRole('textbox', { name: 'Name', exact: true })).toHaveCount(0);
  await dialog.getByLabel('Date', { exact: true }).fill('2026-11-10');
  await dialog.getByRole('button', { name: 'Save workout' }).click();
  await expect(dialog).toBeHidden();
  await expect(workout.getByRole('heading', { name: /10 Nov 2026|Nov 10, 2026/ })).toBeVisible();
  await page.screenshot({ path: `test-results/${testInfo.project.name}-workout.png`, fullPage: true });
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await workout.getByRole('button', { name: 'Delete set 1 of Barbell bench press' }).click();
  await confirmDelete(page);
  await expect(workout.getByRole('cell', { name: '42.5', exact: true })).toHaveCount(0);
  await workout.getByRole('button', { name: 'Remove', exact: true }).click();
  await confirmDelete(page);
  await expect(workout.getByRole('heading', { name: 'No exercises yet' })).toBeVisible();
  await workout.getByRole('button', { name: 'Delete', exact: true }).click();
  await confirmDelete(page);
  await expect(page.getByRole('heading', { name: 'Your workouts start here' })).toBeVisible();
  await page.getByRole('button', { name: /Programs/ }).click();
  await page.getByRole('button', { name: /Push day/ }).click();
  await details.getByRole('button', { name: 'Remove', exact: true }).click();
  await confirmDelete(page);
  await expect(details.getByRole('heading', { name: 'No exercises yet' })).toBeVisible();
  await details.getByRole('button', { name: 'Delete', exact: true }).click();
  await confirmDelete(page);
  await page.getByRole('button', { name: /Exercises/ }).click();
  await page.getByRole('button', { name: 'Delete', exact: true }).click();
  await confirmDelete(page);
  await expect(page.getByRole('heading', { name: 'Your exercises start here' })).toBeVisible();
  expect(errors).toEqual([]);
});

test('login failures stay visible and expired sessions return to sign in', async ({ page, context }) => {
  await page.goto('/');
  await page.getByRole('textbox', { name: 'Username', exact: true }).fill('missing-user');
  await page.getByLabel('Password', { exact: true }).fill('wrong-password');
  await page.getByRole('button', { name: /Sign in/ }).click();
  await expect(page.getByRole('alert')).toContainText('Invalid username or password');
  await context.clearCookies();
  await page.reload();
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
});
