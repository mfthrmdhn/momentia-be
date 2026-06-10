# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: api/register.test.ts >> Register API Tests >> TC-REG-011 - should fail registration with duplicate msisdn
- Location: tests/api/register.test.ts:14:13

# Error details

```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 400
Received: 200
```

# Test source

```ts
  1  | import { test, expect } from '@playwright/test';
  2  | import { assert } from 'superstruct';
  3  | import { UserEndpoints } from '../../src/api/endpoints/user.endpoint';
  4  | import { registerScenarios } from '../../src/scenario/register.scenario';
  5  | 
  6  | let userEndpoints: UserEndpoints;
  7  | 
  8  | test.describe('Register API Tests', () => {
  9  |     test.beforeEach(async ({ request }) => {
  10 |         userEndpoints = new UserEndpoints(request);
  11 |     });
  12 | 
  13 |     for (const scenario of registerScenarios) {
  14 |         test(scenario.description, async () => {
  15 |             const response = await userEndpoints.register(scenario.payload);
  16 |             const responseBody = await response.json();
  17 | 
> 18 |             expect(response.status()).toBe(scenario.expectedStatus);
     |                                       ^ Error: expect(received).toBe(expected) // Object.is equality
  19 |             assert(responseBody, scenario.expectedSchema);
  20 | 
  21 |             expect(responseBody.success).toBe(scenario.expectedSuccess);
  22 |             expect(responseBody.message).toBe(scenario.expectedMessage);
  23 |             if (scenario.expectedError) {
  24 |                 expect(responseBody.error).toMatch(scenario.expectedError);
  25 |             }
  26 |         });
  27 |     }
  28 | });
  29 | 
```