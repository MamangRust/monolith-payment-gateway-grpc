import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:5000';
const TOKEN =
  'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNyIsImF1ZCI6WyJhY2Nlc3MiXSwiZXhwIjoxNzcwNjAyNDc3fQ.vTTFF7CxENCjLmQdZtIv9cBzuSxpBoMJVrV_mBldlXU';

export default function () {
  const params = {
    headers: { Authorization: TOKEN, 'Content-Type': 'application/json' },
  };

  const basicEndpoints = [
    // ===== Card Query =====
    '/api/card-query?page=1&limit=10',
    '/api/card-query/active?page=1&limit=10',
    '/api/card-query/trashed?page=1&limit=10',
    '/api/card-query/user?user_id=11',
    '/api/card-query/card_number/4111111111111111',
    '/api/card-query/11',

    // ===== Card Dashboard =====
    '/api/card-dashboard/dashboard',
    '/api/card-dashboard/dashboard/4111111111111111',

    // ===== Balance Stats =====
    '/api/card-stats-balance/monthly-balance?year=2025&month=1',
    '/api/card-stats-balance/yearly-balance?year=2025',
    '/api/card-stats-balance/monthly-balance-by-card?year=2025&month=1&card_number=4111111111111111',
    '/api/card-stats-balance/yearly-balance-by-card?year=2025&card_number=4111111111111111',

    // ===== Topup Stats =====
    '/api/card-stats-topup/monthly-topup-amount?year=2025&month=1',
    '/api/card-stats-topup/yearly-topup-amount?year=2025',
    '/api/card-stats-topup/monthly-topup-amount-by-card?year=2025&month=1&card_number=4111111111111111',
    '/api/card-stats-topup/yearly-topup-amount-by-card?year=2025&card_number=4111111111111111',

    // ===== Withdraw Stats =====
    '/api/card-stats-withdraw/monthly-withdraw-amount?year=2025&month=1',
    '/api/card-stats-withdraw/yearly-withdraw-amount?year=2025',
    '/api/card-stats-withdraw/monthly-withdraw-amount-by-card?year=2025&month=1&card_number=4111111111111111',
    '/api/card-stats-withdraw/yearly-withdraw-amount-by-card?year=2025&card_number=4111111111111111',

    // ===== Transaction Stats =====
    '/api/card-stats-transaction/monthly-transaction-amount?year=2025&month=1',
    '/api/card-stats-transaction/yearly-transaction-amount?year=2025',
    '/api/card-stats-transaction/monthly-transaction-amount-by-card?year=2025&month=1&card_number=4111111111111111',
    '/api/card-stats-transaction/yearly-transaction-amount-by-card?year=2025&card_number=4111111111111111',

    // ===== Transfer Stats =====
    '/api/card-stats-transfer/monthly-transfer-sender-amount?year=2025&month=1',
    '/api/card-stats-transfer/yearly-transfer-sender-amount?year=2025',

    '/api/card-stats-transfer/monthly-transfer-receiver-amount?year=2025&month=1',
    '/api/card-stats-transfer/yearly-transfer-receiver-amount?year=2025',

    '/api/card-stats-transfer/monthly-transfer-sender-amount-by-card?year=2025&month=1&card_number=4111111111111111',
    '/api/card-stats-transfer/yearly-transfer-sender-amount-by-card?year=2025&card_number=4111111111111111',

    '/api/card-stats-transfer/monthly-transfer-receiver-amount-by-card?year=2025&month=1&card_number=4111111111111111',
    '/api/card-stats-transfer/yearly-transfer-receiver-amount-by-card?year=2025&card_number=4111111111111111',
  ];

  for (let endpoint of basicEndpoints) {
    let res = http.get(`${BASE_URL}${endpoint}`, params);
    check(res, { [`GET ${endpoint} success`]: (r) => r.status === 200 });
  }

  sleep(0.1);
}
